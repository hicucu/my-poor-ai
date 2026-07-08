---
name: which-way-should-i-go
description: Use when starting a project with no existing codebase and only a goal, when a goal exists but the tech stack or overall approach is entirely undecided, or when adding a feature whose domain or paradigm sharply diverges from the current codebase (e.g. adding social features to an e-commerce app, adding a game to a social app, adding real-item trading to a game) — before deciding what to brainstorm.
---

# which-way-should-i-go

## 개요

무엇을 만들지조차 정해지지 않은 상태에서, 서로 다른 세대의 검증된 접근법 3가지를 실시간 웹 조사로 비교해 방향을 정하는 기법. brainstorming의 "2-3가지 접근법 제안" 단계를 학습된 지식만으로 어림잡지 않고, 실제 최신 사례로 대체함.

## 사용 시점

- 코드베이스 없이 목표만 있음 (0에서 시작)
- 목표는 있으나 기술스택/전체 접근법이 전혀 안 정해짐
- 기존 코드베이스와 결이 다른 기능 추가 (쇼핑몰+소셜, 소셜+게임, 게임+실물거래 등)

**사용 안 함:** 스택/방향이 이미 정해져 있고 세부 설계만 남은 경우 — 바로 `my-poor-ai:brainstorming` 진입.

## 핵심 패턴: 3개 세대 렌즈

| 렌즈 | 정의 | 경계 |
| --- | --- | --- |
| A. 정통·지속형 | 죽지 않고 여전히 검증된 주류 방식 | 사장된 기술(jQuery류)은 제외 — "오래됨"이 아니라 "여전히 쓰임"이 기준 |
| B. 모던·엘레강스형 | 현재 가장 세련된 표준 접근법 | 학습 지식만으로 판단 금지 — WebSearch로 현재 시점 확인 |
| C. 신흥·핫형 | 최근 1-2년 급부상한 방식 | 학습 지식 컷오프 이후 정보일 가능성 높음 — WebSearch 필수 |

예시(한국 인터넷 쇼핑몰): A=오픈마켓 입점, B=D2C 자사몰, C=SNS 커머스. 이 예시는 하나의 도메인 사례일 뿐 — 매번 실제 목표 도메인에 맞게 3렌즈를 새로 채움.

## 프로세스

1. 목표/도메인을 한 줄로 확정 (예산·상세 요구사항 질문은 하지 않음 — 방향 탐색에는 도메인 한 줄이면 충분, 세부 스펙은 이후 brainstorming의 몫)
2. `my-poor-ai:dispatching-parallel-agents` 패턴으로 general-purpose 에이전트 3개를 한 메시지에서 동시 디스패치. 각 에이전트에게: 목표/도메인, 자신이 맡은 렌즈의 정의와 경계, WebSearch로 실제 사례 확인 지시, 반환 형식(접근법 이름 / 대표 사례 / 강점 / 리스크 / 최소 요구 스택 / 적합 상황)
3. 3개 결과를 비교표로 종합해 `_workspaces/{branch-slug}/direction-comparison.md`에 저장
4. 사용자에게 제시하고 방향(또는 혼합) 선택받기
5. 선택된 방향을 컨텍스트로 `my-poor-ai:brainstorming` 호출 (writing-plans 직접 호출 금지 — brainstorming이 다음 단계)

## 흔한 실수

| 생각 | 현실 |
| --- | --- |
| "예산·카테고리부터 물어봐야 정확한 비교가 된다" | 방향 탐색은 도메인 한 줄로 충분함. 세부 질문으로 막으면 3렌즈 조사 자체가 시작되지 않음 — 세부 스펙은 brainstorming 단계에서 다룸 |
| "기존 지식으로 개략적 비교가 가능하다" | B·C 렌즈는 정의상 "현재/최근" 기준이라 학습 컷오프 이후 정보를 놓침. 반드시 WebSearch로 검증함 |
| "사용자가 준 예시가 이미 정답이다" | 예시는 하나의 사례일 뿐. 실제 목표 도메인에 3렌즈 프레임을 새로 적용함 |
| "이 3개는 서로 영향을 주므로 병렬화가 부적합하다" | 3렌즈는 서로 참조 없이 독립 조사 가능 — 병렬 디스패치 대상임 |

## 관련 스킬

- **필수 후속:** `my-poor-ai:brainstorming` — 방향 확정 후 세부 설계로 진입
- **디스패치 메커니즘:** `my-poor-ai:dispatching-parallel-agents`
- **구분:** `deep-research`는 하나의 주제에 대한 사실검증형 리포트가 목적. 이 스킬은 "방향 결정을 위한 3-way 비교"가 목적이며 각 렌즈의 심층 사실검증까지는 하지 않음
