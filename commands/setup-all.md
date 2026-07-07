---
description: graphify-setup, generate-claudeignore 커맨드를 순차 실행하는 통합 커맨드. 코드 그래프 도구(graphifyy 또는 codegraph) 선택·설치·그래프 생성·Claude 통합·git hook·.gitignore 설정까지 순차 실행. 신규 프로젝트 onboarding 또는 머신 초기 설정에 사용.
---

# Setup All (graphify-setup + generate-claudeignore 통합 실행)

`/my-poor-ai:graphify-setup`, `/my-poor-ai:generate-claudeignore` 커맨드를 순차 실행하는 통합 진입점.
각 단계는 이미 완료된 경우 자동으로 건너뛰므로 재실행해도 안전함.

> SessionStart 훅 등록은 my-poor-ai 플러그인 설치 시 `hooks/hooks.json`으로 자동 처리되므로 별도 설치 단계가 없음.

## 사용법

```
/my-poor-ai:setup-all              graphify-setup → generate-claudeignore 순차 실행 (각 단계 완료 시 건너뜀)
/my-poor-ai:setup-all check        두 커맨드의 현재 설치 상태만 확인하고 종료
/my-poor-ai:setup-all reset        모든 단계 강제 재실행 (기존 설정 덮어쓰기)
/my-poor-ai:setup-all graphify     graphify-setup만 graphifyy 도구로 실행
/my-poor-ai:setup-all codegraph    graphify-setup만 codegraph 도구로 실행
```

## 실행 순서

### Step 1: 사전 확인

다음 항목 모두 확인:

| 항목                  | 확인 방법                                            | 실패 시                         |
| --------------------- | ---------------------------------------------------- | ------------------------------- |
| 현재 CWD가 git 저장소 | `git rev-parse --is-inside-work-tree`                | graphify-setup 단계 생략 안내   |
| pip 또는 npm 존재     | `pip --version` / `pip3 --version` / `npm --version` | 선택 도구에 따라 종료 또는 경고 |

> my-poor-ai가 brainstorming·writing-plans·test-driven-development 등 필요한 skill을 자체 제공하므로,
> 외부 플러그인 설치 단계는 필요하지 않음.

### Step 2: graphify-setup 실행

`check`/`claudeignore` 인수가 아니면, `/my-poor-ai:graphify-setup` 커맨드를 동일한 인수와 함께 호출.

```
/my-poor-ai:graphify-setup            (인수 없음 → 도구 선택 후 6단계 자동 진행)
/my-poor-ai:graphify-setup graphify   (graphifyy 도구 지정)
/my-poor-ai:graphify-setup codegraph  (codegraph 도구 지정)
/my-poor-ai:graphify-setup check      (--check 모드)
/my-poor-ai:graphify-setup reset      (--reset 모드)
```

graphify-setup의 6단계: 도구 선택 → 패키지 설치 → 그래프 생성 → Claude 통합 → git hook → .gitignore.
도구 선택에 따라 각 단계의 명령이 분기됨 (graphifyy: pip 기반 / codegraph: npm + MCP 기반).

산출 요약을 그대로 콘솔에 출력. 오류 발생 시 에러 보고 후 generate-claudeignore 단계 진행 여부를 사용자에게 확인.

### Step 3: generate-claudeignore 실행

`check` 인수가 아니면, `/my-poor-ai:generate-claudeignore` 커맨드를 호출함.
기존 `.claudeignore`가 있으면 누락 항목만 병합하고, 없으면 신규 생성함.

```
/my-poor-ai:generate-claudeignore     (인수 없음 → 현재 CWD 기준 생성/병합)
/my-poor-ai:generate-claudeignore check  (--check 모드 → 미리보기만)
/my-poor-ai:generate-claudeignore reset  (--reset 모드 → 전체 재생성)
```

### Step 4: 프로젝트 컨텍스트 탐색

`_workspaces/` 디렉토리 생성 (이미 존재하면 생략):

```bash
mkdir -p _workspaces
```

`project-context` 에이전트 호출:

```
{팀_위치}/agents/project-context.md를 읽고 그 지침에 따라 작업한다.
CWD: {현재 프로젝트 루트}
mode: full
output: _workspaces/project-context.md
```

완료 후 `_workspaces/project-context.md` 생성 여부 확인.

### Step 5: 통합 요약

두 커맨드 결과를 합쳐 다음 형식으로 출력. `<도구>` 자리에 선택된 도구명(graphifyy / codegraph)을 기입함.

```
setup-all 완료 — <CWD>
─────────────────────────────────────────────
graphify-setup  [<도구> 사용]
  패키지 설치          {결과}
  그래프 생성          {결과}
  Claude 통합          {결과}
  git hook 등록        {결과}
  .gitignore 설정      {결과}

generate-claudeignore
  .claudeignore        {결과 — 신규 생성 / 병합 N개 항목 추가 / 기존 유지}

다음 단계
  - 수동 확인: /my-poor-ai:graphify-setup check, /my-poor-ai:generate-claudeignore check
```

`{결과}` 자리에 각 단계 실행 결과 (✓ 완료 / — 기존 유지 / ✗ 실패).

## 인수 전파 규칙

**인수 전파**: `check` 또는 `reset` 인수가 setup-all에 전달되면 graphify-setup, generate-claudeignore 호출 시 동일 인수를 그대로 전달함. 도구 지정 인수(`graphify`, `codegraph`, `claudeignore`)는 해당 단계에만 적용되고 다른 단계에는 전달하지 않음.

## 인수 처리 우선순위

1. `check` — 두 커맨드 모두 check 모드, 설치/생성 안 함
2. `reset` — 두 커맨드 모두 reset 모드, 강제 재실행
3. `graphify` — graphify-setup을 graphifyy 도구로 실행 (generate-claudeignore도 실행)
4. `codegraph` — graphify-setup을 codegraph 도구로 실행 (generate-claudeignore도 실행)
5. `claudeignore` — generate-claudeignore만 실행 (graphify-setup 생략)
6. (인수 없음) — 두 커맨드 순차 실행 (graphify-setup 내부에서 도구 선택)

## 활용 예시

### 신규 프로젝트 onboarding

```
/my-poor-ai:setup-all
```

→ 코드 그래프 도구(graphifyy 또는 codegraph) 선택·설치 + 그래프 생성 + Claude 통합 + git hook + .claudeignore 설정을 한 번에. 약 1-2분 소요.

### 머신 변경 후 동기화

```
/my-poor-ai:setup-all check
```

→ 두 자산의 현재 상태만 확인. 어떤 단계가 미설치인지 파악.

### 부분 재실행

```
/my-poor-ai:setup-all claudeignore
```

→ graphify-setup은 건드리지 않고 generate-claudeignore만 재실행. ignore 규칙이 바뀐 경우 유용.

## 에러 핸들링

| 상황                                                | 대응                                        |
| --------------------------------------------------- | ------------------------------------------- |
| graphify-setup 실패 (예: 패키지 설치 권한 부족)     | 에러 보고 + generate-claudeignore 진행 여부 확인 |
| generate-claudeignore 실패 (예: 쓰기 권한 없음)     | 에러 보고 + 통합 요약에 ✗ 표시              |
| pip/npm 모두 미설치                                 | graphify-setup 종료, 설치 안내 출력         |
| 인수 잘못 사용 (예: `/my-poor-ai:setup-all check reset`) | 사용법 출력 후 종료                         |

## 절대 금지

- 두 커맨드를 동시(병렬) 실행 (설정 파일 충돌 가능)
- 사용자 확인 없이 reset 강제 (덮어쓰기는 항상 사용자 의도 확인)

## 참조

- 단독 커맨드: `commands/graphify-setup.md`, `commands/generate-claudeignore.md`
