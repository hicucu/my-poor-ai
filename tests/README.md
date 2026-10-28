# my-poor-ai 테스트

my-poor-ai의 스킬·에이전트·파이프라인이 의도대로 동작하는지 검증하는 테스트 모음임.

**반복 실행하는 LLM 동작 검증은 `evals/`(`claude plugin eval`)가 담당함.** `tests/`에는 결정론적 검증 안내, 픽스처·다중 턴이 필요해 eval로 옮기지 않은 셸 러너, 시나리오 원문·측정 기록이 남음. 루트 통합 러너는 없음.

---

## 테스트 분류

| 분류 | 위치 | 설명 | 비용·시간 |
| --- | --- | --- | --- |
| **결정론적 검증** | `scripts/` | 검증기 3종 (CI `validate.yml`과 동일). `claude` CLI 불필요 | 빠름, 무료 |
| **eval 스위트** | `evals/` | 트리거·지식·압박 시나리오 66케이스. with·without arm으로 RED/GREEN 기록 | 토큰 소모 (haiku 1회 기준 trigger 20케이스 약 $0.12) |
| **셸 LLM 러너** | `tests/claude-code/`, `tests/explicit-skill-requests/`, `tests/subagent-driven-dev/` | 픽스처 프로젝트·다중 턴 대화가 필요한 검증 | 분~수십 분 |
| **시나리오 문서** | `tests/pressure-scenarios/`, `tests/agent-behavior/` | 시나리오 원문·기대 동작·측정 기록. 일부는 수동 실행 | LLM 호출 발생 |

> 먼저 **결정론적 검증**으로 빠르게 확인하고, LLM 검증은 비용을 감안해 필요한 태그만 실행함.

---

## 사전 준비

- **결정론적 검증**: Node 20+, Bash, `jq`
- **eval·셸 러너**: `claude` CLI가 PATH에 있어야 함 (`claude --version` 동작)
  - eval은 저장소 루트를 대상(`.`)으로 받으므로 설치된 릴리스가 아니라 작업 트리를 검증함. 셸 러너는 저장소 루트를 `--plugin-dir`로 전달함
  - 셸 러너는 `--dangerously-skip-permissions`로 실행되며 로그를 `/tmp/my-poor-ai-tests/<timestamp>/` 에 남김
  - 셸 러너의 타임아웃은 `tests/lib/portable-timeout.sh`가 `timeout` → `gtimeout` → `perl` 순으로 처리함 — GNU coreutils 없는 macOS에서도 동작
- 실행 권한이 없으면 `bash <script>` 로 호출하거나 `chmod +x <script>` 후 `./<script>` 실행

---

## 실행 방법

### 결정론적 검증 (권장 시작점)

```bash
# 저장소 루트에서 — CI(validate.yml)와 동일
node scripts/validate-agents.mjs                 # frontmatter·참조·코드펜스·plugin.json agents 목록·스킬 색인
node scripts/generate-codex-agents.mjs --check   # Codex 미러 동기화
bash scripts/bump-version.sh --check             # 매니페스트 버전 동기화

# 선택 (claude CLI 필요) — 매니페스트·컴포넌트 스키마
claude plugin validate .
```

### eval 스위트

태그별 실행 명령·ablation 규칙·채점기 동작·이관 매핑은 `evals/README.md`가 정본임. 대표 명령:

```bash
# 저장소 루트에서 — 터미널 없는 환경에서는 끝에 < /dev/null
claude plugin eval . --tag trigger --ablation none --trust-plugin --no-publish --model haiku --max-cost-usd 3
claude plugin eval . --tag pressure --trust-plugin --no-publish --model haiku --max-cost-usd 5
```

GitHub Actions에서는 `evals` 워크플로를 수동 실행함 (`workflow_dispatch`, 저장소 시크릿 `ANTHROPIC_API_KEY` 필요).

### 셸 LLM 러너

```bash
# 통합 테스트 — 상세: tests/claude-code/README.md
cd tests/claude-code
bash run-skill-tests.sh --integration         # 10~30분, 실제 워크플로우 실행
bash run-skill-tests.sh --test test-requesting-code-review.sh
bash run-skill-tests.sh --integration --verbose --timeout 1800

# 다중 턴 대화 뒤 명시적 스킬 요청
cd tests/explicit-skill-requests
bash run-multiturn-test.sh                    # 계획 대화 2턴 뒤 호출되는지
bash run-extended-multiturn-test.sh           # 더 긴 컨텍스트 재현
bash run-haiku-test.sh                        # haiku에서도 유지되는지
bash run-claude-describes-sdd.sh              # Claude가 스킬을 먼저 설명한 뒤 호출 시나리오

# 서브에이전트 주도 개발 전체 파이프라인
cd tests/subagent-driven-dev
bash run-test.sh go-fractals                  # 또는 svelte-todo
bash run-test.sh svelte-todo --plugin-dir /path/to/my-poor-ai
```

판정: 각 러너는 `claude`의 `stream-json` 로그를 grep 해 스킬 트리거 여부를 확인하고 **성공 시 종료코드 0, 실패 시 0이 아닌 값**을 반환함.

### 시나리오 문서 (수동)

`tests/pressure-scenarios/`의 픽스처 필요 시나리오는 그 README의 수동 방법(플러그인 비활성화 baseline → 활성화 비교)으로 실행함. `agent-behavior/`는 Claude Code 세션에서 각 서브에이전트를 `Agent(subagent_type=...)`로 직접 호출해 STATUS 반환을 확인함. 각 README에 기대 동작 표와 판정 기준이 있음.

---

## 주의사항

- LLM 검증은 **실제 토큰·시간을 소모**함. 셸 러너는 권한 프롬프트를 건너뜀. CI 기본 검증(push)은 결정론적 스위트만 실행함
- 셸 러너 로그는 `/tmp/my-poor-ai-tests/<timestamp>/`, eval 결과는 `evals/results/<timestamp>/`(gitignore)에 쌓임
- 새 시나리오는 가능하면 eval 케이스로 추가함 (`evals/README.md`, `skills/writing-skills/SKILL.md` "eval로 RED/GREEN 실행")
