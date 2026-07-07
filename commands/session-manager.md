---
description: 로컬 Claude Code 세션 전체를 최대 10개 서브에이전트로 병렬 분석하여 목록 조회·이름 변경·삭제를 수행하는 세션 관리 커맨드
model: haiku
---

# 세션 관리자 (Session Manager)

`~/.claude/projects/` 하위의 모든 JSONL 세션을 최대 10개 서브에이전트로 병렬 분석함.

## 사용법

```
/session-manager                              세션 목록 표시 (최신 50개)
/session-manager all                          세션 목록 전체 표시
/session-manager project <이름>              특정 프로젝트 세션만 표시
/session-manager rename <번호> <새 이름>     세션 이름 변경
/session-manager delete <번호...>            번호 기반 삭제 (다중 지원)
/session-manager delete <조건>              조건 기반 삭제
```

삭제 예시:

- `/session-manager delete 1 3 5` — 1, 3, 5번 삭제
- `/session-manager delete 2-7` — 2~7번 범위 삭제
- `/session-manager delete 1 3 5-8 12` — 혼합
- `/session-manager delete 교환이 5번보다 적은 세션` — 조건 기반
- `/session-manager delete 2주 이상 된 세션` — 날짜 조건
- `/session-manager delete web-mono 프로젝트 세션` — 프로젝트 조건

---

## 권한 정책

**원칙**: 파일 삭제·내용 변경 등 되돌리기 어려운 작업이 아닌 모든 작업은 사용자 확인 없이 자동 허가함.

자동 허가 (목록 조회·분석에 필요한 모든 읽기·실행):

- 파일 읽기 (Read 도구, `cat` / `grep` / `find` / `head` / `tail` / `wc` / Python `os.stat` 등)
- Bash / PowerShell 실행 — 읽기·집계·반복 목적
- 서브에이전트 병렬 실행
- `~/.claude/.session-names.json` 읽기 및 쓰기 (이름 변경 전용)
- Python / 파이프 조합 등 데이터 추출·가공

확인 필수 (되돌리기 불가능한 작업):

- JSONL 파일 삭제 (`rm`, `Remove-Item`, `unlink` 등 모든 삭제 수단)
- 확인(`y`/`yes`) 없이는 어떤 경우에도 파일을 삭제하지 않음

---

## Phase 0: 입력 파싱

커맨드 호출 시 전달된 인수를 파싱함.

| 인수 패턴              | 동작                                         |
| ---------------------- | -------------------------------------------- |
| 없음                   | Phase 1 → Phase 2 (목록 표시, 기본 50개)     |
| `all`                  | Phase 1 → Phase 2 (전체 표시)                |
| `project <이름>`       | Phase 1 → Phase 2 (프로젝트 필터)            |
| `rename <번호> <이름>` | Phase 1 → Phase 3-A (목록 분석 후 이름 변경) |
| `delete <번호/조건>`   | Phase 1 → Phase 3-B (목록 분석 후 삭제)      |

인수 없이 호출(`/session-manager`)하면 목록을 먼저 표시하고 사용자 추가 입력을 기다림.
인수와 함께 호출하면 Phase 1 분석 후 해당 작업을 바로 실행함.

---

## Phase 1: 세션 파일 수집

아래 명령으로 모든 JSONL 파일 경로와 파일 통계를 수집함.

```bash
# OS 무관 Python 방식 (권장) — Linux/macOS/Windows 모두 동작
python3 -c "
import os, json, glob
pattern = os.path.expanduser('~/.claude/projects/**/*.jsonl')
files = glob.glob(pattern, recursive=True)
results = []
for p in files:
    try:
        s = os.stat(p)
        results.append({'path': p, 'ctime': s.st_ctime, 'mtime': s.st_mtime, 'size': s.st_size})
    except Exception as e:
        results.append({'path': p, 'error': str(e)})
results.sort(key=lambda x: x.get('mtime', 0), reverse=True)
for r in results:
    print(json.dumps(r))
"
```

각 결과 필드: `path` = 절대 경로, `ctime` = 생성 시각(epoch, `st_ctime`), `mtime` = 마지막 수정 시각(epoch), `size` = 파일 크기(bytes)

- `st_ctime`이 생성 시각을 반환하지 않는 환경(Linux ext4 등)에서는 `mtime`으로 대체
- `error` 필드가 있는 항목은 접근 실패로 간주하고 목록에 `(접근 오류)` 표기

수집 결과에서:

- 총 세션 수 N을 계산함
- 파일 경로에서 **프로젝트명** 추출: `~/.claude/projects/<project>/uuid.jsonl` → `<project>`
- 프로젝트명 디코드 (표시용): `D--example-my-project` → `example-my-project` (드라이브 접두사 `X--` 제거)
- 파일을 **최신 수정일 내림차순** 정렬 후 인덱스 번호(1-based) 부여
- 커스텀 이름 파일 `~/.claude/.session-names.json` 존재 시 읽어 메모리에 보관

---

## Phase 2: 병렬 세션 분석

전체 N개 파일을 **최대 10개 배치**로 균등 분할함 (배치당 ceil(N/10)개).

**필수: 아래 서브에이전트 호출을 단일 응답에 동시에 포함함. 각 Agent 호출 시 `model: "haiku"` 를 지정함.**

각 서브에이전트 프롬프트 (배치 파일 경로 목록을 채워 전달):

```
다음 JSONL 파일 목록을 분석하여 각 세션의 메타데이터를 JSON 배열로 반환하세요.
Bash 도구만 사용하세요. 읽기·grep·파이프 명령은 모두 자동 허가됩니다.

파일 목록:
[배치 파일 경로 목록 — 절대 경로, 한 줄에 하나씩]

각 파일에 대해 아래 bash 명령으로 정보를 추출하세요:

1. 첫 번째 timestamp (생성일):
   grep -m 1 '"timestamp"' FILE | grep -oE '"[0-9]{4}-[0-9TZ:.-]+"' | tr -d '"'

2. 마지막 timestamp (마지막 작업일):
   grep '"timestamp"' FILE | tail -1 | grep -oE '"[0-9]{4}-[0-9TZ:.-]+"' | tr -d '"'

3. 사용자 메시지 수 (교환 횟수):
   grep -c '"role":"user"' FILE

4. 모든 사용자 메시지 텍스트 추출 후 요약 (작업 요약):
   grep '"role":"user"' FILE | python3 -c "
   import sys, json, re
   lines = sys.stdin.read().strip().split('\n')
   texts = []
   for line in lines:
       try:
           d = json.loads(line)
           c = d.get('message', {}).get('content', '')
           if isinstance(c, list):
               c = ' '.join(b.get('text','') for b in c if isinstance(b,dict) and b.get('type')=='text')
           t = re.sub(r'\s+', ' ', str(c)).strip()
           if t: texts.append(t[:200])
       except: pass
   combined = ' | '.join(texts[:5])
   print(combined[:300])
   "
   위 추출 결과를 바탕으로 이 세션에서 수행한 작업을 한국어로 1~2문장으로 요약하세요.
   예: "React 로그인 컴포넌트 리팩터링 및 상태 관리 분리 작업 수행. useAuth 훅으로 인증 로직 추출."

파일 경로에서:
- sessionId = 파일명(UUID, 확장자 제외)
- project = projects/ 바로 아래 디렉토리명

반환 형식 (JSON 배열만, 다른 텍스트 없이):
[
  {
    "sessionId": "00000000-0000-0000-0000-000000000000",
    "project": "D--example-my-project",
    "createdAt": "2026-05-08T01:07:13.044Z",
    "lastActiveAt": "2026-05-08T01:22:51.303Z",
    "exchanges": 15,
    "summary": "세션 관리 커맨드 제작. 병렬 서브에이전트로 JSONL 분석 및 목록·삭제·이름 변경 기능 구현.",
    "filePath": "/절대/경로/uuid.jsonl"
  }
]

오류 발생 파일: { "sessionId": "...", "project": "...", "filePath": "...", "error": "오류 내용" }
```

배치가 1개뿐이면 서브에이전트 없이 직접 처리함.

---

## Phase 3: 결과 병합 및 목록 표시

서브에이전트 결과를 수집하여 **최신 lastActiveAt 내림차순**으로 정렬함.
커스텀 이름(`~/.claude/.session-names.json`)이 있는 세션은 헤더 줄 끝에 `"이름"` 형태로 표시함.

출력 형식:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Claude 세션 관리자   (전체 N개 세션 / M개 프로젝트)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

 #   생성일       마지막 작업    교환   프로젝트
─────────────────────────────────────────────────────────────
  1  2026-05-14  2026-05-14      23   example-project
     세션 관리 커맨드 제작. 병렬 서브에이전트로 JSONL 분석 및
     목록·삭제·이름 변경 기능 구현.

  2  2026-05-13  2026-05-13      18   web-mono  "로그인 개선"
     React 로그인 컴포넌트 리팩터링 및 useAuth 훅 분리.

  3  2026-05-12  2026-05-12       4   docs
     README API 변경사항 반영 및 예시 코드 업데이트.
...
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

/session-manager rename <번호> <새 이름>   이름 변경
/session-manager delete <번호/조건>        삭제
```

- 날짜: `YYYY-MM-DD` 형식
- 분석 실패 세션: summary 자리에 `(분석 오류)` 표기, 목록에는 포함
- 기본 최신 50개 표시. `all` 전달 시 전체 / `project <이름>` 전달 시 필터

---

## Phase 4-A: 이름 변경

`rename <번호> <새 이름>` 처리:

1. `~/.claude/.session-names.json` 읽기 (없으면 `{}`)
2. 해당 번호의 sessionId 확인
3. `"<sessionId>": "<새 이름>"` 추가/업데이트
4. 들여쓰기 2칸 JSON으로 저장

출력:

```
세션 #<번호> 이름 변경 완료
  이전: "<기존 이름>"
  이후: "<새 이름>"
```

번호 범위 초과 또는 이름 누락 시 오류 출력.

---

## Phase 4-B: 삭제

### 번호 기반 삭제

`delete <번호...>` 파싱 규칙:

- 공백/쉼표 구분 다중 번호: `1 3 5`, `1,3,5`, `1, 3, 5`
- 범위: `2-7` → 2,3,4,5,6,7
- 혼합: `1 3 5-8 12`
- 중복 번호는 한 번만 처리

### 조건 기반 삭제

`delete <자연어 조건>` — 번호가 아닌 텍스트가 전달된 경우 조건으로 해석함.

지원 조건 예시:

| 조건 표현               | 필터 기준                   |
| ----------------------- | --------------------------- |
| `교환이 N번보다 적은`   | exchanges < N               |
| `교환이 N번 이하인`     | exchanges ≤ N               |
| `N주/개월/일 이상 된`   | lastActiveAt 기준 경과 기간 |
| `<프로젝트명> 프로젝트` | project 필드 포함 여부      |
| `<날짜> 이전`           | lastActiveAt < 날짜         |

조건에 해당하는 세션 목록을 먼저 표시하고 확인을 구함.

### 삭제 확인 절차 (번호·조건 공통)

삭제 대상이 확정되면 **반드시** 아래 형식으로 표시 후 확인:

```
삭제 대상 세션 (N개):

  #1  2026-05-08  2026-05-08   15교환  example-project
      세션 관리 커맨드 제작. 병렬 서브에이전트로 JSONL 분석...

  #3  2026-04-20  2026-04-21    4교환  docs
      README 업데이트 요청 — API 변경사항 반영.

위 N개 세션을 삭제하시겠습니까? 되돌릴 수 없습니다. [y/N]
```

사용자가 `y` 또는 `yes` 입력 시에만:

1. 대상 JSONL 파일 전체 삭제
2. `~/.claude/.session-names.json`에서 해당 sessionId 항목 제거
3. `N개 세션 삭제됨.` 출력
4. 아래 질문을 출력하고 사용자 입력을 기다림:

```
삭제된 세션을 제외한 목록을 다시 표시할까요? [y/N]
```

- `y` 또는 `yes` 입력 시: 이미 수집된 세션 데이터에서 삭제된 항목을 제거한 후 Phase 3 형식으로 목록을 재출력함. JSONL 파일 재스캔 없이 메모리 내 데이터를 재사용함.
- 그 외 입력: 종료

그 외 입력(`n`, `N`, `no`, Enter 등)은 `삭제 취소됨.` 출력 후 종료.

**절대 규칙: 사용자 명시적 확인(`y`/`yes`) 없이 파일을 삭제하지 않음. 조건이 자명해 보여도 예외 없음.**
