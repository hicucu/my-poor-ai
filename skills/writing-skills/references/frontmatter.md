# SKILL.md frontmatter 필드

> 2026-10-19 확인 — Claude Code 2.1.303의 스킬 frontmatter 스키마 기준. 버전이 바뀌면 재확인함.

## 플랫폼별 인식 범위

| 환경 | 인식하는 필드 |
| --- | --- |
| Agent Skills 표준 (claude.ai 업로드·Claude API) | `name`·`description`·`license`·`compatibility`·`metadata`·`allowed-tools` 6개 |
| Claude Code (개인·프로젝트·플러그인 스킬) | 아래 표 전부 |

다른 환경에도 배포할 스킬은 표준 6필드만으로 동작하게 쓰고, Claude Code 전용 필드는 부가 기능으로만 둠.

## Claude Code 필드

| 필드 | 의미 |
| --- | --- |
| `name` | 표시 이름. 생략 시 디렉토리명. 플러그인 스킬은 이 값이 `/플러그인:이름`의 마지막 세그먼트가 됨 — 이 저장소는 디렉토리명과 같아야 함 (validator 강제) |
| `description` | 목록과 Skill 도구에 노출되는 한 줄. 트리거 조건 중심으로 씀 (SKILL.md "CSO" 절) |
| `when_to_use` | 모델이 이 스킬을 집어야 할 때의 안내. 도구 설명에 합쳐짐 — `description`과 함께 상시 목록에 들어감 |
| `argument-hint` | 슬래시 커맨드 이름 뒤에 표시되는 인수 자리표시자 |
| `disable-model-invocation` | `true`면 모델이 Skill 도구로 호출하지 못함 — 사용자만 슬래시 커맨드로 시작. 프로젝트 밖에 기록하는 등 부작용 있는 스킬에 사용 |
| `user-invocable` | `false`면 슬래시 커맨드를 숨김 — 모델만 호출 |
| `allowed-tools` | 스킬이 활성인 동안 쓸 수 있는 도구. 쉼표 문자열 또는 YAML 리스트 |
| `disallowed-tools` | 스킬이 활성인 동안 제거할 도구. 사용자의 다음 메시지에서 해제됨 |
| `model` | 모델 지정 (`haiku`·`sonnet`·`opus`·`fable`·전체 ID, `inherit`) |
| `effort` | 사고 강도 (`low`·`medium`·`high`·`max` 또는 정수) |
| `context` | `inline`(현재 대화에 전개, 기본) 또는 `fork`(서브에이전트로 실행) |
| `agent` | `context: fork`일 때 띄울 에이전트 타입 |
| `background` | `context: fork` 전용. 기본은 백그라운드 실행 후 작업 알림으로 보고, `false`면 호출자가 결과를 기다림 |
| `paths` | glob 패턴. 모델이 일치하는 파일을 건드릴 때만 로드됨 |
| `hooks` | 스킬이 활성인 동안 등록되는 훅. settings.json `hooks`와 같은 형태 |
| `shell` | `!`-명령 블록용 셸 (`bash`·`powershell`, 기본 bash) |
| `metadata` | 작성자용 자유 맵. 로드된 스킬에 보존됨 |

`arguments`·`version`은 내부용으로 문서화된 필드가 아님 — 인수 안내는 `argument-hint`로 씀.

## 상시 비용과 지연 로드

- `description`·`when_to_use`는 **모든 대화의 스킬 목록에 상시 포함**됨 — 짧게 씀. 워크플로우 요약 금지 규칙(SKILL.md "CSO")은 두 필드 모두에 적용됨
- 본문(SKILL.md)은 호출될 때 로드됨. 본문이 참조하는 보조 파일은 모델이 읽을 때만 로드됨 — 대용량 참조는 보조 파일로 분리함
- `disable-model-invocation: true` 스킬은 모델용 목록에서 빠짐 — 상시 비용이 없음
