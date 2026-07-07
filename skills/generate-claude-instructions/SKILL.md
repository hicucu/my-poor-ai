---
name: generate-claude-instructions
description: "An orchestrator that generates CLAUDE.md and reference documents (DEVELOPMENT.md, LANGUAGE_GUIDELINES_{LANG}.md split per language, AI_BEHAVIOR.md, COMMIT_CONVENTION.md) into the instruction/ folder of the current working directory, based on the user's input (existing instruction files / project folder / none). Must be used for any request to create or update CLAUDE instruction documents, such as \"create a CLAUDE.md\", \"generate instructions\", \"write instruction documents\", \"analyze this project and create instructions\", \"rewrite based on existing instructions\", \"refresh the instructions\", \"regenerate AI behavior guidelines\", \"update the development principles document\", \"rewrite the per-language guidelines\", \"create a global CLAUDE.md\", \"generate reference documents\", \"supplement the instructions\". Also applies to re-run, regenerate, update, supplement, and partial-edit requests."
---

# CLAUDE Instruction Generator — 오케스트레이터

CLAUDE.md와 참조 문서를 사용자의 입력 컨텍스트(기존 지침 파일 또는 프로젝트 폴더)를 기반으로 생성함. 모든 산출물은 현재 작업 디렉토리(CWD)의 `instruction/` 폴더에 저장됨.

## 출력 위치

```
{CWD}/instruction/
├── CLAUDE.md                          ← 메인 지침 (모든 참조 문서 포함)
├── DEVELOPMENT.md                     ← 개발원칙 및 개발론
├── LANGUAGE_GUIDELINES_TYPESCRIPT.md  ← TypeScript 지침 (감지된 언어별 생성)
├── LANGUAGE_GUIDELINES_REACT.md       ← React 지침
├── LANGUAGE_GUIDELINES_NEXTJS.md      ← Next.js 지침
├── AI_BEHAVIOR.md                     ← AI 동작 지침
└── COMMIT_CONVENTION.md               ← 커밋 컨벤션 (commitlint 기준)
```

언어/프레임워크 파일은 프로젝트 스택에서 **감지된 것만** 생성함.
지원 파일명: `LANGUAGE_GUIDELINES_{TYPESCRIPT|REACT|NEXTJS|CSHARP|PYTHON|JAVA}.md`

`CWD`는 현재 Claude가 실행되고 있는 디렉토리. 이 스킬이 위치한 팀 디렉토리(`.claude/`)와는 별개.

## 팀 위치

my-poor-ai 플러그인 루트(`{팀_위치}`)의 `agents/`에서 에이전트 정의 파일을 읽는다. 플러그인이 어디에 설치되든 상대 위치로 작동.

## 경로 정책

**경로 정책**: `{팀_위치}`는 항상 절대 경로로 결정한다(`os.path.abspath()` 등). 에이전트에 전달 시 절대 경로 문자열 사용.

## Phase 0: 컨텍스트 확인

실행 전 다음을 결정함.

### 0-1. 출력 디렉토리 준비

- 현재 작업 디렉토리(CWD)를 확인
- `{CWD}/instruction/` 폴더 존재 여부 확인, 없으면 생성

### 0-2. 입력 참조 파악

사용자 메시지에서 입력 참조를 추출하여 모드를 결정함:

| 입력 유형          | 처리 방식                                  | 모드          |
| ------------------ | ------------------------------------------ | ------------- |
| 파일 경로(들) 명시 | 파일 내용을 읽어 기존 지침/컨벤션 추출     | **참조 모드** |
| 디렉토리 경로 명시 | 디렉토리 트리 + 설정 파일 + 샘플 소스 분석 | **분석 모드** |
| 입력 없음          | 일반 모범 사례 기반으로 작성               | **표준 모드** |
| 혼합               | 모두 활용                                  | **혼합 모드** |

입력 참조가 식별되면 사용자에게 한 문장으로 알림 (예: "프로젝트 폴더 `D:\my-project`를 분석하여 지침을 작성함").

### 0-3. 실행 모드 결정

- `{CWD}/instruction/`에 기존 파일 없음 → **초기 실행**
- 파일 있음 + 사용자가 전체 갱신 요청 → **전체 갱신** (덮어쓰기)
- 파일 있음 + 특정 문서 수정 요청 → **부분 수정** (해당 에이전트 + composer만)

## Phase 1: 병렬 문서 생성 (Fan-out)

**실행 모드: 서브 에이전트 4개 병렬**

4개 Agent를 동시에 `run_in_background: true`로 호출. 모두 완료될 때까지 대기. 모델은 에이전트별로 아래에 명시.

각 에이전트 프롬프트에 포함할 공통 컨텍스트:

- **출력 디렉토리**: `{CWD}/instruction/` (절대 경로)
- **입력 참조**: Phase 0-2에서 식별한 파일/디렉토리 경로 목록 (없으면 "없음")
- **실행 모드**: 초기/전체/부분 중 하나
- **이 팀의 위치**: 이 SKILL.md가 위치한 `.claude/`의 절대 경로 (에이전트 정의 파일을 읽기 위함)

### Agent 1 — 개발원칙 작성 `(model: opus)`

에이전트 정의: `{팀_위치}/agents/dev-principles.md`

프롬프트 골격:

```
{팀_위치}/agents/dev-principles.md를 읽고 그 지침에 따라 작업한다.

출력: {CWD}/instruction/DEVELOPMENT.md
입력 참조: [파일/디렉토리 경로 목록 또는 "없음"]
실행 모드: [초기/전체/부분]

입력 참조가 있으면 모두 읽고 컨텍스트를 반영한다.
디렉토리인 경우 소스코드 구조와 사용 패턴을 분석하여 해당 프로젝트에 최적화된 지침을 작성한다.
```

### Agent 2 — 언어별 지침 작성 `(model: opus)`

에이전트 정의: `{팀_위치}/agents/language-guidelines.md`

단일 파일이 아닌 **언어/프레임워크별 별도 파일**로 출력함.

프롬프트 골격:

```
{팀_위치}/agents/language-guidelines.md를 읽고 그 지침에 따라 작업한다.

출력 디렉토리: {CWD}/instruction/
출력 형식: 감지된 언어/프레임워크마다 LANGUAGE_GUIDELINES_{LANG}.md 파일 생성
  예: LANGUAGE_GUIDELINES_TYPESCRIPT.md, LANGUAGE_GUIDELINES_REACT.md, LANGUAGE_GUIDELINES_NEXTJS.md
입력 참조: [파일/디렉토리 경로 목록 또는 "없음"]
실행 모드: [초기/전체/부분]
```

완료 후 에이전트는 생성된 파일 목록을 오케스트레이터에 보고함.

### Agent 3 — AI 동작 지침 작성 `(model: haiku)`

에이전트 정의: `{팀_위치}/agents/ai-behavior.md`
출력: `{CWD}/instruction/AI_BEHAVIOR.md`
나머지 프롬프트는 Agent 1과 동일 패턴.

### Agent 4 — 커밋 컨벤션 작성 `(model: haiku)`

에이전트 정의: `{팀_위치}/agents/commit-convention.md`
출력: `{CWD}/instruction/COMMIT_CONVENTION.md`
나머지 프롬프트는 Agent 1과 동일 패턴. commitlint type 목록·검증 규칙·예시 포함.

## Phase 2: CLAUDE.md 합성 (Fan-in)

**실행 모드: 서브 에이전트 1개 순차**

Phase 1의 4개 에이전트 모두 완료 확인 후 실행.
Agent 2가 보고한 `LANGUAGE_GUIDELINES_*.md` 파일 목록을 확인하여 composer 프롬프트에 포함함.

`model: "haiku"`, 동기 실행. (합성 작업 — 프로젝트 스캔 불필요)

프롬프트 골격:

```
{팀_위치}/agents/claude-md-composer.md를 읽고 그 지침에 따라 작업한다.

출력: {CWD}/instruction/CLAUDE.md
입력 (모두 읽을 것):
- {CWD}/instruction/DEVELOPMENT.md
- {CWD}/instruction/LANGUAGE_GUIDELINES_*.md  (존재하는 파일 모두 — 언어별 분리 파일)
- {CWD}/instruction/AI_BEHAVIOR.md
- {CWD}/instruction/COMMIT_CONVENTION.md
- [Phase 0-2의 입력 참조 (있으면)]

CLAUDE.md는 사용자가 전역(~/.claude/)으로 복사하여 사용할 것이므로,
참조 문서 경로는 같은 디렉토리에 있다고 가정하고 상대 경로/파일명만으로 작성한다.
언어별 지침 파일은 각 파일명을 명시적으로 나열한다 (glob 패턴 대신).
오늘 날짜(YYYY-MM-DD)를 변경 이력에 기록.
```

## Phase 3: 검증 및 완료 보고

1. 필수 파일 생성 확인:
   - `{CWD}/instruction/DEVELOPMENT.md`
   - `{CWD}/instruction/LANGUAGE_GUIDELINES_*.md` (1개 이상)
   - `{CWD}/instruction/AI_BEHAVIOR.md`
   - `{CWD}/instruction/COMMIT_CONVENTION.md`
   - `{CWD}/instruction/CLAUDE.md`
2. CLAUDE.md가 생성된 언어별 지침 파일을 모두 언급하는지 확인
3. 각 파일 줄 수 간단 표시
4. 완료 보고:

```
지침서 생성 완료
  위치: {CWD}/instruction/
  생성 파일:
    - DEVELOPMENT.md
    - LANGUAGE_GUIDELINES_TYPESCRIPT.md  (감지된 언어별)
    - LANGUAGE_GUIDELINES_REACT.md
    - AI_BEHAVIOR.md
    - COMMIT_CONVENTION.md
    - CLAUDE.md

재실행: "지침서 갱신해줘" 또는 "TypeScript 지침만 다시 써줘"
```

## 에러 핸들링

| 상황                       | 대응                                     |
| -------------------------- | ---------------------------------------- |
| Agent 1개 실패             | 나머지 완료 후 해당 에이전트 1회 재시도  |
| 재시도 후도 실패           | 해당 파일 없이 CLAUDE.md 작성, 누락 명시 |
| composer 실패              | 기존 CLAUDE.md(있으면) 보존, 수동 안내   |
| 입력 참조 파일/폴더 미존재 | 표준 모드로 폴백, 사용자에게 알림        |
| 출력 폴더 쓰기 실패        | 권한 확인 요청                           |
| 언어 감지 실패 (표준 모드) | TypeScript + React 기본 파일 2개 생성    |

## 테스트 시나리오

**시나리오 1 — 표준 모드 (입력 없음):**
"CLAUDE.md 만들어줘" → Phase 0(표준, 초기) → Phase 1(4개 병렬) → Phase 2(composer) → LANGUAGE_GUIDELINES_TYPESCRIPT.md + LANGUAGE_GUIDELINES_REACT.md 기본 생성

**시나리오 2 — 분석 모드 (Next.js 프로젝트):**
"D:\my-project 분석해서 지침 만들어줘" → Phase 0(분석) → Agent 2가 package.json 분석 → LANGUAGE_GUIDELINES_TYPESCRIPT.md + LANGUAGE_GUIDELINES_REACT.md + LANGUAGE_GUIDELINES_NEXTJS.md 생성

**시나리오 3 — 부분 갱신:**
"TypeScript 지침만 다시 써줘" → Phase 0(부분) → Agent 2만 실행 (LANGUAGE_GUIDELINES_TYPESCRIPT.md만) → Phase 2(composer 재실행)

**시나리오 4 — C# 프로젝트:**
"D:\dotnet-api 분석해서 지침 만들어줘" → Agent 2가 \*.csproj 감지 → LANGUAGE_GUIDELINES_CSHARP.md 생성
