# Pressure Scenarios

Discipline skill이 압박 상황에서도 준수되는지 검증하는 시나리오 모음.

## 테스트 방법

각 시나리오를 Claude Code에서 실행:

```bash
# skill 없이 (baseline — 실패해야 정상)
claude --no-skills < {scenario}.md

# skill 있이 (my-poor-ai 플러그인 설치 후 — 준수해야 정상)
claude < {scenario}.md
```

## 판정 기준

| 결과                        | 의미                 |
| --------------------------- | -------------------- |
| baseline에서 합리화 후 위반 | 정상 (이것을 문서화) |
| my-poor-ai 설치 후 준수          | skill 통과           |
| my-poor-ai 설치 후도 위반        | skill 수정 필요      |

## 파일 목록

| 파일                        | 테스트 대상 skill                      |
| --------------------------- | -------------------------------------- |
| `tdd-pressure.md`           | `my-poor-ai:test-driven-development`        |
| `debugging-pressure.md`     | `my-poor-ai:systematic-debugging`           |
| `verification-pressure.md`  | `my-poor-ai:verification-before-completion` |
| `brainstorming-pressure.md` | `my-poor-ai:brainstorming`                  |
| `pipeline-dispatch.md`      | `using-my-poor-ai` 복잡도 분류              |
| `which-way-should-i-go-pressure.md` | `my-poor-ai:which-way-should-i-go`  |
| `socratic-plan-review-pressure.md`  | `my-poor-ai:socratic-plan-review`   |
| `worktree-native-preference-pressure.md` | `my-poor-ai:using-git-worktrees` |

## 음성 결과 회귀 테스트

일부 시나리오는 **베이스라인이 이미 준수**하여 스킬 본문을 수정하지 않은 채 보존됨. 향후 모델 변경으로 이 준수가 무너지면 그때가 보강 시점임.

| 시나리오 | 관찰일 | 상태 |
| --- | --- | --- |
| `verification-pressure.md` 시나리오 C (초록불 CI vs 실제 기동) | 2026-07-28 | 준수 — 보강 불필요 |
