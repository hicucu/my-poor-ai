---
name: which-way-should-i-go
description: Use when starting a project with no existing codebase and only a goal, when a goal exists but the tech stack or overall approach is entirely undecided, or when adding a feature whose domain or paradigm sharply diverges from the current codebase (e.g. adding social features to an e-commerce app, adding a game to a social app, adding real-item trading to a game) — before deciding what to brainstorm.
---

# which-way-should-i-go

## 개요

무엇을 브레인스토밍할지조차 정하기 어려울 때, 서로 다른 세대의 검증된 접근법 3가지를 병렬 웹 조사로 비교해 방향을 먼저 정하는 기법. 학습된 지식으로 어림잡는 대신 현재 시점의 실제 사례로 판단함. 산출물은 심층 리서치 리포트가 아니라 방향 결정에 필요한 최소한의 비교표임.

## 사용 시점

- 코드베이스 없이 목표만 있음 (0에서 시작)
- 목표는 있으나 기술스택·전체 접근법이 미정
- 기존 코드베이스와 결이 다른 도메인의 기능 추가 (커머스 앱에 소셜 기능, 소셜 앱에 게임, 게임에 실물 거래 등)

**사용 안 함:** 방향이 이미 정해져 세부 설계만 남은 경우 — 바로 `my-poor-ai:brainstorming` 진입.

일부만 확정된 상태(예: "언어는 Python 확정, 아키텍처는 미정")는 방향 미정에 해당함. 확정 사항은 사용 안 함 사유가 아니라 조사 제약 조건으로 각 에이전트에 전달함.

## 핵심 패턴: 3세대 렌즈

| 렌즈 | 정의 | 경계 |
| --- | --- | --- |
| A. 정통·지속형 | 오래 검증됐고 여전히 현역인 주류 방식 | 사장된 기술은 제외 — "오래됨"이 아니라 "여전히 쓰임"이 기준 |
| B. 모던·엘레강스형 | 현재 가장 세련된 표준 접근법 | 학습 지식만으로 단정 금지 — 웹 검색으로 현재 시점 확인 |
| C. 신흥·핫형 | 최근 1-2년 급부상한 방식 | 학습 컷오프 이후 정보일 가능성이 높음 — 웹 검색 필수 |

렌즈는 시점 상대적이라 같은 결정도 시기마다 답이 달라짐. 아래는 프레임 이해를 위한 스냅숏일 뿐 — 매번 현재 시점과 실제 도메인으로 새로 채움:

| 시점 · 결정 | A. 정통·지속형 | B. 모던·엘레강스형 | C. 신흥·핫형 |
| --- | --- | --- | --- |
| 2010년대 후반, 웹 프론트엔드 | 서버 렌더링 MPA | Angular | React |
| 2020년경, React 상태 관리 | Redux | Redux Toolkit | zustand |
| 한국 온라인 커머스 진출 | 오픈마켓 입점 | D2C 자사몰 | SNS 커머스 |

마지막 행처럼 렌즈는 기술 선택뿐 아니라 비즈니스 방향 결정에도 동일하게 적용됨.

## 프로세스

1. 목표·도메인을 한 줄로 확정. 이미 주어진 확정 사항(언어·플랫폼 등)은 제약 조건으로 함께 기록. 세부 요구사항 질문은 하지 않음 — 그건 brainstorming의 몫
2. `my-poor-ai:dispatching-parallel-agents` 패턴으로 조사 에이전트 3개를 한 메시지에서 동시 디스패치. 각 에이전트 프롬프트:

   ```
   목표: {도메인 한 줄}. 제약: {확정 사항 또는 "없음"}.
   당신의 렌즈: {A|B|C} — {렌즈 정의와 경계}.
   웹 검색으로 현재 시점의 실제 사례를 확인할 것. 학습 지식만으로 답하지 말 것.
   반환: 접근법 이름 / 대표 사례 2-3개 / 강점 / 리스크 / 최소 요구 스택 / 적합한 상황
   ```

3. 3개 결과를 비교표로 종합해 `_workspaces/{branch-slug}/direction-comparison.md`에 저장
4. 비교표와 권장안(이유 포함)을 제시하고 사용자가 방향을 선택 — 혼합 선택도 가능
5. 선택된 방향을 컨텍스트로 `my-poor-ai:brainstorming` 호출 (writing-plans 직행 금지 — brainstorming이 다음 단계)

## 흔한 실수

| 생각 | 현실 |
| --- | --- |
| "예산·카테고리부터 물어봐야 정확한 비교가 된다" | 방향 탐색은 도메인 한 줄이면 충분함. 세부 질문으로 막으면 조사 자체가 시작되지 않음 — 세부 스펙은 brainstorming에서 다룸 |
| "기존 지식으로 개략적 비교가 가능하다" | B·C 렌즈는 정의상 "현재/최근" 기준이라 학습 컷오프 이후 정보를 놓침. 반드시 웹 검색으로 검증함 |
| "사용자가 언급한 선택지가 곧 3렌즈다" | 사용자의 후보는 출발점일 뿐. 실제 도메인에 3렌즈 프레임을 새로 적용하고 웹 조사로 재확인함 |
| "세 방향은 서로 영향을 주므로 병렬화가 부적합하다" | 3렌즈는 서로 참조 없이 독립 조사 가능 — 병렬 디스패치 대상임 |

## 관련 스킬

- **필수 후속:** `my-poor-ai:brainstorming` — 방향 확정 후 세부 설계로 진입
- **디스패치 메커니즘:** `my-poor-ai:dispatching-parallel-agents`
