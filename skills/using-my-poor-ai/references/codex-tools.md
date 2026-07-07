# Codex 도구 매핑

스킬은 Claude Code 도구 이름을 사용함. 스킬에서 아래 도구를 만나면, 사용 중인 플랫폼의 대응 도구를 사용함:

| 스킬에서 참조하는 도구          | Codex 대응 도구                                                                                                      |
| ------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| `Task` 도구 (서브에이전트 파견) | `spawn_agent` ([서브에이전트 파견에는 멀티 에이전트 지원 필요](#서브에이전트-파견에는-멀티-에이전트-지원-필요) 참조) |
| 여러 `Task` 호출 (병렬)         | 여러 `spawn_agent` 호출                                                                                              |
| Task 결과 반환                  | `wait_agent`                                                                                                         |
| Task 자동 완료                  | `close_agent`로 슬롯 해제                                                                                            |
| `TodoWrite` (작업 추적)         | `update_plan`                                                                                                        |
| `Skill` 도구 (스킬 호출)        | 스킬이 네이티브로 로드됨 — 지침을 그대로 따름                                                                        |
| `Read`, `Write`, `Edit` (파일)  | 네이티브 파일 도구 사용                                                                                              |
| `Bash` (명령 실행)              | 네이티브 셸 도구 사용                                                                                                |

## 서브에이전트 파견에는 멀티 에이전트 지원 필요

Codex 설정(`~/.codex/config.toml`)에 다음을 추가함:

```toml
[features]
multi_agent = true
```

이렇게 하면 `dispatching-parallel-agents`, `subagent-driven-development` 같은 스킬을 위한 `spawn_agent`, `wait_agent`, `close_agent`가 활성화됨.

레거시 참고: `rust-v0.115.0` 이전 Codex 빌드는 파견된 에이전트의 대기를
`wait`로 노출했음. 현재 Codex는 파견된 에이전트에 `wait_agent`를 사용함.
이제 `wait`라는 이름은 code-mode `exec/wait`에 속하며, 이는 yield된 exec
셀을 `cell_id`로 재개하는 도구임. 파견 에이전트의 결과 도구가 아님.

## 환경 감지

워크트리를 생성하거나 브랜치를 마무리하는 스킬은, 진행 전에
읽기 전용 git 명령으로 환경을 감지해야 함:

```bash
GIT_DIR=$(cd "$(git rev-parse --git-dir)" 2>/dev/null && pwd -P)
GIT_COMMON=$(cd "$(git rev-parse --git-common-dir)" 2>/dev/null && pwd -P)
BRANCH=$(git branch --show-current)
```

- `GIT_DIR != GIT_COMMON` → 이미 연결된 워크트리 내부 (생성 건너뜀)
- `BRANCH`가 비어 있음 → detached HEAD (샌드박스에서 브랜치/push/PR 불가)

각 스킬이 이 신호를 어떻게 사용하는지는 `using-git-worktrees` Step 0과
`finishing-a-development-branch` Step 1을 참조함.

## Codex 앱에서의 마무리

샌드박스가 브랜치/push 작업을 차단하면(외부에서 관리되는 워크트리의
detached HEAD), 에이전트는 모든 작업을 커밋하고 사용자에게 앱의
네이티브 컨트롤을 사용하도록 안내함:

- **"Create branch"** — 브랜치 이름을 지정한 뒤, 앱 UI로 커밋/push/PR
- **"Hand off to local"** — 작업을 사용자의 로컬 체크아웃으로 이관

에이전트는 여전히 테스트 실행, 파일 스테이징을 할 수 있고, 사용자가
복사해 쓸 수 있도록 추천 브랜치 이름·커밋 메시지·PR 설명을 출력할 수 있음.
