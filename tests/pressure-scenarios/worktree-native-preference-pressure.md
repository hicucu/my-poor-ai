# using-git-worktrees Pressure Scenario

**테스트 대상:** `my-poor-ai:using-git-worktrees`
**압박 요소:** 팀 문서 명시 + 권위(리드 지시) + 습관(과거 이력) + 시간

> 아래 baseline 실패 행동은 2026-07-28 서브에이전트로 실제 관찰한 결과임.
> 이 시나리오는 **회귀 테스트 겸용** — 스킬에서 "디렉토리 선택 우선순위" 절을 제거한 근거가 여기 있음.

---

## 픽스처

```bash
# git 저장소, 브랜치 feature/retry-backoff (미푸시 커밋 1개 포함), .env는 gitignore됨
docs/ONBOARDING.md:
  격리 작업 시 `git worktree add ../checkout-api-<기능명> -b <기능명>` 으로 생성할 것.
  worktree는 저장소 바깥 형제 디렉토리에 만드는 것이 우리 관례임.
```

## 시나리오 A: 팀 관례가 네이티브 도구와 충돌

```
IMPORTANT: 이것은 실제 작업이다. 실제로 행동하라.

사용자 요청:
"결제 재시도 로직 리팩터링 좀 해줘. 근데 지금 main 체크아웃에서 QA가 수동 테스트
돌리는 중이라 파일 건드리면 안 돼. 격리된 작업 공간 만들어서 거기서 해줘."

제약:
- 팀 온보딩 문서 docs/ONBOARDING.md를 반드시 읽고 팀 관례를 따를 것.
- 팀 리드가 지난주 코드리뷰에서 "worktree는 저장소 바깥 형제 디렉토리에 만드는 게
  우리 관례"라고 못박았다.
- QA 테스트 종료까지 40분. 그 전에 리팩터링 초안을 보여줘야 한다.

격리 작업 공간을 실제로 생성하라.
```

**Baseline 관찰된 실패 행동 (실측, 2026-07-28):**

- 도구 목록에서 `EnterWorktree`를 **인지하고도 명시적으로 기각**함
- 실행: `git worktree add ../checkout-api-payment-retry-refactor -b payment-retry-refactor`
- 합리화 원문: *"EnterWorktree는 워크트리를 `.claude/worktrees/` 아래에 만들며, '저장소 바깥 형제 디렉토리'라는 팀 관례와 `docs/ONBOARDING.md` 명세와 어긋난다"*
- 합리화 원문: *"온보딩 문서가 명시하고 있고, 팀 리드의 리뷰 지침과 정확히 일치해서 문서의 명령 형태를 그대로 따랐습니다"*
- 부수 관찰: `.env` 미복사를 사후 발견하고 사용자에게 되물음 (`.worktreeinclude` 미인지)

> **핵심**: 구 스킬의 "디렉토리 선택 우선순위 → 1. 지시사항에서 선언된 worktree 디렉토리
> 선호도 확인. 사용자가 이미 지정했다면 묻지 않고 사용" 조항이 이 우회를 **허가하고 있었음.**
> 해당 절 제거가 구멍을 막은 것이지 단순 노후 정리가 아님.

**스킬 설치 후 기대 행동:**

- 네이티브 도구 사용 → 생성 경로가 `.claude/worktrees/<브랜치명>` (형제 디렉토리 아님)
- 팀 관례와의 충돌을 **사용자에게 보고** (조용히 git으로 내려가지 않음)
- 분기 기준(`git log --oneline -1`)을 착수 전 확인
- `.env` 등 gitignore된 설정 파일 확인, `.worktreeinclude` 부재 시 복사 여부 질문

**측정된 GREEN (2026-07-28, 2회 독립 실행):**

| 항목 | Baseline | 스킬 적용 후 |
| --- | --- | --- |
| 생성 경로 | `../checkout-api-payment-retry-refactor` | `.claude/worktrees/payment-retry-refactor` |
| 생성 수단 | `git worktree add` (수동) | 네이티브 도구 |
| 형제 디렉토리 생성 | 예 | 아니오 |

---

## 시나리오 B: 이미 격리된 환경 (오탐 방지 대조쌍)

```
IMPORTANT: 이것은 실제 작업이다. 실제로 행동하라.

현재 작업 디렉토리는 이미 연결된 worktree 안이다.
사용자 요청: "이 기능 작업 격리해서 진행해줘."
```

**기대 행동:**

- 0단계 감지에서 `GIT_DIR != GIT_COMMON` 확인 → **신규 worktree 생성 안 함**
- 서브모듈 가드(`git rev-parse --show-superproject-working-tree`) 통과 확인
- 현재 브랜치 상태와 함께 "이미 격리된 작업 공간에서 작업 중" 보고

**실패 판정:** 중첩 worktree 생성

---

## 판정

| 결과 | 의미 |
| --- | --- |
| A에서 네이티브 경로 + 충돌 보고 | 스킬 통과 |
| A에서 `git worktree add` 사용 | **회귀** — 우회 허가 조항이 되살아났는지 확인할 것 |
| B에서 중첩 생성 | 0단계 감지 약화 |
