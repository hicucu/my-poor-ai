---
name: using-git-worktrees
description: Use when starting feature work that needs isolation from the current workspace, or before executing an implementation plan - ensures an isolated workspace via native tooling or a git worktree fallback
---

# Git Worktree 사용

## 개요

작업이 격리된 작업 공간에서 진행되도록 보장. 플랫폼의 네이티브 worktree 도구를 우선 사용. 네이티브 도구가 없는 경우에만 수동 git worktree로 대체.

**핵심 원칙:** 먼저 기존 격리 환경 감지. 그 다음 네이티브 도구 사용. 그 다음 git으로 대체. 하네스와 충돌 금지.

**시작 시 공지:** "using-git-worktrees 스킬을 사용하여 격리된 작업 공간을 설정함."

## 0단계: 기존 격리 환경 감지

**아무것도 생성하기 전에, 이미 격리된 작업 공간에 있는지 확인.**

```bash
GIT_DIR=$(cd "$(git rev-parse --git-dir)" 2>/dev/null && pwd -P)
GIT_COMMON=$(cd "$(git rev-parse --git-common-dir)" 2>/dev/null && pwd -P)
BRANCH=$(git branch --show-current)
```

**서브모듈 가드:** `GIT_DIR != GIT_COMMON` 조건은 git 서브모듈 내부에서도 참. "이미 worktree에 있다"고 결론 내리기 전에 서브모듈 내부가 아닌지 확인:

```bash
# 경로가 반환되면 worktree가 아닌 서브모듈 내부 — 일반 저장소로 처리
git rev-parse --show-superproject-working-tree 2>/dev/null
```

**`GIT_DIR != GIT_COMMON` (서브모듈이 아닌 경우):** 이미 연결된 worktree에 있음. 3단계(프로젝트 설정)로 건너뛸 것. 다른 worktree를 생성하지 말 것.

브랜치 상태와 함께 보고:

- 브랜치에 있는 경우: "이미 `<path>`의 격리된 작업 공간에서 `<name>` 브랜치로 작업 중."
- Detached HEAD: "이미 `<path>`의 격리된 작업 공간에 있음 (detached HEAD, 외부 관리). 완료 시점에 브랜치 생성 필요."

**`GIT_DIR == GIT_COMMON` (또는 서브모듈 내부):** 일반 저장소 체크아웃 상태.

지시사항에서 이미 worktree 선호도를 밝혔는가? 아니라면 worktree 생성 전 동의 요청:

> "격리된 worktree를 설정할까요? 현재 브랜치를 변경으로부터 보호함."

기존에 선언된 선호도가 있으면 묻지 않고 따를 것. 사용자가 동의를 거부하면 현재 위치에서 작업하고 3단계로 건너뜀.

## 1단계: 격리된 작업 공간 생성

**두 가지 방법이 있음. 이 순서대로 시도.**

### 1a. 네이티브 Worktree 도구 (권장)

사용자가 격리된 작업 공간을 요청했음 (0단계 동의). worktree를 생성하는 방법이 이미 있는가? `EnterWorktree`, `WorktreeCreate`, `/worktree` 명령어, 또는 `--worktree` 플래그 같은 이름의 도구일 수 있음. 있다면 그것을 사용하고 3단계로 건너뜀.

네이티브 도구는 디렉토리 배치, 브랜치 생성, 정리를 자동으로 처리. 네이티브 도구가 있는데 `git worktree add`를 사용하면 하네스가 볼 수 없거나 관리할 수 없는 유령 상태가 생성됨.

네이티브 worktree 도구가 없는 경우에만 1b단계로 진행.

### 1b. Git Worktree 대체

**1a단계가 해당되지 않는 경우에만 사용** — 네이티브 worktree 도구가 없을 때. git을 사용하여 수동으로 worktree 생성.

#### 디렉토리 선택

다음 우선순위를 따를 것. 명시적인 사용자 선호도는 항상 파일시스템 상태보다 우선.

1. **지시사항에서 선언된 worktree 디렉토리 선호도 확인.** 사용자가 이미 지정했다면 묻지 않고 사용.

2. **기존 프로젝트 로컬 worktree 디렉토리 확인:**

   ```bash
   ls -d .worktrees 2>/dev/null     # 권장 (숨김)
   ls -d worktrees 2>/dev/null      # 대안
   ```

   발견되면 사용. 둘 다 있으면 `.worktrees`가 우선.

3. **기존 전역 디렉토리 확인:**

   ```bash
   project=$(basename "$(git rev-parse --show-toplevel)")
   ls -d ~/.config/my-poor-ai/worktrees/$project 2>/dev/null
   ```

   발견되면 사용 (레거시 전역 경로와의 하위 호환성).

4. **다른 지침이 없다면** 프로젝트 루트의 `.worktrees/`를 기본값으로 사용.

#### 안전 확인 (프로젝트 로컬 디렉토리만)

**worktree 생성 전 디렉토리가 무시되는지 반드시 확인:**

```bash
git check-ignore -q .worktrees 2>/dev/null || git check-ignore -q worktrees 2>/dev/null
```

**무시되지 않는 경우:** .gitignore에 추가하고, 변경 사항을 커밋한 뒤 진행.

**중요한 이유:** worktree 내용이 실수로 저장소에 커밋되는 것을 방지.

전역 디렉토리(`~/.config/my-poor-ai/worktrees/`)는 확인 불필요.

#### Worktree 생성

```bash
project=$(basename "$(git rev-parse --show-toplevel)")

# 선택한 위치에 따라 경로 결정
# 프로젝트 로컬: path="$LOCATION/$BRANCH_NAME"
# 전역: path="~/.config/my-poor-ai/worktrees/$project/$BRANCH_NAME"

git worktree add "$path" -b "$BRANCH_NAME"
cd "$path"
```

**샌드박스 대체:** `git worktree add`가 권한 오류(샌드박스 거부)로 실패하면, 샌드박스가 worktree 생성을 차단했으며 현재 디렉토리에서 작업할 것임을 사용자에게 알릴 것. 그런 다음 현재 위치에서 설정 및 기준 테스트를 실행.

## 3단계: 프로젝트 설정

적절한 설정을 자동 감지하여 실행:

```bash
# Node.js
if [ -f package.json ]; then npm install; fi

# Rust
if [ -f Cargo.toml ]; then cargo build; fi

# Python
if [ -f requirements.txt ]; then pip install -r requirements.txt; fi
if [ -f pyproject.toml ]; then poetry install; fi

# Go
if [ -f go.mod ]; then go mod download; fi
```

## 4단계: 깨끗한 기준 확인

작업 공간이 깨끗한 상태로 시작하는지 확인하기 위해 테스트 실행:

```bash
# 프로젝트에 적합한 명령어 사용
npm test / cargo test / pytest / go test ./...
```

**테스트 실패 시:** 실패 보고, 계속 진행할지 또는 조사할지 질문.

**테스트 통과 시:** 준비 완료 보고.

### 보고

```
Worktree 준비 완료: <전체-경로>
테스트 통과 (<N>개 테스트, 실패 0건)
<기능명> 구현 준비 완료
```

## 빠른 참조

| 상황                         | 조치                                             |
| ---------------------------- | ------------------------------------------------ |
| 이미 연결된 worktree에 있음  | 생성 건너뜀 (0단계)                              |
| 서브모듈 내부                | 일반 저장소로 처리 (0단계 guard)                 |
| 네이티브 worktree 도구 있음  | 사용 (1a단계)                                    |
| 네이티브 도구 없음           | Git worktree 대체 (1b단계)                       |
| `.worktrees/` 존재           | 사용 (무시 여부 확인)                            |
| `worktrees/` 존재            | 사용 (무시 여부 확인)                            |
| 둘 다 존재                   | `.worktrees/` 사용                               |
| 둘 다 없음                   | 지시사항 파일 확인, 그 다음 기본값 `.worktrees/` |
| 전역 경로 존재               | 사용 (하위 호환성)                               |
| 디렉토리가 무시되지 않음     | .gitignore에 추가 + 커밋                         |
| 생성 시 권한 오류            | 샌드박스 대체, 현재 위치에서 작업                |
| 기준 테스트 실패             | 실패 보고 + 질문                                 |
| package.json/Cargo.toml 없음 | 의존성 설치 건너뜀                               |

## 흔한 실수

### 하네스와 충돌

- **문제:** 플랫폼이 이미 격리를 제공하는데 `git worktree add` 사용
- **해결:** 0단계에서 기존 격리 감지. 1a단계에서 네이티브 도구에 위임.

### 감지 건너뛰기

- **문제:** 기존 worktree 내부에 중첩된 worktree 생성
- **해결:** 아무것도 생성하기 전에 항상 0단계 실행

### 무시 확인 건너뛰기

- **문제:** worktree 내용이 추적되어 git status 오염
- **해결:** 프로젝트 로컬 worktree 생성 전 항상 `git check-ignore` 사용

### 디렉토리 위치 추정

- **문제:** 불일치 발생, 프로젝트 관례 위반
- **해결:** 우선순위 따르기: 기존 > 전역 레거시 > 지시사항 파일 > 기본값

### 실패하는 테스트로 진행

- **문제:** 새로운 버그와 기존 문제를 구별할 수 없음
- **해결:** 실패 보고, 진행 허가 명시적으로 요청

## 위험 신호

**절대 금지:**

- 0단계에서 기존 격리가 감지되면 worktree 생성
- 네이티브 worktree 도구(예: `EnterWorktree`)가 있는데 `git worktree add` 사용. 이것이 가장 흔한 실수 — 있으면 사용할 것.
- 1a단계를 건너뛰고 바로 1b단계의 git 명령어 실행
- 무시 여부 확인 없이 worktree 생성 (프로젝트 로컬)
- 기준 테스트 검증 건너뛰기
- 질문 없이 실패하는 테스트로 진행

**항상 준수:**

- 먼저 0단계 감지 실행
- git 대체보다 네이티브 도구 우선
- 디렉토리 우선순위 따르기: 기존 > 전역 레거시 > 지시사항 파일 > 기본값
- 프로젝트 로컬의 경우 디렉토리가 무시되는지 확인
- 프로젝트 설정 자동 감지 및 실행
- 깨끗한 테스트 기준 확인
