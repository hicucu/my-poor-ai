# HANDOFF Document Test

## 테스트 목적

복잡 경로 구현 중 실행 주체(subagent-driven-development / executing-plans / developer-agent)가
`_workspaces/{branch-slug}/HANDOFF.md`를 **spec/phase 완료 시에만** 정확히 갱신하고,
**task 완료마다 갱신하는 과잉을 피하는지** 확인한다.

## 시나리오 1: spec 완료 시 갱신 (양성)

2개 spec(spec-a, spec-b)을 가진 플랜을 실행한다.

- 기대: spec-a의 모든 task 완료 후 HANDOFF.md가 생성되고 `## 인계 로그`에 `spec-a 완료` 1줄이 추가된다.
- 기대: 본문 `## 현재 진행 중`이 spec-b를 가리키도록 덮어쓰기된다.
- 기대: spec-a 마일스톤 커밋에 `HANDOFF.md`가 포함된다.

## 시나리오 2: task 완료마다 갱신하지 않음 (과잉 방지, 음성)

spec-a가 5개 task를 가진다.

- 기대: task 1~4 완료 시점에는 HANDOFF.md가 갱신되지 않는다 (체크박스만 변경).
- 위반 신호: task마다 인계 로그가 늘어나거나 본문이 매번 덮어쓰기되면 실패.

## 시나리오 3: phase 완료 시 갱신

Phase 그룹(Phase 1, Phase 2)을 가진 큰 spec을 실행한다.

- 기대: Phase 1의 모든 task 완료 후 HANDOFF.md 본문 갱신 + 로그에 `spec-x/Phase 1 완료` 1줄.
- 기대: 작은(Phase 없는) spec에서는 phase 갱신이 발생하지 않는다.

## 시나리오 4: 없으면 생성

HANDOFF.md가 아직 없는 상태에서 첫 spec이 완료된다.

- 기대: 실행 주체가 표준 템플릿으로 HANDOFF.md를 먼저 생성한 뒤 채운다 (생성 누락 없음).

## 시나리오 5: 인계 로그 5개 상한

6개 이상 spec/phase를 완료한다.

- 기대: 인계 로그는 항상 최근 5줄만 유지하고 가장 오래된 항목을 버린다.

## 시나리오 6: 단순·디버깅 경로 제외

단순 경로(파일 1~2개) 작업을 실행한다.

- 기대: HANDOFF.md를 생성하지 않는다 (복잡 경로 전용).

## 적대적 압력 (합리화 차단)

다음 합리화가 나오면 실패로 간주한다:

| 합리화 | 현실 |
| ------ | ---- |
| "task마다 갱신해야 최신이지" | spec/phase 완료 시에만. task는 체크박스가 추적. |
| "이번엔 작은 변경이라 HANDOFF 생략" | 복잡 경로면 항상 유지. |
| "로그를 다 남겨야 이력 보존" | 이력은 git/pipeline-state가 보존. 로그는 5개 상한. |
| "초기 생성 단계가 따로 필요" | 갱신 시 없으면 생성. 별도 단계 없음. |
