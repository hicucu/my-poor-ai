# my-poor-ai — AI 에이전트 작업 지침

## AI 에이전트라면 반드시 먼저 읽을 것

작업 시작 전 이 섹션 완독 필수.

이 저장소는 오픈소스 프로젝트임. 외부 기여(이슈·PR)를 환영하며 절차는 CONTRIBUTING.md 참조. 메인테이너 측 변경은 사용자(소유자)와 AI 에이전트 간의 협업으로 이루어짐.

**에이전트의 역할은 사용자를 실수로부터 보호하는 것임.** 검증되지 않은 변경, 범위를 벗어난 수정, 되돌릴 수 없는 작업은 도움이 아니라 피해임.

---

## 작업 전 필수 확인 사항

변경 시작 전 다음을 반드시 수행함:

1. **변경 범위 명확화** — 요청이 모호하거나 선택지가 많을 경우, 추론하기 전에 먼저 질문
2. **기존 파일 확인** — 수정 전 현재 내용을 읽고 영향 범위 파악
3. **완전한 diff 제시** — 2단계 이상의 변경은 Plan Mode로 작성 후 사용자 승인 획득
4. **되돌릴 수 없는 작업 사전 확인** — 삭제, force push, 덮어쓰기는 반드시 확인 후 진행
5. **동작 증명 후 완료 보고** — "완료"는 테스트·빌드·로그로 증명된 경우에만 사용

---

## 코딩 원칙

### TDD (Test-Driven Development)

- 테스트를 먼저 작성하고, 실패를 확인한 뒤, 최소 코드로 통과시킴
- 테스트 없이 작성된 구현 코드는 삭제 후 재작성
- RED → GREEN → REFACTOR 사이클 준수

### YAGNI (You Aren't Gonna Need It)

- 현재 요구사항에 없는 기능은 추가하지 않음
- "나중에 필요할 것 같아서" 식의 추가 금지
- 요청 범위 외 리팩터링·기능 추가 금지 — 발견한 개선점은 별도 제안으로 분리

### DRY (Don't Repeat Yourself)

- 중복 코드는 추상화로 해결
- 단, 과도한 추상화는 복잡성을 높임 — 실제 중복이 2회 이상일 때 적용

### 근본 원인 해결

- 임시방편(workaround) 대신 근본 원인을 해결
- 임시 수정이 불가피한 경우, 기술 부채로 명시하고 별도 태스크로 추적

---

## 문서 작성 규칙

- my-poor-ai가 생산하는 모든 한글 문서(리포트, 산출물, 커밋 메시지 본문 등)는 문장을 명사형으로 종결
- 변환 예: "-한다/-합니다" → "-함", "-된다/-됩니다" → "-됨", "-이다/-입니다" → "-임"
- 짧은 리스트 항목은 개조식 명사(예: "확인", "보고")로 종결 가능
- 사용자에게 직접 묻는 질문·대화형 문장은 예외

---

## Skills 수정 시 주의사항

skills는 단순한 문서가 아니라 에이전트 동작을 형성하는 코드임.

### 수정 전

- `writing-skills` skill을 사용하여 변경사항을 개발하고 테스트
- 수정하려는 skill의 현재 내용과 의도를 먼저 완전히 파악
- 세심하게 튜닝된 콘텐츠(Red Flags 표, 합리화 목록, 핵심 문구)는 증거 없이 변경 금지

### 수정 후

- 여러 세션에서 적대적 압력 테스트 실행
- 변경 전/후 동작 차이를 비교하여 개선 여부 확인
- 관련 없는 변경을 skill 수정에 묶어서 처리하지 않음

### 신규 Skill 추가

- `writing-skills` skill의 가이드라인 준수
- 범용성 확인: 특정 프로젝트·도메인·도구에만 유용한 skill은 별도 파일로 분리
- 기존 skill과의 중복 여부 확인

---

## 파일 구조

```
my-poor-ai/
├── skills/               # 개별 skill 디렉토리 (19개)
│   ├── brainstorming/
│   ├── writing-plans/
│   ├── subagent-driven-development/
│   ├── executing-plans/
│   ├── test-driven-development/
│   ├── systematic-debugging/
│   ├── verification-before-completion/
│   ├── requesting-code-review/
│   ├── receiving-code-review/
│   ├── dispatching-parallel-agents/
│   ├── finishing-a-development-branch/
│   ├── preventing-github-actions-loops/
│   ├── using-git-worktrees/
│   ├── writing-skills/
│   ├── using-my-poor-ai/
│   ├── feature-pipeline/               # 흡수
│   ├── generate-claude-instructions/   # 흡수
│   ├── socratic-plan-review/           # 흡수
│   └── sync-docs-from-diff/            # 흡수
├── agents/               # 서브에이전트 (24개: project-context + docs-suite 10 + feature-pipeline 9 + subagent-driven 플로우 4) — AGENTS.md 참조
├── commands/             # 슬래시 커맨드 (my-poor-ai·setup·codex-setup·commands 카탈로그 + 세부 커맨드 9개)
├── hooks/                # Claude Code hooks
├── README.md
├── AGENTS.md             # 에이전트 명세
├── CLAUDE.md             # 이 파일
├── CONTRIBUTING.md       # 외부 기여 가이드
└── CODE_OF_CONDUCT.md
```

---

## 금지 사항

| 금지 행위                       | 이유                        |
| ------------------------------- | --------------------------- |
| 요청 범위 외 파일 수정          | 회귀 위험, 사용자 의도 위반 |
| 검증 없이 "완료" 보고           | 신뢰성 훼손                 |
| 불확실한 사실을 단정적으로 서술 | 오류 전파                   |
| 되돌릴 수 없는 작업 무단 실행   | 데이터 손실 위험            |
| 요청 없는 리팩터링·기능 추가    | YAGNI 위반                  |

---

## 참고

- 라이선스 및 컨셉 차용 출처 고지: LICENSE / NOTICE 참조
- Skill 작성 가이드: `skills/writing-skills/SKILL.md`
- 전역 AI 동작 지침: `~/.claude/AI_BEHAVIOR.md`
