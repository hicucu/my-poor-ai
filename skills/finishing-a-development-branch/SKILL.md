---
name: finishing-a-development-branch
description: Use when implementation is complete and all tests pass, and you need to decide how to integrate the work - presents structured options for merge, PR, and cleanup, and guides finishing up development
---

# 개발 브랜치 마무리

## 개요

명확한 옵션을 제시하고 선택한 워크플로우를 처리하여 개발 작업 완료를 안내.

**핵심 원칙:** 테스트 검증 → 환경 감지 → 옵션 제시 → 선택 실행 → 정리.

**시작 시 공지:** "finishing-a-development-branch 스킬을 사용하여 이 작업을 완료함."

## 프로세스

### 1단계: 테스트 검증

**옵션 제시 전, 테스트 통과 여부 확인:**

```bash
# 프로젝트 테스트 스위트 실행
npm test / cargo test / pytest / go test ./...
```

**테스트 실패 시:**

```
Tests failing (<N> failures). Must fix before completing:

[Show failures]

Cannot proceed with merge/PR until tests pass.
```

중단. 2단계로 진행하지 말 것.

**테스트 통과 시:** 2단계로 계속.

### 2단계: 환경 감지

**옵션 제시 전 작업 공간 상태 확인:**

```bash
GIT_DIR=$(cd "$(git rev-parse --git-dir)" 2>/dev/null && pwd -P)
GIT_COMMON=$(cd "$(git rev-parse --git-common-dir)" 2>/dev/null && pwd -P)
```

표시할 메뉴와 정리 방법 결정:

| 상태                                   | 메뉴                           | 정리                   |
| -------------------------------------- | ------------------------------ | ---------------------- |
| `GIT_DIR == GIT_COMMON` (일반 저장소)  | 표준 4가지 옵션                | 정리할 worktree 없음   |
| `GIT_DIR != GIT_COMMON`, 명명된 브랜치 | 표준 4가지 옵션                | 출처 기반 (6단계 참조) |
| `GIT_DIR != GIT_COMMON`, detached HEAD | 축소된 3가지 옵션 (merge 없음) | 정리 없음 (외부 관리)  |

### 3단계: 베이스 브랜치 결정

```bash
# 일반적인 베이스 브랜치 시도
git merge-base HEAD main 2>/dev/null || git merge-base HEAD master 2>/dev/null
```

또는 질문: "이 브랜치는 main에서 분기되었습니다 — 맞습니까?"

### 4단계: 옵션 제시

**일반 저장소 및 명명된 브랜치 worktree — 정확히 4가지 옵션 제시:**

```
Implementation complete. What would you like to do?

1. Merge back to <base-branch> locally
2. Push and create a Pull Request
3. Keep the branch as-is (I'll handle it later)
4. Discard this work

Which option?
```

**Detached HEAD — 정확히 3가지 옵션 제시:**

```
Implementation complete. You're on a detached HEAD (externally managed workspace).

1. Push as new branch and create a Pull Request
2. Keep as-is (I'll handle it later)
3. Discard this work

Which option?
```

**설명 추가 금지** — 옵션은 간결하게 유지.

### 5단계: 선택 실행

#### 옵션 1: 로컬 Merge

```bash
# 안전한 CWD를 위해 메인 저장소 루트 확인
MAIN_ROOT=$(git -C "$(git rev-parse --git-common-dir)/.." rev-parse --show-toplevel)
cd "$MAIN_ROOT"

# 먼저 Merge — 아무것도 제거하기 전에 성공 확인
git checkout <base-branch>
git pull
git merge <feature-branch>

# Merge 결과로 테스트 검증
<test command>

# Merge 성공 후에만: worktree 정리 (6단계), 그 다음 브랜치 삭제
```

이후: worktree 정리 (6단계), 그 다음 브랜치 삭제:

```bash
git branch -d <feature-branch>
```

#### 옵션 2: Push 및 PR 생성

```bash
# 브랜치 Push
git push -u origin <feature-branch>

# PR 생성
gh pr create --title "<title>" --body "$(cat <<'EOF'
## Summary
<2-3 bullets of what changed>

## Test Plan
- [ ] <verification steps>
EOF
)"
```

**worktree 정리 금지** — PR 피드백 반영을 위해 사용자에게 필요함.

#### 옵션 3: 현 상태 유지

보고: "브랜치 <name> 유지. Worktree가 <path>에 보존됨."

**worktree 정리 금지.**

#### 옵션 4: 폐기

**먼저 확인:**

```
This will permanently delete:
- Branch <name>
- All commits: <commit-list>
- Worktree at <path>

Type 'discard' to confirm.
```

정확한 확인 대기.

확인 시:

```bash
MAIN_ROOT=$(git -C "$(git rev-parse --git-common-dir)/.." rev-parse --show-toplevel)
cd "$MAIN_ROOT"
```

이후: worktree 정리 (6단계), 그 다음 강제 브랜치 삭제:

```bash
git branch -D <feature-branch>
```

### 6단계: 작업 공간 정리

**옵션 1과 4에서만 실행.** 옵션 2와 3은 항상 worktree를 보존.

```bash
GIT_DIR=$(cd "$(git rev-parse --git-dir)" 2>/dev/null && pwd -P)
GIT_COMMON=$(cd "$(git rev-parse --git-common-dir)" 2>/dev/null && pwd -P)
WORKTREE_PATH=$(git rev-parse --show-toplevel)
```

**`GIT_DIR == GIT_COMMON`인 경우:** 일반 저장소, 정리할 worktree 없음. 완료.

**worktree 경로가 `.worktrees/`, `worktrees/`, 또는 `~/.config/my-poor-ai/worktrees/` 아래인 경우:** my-poor-ai가 이 worktree를 생성했으므로 정리를 담당.

```bash
MAIN_ROOT=$(git -C "$(git rev-parse --git-common-dir)/.." rev-parse --show-toplevel)
cd "$MAIN_ROOT"
git worktree remove "$WORKTREE_PATH"
git worktree prune  # 자가 복구: 오래된 등록 정리
```

**그 외의 경우:** 호스트 환경(하네스)이 이 작업 공간을 소유. 제거 금지. 플랫폼에 작업 공간 종료 도구가 있으면 사용할 것. 없으면 작업 공간을 그대로 유지.

## 빠른 참조

| 옵션            | Merge | Push | Worktree 유지 | 브랜치 정리 |
| --------------- | ----- | ---- | ------------- | ----------- |
| 1. 로컬 Merge   | 예    | -    | -             | 예          |
| 2. PR 생성      | -     | 예   | 예            | -           |
| 3. 현 상태 유지 | -     | -    | 예            | -           |
| 4. 폐기         | -     | -    | -             | 예 (강제)   |

## 흔한 실수

**테스트 검증 건너뛰기**

- **문제:** 손상된 코드를 Merge하거나 실패하는 PR 생성
- **해결:** 항상 옵션 제시 전 테스트 검증

**개방형 질문**

- **문제:** "다음에 무엇을 해야 하나요?"는 모호함
- **해결:** 정확히 4가지 구조화된 옵션 제시 (detached HEAD의 경우 3가지)

**옵션 2에서 worktree 정리**

- **문제:** PR 반영에 필요한 worktree 제거
- **해결:** 옵션 1과 4에서만 정리

**worktree 제거 전 브랜치 삭제**

- **문제:** worktree가 여전히 브랜치를 참조하고 있어 `git branch -d` 실패
- **해결:** 먼저 Merge, worktree 제거, 그 다음 브랜치 삭제

**worktree 내부에서 git worktree remove 실행**

- **문제:** 제거 중인 worktree 내부에 CWD가 있으면 명령이 자동으로 실패
- **해결:** `git worktree remove` 전에 항상 메인 저장소 루트로 `cd`

**하네스 소유 worktree 정리**

- **문제:** 하네스가 생성한 worktree 제거 시 유령 상태 발생
- **해결:** `.worktrees/`, `worktrees/`, 또는 `~/.config/my-poor-ai/worktrees/` 아래의 worktree만 정리

**폐기 시 확인 없음**

- **문제:** 실수로 작업 삭제
- **해결:** 타이핑으로 "discard" 확인 요구

## 위험 신호

**절대 금지:**

- 실패하는 테스트로 진행
- Merge 결과 테스트 검증 없이 Merge
- 확인 없이 작업 삭제
- 명시적 요청 없이 force-push
- Merge 성공 확인 전 worktree 제거
- 생성하지 않은 worktree 정리 (출처 확인)
- worktree 내부에서 `git worktree remove` 실행

**항상 준수:**

- 옵션 제시 전 테스트 검증
- 메뉴 제시 전 환경 감지
- 정확히 4가지 옵션 제시 (detached HEAD의 경우 3가지)
- 옵션 4에 대해 타이핑 확인 요구
- 옵션 1 및 4에서만 worktree 정리
- worktree 제거 전 메인 저장소 루트로 `cd`
- 제거 후 `git worktree prune` 실행
