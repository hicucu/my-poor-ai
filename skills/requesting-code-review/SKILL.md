---
name: requesting-code-review
description: Use when completing a task or implementing a major feature, or before merging, to verify that the work meets requirements
---

# 코드 리뷰 요청

문제가 연쇄되기 전에 잡아내기 위해 코드 리뷰어 subagent를 파견할 것. 리뷰어는 평가를 위해 정밀하게 작성된 컨텍스트를 받음 — 절대 현재 세션의 히스토리가 아님. 이렇게 하면 리뷰어가 사고 과정이 아닌 작업 결과물에 집중하게 되고, 지속적인 작업을 위한 자신의 컨텍스트도 보존됨.

**핵심 원칙:** 일찍, 자주 리뷰할 것.

## 리뷰 요청 시점

**필수:**

- subagent 주도 개발에서 각 작업 후
- 주요 기능 완료 후
- main에 병합하기 전

**선택적이지만 가치 있음:**

- 막혔을 때 (새로운 시각)
- 리팩터링 전 (기준선 확인)
- 복잡한 버그 수정 후

## 요청 방법

**1. git SHA 가져오기:**

```bash
BASE_SHA=$(git rev-parse HEAD~1)  # 또는 origin/main
HEAD_SHA=$(git rev-parse HEAD)
```

**2. 코드 리뷰어 subagent 파견:**

`general-purpose` 타입의 Task 도구 사용, `code-reviewer.md`의 템플릿 작성

**플레이스홀더:**

- `{DESCRIPTION}` - 빌드한 것의 간략한 요약
- `{PLAN_OR_REQUIREMENTS}` - 무엇을 해야 하는지
- `{BASE_SHA}` - 시작 커밋
- `{HEAD_SHA}` - 종료 커밋

**3. 피드백 반영:**

- Critical 이슈는 즉시 수정
- Important 이슈는 진행 전에 수정
- Minor 이슈는 나중을 위해 기록
- 리뷰어가 틀렸다면 이유와 함께 반박

## 예시

```
[방금 작업 2 완료: 검증 함수 추가]

본인: 진행하기 전에 코드 리뷰를 요청하겠습니다.

BASE_SHA=$(git log --oneline | grep "Task 1" | head -1 | awk '{print $1}')
HEAD_SHA=$(git rev-parse HEAD)

[코드 리뷰어 subagent 파견]
  DESCRIPTION: 4가지 이슈 유형과 함께 verifyIndex() 및 repairIndex() 추가
  PLAN_OR_REQUIREMENTS: _workspaces/{branch-slug}/specs/spec-deployment.md의 작업 2
  BASE_SHA: a7981ec
  HEAD_SHA: 3df7661

[Subagent 반환]:
  강점: 깔끔한 아키텍처, 실제 테스트
  이슈:
    Important: 진행 상황 표시기 누락
    Minor: 보고 간격에 매직 넘버(100) 사용
  평가: 진행 가능

본인: [진행 상황 표시기 수정]
[작업 3으로 계속]
```

## 워크플로우 통합

**Subagent 주도 개발:**

- 각 작업 후 리뷰
- 이슈가 누적되기 전에 잡아내기
- 다음 작업으로 이동하기 전에 수정

**계획 실행:**

- 각 작업 후 또는 자연스러운 체크포인트에서 리뷰
- 피드백 받고, 적용하고, 계속 진행

**임시 개발:**

- 병합 전 리뷰
- 막혔을 때 리뷰

## 위험 신호

**절대 금지:**

- "간단하다"는 이유로 리뷰 건너뛰기
- Critical 이슈 무시
- 수정되지 않은 Important 이슈로 진행
- 유효한 기술적 피드백에 반박

**리뷰어가 틀린 경우:**

- 기술적 근거로 반박할 것
- 작동을 증명하는 코드/테스트를 보여줄 것
- 명확화를 요청할 것

템플릿 위치: requesting-code-review/code-reviewer.md
