---
description: 과거 commit 이력을 분석하여 이전 작업 맥락을 복원하고 이어서 작업할 수 있도록 컨텍스트를 제공하는 커맨드. 자연어 시간 표현("어제", "지난주", 날짜) 또는 commit hash를 인수로 받는다.
---

# 과거 작업 복원 (git-resume)

git commit 이력과 diff를 분석하여 이전 작업 맥락을 재구성하고, 작업을 자연스럽게 이어서 진행할 수 있도록 컨텍스트를 제공함.

## 사용법

```
/git-resume                      오늘·어제 커밋 기준 작업 복원 (기본)
/git-resume yesterday            어제 커밋 기반 작업 복원
/git-resume "2 days ago"         2일 전 커밋 기반 작업 복원
/git-resume "last week"          지난주 커밋 기반 작업 복원
/git-resume 2026-05-20           특정 날짜 커밋 기반 작업 복원
/git-resume <commit-hash>        특정 커밋 기반 작업 복원
/git-resume <hash1>..<hash2>     커밋 범위 기반 작업 복원
```

자연어 시간 표현도 그대로 사용할 수 있음:

```
/git-resume 어제
/git-resume 지난주
/git-resume 3일 전
/git-resume 지난 금요일
```

---

## Phase 0: 사전 확인

### 0-1. git 저장소 확인

```bash
git rev-parse --is-inside-work-tree 2>/dev/null
```

실패하면 아래를 출력하고 종료함.

```
오류: 현재 디렉터리가 git 저장소가 아닙니다.
  해결: 프로젝트 루트로 이동 후 재실행
```

### 0-2. 인수 파싱

인수를 아래 규칙으로 해석함.

| 인수 형태                         | 해석 방식                                            |
| --------------------------------- | ---------------------------------------------------- |
| 없음                              | `--since="1 day ago"` (오늘 + 어제)                  |
| `어제` / `yesterday`              | `--since="yesterday" --until="today"`                |
| `지난주` / `last week`            | `--since="1 week ago"`                               |
| `N일 전` / `N days ago`           | `--since="N days ago"`                               |
| `YYYY-MM-DD`                      | `--since="YYYY-MM-DD" --until="YYYY-MM-DD 23:59:59"` |
| `지난 요일명` / `last <weekday>`  | 해당 요일의 시작~종료                                |
| 7자 이상 hex 문자열 (commit hash) | 단일 커밋 `git show <hash>`                          |
| `<hash1>..<hash2>`                | 커밋 범위 `git log <hash1>..<hash2>`                 |

**검증 (필수)**: 인수는 셸 명령에 치환되므로 위 표의 형태만 허용함. commit hash는 `^[0-9a-f]{7,40}(\.\.[0-9a-f]{7,40})?$` 정규식에 일치해야 하고, 날짜는 `YYYY-MM-DD` 형식만, 자연어는 아래 매핑 표의 고정 문자열로만 변환함. 어느 형태에도 해당하지 않는 인수는 git 명령에 넣지 말고 사용법을 출력한 후 종료함.

한국어 자연어는 아래와 같이 변환함:

| 한국어 표현   | git 인수                                                     |
| ------------- | ------------------------------------------------------------ |
| 어제          | `--since="yesterday" --until="today"`                        |
| 그제 / 그저께 | `--since="2 days ago" --until="1 day ago"`                   |
| 지난주        | `--since="1 week ago"`                                       |
| 이번 주       | `--since="last monday"`                                      |
| N일 전        | N을 숫자로 추출 → `--since="N days ago"`                     |
| 지난 월요일   | `--since="last monday" --until="last monday 23:59:59"`       |
| 지난 화요일   | `--since="last tuesday" --until="last tuesday 23:59:59"`     |
| 지난 수요일   | `--since="last wednesday" --until="last wednesday 23:59:59"` |
| 지난 목요일   | `--since="last thursday" --until="last thursday 23:59:59"`   |
| 지난 금요일   | `--since="last friday" --until="last friday 23:59:59"`       |
| 지난 토요일   | `--since="last saturday" --until="last saturday 23:59:59"`   |
| 지난 일요일   | `--since="last sunday" --until="last sunday 23:59:59"`       |

---

## Phase 1: commit 이력 조회

### 1-1. 현재 사용자 확인

```bash
git config user.name
git config user.email
```

### 1-2. 해당 기간 커밋 목록 조회

```bash
git log <since/until 또는 hash 범위> \
  --author="$(git config user.name)" \
  --oneline \
  --no-merges \
  --decorate
```

커밋이 없으면:

```
해당 기간에 커밋이 없습니다.
  - 다른 기간을 지정하거나
  - /git-resume last week 으로 범위를 넓혀보세요.
```

커밋이 20개를 초과하면 사용자에게 확인:

```
<N>개의 커밋이 발견됐습니다. 최근 20개만 분석할까요? [Y/n]
```

---

## Phase 2: commit별 상세 분석

발견된 각 커밋에 대해 아래를 실행함.

### 2-1. 커밋 상세 조회

```bash
git show <hash> --stat --no-patch
```

변경 파일 목록과 삽입/삭제 줄 수를 파악함.

### 2-2. diff 내용 조회

```bash
git show <hash> -p --unified=3
```

diff가 500줄을 초과하면 파일별로 분할하여 읽음:

```bash
git show <hash> -- <file> -p
```

### 2-3. 커밋 분류

각 커밋을 아래 유형으로 분류함.

| 유형        | 판단 기준                                          |
| ----------- | -------------------------------------------------- |
| 신규 기능   | `feat:` prefix 또는 새 파일 추가 비중이 높음       |
| 버그 수정   | `fix:` prefix 또는 기존 로직 수정                  |
| 문서 작업   | `docs:` prefix 또는 `.md` 파일만 변경              |
| 리팩터링    | `refactor:` prefix                                 |
| 미완성(WIP) | `wip`, `WIP`, `tmp`, `temp`, `[skip`, `draft` 포함 |
| 설정 변경   | `chore:`, `ci:`, `build:` prefix                   |

---

## Phase 3: 작업 컨텍스트 분석

수집된 커밋 데이터를 종합하여 아래 항목을 분석함.

### 3-1. 핵심 작업 파악

- 가장 많이 변경된 파일/디렉터리 → 작업의 중심 영역
- 커밋 메시지 패턴 → 어떤 기능/이슈를 다루고 있었는지
- 추가된 함수·클래스·컴포넌트 → 새로 만들기 시작한 것
- 삭제된 코드 → 무엇을 제거하거나 교체하고 있었는지

### 3-2. 미완성 작업 탐지

아래 패턴이 있으면 미완성 작업으로 표시함.

```bash
# TODO/FIXME/HACK 추가 여부
git show <hash> -p | grep -E "^\+.*(TODO|FIXME|HACK|XXX|WIP)"

# 주석 처리된 코드 추가
git show <hash> -p | grep -E "^\+\s*//"
```

마지막 커밋 이후 스테이징되지 않은 변경사항 확인:

```bash
git status --short
git diff --stat
```

### 3-3. 다음 작업 단계 추론

커밋 흐름과 마지막 상태를 바탕으로 아래를 추론함.

- 마지막으로 작업 중이던 파일
- 커밋 메시지가 암시하는 다음 단계
- 미완성으로 표시된 항목

---

## Phase 4: 작업 재개 요약 출력

아래 형식으로 분석 결과를 출력함.

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  작업 복원 — <기간 또는 커밋 범위>
  저장소: <현재 디렉터리>  |  브랜치: <현재 브랜치>
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

## 커밋 이력 (<N>개)

  <hash> <날짜 시각> <커밋 메시지>
  <hash> <날짜 시각> <커밋 메시지>
  ...

## 작업 요약

  주요 작업 영역:  <디렉터리/파일 목록>
  작업 유형:       <신규 기능 / 버그 수정 / 리팩터링 / 문서 / 혼합>
  변경 규모:       <+N줄 / -N줄, N개 파일>

## 무엇을 하고 있었나

  <커밋 메시지·diff 분석 기반 자연어 설명 3~5문장>
  예) "UserService에 OAuth 로그인 기능을 추가하고 있었습니다.
      refreshToken 처리까지 완료했고, 마지막 커밋에서
      access token 만료 처리 로직을 작성 중이었습니다."

## 미완성 항목

  <없으면 이 섹션 생략>
  - TODO: <내용> — <파일:줄번호>
  - WIP: <내용>
  - 스테이징 안 된 변경: <파일 목록>

## 권장 다음 단계

  1. <구체적인 다음 작업>
  2. <그 다음 작업>
  3. <확인이 필요한 사항>

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

출력 후 사용자에게 물음:

```
어떤 작업부터 이어서 진행할까요?
  [1] <권장 다음 단계 1>
  [2] <권장 다음 단계 2>
  [3] 직접 입력
```

사용자 선택에 따라 해당 작업을 바로 시작함.

---

## 특수 케이스 처리

### 커밋이 매우 많은 경우

20개를 초과하면 최근 20개만 분석하되, 나머지 커밋의 메시지 목록만 별도로 보여줌.

### 바이너리 파일 변경

이미지·DB 파일 등 바이너리 변경은 파일명과 크기 변화만 표시하고 diff는 생략함.

### merge 커밋

`--no-merges` 플래그로 기본적으로 제외함. 명시적으로 merge commit을 포함하려면 `--merges` 인수를 추가함.

### 다른 작성자의 커밋 포함

기본은 `git config user.name` 기준이지만, `--all-authors` 인수를 주면 모든 작성자의 커밋을 포함함.

---

## 활용 예시

### 예시 1 — 어제 하던 작업 이어서

```
/git-resume yesterday
```

→ 어제의 커밋 분석 → "UserService OAuth 구현 중이었음" → 다음 단계 제안

### 예시 2 — 특정 날짜로 돌아가기

```
/git-resume 2026-05-19
```

→ 해당 날짜 커밋 분석 → 그 시점에 무엇을 하고 있었는지 파악

### 예시 3 — 특정 커밋부터 이어서

```
/git-resume abc1234
```

→ 해당 커밋의 변경 내용 분석 → 그 작업의 맥락 파악

### 예시 4 — 지난주 전체 작업 흐름 파악

```
/git-resume last week
```

→ 지난주 커밋 전체 분석 → 주간 작업 흐름 요약 → 이번 주 이어서 할 작업 제안
