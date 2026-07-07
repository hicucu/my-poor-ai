# Agent Behavior Tests

각 subagent가 올바른 STATUS를 반환하고, 파이프라인 상태가 정확하게 추적되는지 검증합니다.

## 테스트 구조

```
agent-behavior/
  brainstorming-agent-test.md   — design.md 생성 + DONE 반환
  planning-agent-test.md        — specs/ 생성 + file-manifest.json
  developer-agent-test.md       — TDD 준수 + STATUS 반환
  review-agent-test.md          — 4개 병렬 리뷰 + aggregator
  pipeline-state-test.md        — pipeline-state.md 추적 정확성
```

## 테스트 방법

각 에이전트를 Claude Code에서 직접 호출:

```
Agent(subagent_type="brainstorming-agent"):
  요구사항: {테스트 케이스 요구사항}
  프로젝트 경로: {테스트 프로젝트 경로}
  branch-slug: test-feature
```

## STATUS 판정 기준

| STATUS             | 의미          | 허용 조건                        |
| ------------------ | ------------- | -------------------------------- |
| DONE               | 완료          | 모든 완료 기준 충족              |
| DONE_WITH_CONCERNS | 완료+우려     | 우려 내용이 구체적으로 기술됨    |
| NEEDS_CONTEXT      | 컨텍스트 부족 | 질문이 구체적이고 블로킹됨       |
| BLOCKED            | 진행 불가     | 이유가 명확하고 시도한 것 기술됨 |
