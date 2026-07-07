# Gemini CLI 도구 매핑

스킬은 Claude Code 도구 이름을 사용함. 스킬에서 아래 도구를 만나면, 사용 중인 플랫폼의 대응 도구를 사용함:

| 스킬에서 참조하는 도구          | Gemini CLI 대응 도구                                         |
| ------------------------------- | ------------------------------------------------------------ |
| `Read` (파일 읽기)              | `read_file`                                                  |
| `Write` (파일 생성)             | `write_file`                                                 |
| `Edit` (파일 편집)              | `replace`                                                    |
| `Bash` (명령 실행)              | `run_shell_command`                                          |
| `Grep` (파일 내용 검색)         | `grep_search`                                                |
| `Glob` (이름으로 파일 검색)     | `glob`                                                       |
| `TodoWrite` (작업 추적)         | `write_todos`                                                |
| `Skill` 도구 (스킬 호출)        | `activate_skill`                                             |
| `WebSearch`                     | `google_web_search`                                          |
| `WebFetch`                      | `web_fetch`                                                  |
| `Task` 도구 (서브에이전트 파견) | `@agent-name` ([서브에이전트 지원](#서브에이전트-지원) 참조) |

## 서브에이전트 지원

Gemini CLI는 `@` 문법으로 서브에이전트를 네이티브로 지원함. 내장 `@generalist` 에이전트를 사용해 어떤 작업이든 파견함 — 이 에이전트는 모든 도구에 접근할 수 있고 제공된 프롬프트를 따름.

스킬이 특정 이름의 에이전트 유형을 파견하라고 할 때는, 스킬의 프롬프트 템플릿에서 채운 전체 프롬프트와 함께 `@generalist`를 사용함:

| 스킬 지침                                       | Gemini CLI 대응                                                                 |
| ----------------------------------------------- | ------------------------------------------------------------------------------- |
| `Task tool (my-poor-ai:implementer)`                 | 채워진 `implementer-prompt.md` 템플릿과 함께 `@generalist`                      |
| `Task tool (my-poor-ai:spec-reviewer)`               | 채워진 `spec-reviewer-prompt.md` 템플릿과 함께 `@generalist`                    |
| `Task tool (my-poor-ai:code-reviewer)`               | `@code-reviewer` (번들 에이전트) 또는 채워진 리뷰 프롬프트와 함께 `@generalist` |
| `Task tool (my-poor-ai:code-quality-reviewer)`       | 채워진 `code-quality-reviewer-prompt.md` 템플릿과 함께 `@generalist`            |
| `Task tool (general-purpose)` (인라인 프롬프트) | 인라인 프롬프트와 함께 `@generalist`                                            |

### 프롬프트 채우기

스킬은 `{WHAT_WAS_IMPLEMENTED}`나 `[FULL TEXT of task]` 같은 플레이스홀더가 있는 프롬프트 템플릿을 제공함. 모든 플레이스홀더를 채우고 완성된 프롬프트를 `@generalist`에 메시지로 전달함. 프롬프트 템플릿 자체에 에이전트의 역할, 리뷰 기준, 기대 출력 형식이 담겨 있으며 — `@generalist`가 이를 따름.

### 병렬 파견

Gemini CLI는 서브에이전트 병렬 파견을 지원함. 스킬이 독립적인 서브에이전트 작업 여러 개를 병렬로 파견하라고 하면, 그 `@generalist` 또는 이름이 지정된 서브에이전트 작업을 모두 같은 프롬프트에서 함께 요청함. 의존성이 있는 작업은 순차로 유지하되, 단지 이력을 단순하게 보이게 하려고 독립적인 서브에이전트 작업을 직렬화하지 않음.

## 추가 Gemini CLI 도구

다음 도구는 Gemini CLI에서 사용할 수 있지만 대응하는 Claude Code 도구가 없음:

| 도구                                 | 용도                                        |
| ------------------------------------ | ------------------------------------------- |
| `list_directory`                     | 파일과 하위 디렉토리 나열                   |
| `save_memory`                        | GEMINI.md에 사실을 세션 간 영구 저장        |
| `ask_user`                           | 사용자에게 구조화된 입력 요청               |
| `tracker_create_task`                | 풍부한 작업 관리 (생성, 갱신, 나열, 시각화) |
| `enter_plan_mode` / `exit_plan_mode` | 변경 전 읽기 전용 조사 모드로 전환          |
