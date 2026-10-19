# Claude Search Optimization (CSO) — 상세

SKILL.md "CSO" 절의 description 규칙을 보완하는 참조임. 핵심 규칙(description = 사용 시점, 워크플로우 요약 금지)은 SKILL.md에 있음.

## 키워드 커버리지

Claude가 검색할 단어 사용:

- 에러 메시지: "Hook timed out", "ENOTEMPTY", "race condition"
- 증상: "flaky", "hanging", "zombie", "pollution"
- 동의어: "timeout/hang/freeze", "cleanup/teardown/afterEach"
- 도구: 실제 명령어, 라이브러리 이름, 파일 유형

## 설명적 이름 짓기

**능동태, 동사 시작:**

- ✅ `creating-skills` not `skill-creation`
- ✅ `condition-based-waiting` not `async-test-helpers`

**행동이나 핵심 통찰에 따라 이름 짓기:**

- ✅ `condition-based-waiting` > `async-test-helpers`
- ✅ `using-skills` not `skill-usage`
- ✅ `flatten-with-flags` > `data-structure-refactoring`
- ✅ `root-cause-tracing` > `debugging-techniques`

**동명사(-ing)는 프로세스에 잘 맞음:** `creating-skills`, `testing-skills`, `debugging-with-logs` — 능동적이고, 수행 중인 행동을 설명

## 토큰 효율성

**상시 비용이 드는 곳:**

- 세션 시작 훅이 주입하는 블록 (이 플러그인의 `skills/using-my-poor-ai/router.md`) — 모든 대화에 들어감. 판단에 필요한 것만 두고 상세는 스킬로 위임
- 모델이 호출할 수 있는 모든 스킬의 `description`·`when_to_use` — 모든 대화의 스킬 목록에 들어감 (`references/frontmatter.md`)

**호출 시 비용이 드는 곳:** SKILL.md 본문. 500줄 미만 유지 (`anthropic-best-practices.md`), 넘으면 보조 파일로 분리.

**기법:**

**상세 내용을 도구 도움말로 이동:**

```bash
# ❌ BAD: Document all flags in SKILL.md
search-conversations supports --text, --both, --after DATE, --before DATE, --limit N

# ✅ GOOD: Reference --help
search-conversations supports multiple modes and filters. Run --help for details.
```

**교차 참조 사용:**

```markdown
# ❌ BAD: Repeat workflow details

When searching, dispatch subagent with template...
[20 lines of repeated instructions]

# ✅ GOOD: Reference other skill

Always use subagents (50-100x context savings). REQUIRED: Use [other-skill-name] for workflow.
```

**예시 압축:**

```markdown
# ❌ BAD: Verbose example (42 words)

your human partner: "How did we handle authentication errors in React Router before?"
You: I'll search past conversations for React Router authentication patterns.
[Dispatch subagent with search query: "React Router authentication error handling 401"]

# ✅ GOOD: Minimal example (20 words)

Partner: "How did we handle auth errors in React Router?"
You: Searching...
[Dispatch subagent → synthesis]
```

**중복 제거:**

- 교차 참조된 스킬에 있는 내용 반복 금지
- 명령어에서 명백한 것 설명 금지
- 동일 패턴의 여러 예시 포함 금지

**검증:**

```bash
wc -l skills/<name>/SKILL.md    # 본문 500줄 미만
```

## 다른 스킬 교차 참조

스킬 이름만 사용하고, 명시적인 필수 표시 사용:

- ✅ 좋음: `**필수 서브 스킬:** my-poor-ai:test-driven-development 사용`
- ✅ 좋음: `**필수 배경:** my-poor-ai:systematic-debugging을 반드시 이해해야 함`
- ❌ 나쁨: `See skills/testing/test-driven-development` (필수 여부 불명확)
- ❌ 나쁨: `@skills/testing/test-driven-development/SKILL.md` (강제 로드, 컨텍스트 소모)

**@ 링크를 사용하지 않는 이유:** `@` 구문은 필요하기 전에 파일을 즉시 로드하여 컨텍스트를 소모함. 같은 스킬의 보조 파일도 `@` 없이 경로로만 가리킴.

## 발견 워크플로우

미래의 Claude가 스킬을 찾는 방법:

1. **문제 발생** ("tests are flaky")
2. **스킬 찾기** (description이 일치)
3. **개요 스캔** (관련성 있는가?)
4. **패턴 읽기** (빠른 참조 표)
5. **예시 로드** (구현 시에만)

**이 흐름에 최적화** — 검색 가능한 용어를 초반에, 자주 배치.
