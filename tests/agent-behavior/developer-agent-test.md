# Developer Agent Test

## 테스트 목적

`developer-agent`가 TDD(RED-GREEN-REFACTOR)를 실제로 준수하는지, 올바른 STATUS를 반환하는지 확인.

## 호출 방법

```
Agent(subagent_type="developer-agent"):
  spec-path: _workspaces/feature-tags/specs/spec-a.md
  branch-slug: feature-tags
  프로젝트 경로: {테스트용 프로젝트 경로}
  컨텍스트: branch feature-tags 기준, 의존 스펙 없음
```

## TDD 준수 검증

TDD 순서를 역추적할 수 있어야 합니다:

- [ ] 커밋 이력에서 "test:" 또는 failing test 커밋이 구현 커밋보다 앞에 있음
- [ ] 최종 테스트 실행 결과 전부 GREEN
- [ ] 스펙 범위 밖 코드 추가 없음

```bash
# TDD 순서 확인
git log --oneline feature-tags | head -10
# test 커밋이 feat 커밋보다 앞에 있어야 함

# 테스트 통과 확인
npm test  # 또는 프로젝트 테스트 명령어
```

## STATUS 검증

| 상황          | 기대 STATUS                           |
| ------------- | ------------------------------------- |
| 정상 완료     | DONE                                  |
| 완료+우려사항 | DONE_WITH_CONCERNS (우려 내용 구체적) |
| 스펙 불명확   | NEEDS_CONTEXT (질문 구체적)           |
| 진행 불가     | BLOCKED (이유+시도 기술)              |

## 실패 판정

- 구현 커밋이 테스트 커밋보다 앞에 있음 → TDD 미준수
- 테스트 실패 상태로 DONE 반환 → verification 미수행
- 스펙 범위 밖 리팩터링 포함 → YAGNI 위반
