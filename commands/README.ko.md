# Commands

[English](README.md) | **한국어**

`/my-poor-ai:{커맨드}` 형태로 실행하는 슬래시 커맨드 12개. `/my-poor-ai:commands`가 카탈로그 진입점이며, 나머지는 아래 표에서 빠르게 확인함.

| 커맨드 | 역할 |
| --- | --- |
| [`my-poor-ai.md`](my-poor-ai.md) | 요청을 알맞은 파이프라인(DEBUG/SIMPLE/FULL)으로 라우팅하는 진입점 — setup 없이도 사용 가능 |
| [`commands.md`](commands.md) | 커맨드 목록 표시, `/my-poor-ai:commands` 카탈로그 자체 |
| [`setup.md`](setup.md) | `SessionStart` 훅을 `~/.claude/settings.json`에 자동 등록 |
| [`codex-setup.md`](codex-setup.md) | my-poor-ai 에이전트와 multi-agent 기능을 `~/.codex/config.toml`에 등록 |
| [`roles.md`](roles.md) | 역할 프리셋 카탈로그 — 역할명(Architect/Builder/Debugger/Reviewer/Docs)을 스킬 번들로 라우팅 |
| [`code-review.md`](code-review.md) | 4개 전문 에이전트(Architecture/Security/Performance/Style) 병렬 리뷰 + 통합 리포트를 산출하는 단독 커맨드 |
| [`detect-stack.md`](detect-stack.md) | 마커 파일 스캔으로 기술 스택을 감지, `stack-profile.json` 생성 (feature-pipeline 전체 실행 없이) |
| [`git-resume.md`](git-resume.md) | 자연어 시간 표현("어제", "지난주") 또는 commit hash로 과거 작업 맥락 복원 |
| [`generate-claudeignore.md`](generate-claudeignore.md) | 감지된 스택과 실제 파일 기반으로 `.claudeignore` 생성, 기존 파일에는 누락 항목 병합 |
| [`graphify-setup.md`](graphify-setup.md) | 코드 그래프 도구(`graphifyy` 또는 `codegraph`) 설치·설정 원스톱 — 패키지 설치, 그래프 생성, Claude Code 통합, git hook 등록 |
| [`session-manager.md`](session-manager.md) | 로컬 Claude Code 세션 전체를 최대 10개 서브에이전트로 병렬 분석, 목록 조회·이름 변경·삭제 |
| [`weekly-commits.md`](weekly-commits.md) | 지정한 GitHub ID/이름으로 이번 주 commit 내역을 markdown 표로 출력, 모노레포는 프로젝트별 분리 |
