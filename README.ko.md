# My Poor AI — AI 코딩 에이전트를 위한 엔지니어링 규율

[English](README.md) | **한국어**

[![validate](https://github.com/hicucu/my-poor-ai/actions/workflows/validate.yml/badge.svg)](https://github.com/hicucu/my-poor-ai/actions/workflows/validate.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**my-poor-ai**는 Claude Code가 바이브 코딩 대신 실제 엔지니어링 프로세스 — 테스트 주도 개발, 근본 원인 디버깅, 설계 리뷰, 검증된 완료 — 를 따르게 만듦. 모든 요청을 오케스트레이터가 분류해 알맞은 파이프라인으로 라우팅하고, 전문 서브에이전트를 파견하며, 증거 없이는 "완료" 보고를 거부함.

## 왜 My Poor AI인가

AI 코딩 에이전트는 빠르지만 규율이 없음: 근본 원인 대신 증상을 고치고, 압박 속에서 테스트를 건너뛰고, 검증 없이 완료를 선언함. my-poor-ai는 **스킬 19개**(에이전트가 따라야 할 프로세스 규칙), **서브에이전트 24개**(단일 책임 워커), **슬래시 커맨드 13개**를 오케스트레이터로 엮어 요청마다 맞는 파이프라인을 강제함.

## 빠른 시작

### 1. 마켓플레이스 등록 (최초 1회)

```
/plugin marketplace add hicucu/my-poor-ai
```

### 2. 플러그인 설치

```
/plugin install my-poor-ai@hicucu
/reload-plugins
```

### 3. SessionStart 훅 등록

`~/.claude/settings.json`의 `hooks` 섹션에 아래를 추가함:

```json
"SessionStart": [
  {
    "hooks": [
      {
        "type": "command",
        "command": "bash \"${CLAUDE_PLUGIN_ROOT}/hooks/session-start\"",
        "timeout": 10000
      }
    ]
  }
]
```

등록되면 매 세션 시작(`/clear`, `/compact`, 신규 세션)마다 `using-my-poor-ai` 스킬 컨텍스트가 자동 주입됨.

## 동작 방식

모든 요청은 세 파이프라인 중 하나로 분류됨:

| 경로       | 트리거                          | 파이프라인                                          |
| ---------- | ------------------------------- | ---------------------------------------------------- |
| **DEBUG**  | 버그·에러·예상치 못한 동작      | GOAL.md → systematic-debugging → verification         |
| **SIMPLE** | 파일 1–2개, 설계 불필요, 10분 내 | GOAL.md → TDD → verification                          |
| **FULL**   | 신규 기능, 복잡한 변경          | brainstorming → planning → 병렬 개발 → 4종 병렬 리뷰  |

FULL 경로는 5단계 멀티에이전트 파이프라인: brainstorming 에이전트가 설계 문서를 작성하고(사용자 승인 게이트), planning 에이전트가 TDD 태스크 스펙으로 분해하며, developer 에이전트들이 그룹 단위로 병렬 구현하고, 리뷰 오케스트레이터가 4종 리뷰어(아키텍처/보안/성능/스타일)를 병렬 파견해 결과를 통합하고 issue-fixer를 병렬 파견함.

## 핵심 구성

- **스킬 19개** — TDD, 체계적 디버깅, 브레인스토밍, 플랜 작성, 코드 리뷰(요청·수신), 멀티에이전트 파이프라인, 문서 동기화, 워크트리 격리, 스킬 작성법 등
- **서브에이전트 24개** — project-context 캡처, docs-suite 10개, feature-pipeline 9개, subagent-driven 플로우 4개; 각각 단일 책임과 명시적 입출력 계약 보유 (`AGENTS.md` 참조)
- **슬래시 커맨드 13개** — `/my-poor-ai:code-review`, `/my-poor-ai:detect-stack`, `/my-poor-ai:roles`, 세션 관리·스택 감지·셋업 유틸리티
- **세션 인계** — spec/phase 완료 시 `HANDOFF.md`에 서술형 맥락을 기록해 새 세션이 파이프라인 중간부터 이어받음; `GOAL.md`는 목표·성공 기준을 완료 게이트로 추적
- **멀티플랫폼** — Claude Code 우선; Codex용 에이전트 정의 자동 생성(`.codex/agents/`), Copilot CLI·Gemini CLI 도구 매핑, OpenCode 테스트 스위트

## 역할 프리셋

어떤 스킬로 시작할지 모르겠다면 `/my-poor-ai:roles`가 역할을 스킬 번들로 매핑함:

| 역할          | 번들                                                                 |
| ------------- | -------------------------------------------------------------------- |
| **Architect** | brainstorming → writing-plans → socratic-plan-review                  |
| **Builder**   | test-driven-development → subagent-driven-development → finishing     |
| **Debugger**  | systematic-debugging → verification-before-completion                 |
| **Reviewer**  | requesting-code-review / receiving-code-review / `/my-poor-ai:code-review` |
| **Docs**      | sync-docs-from-diff / generate-claude-instructions                    |

## 바이브가 아니라 검증

my-poor-ai는 자신의 규율을 스스로에게도 적용함:

- **push마다 CI 검증** — `validate-agents.mjs`가 100+ 마크다운 파일의 frontmatter 계약(name/model/도구 화이트리스트), 참조 해소, 코드펜스 균형을 검사; `generate-codex-agents.mjs --check`가 에이전트 정의 24개와 Codex 미러 간 드리프트를 차단
- **행동 테스트** — 스킬은 실제 에이전트 대상 RED–GREEN–PRESSURE 실행으로 검증됨; 워크트리 격리 스킬은 **50/50 실행 무실패** 기록 (GREEN 20 + PRESSURE 20 + 전체 스킬 텍스트 10)
- **적대적 압박 시나리오** — 규율 스킬(TDD, 디버깅)이 시간 압박·매몰 비용·권위 압박, 즉 에이전트가 지름길을 합리화하는 바로 그 조건에서 버티는지 전용 테스트 스위트로 검증

## 저장소 구조

```
my-poor-ai/
├── .claude-plugin/        # 마켓플레이스 + 플러그인 매니페스트
├── .codex/agents/         # Codex용 에이전트 정의 자동 생성물 (수동 편집 금지)
├── agents/                # 서브에이전트 정의 24개 (단일 소스)
├── commands/              # 슬래시 커맨드 13개
├── hooks/                 # SessionStart 훅 (Claude Code + Cursor)
├── skills/                # 스킬 디렉토리 19개
├── scripts/               # CI 검증기 + Codex 미러 생성기
├── tests/                 # 결정론적 + LLM 행동 + 압박 시나리오 스위트
├── docs/                  # 권장 MCP 조합, 등재 초안
├── AGENTS.md              # 에이전트 입출력 계약·불변식
├── CLAUDE.md              # 이 저장소에서 AI 에이전트의 작업 협약
├── CHANGELOG.md
├── CONTRIBUTING.md        # 기여 가이드
└── SECURITY.md
```

## 함께 쓰면 좋은 것

my-poor-ai는 순수 지침이라 외부 연동을 내장하지 않음. [docs/recommended-mcp.md](docs/recommended-mcp.md)가 파이프라인 단계별로 능력을 보강하는 MCP 서버를 정리함 (설계 시 문서 조회, 검증 시 브라우저 자동화, 리뷰 흐름의 GitHub 연동). 전부 없어도 동작함.

약속이 아니라 증거를 원한다면 [examples/go-fractals/](examples/go-fractals/) 참조 — 이 파이프라인이 **완전 무인으로** 만든 실동작 Go CLI(Sierpinski + Mandelbrot ASCII 렌더러)임. 플랜 태스크 10개에 태스크당 TDD 커밋 1개, 자체 테스트·리뷰 단계가 잡아낸 수정 커밋까지 포함. 전체 커밋 이력과 검증 결과는 [PROVENANCE.md](examples/go-fractals/PROVENANCE.md)에 있음. 직접 재현: `bash tests/subagent-driven-dev/run-test.sh go-fractals` (`claude` CLI 호출, 실제 토큰과 10–30분 소요).

## 기원

> 시작은 "문장 종결 어미를 명사로 해"였다.

한 줄짜리 문체 요청이 저장소 전체 감사, 죽은 코드 정리, 단일 소스 에이전트 정의 파이프라인, 그리고 이 공개 릴리스로 이어짐. 엔지니어링 규율에 관한 프로젝트다운 기원임.

## 철학

테스트 우선 개발. 추측 대신 체계적 프로세스. 증거 기반 완료 검증. 복잡성을 전문 에이전트에 위임.

## 기여

이슈와 PR을 환영함 — [CONTRIBUTING.md](CONTRIBUTING.md) 참조. 모든 변경은 CI 검증기를 통과해야 하고, 스킬 변경은 행동 압박 테스트가 필요함 (`skills/writing-skills/` 참조).

## 라이선스

MIT — [LICENSE](LICENSE) 참조. [Superpowers](https://github.com/obra/superpowers) 컨셉 차용 고지는 [NOTICE](NOTICE) 참조.
