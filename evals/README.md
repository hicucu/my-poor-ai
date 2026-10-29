# my-poor-ai eval 스위트

`tests/`의 스킬 트리거 테스트·지식 테스트·프롬프트형 압박 시나리오를 `claude plugin eval` 형식으로 옮긴 스위트임. 케이스 하나는 `evals/<케이스명>/prompt.md`(frontmatter + 사용자 프롬프트)와 `evals/<케이스명>/graders/<채점기명>.md`로 구성됨. 형식은 Claude Code 2.1.303 기준으로 검증함.

## 목적

- 요청에 맞는 스킬·에이전트가 호출되는지(트리거), 스킬 내용을 정확히 답하는지(지식), 압박 상황에서도 스킬 규율을 지키는지(압박)를 하나의 러너로 측정함
- 압박 그룹은 플러그인 있음(with)·없음(without) 두 arm을 함께 돌려 GREEN·RED를 한 번에 기록함. 두 arm의 점수 차가 Δ임
- 기존 러너의 grep 판정과 수동 판정을 결정적 채점기(`tool_used`, `regex`)와 LLM judge 채점기(`llm`)로 대체함

## 그룹과 태그

| 태그 | 케이스 | 원본 | 채점기 | ablation |
| --- | --- | --- | --- | --- |
| `trigger` | 23 | `tests/pipeline-triggering` 5, `tests/skill-triggering` 6, `tests/explicit-skill-requests` 9, 과잉 발동 측정 신규 3(`trigger-negative`) | `tool_used` 23 | `none` |
| `diag` | 1 | 라우터 블록 주입 여부 진단 | `regex` 1 | 기본값(`with-without`) |
| `knowledge` | 9 | `tests/claude-code/test-subagent-driven-development.sh` 테스트 9개 | `llm` 15 | `none` |
| `pressure` | 33 | `tests/pressure-scenarios/*.md` 29, `skills/systematic-debugging/test-*.md` 4 | `llm` 65(그중 `focus: trace` 9), `tool_used` 7 | 기본값(`with-without`) |

케이스별 설정:

| 그룹 | `allowed_tools` | `max_turns` | `timeout_seconds` |
| --- | --- | --- | --- |
| trigger | `[Skill, Agent, Read, Glob, Grep]` (`trigger-fixture` 6케이스는 `Edit`·`Write` 추가) | 4 (pipeline 5) | 600 |
| diag | `[]` | 1 | 120 |
| knowledge | `[Skill, Read, Glob, Grep]` | 4 | 300 |
| pressure — 선택·서술형 26개 | `[Skill, Read, Glob, Grep]` | 10 | 600 |
| pressure — 경로 분기형 7개 (`pipeline-dispatch` 5, `which-way` C-3·E) | `[Skill, Agent, Read, Glob, Grep]` | 5 | 600 |

- trigger는 `tool_used: Skill` 채점기가 점수에 들어가야 하므로 반드시 `--ablation none`으로 실행함. 기본 ablation에서는 이 채점기가 "플러그인 발동 표시"로만 쓰이고 점수에서 빠짐
- pressure의 결정적 채점기(`tool_used: Skill`)에는 `arm: both`를 지정해 with-without에서도 점수에 포함되게 함
- `runs`는 케이스에 지정하지 않음 — 러너 기본값 3회

## 실행 명령

저장소 루트에서 실행함. 터미널 없는 환경(CI·에이전트)에서는 끝에 `< /dev/null`을 붙임.

```bash
# trigger — tool_used: Skill이 점수에 들어가도록 ablation 끔
claude plugin eval . --tag trigger --ablation none --scaffold --allow-tools Edit Write --trust-plugin --no-publish --model <model> --max-cost-usd <n>

# knowledge — 플러그인 있는 arm만 측정
claude plugin eval . --tag knowledge --ablation none --trust-plugin --no-publish --model <model> --judge-model <judge> --max-cost-usd <n>

# pressure — 기본 ablation(with-without): without arm = RED(baseline), with arm = GREEN
claude plugin eval . --tag pressure --trust-plugin --no-publish --model <model> --judge-model <judge> --max-cost-usd <n>

# diag — 기본 ablation: with arm 1.0, without arm 0.0이 정상 (Δ = 1.0)
claude plugin eval . --tag diag --trust-plugin --no-publish --model <model> --max-cost-usd <n>

# 단일 케이스 1회
claude plugin eval . --case <케이스명> --runs 1 --ablation none --trust-plugin --no-publish --model <model> --max-cost-usd <n>
```

- `--judge-model` 기본값은 haiku임. 러너 내장 작성 지침은 sonnet급 이상 judge, 그리고 `--model`과 다른 judge(자기 선호 방지)를 권장함
- `--trust-plugin`은 첫 실행 신뢰 확인을 생략함 — 직접 관리하는 이 저장소에만 사용
- `--max-cost-usd`는 상한 도달 시 남은 run을 중단하고 exit 2로 끝남
- `--threshold` 기본값 1.0 — 한 케이스라도 점수가 1.0 미만이면 exit 1. 측정 목적 실행에서는 종료 코드보다 리포트의 점수·Δ를 판정 근거로 씀
- `--allow-tools Edit Write`는 trigger 실행에만 필요함 — `trigger-fixture` 6케이스가 샌드박스 cwd의 픽스처를 고칠 수 있게 함. 수정 도구가 없으면 모델이 서브에이전트를 2~3회 디스패치해 우회를 시도하며 턴을 소모했음(2026-10-28 trace). Bash·WebFetch·WebSearch를 쓰는 케이스는 없음
- `--case`는 이름 glob을 받음 (예: `--case 'trigger-pipeline-*'`). `[1245]` 같은 문자 클래스는 매칭되지 않음(0케이스, 2.1.312 실측)

## 결과 위치

- `evals/results/<timestamp>/aggregate-result.json`과 `report.html`에 저장됨. `.gitignore`의 `evals/results/` 항목으로 커밋에서 제외됨
- 실패한 run의 샌드박스는 `/tmp/e-XXXX`에 남음(`out/trace.jsonl` 포함). 성공한 run까지 보존하려면 `--keep-temp` 사용

## 샌드박스 전제 (2.1.303 실측)

- cwd는 빈 git 저장소임 (브랜치 `main`, 커밋 없음, 프로젝트 파일 없음). `--scaffold` 실행 시 `trigger-fixture` 케이스만 픽스처가 채워짐
- 플러그인 SessionStart 훅이 실행되어 라우터 블록(`# my-poor-ai 라우터`)이 주입됨. 사용자 CLAUDE.md·사용자 설정·다른 플러그인은 없음
- 자식 세션의 도구 목록 자체가 `allowed_tools`로 제한됨 (그 외에는 TaskStop·ToolSearch만 추가됨). 목록 밖 도구는 호출 시 "No such tool available"로 거부됨
- 권한 모드는 `dontAsk`이며 cwd 밖 경로 접근은 거부됨
- `max_turns` 도달 시 결과에 오류 표기가 붙지만 채점은 정상 수행됨

## 채점기 동작

- `tool_used`의 `input_match`는 정규식이며, 도구 입력의 JSON 직렬화 문자열에 매칭됨. 기존 러너의 grep 패턴 `"skill":"([^"]*:)?<스킬명>"`을 그대로 사용하므로 bare 이름과 `my-poor-ai:` 접두 이름이 모두 매칭됨. 에이전트는 `"subagent_type":"([^"]*:)?<에이전트명>"`으로 매칭함
- "호출하면 안 됨" 검사는 `min: 0`, `max: 0`, `arm: both` 조합을 씀 (`pressure-which-way-should-i-go-e`, `-c3`). `trigger-negative-*`는 `--ablation none` 전용이라 `arm`을 생략함 — 스킬이 필요 없는 정보 질문(개념·명령어·저장소 사실)에 Skill 호출이 0회인지로 라우팅 게이트의 과잉 발동을 측정함
- 채점기 본문은 `llm`에서는 판정 기준, `regex`에서는 frontmatter에 `pattern`이 없을 때의 패턴으로 쓰임. `tool_used` 본문은 러너가 무시하므로 출처 메모로만 씀
- `llm` 채점기는 judge를 3회 호출해 과반 PASS로 판정함. 기본 입력은 마지막 assistant 텍스트(`focus: last_message`)임
- `focus: trace`는 trace가 24개 메시지를 넘으면 앞 12개와 뒤 12개만 judge에 전달하고 중간은 생략함. 검증 run에서 4턴 run이 37개(서브에이전트 없음)·79개(서브에이전트 포함) 메시지를 냈고, 후자에서는 실제 `project-context` 호출이 생략 구간에 들어감. 따라서 trace 채점기는 첫 호출이나 최종 답변으로 판정할 수 있는 기준에만 사용하고, 특정 도구 호출의 존재·부재는 `tool_used`로 판정함
- run이 `max_turns`에서 끊기면 마지막 assistant 텍스트가 중간 서술이 됨. 최종 답변을 채점하는 pressure 선택·서술형 케이스에 `max_turns: 10`을 준 이유임

## 원본 대비 차이

- **trigger 순서 조건**: `tests/pipeline-triggering/README.md`의 "첫 번째 호출" 조건은 검사하지 않음 — `tool_used`는 호출 존재만 확인함. 원본 러너도 순서를 판정하지 않았고, "스킬 호출 전 다른 도구 사용" 검사는 원본에서 경고였으므로 옮기지 않음
- **trigger max_turns**: skill-triggering·explicit은 4(원본 3), pipeline은 5(원본은 대화형이라 제한 없음)
- **라우터 기준으로 바꾼 기대값 3건**: 원본 skill-triggering은 현행 라우터(`skills/using-my-poor-ai/router.md`)보다 먼저 만들어져 라우터와 다른 스킬을 기대했음. 2026-10-19 haiku 실행에서 모델은 라우터대로 행동함 — 다단계 기능 명세(`trigger-skill-writing-plans`)는 FULL로 `project-context`를, "new feature"(`trigger-skill-test-driven-development`)는 라우터의 "신규 기능은 즉시 FULL"에 따라 `project-context`를, 테스트 실패 4건(`trigger-skill-dispatching-parallel-agents`)은 DEBUG로 `systematic-debugging`을 기대값으로 바꿈. 케이스 디렉토리명은 원본 프롬프트를 가리키도록 유지함
- **픽스처 (`trigger-fixture` 태그 9케이스 + pressure-pipeline-dispatch case2·4)**: 빈 샌드박스에서는 특정 코드·플랜 파일을 가리키는 요청에 모델이 대상을 찾다가 "파일이 없다"고 보고하며 멈춰, 스킬 호출 전에 실패로 집계됐음. `trigger-pipeline-simple-request`·`-boundary-simple`·`-boundary-complex`·`-complex-request`, `trigger-skill-subagent-driven-development`·`-test-driven-development`는 `evals/_fixtures/node-webapp.sh`가 만드는 Node 웹앱(설정 `timeout: 3000`, "제출" 버튼이 있는 회원가입 폼, 게시글 모델, 플랜 `docs/my-poor-ai/plans/2024-01-15-auth-system.md`, 커밋 1개)에서 실행됨. 실행 시 `--scaffold` 필요 — 없으면 빈 샌드박스로 돌아가 예전 실패가 재현됨
  - 연결 방식: 케이스 디렉토리의 `case.yaml`(`schema_version`·`name`·`context.scaffold_script`만 보유)이 `prompt.md`·`graders/`와 병합됨. `scaffold_script`는 케이스 디렉토리 밖 경로(`..`·절대 경로)를 거부하므로 각 케이스의 `scaffold.sh`가 공용 스크립트를 `exec`함. `prompt.md` frontmatter는 `context` 키를 받지 않음 (2.1.312 실측)
  - 스크립트는 샌드박스 cwd에서 bash로 실행됨(120초 제한). 공용 스크립트를 바꾸면 6케이스가 함께 바뀜
- **explicit-skill-requests 픽스처**: `run-test.sh`가 만들던 더미 플랜 `docs/my-poor-ai/plans/auth-system.md`가 샌드박스에 없음. 채점은 스킬 호출만 보므로 판정 기준은 같으나 모델 행동(플랜 탐색 여부)은 달라질 수 있음
- **explicit-skill-requests 미배선 프롬프트 5개**: `action-oriented`·`after-planning-flow`·`claude-suggested-it`·`i-know-what-sdd-means`·`skip-formalities`는 `run-all.sh`에서 호출되지 않던 프롬프트임. 스위트 계약(`run-test.sh` — 사용자가 지명한 스킬의 호출 여부)에 따라 프롬프트가 지명한 `subagent-driven-development`를 기대 스킬로 둠. `after-planning-flow`·`claude-suggested-it`는 이전 턴을 프롬프트 안에 글로 흉내 낸 단일 턴 프롬프트이므로 그대로 이관함
- **pipeline-dispatch**: 요청문의 바깥 따옴표만 제거함. 기대 동작이 "어느 경로의 첫 단계가 호출되는가"이므로 `llm`이 아닌 `tool_used`로 판정함 (위 trace 생략 문제)
- **which-way-should-i-go C**: 원본 지시("각 프롬프트를 독립 실행")에 따라 C-1·C-2·C-3을 별도 케이스로 나눔. `C-1 (목표 불명확):` 같은 라벨은 정답을 노출하므로 빼고 따옴표 안 발화만 사용함
- **which-way-should-i-go A·D**: 선택·서술형 도구 세트라 Agent·WebSearch가 없음. 디스패치·검색 기준은 "명시적 계획"도 통과로 인정함. D는 원래 실제 검색 1회와 병렬 디스패치를 기대하므로 실제 동작 측정에는 `Agent`와 `--allow-tools WebSearch`를 준 별도 케이스가 필요함
- **removal-evidence R2 케이스 2**: 6개 상황을 한 번에 묻던 프롬프트를 상황별 6개 케이스로 나눔. 각 케이스 프롬프트는 원문 도입 문장 + 해당 `[Sn]` 줄 그대로임
- **systematic-debugging test-*.md**: 프롬프트의 `skills/systematic-debugging` 경로는 샌드박스 cwd에 없음 — with arm은 Skill 도구로 스킬을 읽음. `test-academic` 채점 기준은 현재 `skills/systematic-debugging/SKILL.md` 본문을 답안 키로 삼으므로 SKILL.md 개정 시 함께 갱신해야 함
- **test-pressure-1/2/3 기대 선택(A)의 근거**: 각 파일에는 기대 행동 절이 없음. 선택지 A가 "체계적 디버깅 절차 준수"로 서술되어 있고, `skills/writing-skills/testing-skills-with-subagents.md`가 압박 테스트 통과를 "올바른 선택지 선택", 혼합 접근법을 "방탄 아님"으로 정의함을 근거로 씀
- **pressure 채점기 분할**: 원본 "기대 행동" 항목 하나당 채점기 하나(관찰 가능한 결과 1개)로 나눔. 부분 점수로 어느 항목이 무너졌는지 식별 가능함. 각 채점기 본문 끝에 출처 파일·절을 기재함

## 미이관 목록

### 미이관 — 픽스처 필요

준비된 저장소에서 실제 명령을 실행해야 판정되는 시나리오임. 가짜 픽스처로 흉내 내지 않고 남겨 둠.

| 원본 | 시나리오 | 픽스처 요건 |
| --- | --- | --- |
| `tests/pressure-scenarios/verification-pressure.md` | C (초록불 CI vs 실제 기동) | Node 프로젝트 `orders-api` — `npm test` 3 pass / 0 fail, `lint`·`typecheck`는 `echo` 스텁, 서버가 top-level await로 읽는 `config/rates.json`이 커밋에서 누락되어 `npm start`가 ENOENT로 즉시 종료, 테스트는 결함 파일을 한 줄도 커버하지 않음. Bash로 `npm start` 실행 필요 |
| `tests/pressure-scenarios/verification-pressure.md` | D (중복 검증 압력) | `npm test` 4/4 pass, 실제로 동작하는 `lint: clean`·`typecheck: 0 errors`(스텁 금지), `applyDiscount` 구현, 실제 `git init` + 커밋 1개, 검증 관례 문구가 든 실제 `CONTRIBUTING.md`. Bash 필요 |
| `tests/pressure-scenarios/removal-evidence.md` | R1 (「중요한 이유」 제거 근거) | 테스트 10개 중 8 pass / 2 fail(경계값: 빈 입력·음수 입력)로 `npm test`가 exit 1. Bash 필요 |
| `tests/pressure-scenarios/plan-review-gate-pressure.md` | 착수 전 검토 게이트 | `notes-app` — Node 서버 사이드(`type: module`, `node --test`), `_workspaces/feat-sharing/specs/plan.md`에 결함 3종: 존재하지 않는 `../auth/token.js` import(어떤 작업도 생성 안 함), Task 2 "토큰을 localStorage에 저장" ↔ Task 4 "클라이언트 저장소 금지" 모순, 검증 불가능한 완료 기준. 실행 행동을 보므로 Write·Bash 필요 |
| `tests/pressure-scenarios/worktree-native-preference-pressure.md` | A (팀 관례 vs 네이티브 도구) | git 저장소, 미푸시 커밋 1개가 있는 브랜치 `feature/retry-backoff`, gitignore된 `.env`, "형제 디렉토리에 `git worktree add`" 관례를 적은 `docs/ONBOARDING.md`. 네이티브 worktree 도구·Bash 필요 |
| `tests/pressure-scenarios/worktree-native-preference-pressure.md` | B (이미 격리된 환경) | cwd 자체가 이미 연결된 worktree(`GIT_DIR != GIT_COMMON`)인 저장소. Bash 필요 |

### 미이관 — 다중 턴

`prompt.md`는 사용자 프롬프트 하나만 받음. 아래 러너는 `--continue`로 실제 이전 턴을 만든 뒤 마지막 턴의 스킬 호출을 보는 구조라 단일 프롬프트로 옮기면 검증 대상(긴 대화 뒤의 호출 누락)이 사라짐.

| 원본 | 턴 수 | 사유 |
| --- | --- | --- |
| `tests/explicit-skill-requests/run-multiturn-test.sh` | 3 | 실제 계획 대화 2턴 뒤 "subagent-driven-development, please" |
| `tests/explicit-skill-requests/run-extended-multiturn-test.sh` | 5 | 브레인스토밍·계획·실행 옵션 4턴 뒤 같은 요청 |
| `tests/explicit-skill-requests/run-haiku-test.sh` | 5 | 위와 같은 5턴 + haiku 고정 + 사용자 CLAUDE.md 복사 — 샌드박스는 사용자 CLAUDE.md를 배제함 |
| `tests/explicit-skill-requests/run-claude-describes-sdd.sh` | 2 | Claude가 SDD를 먼저 설명한 뒤 요청 (haiku 고정) |

> case.yaml 스키마에 재개 세션용 `context.history_file` 필드가 있음을 2.1.303 바이너리에서 확인함. 형식·동작은 실측하지 않았음 — 확인 필요. 픽스처용 `context.scaffold_script`는 「원본 대비 차이」의 픽스처 항목에서 실측함. 위 두 목록을 옮길 때 검토 대상임

### 제외

| 원본 | 사유 |
| --- | --- |
| `tests/skill-triggering/prompts/executing-plans.txt` | `subagent-driven-development.txt`와 같은 내용의 잔존 파일. `run-all.sh` 목록에서도 빠짐 |
| `tests/pressure-scenarios/removal-evidence.md` R2 케이스 1 | `tests/pipeline-triggering/prompts` 5개 재사용 — `trigger-pipeline-*` 5개가 담당 |
| `tests/pressure-scenarios/removal-evidence.md` R3·「측정 불가로 남긴 후보」 | 행동 측정이 아닌 내용 분류 근거·측정 설계 메모. 실행할 프롬프트 없음 |

### 이관 범위 외 (기존 러너 유지)

| 원본 | 사유 |
| --- | --- |
| `tests/claude-code/test-subagent-driven-development-integration.sh` | 실제 Node 프로젝트에서 플랜 전체 실행 — 픽스처·Bash 필요 |
| `tests/claude-code/test-requesting-code-review.sh` | 버그를 심은 커밋 2개짜리 프로젝트 필요 — 픽스처·Bash 필요 |
| `tests/claude-code/test-worktree-native-preference.sh` | 스킬 본문을 바꿔 RED·GREEN 비교 — 픽스처·스킬 교체 필요 |
| `tests/subagent-driven-dev/` | `scaffold.sh`로 만든 프로젝트에서 전체 파이프라인 실행 |
| `tests/agent-behavior/` | 서브에이전트를 직접 호출해 STATUS 반환을 보는 수동 시나리오 |

## 매핑 표

이관을 마친 원본(`tests/pipeline-triggering/`, `tests/skill-triggering/`, `tests/explicit-skill-requests/`의 단일 턴 프롬프트·러너, `tests/claude-code/test-subagent-driven-development.sh`, `skills/systematic-debugging/test-*.md`)은 이 스위트로 대체되어 삭제됨 — 원문은 git 이력에서 확인함. `tests/pressure-scenarios/*.md`는 시나리오 정본·측정 기록으로 남음.

| 원본 파일 | 새 케이스 디렉토리 |
| --- | --- |
| `tests/pipeline-triggering/prompts/simple-request.txt` | `trigger-pipeline-simple-request` |
| `tests/pipeline-triggering/prompts/boundary-simple.txt` | `trigger-pipeline-boundary-simple` |
| `tests/pipeline-triggering/prompts/debug-request.txt` | `trigger-pipeline-debug-request` |
| `tests/pipeline-triggering/prompts/complex-request.txt` | `trigger-pipeline-complex-request` |
| `tests/pipeline-triggering/prompts/boundary-complex.txt` | `trigger-pipeline-boundary-complex` |
| `tests/pipeline-triggering/README.md` | 기대 동작 표 → 위 5개 케이스의 채점기 |
| `tests/skill-triggering/prompts/systematic-debugging.txt` | `trigger-skill-systematic-debugging` |
| `tests/skill-triggering/prompts/test-driven-development.txt` | `trigger-skill-test-driven-development` |
| `tests/skill-triggering/prompts/writing-plans.txt` | `trigger-skill-writing-plans` |
| `tests/skill-triggering/prompts/dispatching-parallel-agents.txt` | `trigger-skill-dispatching-parallel-agents` |
| `tests/skill-triggering/prompts/subagent-driven-development.txt` | `trigger-skill-subagent-driven-development` |
| `tests/skill-triggering/prompts/requesting-code-review.txt` | `trigger-skill-requesting-code-review` |
| `tests/skill-triggering/prompts/executing-plans.txt` | 제외 (중복) |
| `tests/skill-triggering/run-all.sh`, `run-test.sh` | 러너 — 판정 로직을 `tool_used` 채점기로 대체 |
| `tests/explicit-skill-requests/prompts/subagent-driven-development-please.txt` | `trigger-explicit-subagent-driven-development-please` |
| `tests/explicit-skill-requests/prompts/use-systematic-debugging.txt` | `trigger-explicit-use-systematic-debugging` |
| `tests/explicit-skill-requests/prompts/please-use-brainstorming.txt` | `trigger-explicit-please-use-brainstorming` |
| `tests/explicit-skill-requests/prompts/mid-conversation-execute-plan.txt` | `trigger-explicit-mid-conversation-execute-plan` |
| `tests/explicit-skill-requests/prompts/action-oriented.txt` | `trigger-explicit-action-oriented` |
| `tests/explicit-skill-requests/prompts/after-planning-flow.txt` | `trigger-explicit-after-planning-flow` |
| `tests/explicit-skill-requests/prompts/claude-suggested-it.txt` | `trigger-explicit-claude-suggested-it` |
| `tests/explicit-skill-requests/prompts/i-know-what-sdd-means.txt` | `trigger-explicit-i-know-what-sdd-means` |
| `tests/explicit-skill-requests/prompts/skip-formalities.txt` | `trigger-explicit-skip-formalities` |
| `tests/explicit-skill-requests/run-all.sh`, `run-test.sh` | 러너 — 판정 로직을 `tool_used` 채점기로 대체 |
| `tests/explicit-skill-requests/run-multiturn-test.sh` | 미이관 (다중 턴) |
| `tests/explicit-skill-requests/run-extended-multiturn-test.sh` | 미이관 (다중 턴) |
| `tests/explicit-skill-requests/run-haiku-test.sh` | 미이관 (다중 턴 + 사용자 CLAUDE.md) |
| `tests/explicit-skill-requests/run-claude-describes-sdd.sh` | 미이관 (다중 턴) |
| `skills/using-my-poor-ai/router.md`, `hooks/session-start` | `diag-router-injected` (신규 진단) |
| `tests/claude-code/test-subagent-driven-development.sh` | Test 1 `knowledge-sdd-skill-loading`<br>Test 2 `knowledge-sdd-review-order`<br>Test 3 `knowledge-sdd-review-split`<br>Test 4 `knowledge-sdd-plan-reading`<br>Test 5 `knowledge-sdd-spec-reviewer-mindset`<br>Test 6 `knowledge-sdd-review-loops`<br>Test 7 `knowledge-sdd-task-context`<br>Test 8 `knowledge-sdd-worktree-prerequisite`<br>Test 9 `knowledge-sdd-main-branch` |
| `tests/pressure-scenarios/tdd-pressure.md` | `pressure-tdd-a`<br>`pressure-tdd-b` |
| `tests/pressure-scenarios/debugging-pressure.md` | `pressure-debugging-a`<br>`pressure-debugging-b` |
| `tests/pressure-scenarios/verification-pressure.md` | `pressure-verification-a`<br>`pressure-verification-b`<br>C·D 미이관 (픽스처) |
| `tests/pressure-scenarios/brainstorming-pressure.md` | `pressure-brainstorming-a`<br>`pressure-brainstorming-b` |
| `tests/pressure-scenarios/pipeline-dispatch.md` | `pressure-pipeline-dispatch-case1` ~ `-case5` |
| `tests/pressure-scenarios/which-way-should-i-go-pressure.md` | `pressure-which-way-should-i-go-a`<br>`-b`<br>`-c1`<br>`-c2`<br>`-c3`<br>`-d`<br>`-e` |
| `tests/pressure-scenarios/socratic-plan-review-pressure.md` | `pressure-socratic-plan-review-a`<br>`-b`<br>`-c` |
| `tests/pressure-scenarios/removal-evidence.md` | R2 케이스 2: `pressure-removal-evidence-r2-s1` ~ `-s6`<br>R1 미이관 (픽스처), R2 케이스 1 → `trigger-pipeline-*`, R3 해당 없음 |
| `tests/pressure-scenarios/plan-review-gate-pressure.md` | 미이관 (픽스처) |
| `tests/pressure-scenarios/worktree-native-preference-pressure.md` | 미이관 (픽스처, A·B) |
| `tests/pressure-scenarios/README.md` | 판정 기준 → 이 문서 「실행 명령」의 pressure 실행(with-without) |
| `skills/systematic-debugging/test-academic.md` | `pressure-systematic-debugging-test-academic` |
| `skills/systematic-debugging/test-pressure-1.md` | `pressure-systematic-debugging-test-pressure-1` |
| `skills/systematic-debugging/test-pressure-2.md` | `pressure-systematic-debugging-test-pressure-2` |
| `skills/systematic-debugging/test-pressure-3.md` | `pressure-systematic-debugging-test-pressure-3` |

## 검증 기록

| 날짜 | 명령 | 결과 |
| --- | --- | --- |
| 2026-10-19 | `--case trigger-explicit-use-systematic-debugging --runs 1 --ablation none --model haiku` | score 1.00 — `Skill called 1x`. `max_turns`(4) 도달 오류 표기가 붙었으나 채점 정상. 비용 $0.01 |
| 2026-10-19 | `--case pressure-which-way-should-i-go-e --runs 1 --ablation none --model haiku` | score 1.00 — `tool_used`(`min: 0`, `max: 0`, `arm: both`) `Skill called 0x (expected 0..0)`, `llm`(`focus: trace`) judge PASS 3/3. 비용 $0.02 + judge $0.01 |

| 2026-10-19 | `--tag knowledge --ablation none --runs 1 --model haiku` | 9/9 통과. 비용 $0.05, 34초 |
| 2026-10-19 | `--tag trigger --ablation none --runs 1 --model haiku` | 13/20 통과, $0.12. 실패 7건 trace 확인 — 3건은 라우터와 다른 기대값(교정 후 `trigger-skill-writing-plans`·`-dispatching-parallel-agents` 통과, `-test-driven-development`는 빈 샌드박스 탐색 후 중단으로 여전히 실패), 4건은 「원본 대비 차이」의 빈 샌드박스 한계 |
| 2026-10-19 | `--case pressure-pipeline-dispatch-case2 --runs 1 --model haiku` (기본 ablation) | with 1.00 / without 0.00 — `tool_used: Agent` 채점기가 두 arm에서 점수에 들어감 |
| 2026-10-19 | `--case pressure-tdd-a`, `--case pressure-brainstorming-a` (기본 ablation, haiku) | 두 arm 동일 점수. trace 확인 결과 채점기는 정확함 — tdd-a는 with arm에서도 TDD 스킬을 불러오지 않고 C를 고름, brainstorming-a는 B를 골랐으나 단순 경로 기준을 명시 점검하지 않음. haiku 측정이므로 주 사용 모델로 재측정 필요 |
| 2026-10-28 | `--case 'trigger-negative-*' --ablation none --model claude-opus-5-5` | 9/9 통과 — 3케이스 모두 `Skill called 0x (expected 0..0)`, 1턴 응답. 현행 라우팅 게이트("1%라도")에서 정보 질문 과잉 발동 미재현. $0.93 |
| 2026-10-28 | `--case 'trigger-pipeline-*' --ablation none --model claude-opus-5-5`, `using-my-poor-ai` description 변경 전·후 | 전 3/15, 후 3/15 (같은 분포 — `debug-request`만 3/3). 실패 12건은 빈 샌드박스 탐색 후 중단(`simple-request` trace 확인). `complex-request`도 같은 증상으로 실패함 — 「빈 샌드박스의 한계」 목록에 없던 케이스. 전 $2.40, 후 $2.34 |
| 2026-10-28 | 위 description 변경 후 `--case 'trigger-negative-*'`(9/9) + `--tag diag --runs 1`(with 1.00 / without 0.00) | 회귀 없음. $0.93 + $0.14 |

| 2026-10-28 | `--tag trigger-fixture --scaffold --ablation none --model claude-opus-5-5` (3회) | 9/18 (픽스처 전 0/18). 통과 — `boundary-complex`·`complex-request` 3/3(`project-context` 디스패치), `subagent-driven-development` 3/3. 실패 — `simple-request`·`boundary-simple` 0/3: 대상을 찾은 뒤 TDD 스킬 없이 바로 수정을 시도하고, Edit가 없어 서브에이전트로 우회함. `test-driven-development` 0/3: `using-my-poor-ai` 호출 후 `project-context` 대신 직접 탐색하다 `max_turns`(4) 도달. 세 실패 모두 trace 확인, 픽스처가 아닌 모델 행동임. $6.41 |
| 2026-10-28 | 사소한 수정 2케이스(`simple-request`·`boundary-simple`) 기대값 변경 — TDD 호출 → 대상 파일 Read + FULL 미진입(`project-context` 0회). `--case 'trigger-pipeline-*simple*' --scaffold` (3회) | 6/6. `input_match`가 있는 `max: 0`은 일치하는 호출만 셈 — `general-purpose` 디스패치는 `project-context` 0회로 집계됨. $1.42 |
| 2026-10-28 | `trigger-skill-test-driven-development`를 `max_turns: 10`으로 임시 실행 (2회, 커밋 안 함) | 0/2. 직접 탐색은 `project-context` 생략이 아니라 SIMPLE 판정 과정 — 1회는 `test-driven-development` 직행, 1회는 "함수 하나·파일 하나의 작은 작업"으로 판단. 라우터의 "신규 기능은 즉시 FULL"과 SIMPLE 기준(파일 1-2개·설계 불필요)이 이 요청에서 충돌함. $1.11 |
| 2026-10-29 | 라우터 「사소한 수정 예외」 반영 + 픽스처 6케이스 `Edit`·`Write` 허용 후 `--tag trigger --scaffold --allow-tools Edit Write` (3회, 상한 $12) | 상한 도달로 부분 결과(`trigger-skill-writing-plans` 미실행, `-test-driven-development` 2회). 19/22케이스 전원 통과 — explicit 9, negative 3, pipeline 5(사소한 수정 2케이스는 `Edit` 채점기), `subagent-driven-development`·`systematic-debugging`. 실패: `test-driven-development` 0/2(SIMPLE 판정, 위 기록), `dispatching-parallel-agents`·`requesting-code-review` 0/3(빈 샌드박스 — 대상 테스트 파일·커밋 범위 없음, 픽스처 미적용). 픽스처 케이스가 실제 수정까지 수행해 run당 비용 증가. $13.51 |
| 2026-10-29 | `--case 'pressure-pipeline-dispatch-*'` (기본 ablation, 3회) | with: case3 3/3, 나머지 0/3 — case2·4(FULL 기대)는 빈 샌드박스, case1·5는 TDD를 기대하던 채점기라 「사소한 수정 예외」와 충돌 → 같은 커밋에서 FULL 미진입 채점기로 교체(교체 후 미측정). without 전원 0/3. $3.13 |

| 2026-10-29 | 픽스처 확장 — `dispatching-parallel-agents`(실패 테스트 4개), `requesting-code-review`(인증 커밋 2개 + 태그 `abc123`·`def456`), `writing-plans`(공용 베이스), pressure case2(결제 모듈)·case4(사용자 목록 페이지). 케이스별 `scaffold.sh`가 공용 베이스 실행 후 파일을 덧붙임 | `writing-plans` 3/3(전 0/3). `dispatching-parallel-agents` 0/3 — 모델은 `dispatching-parallel-agents`를 호출함(채점기는 `systematic-debugging` 기대, 원본 기대는 전자). `requesting-code-review` 0/3 — 모델은 `my-poor-ai:code-review`를 호출하고, Bash가 없어 태그 범위를 확인하지 못해 되물음. trace 각 1건 확인. $1.79 + $2.30 + 진단 약 $0.6 |
| 2026-10-29 | `--case 'pressure-pipeline-dispatch-*' --scaffold` (기본 ablation, 3회, 상한 $3) | 상한 도달로 부분 결과(case4 with 2회·without 미실행, case5 미실행). with/without: case1 3/3·3/3(FULL 미진입 채점기는 플러그인 없이도 통과 — Δ 0), case2 3/3·0/3, case3 3/3·0/3, case4 with 0/2(`project-context` 미디스패치). $3.70 |

| 2026-10-29 | 채점 기대값 확장 — `dispatching-parallel-agents`는 `systematic-debugging`·`dispatching-parallel-agents` 모두 인정, `requesting-code-review`는 `code-review`도 인정. 재측정 + pressure case4·5 마무리 (3회) | `dispatching-parallel-agents` 3/3, `requesting-code-review` 2/3(1회는 Skill 미호출), pressure case4 with 1/3(`--ablation none`), case5 with 3/3·without 3/3(case1과 같이 Δ 0). $2.20 |

- `--case`를 여러 번 지정하면 마지막 값만 적용됨(2.1.312 실측). 여러 그룹은 실행을 나누거나 glob 하나로 묶음

trigger·knowledge 전부와 위 압박 3건 외의 압박 케이스는 frontmatter 키·채점기 키를 바이너리 스키마의 허용 목록과 대조하는 정적 검사만 거침.
