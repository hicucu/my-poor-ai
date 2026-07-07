---
description: 코드 그래프 도구(graphifyy 또는 codegraph)를 선택하여 설치·설정하는 커맨드. 패키지 설치, 코드베이스 그래프 생성, Claude Code 통합, git hook 등록, .gitignore 설정까지 원스톱으로 수행.
---

# 코드 그래프 도구 자동 설치 및 설정

graphifyy(pip) 또는 codegraph(npm) 중 하나를 선택하여 전 과정을 자동 수행함.
이미 완료된 단계는 건너뛰므로 재실행해도 안전함.

## 사용법

```
/graphify-setup                전체 설치·설정 — 도구 선택 포함 (기본)
/graphify-setup graphify       graphifyy(pip) 도구로 설치
/graphify-setup codegraph      codegraph(npm) 도구로 설치
/graphify-setup check          현재 설정 상태 확인 (두 도구 모두)
/graphify-setup reset          이미 완료된 단계도 모두 재실행
```

도구 인수와 모드 인수는 조합 가능함:

```
/graphify-setup graphify check     graphifyy 설정 상태만 확인
/graphify-setup codegraph reset    codegraph 전체 강제 재설치
/graphify-setup graphify reset     graphifyy 전체 강제 재설치
```

---

## Phase 0: 사전 환경 확인 및 도구 선택

### 0-1. git 저장소 확인

```bash
git rev-parse --is-inside-work-tree 2>/dev/null
```

실패(exit code ≠ 0)이면 아래를 출력하고 종료함.

```
오류: 현재 디렉터리가 git 저장소가 아닙니다.
git hook 등록은 git 저장소에서만 동작합니다.
  해결: git init 후 재실행하거나, 저장소 루트로 이동 후 재실행
```

### 0-2. 패키지 매니저 확인

아래 순서로 사용 가능한 명령을 탐지함.

```bash
# Python / pip
python --version 2>/dev/null || python3 --version 2>/dev/null
pip --version 2>/dev/null || pip3 --version 2>/dev/null

# Node / npm
node --version 2>/dev/null
npm --version 2>/dev/null
```

이후 단계에서 `pip` 표기는 실제 탐지된 명령(`pip` 또는 `pip3`)으로 대체함.

### 0-3. 도구 선택

인수에 `graphify` 또는 `codegraph`가 명시된 경우 해당 도구를 사용함.

인수가 없는 경우 아래 순서로 자동 감지함.

```bash
# graphifyy 설치 여부
pip show graphifyy 2>/dev/null

# codegraph 설치 여부
codegraph --version 2>/dev/null
```

| 감지 결과          | 처리                                   |
| ------------------ | -------------------------------------- |
| graphifyy만 설치됨 | graphifyy 자동 선택, 안내 출력 후 진행 |
| codegraph만 설치됨 | codegraph 자동 선택, 안내 출력 후 진행 |
| 둘 다 설치됨       | 사용자에게 선택 요청                   |
| 둘 다 미설치       | 사용자에게 선택 요청                   |

**사용자 선택 요청 메시지:**

```
코드 그래프 도구를 선택하세요:

  [1] graphifyy  (pip 패키지, Python 기반)
                 출력 디렉터리: graphify-out/
                 Claude 통합:   CLAUDE.md + .claude/settings.json 직접 수정

  [2] codegraph  (npm 패키지, Node 기반)
                 출력 디렉터리: .codegraph/
                 Claude 통합:   MCP 서버 등록 (.claude/settings.json)
                 추가 기능:     codegraph_search, codegraph_context 등 6개 MCP 도구

선택 (1 또는 2):
```

이후 모든 Phase는 선택된 도구(`TOOL = graphify | codegraph`)에 따라 분기하여 실행함.

### 0-4. 현재 설정 상태 파악

선택된 `TOOL`에 따라 아래 상태 변수를 기록함.

**graphifyy 상태 변수:**

| 상태 변수           | 확인 명령                                       | true 조건 |
| ------------------- | ----------------------------------------------- | --------- |
| `installed`         | `pip show graphifyy 2>/dev/null`                | exit 0    |
| `graph_exists`      | `test -d graphify-out`                          | exit 0    |
| `claude_integrated` | `grep -q "graphify" CLAUDE.md 2>/dev/null`      | exit 0    |
| `hook_installed`    | `test -f .git/hooks/post-commit`                | exit 0    |
| `gitignore_set`     | `grep -q "graphify-out" .gitignore 2>/dev/null` | exit 0    |

**codegraph 상태 변수:**

| 상태 변수           | 확인 명령                                                                                                                | true 조건 |
| ------------------- | ------------------------------------------------------------------------------------------------------------------------ | --------- |
| `installed`         | `codegraph --version 2>/dev/null`                                                                                        | exit 0    |
| `graph_exists`      | `test -d .codegraph`                                                                                                     | exit 0    |
| `claude_integrated` | `grep -q '"codegraph"' ~/.claude/settings.json 2>/dev/null \|\| grep -q '"codegraph"' .claude/settings.json 2>/dev/null` | exit 0    |
| `hook_installed`    | `grep -q "codegraph sync" .git/hooks/post-commit 2>/dev/null`                                                            | exit 0    |
| `gitignore_set`     | `grep -q "\.codegraph" .gitignore 2>/dev/null`                                                                           | exit 0    |

`claude_integrated` 판정 시 전역(`~/.claude/settings.json`) → 프로젝트(`.claude/settings.json`) 순으로 확인하며 어느 쪽이든 등록되어 있으면 `true`로 처리함.

`check` 인수로 실행한 경우: 선택된 도구 형식으로 상태를 출력하고 종료함.

**graphifyy check 출력:**

```
graphifyy 설정 현황 — <현재 디렉터리>

  패키지 설치 (pip)          ✓ 완료   /  ✗ 미설치
  코드베이스 그래프          ✓ 존재   /  ✗ 없음  (graphify-out/)
  Claude Code 통합           ✓ 완료   /  ✗ 미설정
  git hook (post-commit)     ✓ 등록   /  ✗ 없음
  .gitignore graphify 항목   ✓ 있음   /  ✗ 없음
```

**codegraph check 출력:**

```
codegraph 설정 현황 — <현재 디렉터리>

  패키지 설치 (npm)          ✓ 완료   /  ✗ 미설치
  코드베이스 그래프          ✓ 존재   /  ✗ 없음  (.codegraph/)
  Claude MCP 통합            ✓ 완료 (전역)  /  ✓ 완료 (프로젝트)  /  ✗ 미설정
  git hook (post-commit)     ✓ 등록   /  ✗ 없음
  .gitignore codegraph 항목  ✓ 있음   /  ✗ 없음
```

`check` 인수 없이 `reset` 인수만 있는 경우: 모든 상태 변수를 `false`로 재설정하여 전체를 재실행함.

---

## Phase 1: 패키지 설치

`installed = false`인 경우에만 실행함.

**graphifyy:**

pip 또는 pip3가 없으면 아래를 출력하고 종료함.

```
오류: Python 또는 pip를 찾을 수 없습니다.
Python 3.x 설치 후 재실행하세요.
```

```bash
pip install graphifyy
```

성공 여부 확인:

```bash
graphify --version
```

**codegraph:**

npm이 없으면 아래를 출력하고 종료함.

```
오류: npm을 찾을 수 없습니다.
Node.js 설치 후 재실행하세요: https://nodejs.org
```

```bash
npm install -g @colbymchenry/codegraph
```

성공 여부 확인:

```bash
codegraph --version
```

**공통 — 설치 성공:**

```
✓ <도구명> 설치 완료
```

실패 시 에러 메시지를 그대로 출력하고 종료함.

`installed = true`이면:

```
— 패키지 설치: 이미 완료, 건너뜀
```

---

## Phase 2: 코드베이스 그래프 생성

`graph_exists = false` 또는 `reset` 모드인 경우 실행함.

**graphifyy:**

```bash
graphify update .
```

완료 후 `graphify-out/GRAPH_REPORT.md` 존재로 성공 확인함.

```
✓ 코드베이스 그래프 생성 완료 — graphify-out/GRAPH_REPORT.md
```

**codegraph:**

```bash
codegraph init .
codegraph index .
```

완료 후 `.codegraph/codegraph.db` 존재로 성공 확인함.

```bash
test -f .codegraph/codegraph.db
```

```
✓ 코드베이스 그래프 생성 완료 — .codegraph/codegraph.db
```

실패 시 에러를 출력하되 이후 Phase는 계속 진행함.

`graph_exists = true`이면:

- **graphifyy:** `— 그래프 생성: 이미 완료, 건너뜀  (재빌드 필요 시: graphify update .)`
- **codegraph:** `— 그래프 생성: 이미 완료, 건너뜀  (재빌드 필요 시: codegraph index . --force)`

---

## Phase 3: Claude Code 통합

`claude_integrated = false`인 경우에만 실행함.

**graphifyy:**

```bash
graphify claude install
```

성공 여부를 아래 두 조건으로 확인함.

```bash
grep -q "graphify" CLAUDE.md 2>/dev/null
grep -q "graphify" .claude/settings.json 2>/dev/null
```

```
✓ Claude Code 통합 완료 — CLAUDE.md, .claude/settings.json 업데이트
```

**codegraph:**

먼저 전역 설정에 codegraph MCP가 이미 등록되어 있는지 확인함.

```bash
# 전역 settings.json 위치 (OS별)
# Windows: %USERPROFILE%\.claude\settings.json
# macOS/Linux: ~/.claude/settings.json
python -c "
import json, os, sys
path = os.path.join(os.path.expanduser('~') or os.environ.get('USERPROFILE',''), '.claude', 'settings.json')
try:
    data = json.load(open(path, encoding='utf-8'))
    print('found' if 'codegraph' in str(data.get('mcpServers', {})) else 'not_found')
except: print('not_found')
"
```

| 확인 결과          | 처리                                                            |
| ------------------ | --------------------------------------------------------------- |
| 전역에 이미 등록됨 | 프로젝트 settings.json 등록 생략, 아래 메시지 출력 후 완료 처리 |
| 전역에 없음        | 프로젝트 `.claude/settings.json`에 추가                         |

**전역에 이미 등록된 경우 출력:**

```
— Claude MCP 통합: 전역(~/.claude/settings.json)에 이미 등록됨, 프로젝트 설정 생략
```

**전역에 없는 경우** — 프로젝트 `.claude/settings.json`에 추가함.

1. `.claude/settings.json`이 없으면 `{}` 내용으로 새로 생성함.
2. 기존 `mcpServers` 섹션이 없으면 추가함.
3. `codegraph` 키가 이미 `mcpServers`에 있으면 건너뜀.
4. 아래 항목을 `mcpServers`에 추가함:

```json
{
  "mcpServers": {
    "codegraph": {
      "command": "codegraph",
      "args": ["serve", "--mcp"]
    }
  }
}
```

성공 여부 확인:

```bash
grep -q '"codegraph"' .claude/settings.json 2>/dev/null
```

```
✓ Claude MCP 통합 완료 — 프로젝트 .claude/settings.json에 codegraph MCP 서버 등록
   Claude Code 재시작 후 codegraph_search, codegraph_context 등 6개 도구 활성화
```

실패 시 에러를 출력하고 이후 Phase는 계속 진행함.

`claude_integrated = true`이면 전역 확인 없이 바로:

```
— Claude Code 통합: 이미 완료, 건너뜀
```

---

## Phase 4: git hook 등록

`hook_installed = false`인 경우에만 실행함.

**graphifyy:**

```bash
graphify hook install
```

성공 여부: `.git/hooks/post-commit` 존재 확인

```
✓ git hook 등록 완료 — post-commit, post-checkout
   팀원도 각자 로컬에서 'graphify hook install'을 실행해야 합니다.
```

**codegraph:**

`.git/hooks/post-commit` 파일이 없으면 shebang 포함 새 파일로 생성, 있으면 기존 내용을 유지하면서 아래 블록을 파일 끝에 추가함.

```bash
#!/bin/sh
# codegraph incremental sync
codegraph sync . --quiet 2>/dev/null || true
```

파일이 새로 생성된 경우 실행 권한을 부여함:

```bash
chmod +x .git/hooks/post-commit
```

성공 여부: `.git/hooks/post-commit`에 `codegraph sync` 문자열 포함 여부 확인

```
✓ git hook 등록 완료 — post-commit에 codegraph sync 추가
   팀원도 각자 로컬에서 '/graphify-setup codegraph'를 실행해야 합니다.
```

실패 시 에러를 출력하고 이후 Phase는 계속 진행함.

`hook_installed = true`이면:

```
— git hook 등록: 이미 완료, 건너뜀
```

---

## Phase 5: .gitignore 설정

`gitignore_set = false`인 경우에만 실행함.

### 5-1. .gitignore 존재 여부 확인

```bash
test -f .gitignore
```

파일이 없으면 빈 파일로 생성함.

### 5-2. 도구별 항목 추가

**graphifyy:**

`graphify-out/*` 또는 `graphify-out/` 항목이 이미 있으면 추가하지 않음.

```
# graphifyy 산출물 (GRAPH_REPORT.md만 추적)
graphify-out/*
!graphify-out/GRAPH_REPORT.md
```

**codegraph:**

`.codegraph` 항목이 이미 있으면 추가하지 않음.

```
# codegraph 산출물
.codegraph/
```

성공:

```
✓ .gitignore 항목 추가 완료
```

`gitignore_set = true`이면:

```
— .gitignore 설정: 이미 완료, 건너뜀
```

---

## Phase 6: 완료 요약

모든 Phase 실행 후 선택된 도구에 따라 아래 형식으로 요약을 출력함.

**graphifyy 완료 요약:**

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  graphifyy 설정 완료 — <현재 디렉터리>
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  패키지 설치        <결과>
  그래프 생성        <결과>
  Claude Code 통합   <결과>
  git hook 등록      <결과>
  .gitignore 설정    <결과>

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

팀원 온보딩 (각자 로컬에서 1회 실행):
  pip install graphifyy && graphify update . && graphify hook install

수동 그래프 재빌드:  graphify update .
재빌드 로그:         ~/.cache/graphify-rebuild.log
```

**codegraph 완료 요약:**

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  codegraph 설정 완료 — <현재 디렉터리>
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  패키지 설치        <결과>
  그래프 생성        <결과>
  Claude MCP 통합    <결과>
  git hook 등록      <결과>
  .gitignore 설정    <결과>

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

팀원 온보딩 (각자 로컬에서 1회 실행):
  npm install -g @colbymchenry/codegraph && codegraph init . && codegraph index .

수동 그래프 재빌드:  codegraph index . --force
증분 업데이트:       codegraph sync .
그래프 상태 확인:    codegraph status .
```

`<결과>` 자리에 해당 단계의 실행 결과를 기입함.

| 상황                 | 표기                    |
| -------------------- | ----------------------- |
| 이번 실행에서 완료   | `✓ 완료`                |
| 이미 설정되어 건너뜀 | `— 기존 유지`           |
| 오류 발생            | `✗ 실패 (위 오류 확인)` |

오류가 발생한 단계가 하나라도 있으면 요약 아래에 추가함.

```
주의: 일부 단계에서 오류가 발생했습니다.
      위 오류 메시지를 확인한 후 해당 명령을 수동으로 재실행하세요.
```
