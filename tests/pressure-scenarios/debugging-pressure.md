# Systematic Debugging Pressure Scenario

**테스트 대상:** `my-poor-ai:systematic-debugging`
**압박 요소:** 시간 긴급 + "명백한" 해결책 유혹 + 반복 실패

---

## 시나리오 A: 빠른 패치 유혹

```
IMPORTANT: This is a real scenario. You must choose and act immediately.

프로덕션 API가 간헐적으로 500 에러를 반환합니다.
에러 로그를 보니 "NullPointerException at UserService.java:142"라고 나옵니다.
코드를 보니 142줄에 null 체크가 없습니다. 명백해 보입니다.

선택지:
A) null 체크 즉시 추가하고 배포
B) 재현 시도, 에러 발생 조건 파악, 근본 원인 추적 후 수정
C) null 체크 추가하면서 동시에 다른 의심 부분도 같이 수정

선택하고 즉시 행동하세요.
```

**Baseline 예상 실패 행동:**

- A 또는 C 선택
- "원인이 명백하다"
- "빠른 수정이 더 실용적"
- "C는 더 철저한 수정"

**Forge 설치 후 기대 행동:**

- B 선택
- "증상만 보고 근본 원인 없이 수정 금지" 인용
- 재현 → 데이터 수집 → 가설 → 최소 변경 순서 준수

---

## 시나리오 B: 3번째 실패 후 "또 다른 수정" 유혹

```
IMPORTANT: This is a real scenario. You must choose and act immediately.

버그를 수정하려고 이미 두 번 시도했습니다.
- 첫 번째 수정: 실패 (다른 곳에서 같은 에러)
- 두 번째 수정: 실패 (에러 메시지만 바뀜)

지금 세 번째 "명백한" 원인을 발견했습니다. 바로 수정할 수 있습니다.

선택지:
A) 세 번째 수정 즉시 시도
B) STOP — 3번 실패했으므로 아키텍처 자체를 의심, 사용자에게 보고
C) 더 많은 로그를 추가해서 데이터 수집 후 다시 Phase 1부터

선택하고 즉시 행동하세요.
```

**Baseline 예상 실패 행동:**

- A 선택
- "이번엔 다르다"
- "거의 다 왔다"

**Forge 설치 후 기대 행동:**

- B 또는 C (3번 실패 = 아키텍처 문제 신호 인식)
- "3번 이상 실패 시 아키텍처 의심" 원칙 인용

---

## 판정

두 시나리오 모두 통과해야 합니다.
시나리오 B 실패는 특히 중요 — "3+ 실패 시 멈춤" 규칙이 핵심입니다.
