---
description: 프로젝트 마커 파일을 스캔하여 기술 스택을 자동 감지하고 stack-profile.json을 생성하는 단독 커맨드. feature-pipeline 파이프라인 없이 스택 정보만 필요할 때 사용.
model: opus
---

# Detect Stack (스택 감지 단독 커맨드)

대상 디렉토리의 기술 스택을 마커 파일 기반으로 감지하여 `stack-profile.json` 결과를 산출함.
feature-pipeline의 Phase 1.0 로직만 추출한 단독 실행 커맨드.

## 사용법

```
/my-poor-ai:detect-stack                          현재 디렉토리 감지, _workspaces/stack-profile.json 작성
/my-poor-ai:detect-stack <경로>                   지정 경로 감지, <경로>/_workspaces/stack-profile.json 작성
/my-poor-ai:detect-stack --inline                 파일 작성 없이 결과 출력만
/my-poor-ai:detect-stack <경로> --inline          지정 경로 감지, 출력만
/my-poor-ai:detect-stack --out <파일>             출력 경로 직접 지정 (CWD 기준 상대)
```

## 실행 절차

### Step 1: 입력 인자 해석

- 첫 인수가 디렉토리 경로면 → 대상 루트로 사용 (기본: CWD)
- `--inline` 플래그가 있으면 → 파일 작성 안 함, 콘솔 출력만
- `--out <파일>` 옵션이 있으면 → 해당 경로에 작성 (기본: `<대상>/_workspaces/stack-profile.json`)

### Step 2: feature-planner의 Phase 1.0 매트릭스 적용

`{팀_위치}/agents/feature-planner.md`를 읽고 "Phase 1.0 — 스택 감지 매트릭스" 섹션의 절차만 수행함.

수행 범위:

1. 대상 루트의 마커 파일 스캔 (`package.json`, `pyproject.toml`, `*.csproj`, `go.mod`, `pom.xml`, `Cargo.toml`, `composer.json` 등)
2. 매트릭스 표에서 `primary`/`subtype` 결정
3. `language`/`packageManager`/`moduleSystem` 결정
4. `testFramework` 우선순위 검사
5. `businessLogicMarkers`/`uiMarkers` 휴리스틱 적용 (실제 존재 디렉토리만)
6. `detectionEvidence` 3~5개 기재
7. 마커 미발견 시 `primary: "unknown"`, `fallbackUsed: true`

수행하지 않을 범위 (feature-pipeline 파이프라인에서만 수행):

- 요구사항 분해, plan.md 작성, file-manifest.json 작성
- 코드 파일 수정·생성
- 사용자에게 파일 목록 승인 요청

### Step 3: 산출

`--inline` 미지정 시:

- 출력 경로(`<대상>/_workspaces/stack-profile.json` 또는 `--out` 지정 경로)에 JSON 작성
- 파일이 이미 존재하면 덮어쓰되, 직전 내용을 `stack-profile.prev.json`으로 백업

`--inline` 지정 시:

- 파일 작성 생략, 결과만 콘솔 출력

### Step 4: 결과 요약 출력

콘솔에 다음 형식으로 요약 출력 (항상 출력, `--inline` 여부 무관):

```
스택 감지 결과 ({대상 경로})
─────────────────────────────
primary/subtype : node/express
language        : typescript
packageManager  : pnpm
moduleSystem    : esm
testFramework   : vitest (vitest.config.ts)
linter/formatter: eslint / prettier
buildTool       : tsc
businessLogic   : src/services, src/middleware
ui              : (없음)
근거            : package.json:express@4.19.2, tsconfig.json, vitest.config.ts, pnpm-lock.yaml
fallbackUsed    : false
산출 파일       : ./_workspaces/stack-profile.json
```

`fallbackUsed: true`일 때:

```
스택 감지 실패 — 마커 파일을 식별할 수 없습니다.
다음 중 하나를 확인하세요:
  1. 대상 디렉토리가 맞는지: /my-poor-ai:detect-stack <올바른 경로>
  2. 마커 파일이 비표준 위치에 있는지
  3. 스택을 직접 명시: /my-poor-ai:detect-stack --primary node --subtype express
```

## 인수 처리 우선순위

1. `--inline` (boolean) — 출력 모드
2. `--out <경로>` — 출력 파일 경로
3. `--primary <값>` `--subtype <값>` — 수동 보정 (자동 감지 우회)
4. 첫 번째 위치 인수 — 대상 디렉토리

여러 인수 동시 사용 예시:

```
/my-poor-ai:detect-stack ./services/api --out ./services/api/_workspaces/profile.json
/my-poor-ai:detect-stack . --primary python --subtype fastapi --out ./_workspaces/manual-profile.json
```

## 활용 예시

### 예시 1: 빠른 확인

```
/my-poor-ai:detect-stack
```

→ 현재 디렉토리 스캔, 결과만 콘솔에 표시 (`_workspaces/stack-profile.json`도 생성). 다른 작업으로 넘어가기 전 프로젝트 스택을 빠르게 파악.

### 예시 2: 다른 도구의 입력으로 사용

```
/my-poor-ai:detect-stack . --out ./.cache/stack.json
```

→ 빌드 스크립트·CI·다른 에이전트가 참조할 수 있는 위치에 저장.

### 예시 3: feature-pipeline 사전 검증

```
/my-poor-ai:detect-stack
# 결과 확인 후
/feature-pipeline 기능 요청...
```

→ feature-pipeline이 같은 결과를 도출할지 미리 확인. 잘못 감지되면 수동 보정 후 진행.

## 에러 핸들링

| 상황                              | 대응                                                                |
| --------------------------------- | ------------------------------------------------------------------- |
| 대상 경로 미존재                  | 에러 메시지 출력 후 종료                                            |
| 마커 미발견                       | `fallbackUsed: true`로 작성, 안내 메시지 출력                       |
| 여러 framework 동시 발견          | `detectionEvidence`에 모두 기록, 가장 유력한 subtype 선택 + 안내    |
| 출력 경로의 부모 디렉토리 미존재  | 자동 생성 후 작성                                                   |
| `--primary` 수동 지정 시 매트릭스 외 값 | 그대로 작성 (사용자 책임), `detectionEvidence: ["manual"]` 기재 |

## 절대 금지

- 코드 파일(*.ts, *.py, *.cs 등) 수정·생성
- plan.md, file-manifest.json, review-report.md 작성 (파이프라인 산출물)
- 사용자에게 파일 목록 승인 요청 (이 커맨드는 게이트 없음)
- 절대 경로/`~/` 사용 (모든 경로는 CWD 기준 상대)

## 참조

- 매트릭스 정의: `feature-planner.md` "Phase 1.0 — 스택 감지 매트릭스" 섹션 (진실의 원천)
- 스키마: 동일 파일의 "stack-profile.json 스키마" 섹션
