---
name: writing-skills
description: Use when creating a new skill, editing an existing skill, or validating a skill before deployment
---

# 스킬 작성

## 개요

**스킬 작성은 프로세스 문서에 적용된 TDD임.**

**개인 스킬은 에이전트별 디렉토리에 위치 (Claude Code의 경우 `~/.claude/skills`, Codex의 경우 `~/.agents/skills/`)**

테스트 케이스(eval 케이스 또는 subagent를 활용한 압박 시나리오)를 작성하고, 실패를 관찰하고(기준 동작), 스킬(문서)을 작성하고, 테스트 통과를 관찰하고(에이전트 준수), 리팩터링(허점 제거)함.

**핵심 원칙:** 스킬 없이 에이전트가 실패하는 것을 직접 보지 않았다면, 스킬이 올바른 것을 가르치는지 알 수 없음.

**필수 배경:** 이 스킬을 사용하기 전에 my-poor-ai:test-driven-development를 반드시 이해해야 함. 해당 스킬이 근본적인 RED-GREEN-REFACTOR 사이클을 정의함. 이 스킬은 TDD를 문서에 적용.

**공식 가이드:** Anthropic의 공식 스킬 작성 모범 사례는 anthropic-best-practices.md를 참조. 이 문서는 이 스킬의 TDD 중심 접근 방식을 보완하는 추가 패턴과 가이드라인을 제공.

## 스킬이란?

**스킬**은 검증된 기법, 패턴, 도구에 대한 참조 가이드. 스킬은 미래의 Claude 인스턴스가 효과적인 접근 방식을 찾고 적용하는 데 도움을 줌.

**스킬은:** 재사용 가능한 기법, 패턴, 도구, 참조 가이드

**스킬은 아님:** 한 번 문제를 해결한 방법에 대한 서술

## 스킬을 위한 TDD 매핑

| TDD 개념                | 스킬 생성                                          |
| ----------------------- | -------------------------------------------------- |
| **테스트 케이스**       | eval 케이스 또는 subagent를 활용한 압박 시나리오   |
| **프로덕션 코드**       | 스킬 문서 (SKILL.md)                               |
| **테스트 실패 (RED)**   | 스킬 없이 에이전트가 규칙 위반 (기준)              |
| **테스트 통과 (GREEN)** | 스킬이 있을 때 에이전트 준수                       |
| **리팩터링**            | 준수를 유지하면서 허점 제거                        |
| **테스트 먼저 작성**    | 스킬 작성 전 기준 시나리오 실행                    |
| **실패 관찰**           | 에이전트가 자연스럽게 사용하는 합리화 문서화       |
| **최소 코드**           | 해당 특정 위반을 다루는 스킬 작성                  |
| **통과 관찰**           | 에이전트가 이제 준수하는지 확인                    |
| **리팩터링 사이클**     | 새로운 합리화 발견 → 보완 → 재검증                 |

전체 스킬 생성 프로세스는 RED-GREEN-REFACTOR를 따름.

## 스킬 생성 시점

**생성 시:**

- 직관적으로 명확하지 않았던 기법
- 여러 프로젝트에 걸쳐 다시 참조할 내용
- 패턴이 광범위하게 적용됨 (프로젝트 특화 아님)
- 다른 사람들에게도 도움이 될 내용

**생성 안 함:**

- 일회성 해결책
- 다른 곳에 잘 문서화된 표준 관행
- 프로젝트별 관례 (CLAUDE.md에 포함)
- 기계적 제약 (정규식/검증으로 강제할 수 있다면 자동화할 것 — 문서는 판단이 필요한 경우에 사용)

## 스킬 유형

- **기법** — 구체적인 방법과 따를 단계 (condition-based-waiting, root-cause-tracing)
- **패턴** — 문제를 바라보는 방식 (flatten-with-flags, test-invariants)
- **참조** — API 문서, 구문 가이드, 도구 문서 (office docs)

## 디렉토리 구조

```
skills/
  skill-name/
    SKILL.md              # 주요 참조 (필수)
    references/           # 대용량 참조 (필요한 경우에만)
    supporting-file.*     # 재사용 도구 (필요한 경우에만)
```

**플랫 네임스페이스** — 검색 가능한 하나의 네임스페이스에 모든 스킬

**별도 파일 사용 시:** 대용량 참조 (100줄 이상 — API 문서, 포괄적인 구문), 재사용 가능한 도구 (스크립트, 유틸리티, 템플릿)

**인라인 유지:** 원칙과 개념, 코드 패턴 (50줄 미만), 그 외 모든 것

본문 골격·순서도·코드 예시·파일 구성·안티 패턴은 `references/authoring-patterns.md` 참조.

## Frontmatter

- 필수: `name`과 `description` — Agent Skills 표준(claude.ai·API 업로드)에서 필수이며, 이 저장소는 validator가 `name` = 디렉토리명을 강제함
- 최대 1024자
- `name`: 소문자·숫자·하이픈만 사용 (괄호, 특수 문자 없음)
- `description`: 3인칭, 언제 사용할지만 설명 (무엇을 하는지 아님)
  - 트리거 조건에 집중하기 위해 "Use when..."으로 시작
  - 구체적인 증상, 상황, 맥락 포함
  - **스킬의 프로세스나 워크플로우 요약 절대 금지** (아래 CSO 절 참조)
  - 가능하면 500자 이하 유지
- Claude Code는 `when_to_use`·`disable-model-invocation`·`allowed-tools`·`context: fork` 등을 추가로 받음. `when_to_use`도 `description`과 함께 상시 목록에 들어가므로 같은 규칙을 적용함. 전체 필드와 플랫폼별 인식 범위는 `references/frontmatter.md` 참조

## Claude Search Optimization (CSO)

**발견에 중요:** 미래의 Claude가 스킬을 찾을 수 있어야 함

**목적:** Claude는 주어진 작업에 어떤 스킬을 불러올지 결정하기 위해 description을 읽음. "지금 이 스킬을 읽어야 하는가?"에 답하도록 만들 것.

**형식:** 트리거 조건에 집중하기 위해 "Use when..."으로 시작

**중요: Description = 사용 시점, 스킬이 하는 일이 아님**

description은 트리거 조건만 설명해야 함. 스킬의 프로세스나 워크플로우를 요약하지 말 것.

**왜 중요한가:** 테스트에서 description이 스킬의 워크플로우를 요약할 때, Claude가 전체 스킬 내용을 읽는 대신 description을 따를 수 있음을 확인. "code review between tasks"라는 description이 스킬의 순서도에서 두 번의 리뷰(스펙 준수 후 코드 품질)를 명확히 보여줌에도 불구하고 Claude가 한 번만 리뷰하게 만듦.

description을 "Use when executing implementation plans with independent tasks" (워크플로우 요약 없음)으로 변경했을 때, Claude가 순서도를 올바르게 읽고 두 단계 리뷰 프로세스를 따름.

**함정:** 워크플로우를 요약하는 description은 Claude가 취할 지름길을 만듦. 스킬 본문이 Claude가 건너뛰는 문서가 됨.

```yaml
# ❌ BAD: Summarizes workflow - Claude may follow this instead of reading skill
description: Use when executing plans - dispatches subagent per task with code review between tasks

# ❌ BAD: Too much process detail
description: Use for TDD - write test first, watch it fail, write minimal code, refactor

# ✅ GOOD: Just triggering conditions, no workflow summary
description: Use when executing implementation plans with independent tasks in the current session

# ✅ GOOD: Triggering conditions only
description: Use when implementing any feature or bugfix, before writing implementation code
```

**내용:**

- 이 스킬이 적용됨을 알리는 구체적인 트리거, 증상, 상황 사용
- _언어별 증상_ (setTimeout, sleep)이 아닌 _문제_ (race condition, 불일치 동작) 설명
- 스킬 자체가 특정 기술에 한정되지 않는 한 트리거는 기술 중립적으로 유지
- 기술 특화 스킬이면 트리거에서 명시적으로 밝힐 것
- 3인칭으로 작성 (시스템 프롬프트에 주입됨)
- **스킬의 프로세스나 워크플로우 절대 요약 금지**

키워드 커버리지·이름 짓기·토큰 효율성·교차 참조 규칙은 `references/cso.md` 참조.

## 철칙 (TDD와 동일)

```
실패하는 테스트 없이는 스킬도 없다
```

새로운 스킬과 기존 스킬 편집 모두에 적용.

테스트 전에 스킬 작성? 삭제. 처음부터 다시.
테스트 없이 스킬 편집? 동일한 위반.

**예외 없음:**

- "간단한 추가"라도 예외 없음
- "섹션만 추가"라도 예외 없음
- "문서 업데이트"라도 예외 없음
- 테스트되지 않은 변경 사항을 "참조"로 유지 금지
- 테스트 실행 중 "적용" 금지
- 삭제는 삭제를 의미함

**필수 배경:** my-poor-ai:test-driven-development 스킬이 이것이 중요한 이유를 설명함. 동일한 원칙이 문서에도 적용됨.

## eval로 RED/GREEN 실행

플러그인 스킬은 `claude plugin eval`(Claude Code 2.1.269+)로 RED/GREEN을 반복 실행하고 점수화함. 이 저장소의 케이스는 `evals/`에 있음.

| TDD 단계 | eval 대응 |
| --- | --- |
| 테스트 케이스 | `evals/<case>/prompt.md` (frontmatter `max_turns`·`allowed_tools`·`tags`·`runs`) + `evals/<case>/graders/*.md` |
| RED (스킬 없이 실패) | `--ablation with-without`(기본)의 no-plugin 기준 arm |
| GREEN (스킬로 통과) | 같은 실행의 with arm — 두 arm의 점수 차이가 스킬의 효과 |
| 트리거 확인 | `type: tool_used` + `tool: Skill` grader. ablation 실행에서는 점수가 아니라 "플러그인 발동 지표"로만 집계되므로 트리거 케이스는 `--ablation none`으로 따로 실행 |
| 행동 판정 | `type: llm` grader(본문에 기대 행동 기준) 또는 `type: regex` |

```bash
claude plugin eval . --tag trigger --ablation none --no-publish      # 트리거
claude plugin eval . --tag pressure --runs 3 --no-publish             # 압박 시나리오 RED/GREEN
```

샌드박스는 사용자 CLAUDE.md·설정을 배제하고 플러그인 SessionStart 훅은 실행함. 다중 턴 압박처럼 케이스로 표현하기 어려운 시나리오는 서브에이전트 수동 테스트(testing-skills-with-subagents.md)로 보완함.

스킬 유형(규율·기법·패턴·참조)별 테스트 방법과 성공 기준은 `references/testing-by-skill-type.md` 참조.

## 테스트 건너뛰기에 대한 흔한 합리화

| 변명                         | 현실                                                                      |
| ---------------------------- | ------------------------------------------------------------------------- |
| "내 스킬은 명백히 명확하다"  | 나에게 명확 ≠ 다른 에이전트에게 명확. 테스트할 것.                        |
| "그냥 참조일 뿐이다"         | 참조에도 공백과 불명확한 섹션이 있을 수 있음. 검색을 테스트할 것.         |
| "테스트는 과도하다"          | 테스트되지 않은 스킬에는 문제가 있음. 항상. 15분 테스트가 몇 시간을 아낌. |
| "문제가 생기면 테스트하겠다" | 문제 = 에이전트가 스킬을 사용할 수 없음. 배포 전에 테스트할 것.           |
| "테스트는 너무 번거롭다"     | 테스트는 프로덕션에서 나쁜 스킬을 디버깅하는 것보다 덜 번거로움.          |
| "품질이 좋다고 확신한다"     | 과신은 문제를 보장함. 어쨌든 테스트할 것.                                 |
| "학문적 검토로 충분하다"     | 읽기 ≠ 사용하기. 적용 시나리오를 테스트할 것.                             |
| "테스트할 시간이 없다"       | 테스트되지 않은 스킬을 배포하면 나중에 수정하는 데 더 많은 시간이 소요됨. |

**이 모든 것의 의미: 배포 전 테스트. 예외 없음.**

## 합리화에 대한 스킬 방탄화

규율을 강화하는 스킬(TDD 등)은 합리화에 저항해야 함. 에이전트는 압박 하에서 허점을 찾을 만큼 영리함.

**심리학 노트:** 설득 기법이 작동하는 이유를 이해하면 체계적으로 적용하는 데 도움이 됨. 권위, 헌신, 희소성, 사회적 증거, 통일성 원칙에 관한 연구 기반(Cialdini, 2021; Meincke et al., 2025)은 persuasion-principles.md 참조.

### 모든 허점을 명시적으로 제거

규칙을 명시하는 것만으로는 부족 — 특정 우회책을 금지:

<Bad>
```markdown
Write code before test? Delete it.
```
</Bad>

<Good>
```markdown
Write code before test? Delete it. Start over.

**No exceptions:**

- Don't keep it as "reference"
- Don't "adapt" it while writing tests
- Don't look at it
- Delete means delete

````
</Good>

### "정신 vs 글자" 주장 다루기

초반에 기초 원칙 추가:

```markdown
**규칙의 글자를 위반하는 것은 규칙의 정신을 위반하는 것임.**
````

이것으로 "정신을 따르고 있다"는 합리화 전체 범주를 차단.

### 합리화 표 구축

기준 테스트에서 합리화를 포착 (아래 테스트 섹션 참조). 에이전트가 하는 모든 변명을 표에 기록:

```markdown
| 변명                                    | 현실                                                                                   |
| --------------------------------------- | -------------------------------------------------------------------------------------- |
| "테스트하기엔 너무 간단하다"            | 간단한 코드도 깨짐. 테스트에 30초면 됨.                                                |
| "나중에 테스트하겠다"                   | 즉시 테스트 통과는 아무것도 증명하지 않음.                                             |
| "테스트 후 작성도 같은 목표를 달성한다" | 테스트 후 작성 = "이것이 무엇을 하는가?" 테스트 선 작성 = "이것이 무엇을 해야 하는가?" |
```

### 위험 신호 목록 생성

에이전트가 합리화할 때 자가 점검을 쉽게 만들기:

```markdown
## 위험 신호 - 중단하고 처음부터 다시 시작

- 테스트 전 코드
- "이미 수동으로 테스트했다"
- "테스트 후 작성도 같은 목적을 달성한다"
- "의식이 아닌 정신에 관한 것이다"
- "이것은 다르다 왜냐하면..."

**이 모든 것의 의미: 코드 삭제. TDD로 처음부터 다시 시작.**
```

### 위반 증상에 대한 CSO 업데이트

description에 규칙을 곧 위반할 것 같을 때의 증상 추가:

```yaml
description: use when implementing any feature or bugfix, before writing implementation code
```

## RED-GREEN-REFACTOR for Skills

TDD 사이클 따르기:

### RED: 실패하는 테스트 작성 (기준)

스킬 없이 압박 시나리오 실행 (eval 기준 arm 또는 subagent). 정확한 동작 문서화:

- 어떤 선택을 했는가?
- 어떤 합리화를 사용했는가 (그대로)?
- 어떤 압박이 위반을 유발했는가?

이것이 "테스트 실패 관찰" — 스킬 작성 전에 에이전트가 자연스럽게 무엇을 하는지 반드시 봐야 함.

### GREEN: 최소 스킬 작성

그 특정 합리화를 다루는 스킬 작성. 가상의 경우를 위한 추가 내용 금지.

스킬이 있는 동일한 시나리오 실행. 에이전트가 이제 준수해야 함.

### REFACTOR: 허점 제거

에이전트가 새로운 합리화를 찾았는가? 명시적 반론 추가. 방탄화될 때까지 재테스트.

**테스트 방법론:** 압박 시나리오 작성 방법, 압박 유형 (시간, 매몰 비용, 권위, 피로), 체계적인 허점 제거, 메타 테스트 기법은 testing-skills-with-subagents.md 참조.

## 중단: 다음 스킬로 넘어가기 전

**어떤 스킬을 작성한 후에도, 반드시 중단하고 배포 프로세스를 완료해야 함.**

**금지:**

- 각 스킬을 테스트하지 않고 여러 스킬을 일괄 생성
- 현재 스킬 검증 전에 다음 스킬로 이동
- "일괄 처리가 더 효율적이다"는 이유로 테스트 건너뛰기

**배포 체크리스트는 각 스킬에 대해 필수.**

테스트되지 않은 스킬 배포 = 테스트되지 않은 코드 배포. 품질 기준 위반.

## 스킬 생성 체크리스트

**필수:** `references/checklist.md`를 읽고 각 항목(RED·GREEN·REFACTOR·품질 확인·배포)을 할 일로 만들어 추적함 — 작업 추적 도구가 있으면 그 도구, 없으면 응답 내 체크리스트.

## 참조 파일

| 파일 | 내용 |
| --- | --- |
| `references/frontmatter.md` | Claude Code frontmatter 전 필드, 플랫폼별 인식 범위, 상시 비용 |
| `references/cso.md` | 키워드 커버리지, 이름 짓기, 토큰 효율성, 교차 참조, 발견 워크플로우 |
| `references/authoring-patterns.md` | 본문 골격, 순서도, 코드 예시, 파일 구성, 안티 패턴 |
| `references/testing-by-skill-type.md` | 스킬 유형별 테스트 방법과 성공 기준 |
| `references/checklist.md` | 스킬 생성 체크리스트 |
| `testing-skills-with-subagents.md` | 서브에이전트 수동 압박 테스트 방법론 |
| `persuasion-principles.md` | 설득 원칙 연구 기반 |
| `anthropic-best-practices.md` | Anthropic 공식 스킬 작성 모범 사례 (번역) |

## 결론

**스킬 생성은 프로세스 문서를 위한 TDD임.**

동일한 철칙: 실패하는 테스트 없이는 스킬도 없음.
동일한 사이클: RED (기준) → GREEN (스킬 작성) → REFACTOR (허점 제거).
동일한 이점: 높은 품질, 적은 놀람, 방탄화된 결과.

코드에 TDD를 따른다면, 스킬에도 따를 것. 동일한 규율을 문서에 적용하는 것임.
