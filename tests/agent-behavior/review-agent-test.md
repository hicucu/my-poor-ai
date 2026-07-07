# Review Agent Test

## 테스트 목적

`review-agent`가 4개 전문 reviewer를 **병렬**로 실행하고, aggregator로 통합 보고서를 생성하며, Critical 이슈 발견 시 `issue-fixer`를 올바르게 호출하는지 확인.

## 호출 방법

```
Agent(subagent_type="review-agent"):
  branch-slug: feature-tags
  base-branch: main
  프로젝트 경로: {테스트용 프로젝트 경로}
```

## 병렬 실행 검증

4개 reviewer가 순차가 아닌 병렬로 실행되었는지 확인:

- [ ] `_workspaces/review-feature-tags/reviews/architecture.md` 생성됨
- [ ] `_workspaces/review-feature-tags/reviews/security.md` 생성됨
- [ ] `_workspaces/review-feature-tags/reviews/performance.md` 생성됨
- [ ] `_workspaces/review-feature-tags/reviews/style.md` 생성됨
- [ ] `_workspaces/review-feature-tags/review-report.md` 생성됨 (aggregator)
- [ ] review-report.md에 파일별 이슈 그룹화 포함
- [ ] review-report.md에 Critical/High 우선순위 표 포함

## STATUS 검증

| 조건          | 기대 STATUS                            |
| ------------- | -------------------------------------- |
| Critical 0건  | APPROVED                               |
| Critical 1건+ | NEEDS_FIXES → issue-fixer 호출됨 |

## issue-fixer 연동 검증 (Critical 있는 경우)

- [ ] issue-fixer가 Critical 이슈 파일별로 호출됨
- [ ] 수정 후 재리뷰하지 않고 오케스트레이터에 반환 (재리뷰는 오케스트레이터 담당)
- [ ] FIXED STATUS에 수정된 파일 목록 포함

## 실패 판정

- 4개 파일 중 하나라도 없음 → reviewer 호출 오류
- review-report.md에 파일별 그룹화 없음 → aggregator 오류
- Critical 있는데 issue-fixer 미호출 → 연동 오류
