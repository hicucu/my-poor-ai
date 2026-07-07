# Pipeline Triggering Tests

my-poor-ai의 3방향 파이프라인 분기가 올바르게 트리거되는지 검증합니다.

## 테스트 방법

```bash
# Claude Code에서 각 프롬프트 파일을 입력으로 실행
claude < prompts/simple-request.txt
claude < prompts/complex-request.txt
claude < prompts/debug-request.txt
```

## 기대 동작

| 프롬프트               | 기대 경로   | 첫 번째 호출 에이전트/skill     |
| ---------------------- | ----------- | ------------------------------- |
| `simple-request.txt`   | 단순 경로   | `my-poor-ai:test-driven-development` |
| `complex-request.txt`  | 복잡 경로   | `project-context`         |
| `debug-request.txt`    | 디버깅 경로 | `my-poor-ai:systematic-debugging`    |
| `boundary-simple.txt`  | 단순 경로   | `my-poor-ai:test-driven-development` |
| `boundary-complex.txt` | 복잡 경로   | `project-context`         |

## 판정

- 올바른 경로 선택 + 첫 번째 에이전트/skill 일치 → PASS
- 경로 오분류 또는 첫 단계 건너뜀 → FAIL → `using-my-poor-ai/SKILL.md` 분류 기준 수정
