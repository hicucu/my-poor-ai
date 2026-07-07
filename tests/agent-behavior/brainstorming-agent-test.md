# Brainstorming Agent Test

## 테스트 목적

`brainstorming-agent`가 사용자에게 직접 질문하지 않고 설계안을 생성하여 `design.md`를 저장하는지 확인.

## 호출 방법

```
Agent(subagent_type="brainstorming-agent"):
  요구사항: 블로그 포스트에 태그 기능 추가. 태그로 검색 가능해야 함.
  프로젝트 경로: {테스트용 프로젝트 경로}
  branch-slug: feature-tags
```

## 검증 항목

- [ ] `_workspaces/feature-tags/design.md` 파일 생성됨
- [ ] design.md에 2-3가지 접근법 포함
- [ ] 하나의 접근법에 명확한 추천 명시
- [ ] 변경 파일 목록 포함
- [ ] STATUS: DONE 반환
- [ ] 사용자에게 직접 질문 없음 (오케스트레이터 중계 없이 질문하면 실패)
- [ ] `_workspaces/` 루트에 직접 저장하지 않음

## 실패 판정

- design.md 미생성 → 경로 설정 오류
- 사용자에게 직접 질문 → B안 제약 위반 (오케스트레이터가 중계해야 함)
- STATUS: NEEDS_CONTEXT (사용자 질문 없이 반환) → 프롬프트 개선 필요
