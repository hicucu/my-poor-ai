# 권장 MCP 서버 조합

my-poor-ai는 순수 지침(스킬·에이전트) 플러그인이라 외부 도구 연동을 내장하지 않음. 아래 MCP 서버를 함께 쓰면 파이프라인 단계별 능력이 보강됨. 전부 선택 사항 — 없어도 모든 스킬이 동작함.

| 파이프라인 단계 | MCP 서버 | 보강되는 것 |
| --- | --- | --- |
| brainstorming / feature-planner | [Context7](https://github.com/upstash/context7) | 라이브러리 최신 문서 조회 — 설계 시 API 가정 오류 감소 |
| test-writer / verification | [Playwright MCP](https://github.com/microsoft/playwright-mcp) | 브라우저 기반 E2E 검증 — "동작 증명 후 완료 보고" 원칙의 웹 프런트엔드 커버 |
| systematic-debugging | 프로젝트 DB용 MCP (예: Postgres MCP) | 근본 원인 추적 시 데이터 상태 직접 관찰 |
| review-agent (security) | [GitHub MCP](https://github.com/github/github-mcp-server) | PR·이슈 컨텍스트 연동, 리뷰 결과의 PR 코멘트 반영 |
| 전 단계 공통 | 웹 검색 계열 MCP | 스택 관련 최신 정보 확인 (Claude Code 내장 WebSearch로 대체 가능) |

## 연동 원칙

- 스킬은 MCP 도구를 직접 지시하지 않음 — 에이전트가 사용 가능한 도구 중에서 판단해 사용함. MCP를 설치하면 별도 설정 없이 해당 단계 품질이 올라가는 구조.
- MCP 도구명이 필요한 커스텀 지침은 사용자 CLAUDE.md에 두는 것을 권장 (플러그인 업데이트와 독립적으로 유지됨).
