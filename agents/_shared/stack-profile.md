# 공유 지침 — 스택 감지와 `stack-profile.json` (공유 모듈)

> 이 파일은 실행 에이전트가 아니라 **스택 감지 지침 모듈**임. 산출자는 `project-context`(복잡 경로 Phase 0) 하나이며 정의에서 `@include`함. 소비자(리뷰어·aggregator·issue-fixer·developer-agent·writing-plans 등)는 이 문서의 스키마로 `_workspaces/stack-profile.json`을 읽음. 수정 후에는 `node scripts/generate-codex-agents.mjs` 재실행 필수 (Codex 미러에 인라인되는 내용이므로).

## 산출 위치와 재사용

- 경로: **`_workspaces/stack-profile.json`** — 프로젝트 공용 파일. 기능·브랜치별 작업 디렉토리(`{workspaceDir}`) 안에 두지 않음
- 재사용: 24시간 이내에 생성되었고 그 뒤로 아래 매트릭스의 마커 파일이 바뀌지 않았으면 재사용. 그 외에는 다시 감지해 덮어씀
- 작업 디렉토리 경로는 이 파일에 기록하지 않음 — 공용 파일이 특정 작업에 묶이면 다른 경로의 소비자가 엉뚱한 디렉토리를 읽음

## 스택 감지 매트릭스

프로젝트 루트와 직속 하위(`src/`, `app/`, `lib/`, `cmd/` 등)에서 다음 마커를 순차 검사함.

| 마커 파일 (정확 이름)                              | primary  | 추가 시그널 (subtype 결정)                                                                                                |
| -------------------------------------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------- |
| `package.json`                                     | `node`   | dependencies: `react` → react / `next` → next / `express` → express / `@nestjs/core` → nestjs / `vue` → vue / `vite` 빌드 |
| `pyproject.toml` / `requirements.txt` / `setup.py` | `python` | dependencies: `fastapi` → fastapi / `django` → django / `flask` → flask                                                   |
| `*.csproj` / `*.sln`                               | `dotnet` | `Microsoft.AspNetCore.App` → aspnetcore / `Microsoft.NET.Sdk.Web` → web / `Microsoft.NET.Sdk` → library                   |
| `go.mod`                                           | `go`     | imports: `gin-gonic/gin` → gin / `labstack/echo` → echo / `net/http` only → stdlib                                        |
| `pom.xml` / `build.gradle(.kts)`                   | `jvm`    | dependencies: `spring-boot` → spring / `quarkus` → quarkus                                                                |
| `Cargo.toml`                                       | `rust`   | dependencies: `actix-web` → actix / `axum` → axum / `tokio` 단독 → tokio                                                  |
| `composer.json`                                    | `php`    | require: `laravel/framework` → laravel / `symfony/*` → symfony                                                            |

마커 미발견 시 `primary: "unknown"`, `fallbackUsed: true` 설정 — 추측으로 스택을 단정하지 않음. 이후 처리(사용자에게 스택 명시 요청)는 호출한 에이전트의 정의를 따름.

### language / packageManager / moduleSystem 결정

- Node: `tsconfig.json` 존재 → typescript, 없으면 javascript / lock file로 npm·pnpm·yarn 구분 / `"type": "module"` → esm, 없으면 cjs
- Python: poetry.lock → poetry, Pipfile → pipenv, 없으면 pip
- .NET: 항상 csharp (또는 fsproj면 fsharp) / nuget
- 그 외: 표준 매니저 (cargo, go, maven, gradle)

### testFramework 결정 (우선순위)

| 스택   | 우선 검사                                                                                    | fallback   |
| ------ | -------------------------------------------------------------------------------------------- | ---------- |
| node   | `vitest.config.*` → vitest / `jest.config.*` 또는 package.json jest 키 → jest                | jest       |
| python | `pytest.ini` / `pyproject.toml [tool.pytest]` → pytest / `unittest` 디렉토리 패턴 → unittest | pytest     |
| dotnet | csproj `<PackageReference Include="xunit">` → xunit / `nunit` → nunit / `MSTest` → mstest    | xunit      |
| go     | 항상 `go-test` (표준)                                                                        | go-test    |
| jvm    | gradle/maven 의존성: `junit-jupiter` → junit5 / `junit:junit` → junit4 / `kotest` → kotest   | junit5     |
| rust   | 항상 `cargo-test` (표준)                                                                     | cargo-test |
| php    | `phpunit.xml` → phpunit / `pest.config.php` → pest                                           | phpunit    |

### businessLogicMarkers / uiMarkers 휴리스틱

| primary | businessLogicMarkers 후보                                                           | uiMarkers 후보                                                           |
| ------- | ----------------------------------------------------------------------------------- | ------------------------------------------------------------------------ |
| node    | `src/services`, `src/hooks`, `src/lib`, `src/utils`, `src/domain`, `src/middleware` | `src/components`, `src/pages`, `src/views`, `app/`                       |
| python  | `app/services`, `app/use_cases`, `app/domain`, `*/services.py`                      | `app/templates`, `app/views`, `app/routers` (FastAPI 라우터는 검토 필요) |
| dotnet  | `Services/`, `Application/`, `Domain/`, `*Service.cs`                               | `Controllers/`, `Views/`, `Pages/`, `Components/`                        |
| go      | `internal/service`, `pkg/`, `internal/domain`                                       | `cmd/`, `internal/handler`, `web/`                                       |
| jvm     | `src/main/java/**/service`, `src/main/kotlin/**/service`, `domain`                  | `controller`, `resources/templates`                                      |

실제 존재하는 디렉토리만 포함. 휴리스틱이 맞지 않으면 사용자에게 보정 요청.

## stack-profile.json 스키마

```json
{
  "primary": "node",
  "subtype": "express",
  "language": "typescript",
  "packageManager": "pnpm",
  "moduleSystem": "esm",
  "sourceLayout": { "root": "src", "tests": "src/__tests__" },
  "testFramework": { "name": "vitest", "configFile": "vitest.config.ts" },
  "linter": "eslint",
  "formatter": "prettier",
  "buildTool": "tsc",
  "conventions": {
    "fileNaming": "kebab-case",
    "testPattern": "*.test.ts",
    "componentExt": null
  },
  "businessLogicMarkers": ["src/services", "src/domain"],
  "uiMarkers": [],
  "detectionEvidence": [
    "package.json:express@4.18",
    "tsconfig.json",
    "vitest.config.ts"
  ],
  "fallbackUsed": false
}
```

- 전체 필드를 작성하고 모르는 값은 `null`로 둠 — 부분 작성 금지
- `detectionEvidence`에는 판정 근거(파일명 + 발견한 키워드)를 3~5개 기재. 사용자가 검증할 수 있도록

## 경로 정책

- `testFramework.configFile`, `sourceLayout.root`, `sourceLayout.tests`, `businessLogicMarkers`, `uiMarkers`의 모든 경로는 **프로젝트 루트 기준 상대 경로**로 작성함
- 절대 경로 또는 `~/` 사용 금지. 다른 프로젝트로 복사된 후에도 그대로 해석되어야 함
- `detectionEvidence`의 항목은 `파일명:키워드` 형식의 짧은 식별자 (예: `package.json:express@4.19.2`). 전체 경로 불필요
