# my-poor-ai 테스트

my-poor-ai의 스킬·에이전트·파이프라인이 의도대로 동작하는지 검증하는 테스트 모음입니다.

**루트 통합 러너는 없습니다.** 각 스위트는 `tests/<스위트>/` 아래에서 독립적으로 실행합니다.

---

## 테스트 분류

| 분류                | 설명                                                          | 비용·시간             |
| ------------------- | ------------------------------------------------------------- | --------------------- |
| **결정론적 테스트** | Node/셸 단위 테스트. `claude` CLI 불필요, 로컬에서 안전하게 반복 | 빠름, 무료            |
| **LLM 동작 테스트** | `claude` CLI를 실제 호출해 스킬 트리거·동작을 검증              | 토큰·시간 소모 (분~수십 분) |
| **시나리오 문서**   | 러너 없이 프롬프트/시나리오를 Claude에 직접 입력해 수동 판정     | LLM 호출 발생          |

> 먼저 **결정론적 테스트**로 빠르게 확인하고, LLM 동작 테스트는 비용을 감안해 선택적으로 실행하세요.

---

## 사전 준비

- **결정론적 테스트**: Bash
- **LLM 동작 테스트**: `claude` CLI가 PATH에 있어야 함 (`claude --version` 동작)
  - 대부분의 러너가 저장소 루트를 `--plugin-dir`로 자동 전달하므로 별도 플러그인 설치는 불필요
  - 러너는 `--dangerously-skip-permissions`로 실행되며 로그를 `/tmp/my-poor-ai-tests/<timestamp>/` 에 남김
- 실행 권한이 없으면 `bash <script>` 로 호출하거나 `chmod +x <script>` 후 `./<script>` 실행

---

## 스위트 목록

| 디렉토리                   | 분류        | 검증 대상                                      | 진입점                       |
| -------------------------- | ----------- | ---------------------------------------------- | ---------------------------- |
| `opencode/`                | 결정론적    | OpenCode 플러그인 로딩·우선순위·캐싱·도구       | `./run-tests.sh`             |
| `claude-code/`             | LLM 동작    | 스킬 로딩 및 준수 (헤드리스 `claude -p`)        | `./run-skill-tests.sh`       |
| `skill-triggering/`        | LLM 동작    | 자연어 프롬프트로 스킬 자동 트리거             | `./run-all.sh`               |
| `explicit-skill-requests/` | LLM 동작    | 사용자가 스킬명을 직접 지명했을 때 호출 여부    | `./run-all.sh`               |
| `subagent-driven-dev/`     | LLM 동작    | 서브에이전트 주도 개발 전체 파이프라인          | `./run-test.sh <테스트명>`   |
| `pipeline-triggering/`     | 시나리오    | 단순/복잡/디버깅 3방향 파이프라인 분기          | `tests/pipeline-triggering/README.md` |
| `pressure-scenarios/`      | 시나리오    | 압박 상황에서 discipline 스킬 준수             | `tests/pressure-scenarios/README.md`  |
| `agent-behavior/`          | 시나리오    | 각 서브에이전트의 STATUS 반환·상태 추적         | `tests/agent-behavior/README.md`      |

---

## 실행 방법

### 결정론적 테스트 (권장 시작점)

```bash
# OpenCode 플러그인 스위트
cd tests/opencode
bash run-tests.sh                     # 로딩·우선순위·캐싱·도구 전체
```

### LLM 동작 테스트

```bash
# Claude Code 스킬 테스트 — 상세: tests/claude-code/README.md
cd tests/claude-code
bash run-skill-tests.sh                       # 빠른 테스트 전체 (~2분)
bash run-skill-tests.sh --integration         # 통합 테스트 (10~30분, 실제 워크플로우 실행)
bash run-skill-tests.sh --test test-requesting-code-review.sh
bash run-skill-tests.sh --verbose             # 전체 출력
bash run-skill-tests.sh --timeout 1800        # 타임아웃(초) 조정

# 스킬 트리거 (자연어 프롬프트 → 스킬 자동 호출)
cd tests/skill-triggering
bash run-all.sh                               # 전체
bash run-test.sh <skill-name> <prompt-file>   # 개별

# 명시적 스킬 요청 (사용자가 스킬명을 직접 지명)
cd tests/explicit-skill-requests
bash run-all.sh                               # 전체
bash run-test.sh <skill-name> <prompt-file>   # 단일 프롬프트
bash run-multiturn-test.sh                    # 다중 턴 대화 후에도 스킬 호출되는지
bash run-extended-multiturn-test.sh           # 더 긴 컨텍스트 재현
bash run-haiku-test.sh                        # 저비용(haiku) 모델에서도 스킬 강제 유지되는지
bash run-claude-describes-sdd.sh              # Claude가 스킬을 먼저 설명한 뒤 호출 시나리오

# 서브에이전트 주도 개발 전체 파이프라인
cd tests/subagent-driven-dev
bash run-test.sh go-fractals                  # 또는 svelte-todo
bash run-test.sh svelte-todo --plugin-dir /path/to/my-poor-ai
```

판정: 각 러너는 `claude`의 `stream-json` 로그를 grep 해 스킬 트리거 여부를 확인하고 **성공 시 종료코드 0, 실패 시 0이 아닌 값**을 반환합니다.

### 시나리오 문서 (수동/반자동)

러너가 없는 스위트는 프롬프트/시나리오 파일을 Claude에 직접 입력해 기대 동작과 대조합니다.

```bash
# 파이프라인 분기
claude < tests/pipeline-triggering/prompts/simple-request.txt    # → 단순 경로
claude < tests/pipeline-triggering/prompts/complex-request.txt   # → 복잡 경로
claude < tests/pipeline-triggering/prompts/debug-request.txt     # → 디버깅 경로

# 압박 시나리오 (baseline은 위반해야 정상, my-poor-ai 설치 후엔 준수해야 정상)
claude --no-skills < tests/pressure-scenarios/tdd-pressure.md    # baseline
claude          < tests/pressure-scenarios/tdd-pressure.md       # 스킬 적용
```

`agent-behavior/`는 Claude Code 세션에서 각 서브에이전트를 `Agent(subagent_type=...)`로 직접 호출해 STATUS 반환을 확인합니다. 각 스위트의 `README.md`에 기대 동작 표와 판정 기준이 있습니다.

---

## 주의사항

- LLM 동작 테스트는 **실제 토큰·시간을 소모**하고 권한 프롬프트를 건너뜁니다. CI/반복 검증에는 결정론적 스위트를 우선 사용하세요.
- 로그는 `/tmp/my-poor-ai-tests/<timestamp>/` 에 쌓이므로 디스크 정리 시 참고하세요.
- 새 테스트 추가 방법은 `tests/claude-code/README.md` 의 "Adding New Tests" 절을 참고하세요.
