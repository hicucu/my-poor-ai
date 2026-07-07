# Copilot CLI 도구 매핑

스킬은 Claude Code 도구 이름을 사용함. 스킬에서 아래 도구를 만나면, 사용 중인 플랫폼의 대응 도구를 사용함:

| 스킬에서 참조하는 도구           | Copilot CLI 대응 도구                                            |
| -------------------------------- | ---------------------------------------------------------------- |
| `Read` (파일 읽기)               | `view`                                                           |
| `Write` (파일 생성)              | `create`                                                         |
| `Edit` (파일 편집)               | `edit`                                                           |
| `Bash` (명령 실행)               | `bash`                                                           |
| `Grep` (파일 내용 검색)          | `grep`                                                           |
| `Glob` (이름으로 파일 검색)      | `glob`                                                           |
| `Skill` 도구 (스킬 호출)         | `skill`                                                          |
| `WebFetch`                       | `web_fetch`                                                      |
| `Task` 도구 (서브에이전트 파견)  | `agent_type: "general-purpose"` 또는 `"explore"`를 지정한 `task` |
| 여러 `Task` 호출 (병렬)          | 여러 `task` 호출                                                 |
| Task 상태/출력                   | `read_agent`, `list_agents`                                      |
| `TodoWrite` (작업 추적)          | 내장 `todos` 테이블을 사용하는 `sql`                             |
| `WebSearch`                      | 대응 도구 없음 — 검색 엔진 URL로 `web_fetch` 사용                |
| `EnterPlanMode` / `ExitPlanMode` | 대응 도구 없음 — 메인 세션에 머무름                              |

## 비동기 셸 세션

Copilot CLI는 지속형 비동기 셸 세션을 지원하며, 이에 직접 대응하는 Claude Code 도구는 없음:

| 도구                   | 용도                                   |
| ---------------------- | -------------------------------------- |
| `bash` (`async: true`) | 장시간 실행 명령을 백그라운드에서 시작 |
| `write_bash`           | 실행 중인 비동기 세션에 입력 전송      |
| `read_bash`            | 비동기 세션의 출력 읽기                |
| `stop_bash`            | 비동기 세션 종료                       |
| `list_bash`            | 활성 셸 세션 전체 나열                 |

## 추가 Copilot CLI 도구

| 도구                                    | 용도                                                |
| --------------------------------------- | --------------------------------------------------- |
| `store_memory`                          | 코드베이스에 대한 사실을 이후 세션을 위해 영구 저장 |
| `report_intent`                         | 현재 의도를 UI 상태 표시줄에 갱신                   |
| `sql`                                   | 세션의 SQLite 데이터베이스 조회 (todos, 메타데이터) |
| `fetch_copilot_cli_documentation`       | Copilot CLI 문서 조회                               |
| GitHub MCP 도구 (`github-mcp-server-*`) | 네이티브 GitHub API 접근 (이슈, PR, 코드 검색)      |
