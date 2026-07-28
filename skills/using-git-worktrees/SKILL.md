---
name: using-git-worktrees
description: Use when starting feature work that needs isolation from the current workspace, before executing an implementation plan, or when a team convention or instruction file specifies a worktree directory layout - ensures an isolated workspace via native tooling, with a git worktree fallback only when no native tool exists
---

# Git Worktree 사용

## 개요

작업이 격리된 작업 공간에서 진행되도록 보장. 플랫폼의 네이티브 worktree 도구를 우선 사용. 네이티브 도구가 없는 경우에만 수동 git worktree로 대체.

**핵심 원칙:** 기존 격리 감지 → 네이티브 도구 → (없을 때만) git 대체. 하네스와 충돌 금지.

**시작 시 공지:** "using-git-worktrees 스킬을 사용하여 격리된 작업 공간을 설정함."

## 0단계: 기존 격리 환경 감지

**아무것도 생성하기 전에, 이미 격리된 작업 공간에 있는지 확인.**

```bash
GIT_DIR=$(cd "$(git rev-parse --git-dir)" 2>/dev/null && pwd -P)
GIT_COMMON=$(cd "$(git rev-parse --git-common-dir)" 2>/dev/null && pwd -P)
git rev-parse --show-superproject-working-tree 2>/dev/null   # 값이 있으면 서브모듈
```

**서브모듈 가드:** `GIT_DIR != GIT_COMMON`은 서브모듈 내부에서도 참임. 위 세 번째 명령이 경로를 반환하면 worktree가 아니라 서브모듈이므로 일반 저장소로 처리.

| 판정 | 조치 |
| --- | --- |
| `GIT_DIR != GIT_COMMON` (서브모듈 아님) | 이미 격리됨 — **생성하지 않고** 3단계로 |
| `GIT_DIR == GIT_COMMON` 또는 서브모듈 | 일반 체크아웃 — 1단계로 |

이미 격리된 경우 브랜치 상태와 함께 보고함. Detached HEAD면 "완료 시점에 브랜치 생성 필요"를 함께 알림.

일반 체크아웃이고 지시사항에 worktree 선호가 선언되어 있지 않다면, 생성 전 동의를 구함:

> "격리된 worktree를 설정할까요? 현재 브랜치를 변경으로부터 보호함."

거부 시 현재 위치에서 작업하고 3단계로.

## 1단계: 네이티브 도구 우선

**사용 가능한 도구 목록을 실제로 확인.** `EnterWorktree`, `WorktreeCreate`, `/worktree`, `--worktree` 플래그 같은 이름의 도구가 있으면 **그것을 사용하고 2단계로 건너뜀.**

네이티브 도구는 디렉토리 배치·브랜치 생성·정리·잠금을 하네스가 추적하는 방식으로 처리함. 네이티브 도구가 있는데 `git worktree add`를 쓰면 하네스가 보지 못하는 유령 상태가 생기고, 자동 정리·세션 재개·트랜스크립트 이동이 전부 깨짐.

### 디렉토리 배치 관례는 우회 근거가 아님

**팀 문서·온보딩 가이드·리드의 지시가 다른 디렉토리 배치(예: "저장소 바깥 형제 디렉토리")를 명시하더라도, 네이티브 도구를 우회하지 않음.**

배치 관례와 도구 선택은 별개의 문제임. 배치는 네이티브 도구의 설정(`worktree` 설정, `WorktreeCreate` 훅)으로 표현하는 것이지, 도구를 버리는 근거가 아님. 관례가 네이티브 배치와 충돌하면 **사용자에게 그 충돌을 보고**하고 판단을 받음 — 조용히 git으로 내려가지 않음.

### 함정 1: 베이스 브랜치 기본값

네이티브 worktree는 보통 **현재 HEAD가 아니라 저장소 기본 브랜치**에서 분기함(`worktree.baseRef` 기본값 `fresh`). 수동 `git worktree add`의 기본값(HEAD 분기)과 반대임.

진행 중인 미푸시 작업 위에서 격리하면, 그 작업이 **없는** 상태로 시작하게 됨. 착수 전 확인:

```bash
git log --oneline -1        # 새 작업 공간이 어디서 분기했는가
git status -sb
```

기대와 다르면 베이스를 명시하거나(`worktree.baseRef: "head"`) 사용자에게 확인함. 잘못된 베이스는 조용히 진행되다 병합 시점에 드러남.

### 함정 2: gitignore된 설정 파일

worktree는 새 체크아웃이라 `.env`, `.env.local` 같은 **비추적 파일이 없음**. 설정 누락을 엉뚱한 원인으로 진단하게 되는 흔한 함정임.

프로젝트 루트에 `.worktreeinclude`(gitignore 문법)를 두면 매칭되면서 gitignore된 파일이 새 worktree로 자동 복사됨. 없으면 필요한 파일을 확인하고 사용자에게 복사 여부를 물음 — 시크릿을 임의로 복사하지 않음.

## 1b단계: Git 대체 (네이티브 도구가 없을 때만)

```bash
git worktree add "<경로>/<브랜치명>" -b "<브랜치명>"
cd "<경로>/<브랜치명>"
```

경로는 지시사항에 선언된 관례를 따르고, 없으면 프로젝트 루트의 `.worktrees/`를 사용함. 프로젝트 내부에 만들 때는 해당 디렉토리가 무시되는지 먼저 확인:

```bash
git check-ignore -q .worktrees || echo "먼저 .gitignore에 추가하고 커밋할 것"
```

**샌드박스 대체:** 권한 오류로 실패하면 샌드박스가 생성을 차단한 것임. 사용자에게 알리고 현재 디렉토리에서 작업함.

## 2단계: 환경 준비

새 체크아웃이므로 의존성이 없음. 프로젝트에 맞는 설치를 실행함 (`npm install` / `cargo build` / `pip install -r requirements.txt` / `go mod download` 등). 함정 2의 설정 파일도 이 시점에 확인함.

## 3단계: 깨끗한 기준 확인

작업 공간이 깨끗한 상태에서 시작하는지 테스트로 확인함.

**테스트 실패 시:** 실패를 보고하고 계속할지 조사할지 질문. 기준을 확인하지 않으면 새 버그와 기존 실패를 구별할 수 없음.

### 보고

```
작업 공간 준비 완료: <전체-경로>  (생성 수단: <네이티브 도구명 | git 대체>)
분기 기준: <브랜치> @ <short-sha>
테스트 통과 (<N>개, 실패 0건)
```

## 빠른 참조

| 상황 | 조치 |
| --- | --- |
| 이미 연결된 worktree에 있음 | 생성 건너뜀 (0단계) |
| 서브모듈 내부 | 일반 저장소로 처리 (0단계 guard) |
| 네이티브 worktree 도구 있음 | 사용 (1단계) |
| 팀 문서가 다른 배치를 명시 | 네이티브 사용 + 충돌을 사용자에게 보고 |
| 네이티브 도구 없음 | git 대체 (1b단계) |
| 프로젝트 내부에 생성 | `git check-ignore`로 무시 확인 |
| 분기 기준이 기대와 다름 | 베이스 명시 또는 사용자 확인 |
| `.env` 등 설정 누락 | `.worktreeinclude` 확인, 없으면 복사 여부 질문 |
| 기준 테스트 실패 | 실패 보고 + 질문 |

## 합리화 표

| 변명 | 현실 |
| --- | --- |
| "팀 문서에 `git worktree add` 절차가 명시돼 있다" | 문서가 정하는 건 **배치**지 도구가 아님. 배치는 네이티브 도구 설정으로 표현함 |
| "네이티브 도구는 `.claude/worktrees/`에 만들어서 우리 관례와 어긋난다" | 어긋남을 발견했으면 보고할 것. 조용히 git으로 내려가는 건 하네스에 유령 상태를 만드는 선택임 |
| "리드가 형제 디렉토리라고 못박았다" | 권위는 배치에 대한 것임. 도구 우회 승인이 아님 |
| "예전에도 `git worktree add`로 했다" | 이력은 근거가 아님. 도구 목록을 지금 확인할 것 |
| "40분밖에 없어서 빠른 쪽으로" | 네이티브 도구가 더 빠름. 정리·재개까지 하네스가 처리함 |
| "worktree 만들었으니 바로 작업 시작" | 분기 기준과 설정 파일을 확인하지 않으면 잘못된 베이스에서 작업하게 됨 |

## 위험 신호

**절대 금지:**

- 0단계에서 기존 격리가 감지됐는데 worktree 생성
- 네이티브 도구가 있는데 `git worktree add` 실행 — **가장 흔한 실패**
- 팀 문서·관례·권위를 근거로 네이티브 도구를 조용히 기각
- 도구 목록을 확인하지 않고 git 명령부터 실행
- 분기 기준(`git log -1`) 미확인으로 착수
- 기준 테스트 없이 진행

**항상 준수:**

- 0단계 감지 먼저
- 사용 가능한 도구 목록을 실제로 확인한 뒤 선택
- 네이티브와 관례가 충돌하면 사용자에게 보고
- 분기 기준과 설정 파일을 착수 전 확인
