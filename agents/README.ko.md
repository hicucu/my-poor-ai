# Agents

[English](README.md) | **한국어**

단일 책임 원칙을 따르는 서브에이전트 정의 16개, 각각 명확한 입출력 계약을 가짐. 각 에이전트 파일의 YAML frontmatter(`name`, `description`, `model`, `tools`)가 정본이며, 이 문서는 빠른 색인 용도임. 호출 순서·Phase 다이어그램·전체 I/O 계약은 저장소 루트의 [`AGENTS.md`](../AGENTS.md) 참조.

이 에이전트들은 사용자가 직접 호출하지 않음 — `using-my-poor-ai` 복잡 경로(FULL)에서 main agent가 `subagent_type`으로 스폰하고, 리뷰 그룹은 `review-agent`가, 나머지는 `sync-docs-from-diff`·`code-review` 스킬이 스폰함.

## 공통 인프라 (1개)

| 에이전트 | 역할 |
| --- | --- |
| [`project-context.md`](project-context.md) | 기능 개발 시작 전 프로젝트 구조·스택·컨벤션·최근 커밋과 공용 `stack-profile.json` 캡처. 24시간 캐시. |

## docs-suite — `sync-docs-from-diff` 서브그룹 (5개)

브랜치 diff를 분석해 문서 패치를 제안(직접 수정 금지)하고, 사용자 승인 후 반영 결과를 검증함.

| 에이전트 | 역할 |
| --- | --- |
| [`change-analyzer.md`](change-analyzer.md) | 베이스 브랜치~HEAD 간 git diff·커밋을 구조화된 변경 분석 보고서로 정리 |
| [`readme-updater.md`](readme-updater.md) | 사용자 대면 동작이 바뀌었을 때 루트 `README.md` 갱신 제안 |
| [`docs-updater.md`](docs-updater.md) | `./docs/` 하위 가이드·튜토리얼·API 레퍼런스·아키텍처 문서 갱신 제안 |
| [`inline-doc-updater.md`](inline-doc-updater.md) | 변경 파일 근처 인라인 문서(컴포넌트 README, 모듈 노트) 갱신 제안 |
| [`doc-sync-validator.md`](doc-sync-validator.md) | 패치 적용 후 누락·표현 불일치·링크/시그니처 불일치 여부 최종 검증 |

## 리뷰 그룹 (6개)

`review-agent`(복잡 경로 Phase 4)가 스폰함. `/my-poor-ai:code-review`는 `architecture-reviewer`·`performance-reviewer`·`review-aggregator`를 재사용함.

| 에이전트 | 역할 |
| --- | --- |
| [`architecture-reviewer.md`](architecture-reviewer.md) | 레이어 위반·결합도·SRP·추상화, 마이그레이션 없는 스키마 변경 관점 리뷰 |
| [`security-reviewer.md`](security-reviewer.md) | OWASP Top 10·인증/인가·입력 검증·시크릿 노출 관점 리뷰. `review-agent`가 네이티브 `/security-review`를 쓸 수 없을 때만 사용 |
| [`performance-reviewer.md`](performance-reviewer.md) | N+1 쿼리·동기 블로킹·메모리 누수·캐싱 관점 리뷰 |
| [`style-reviewer.md`](style-reviewer.md) | 네이밍·중복·매직 넘버·언어 관용 규칙 관점 리뷰 |
| [`review-aggregator.md`](review-aggregator.md) | 전달받은 축별 리뷰 결과를 파일 단위 `review-report.md`로 통합 |
| [`issue-fixer.md`](issue-fixer.md) | 스택 컨벤션에 맞춰 단일 파일의 리뷰 이슈 수정 |

## subagent-driven 플로우 그룹 (4개)

`using-my-poor-ai` 복잡 경로(FULL)에서 main agent가 `subagent_type`으로 직접 스폰함.

| 에이전트 | 역할 |
| --- | --- |
| [`brainstorming-agent.md`](brainstorming-agent.md) | 요구사항 분석 후 설계안 2~3개 비교, 사용자 승인용 `design.md` 작성 |
| [`planning-agent.md`](planning-agent.md) | 승인된 `design.md`를 TDD 태스크 스펙(`specs/*.md`, `file-manifest.json`)으로 분해. DB 스키마 변경이 있으면 마이그레이션 방식을 먼저 확인 |
| [`developer-agent.md`](developer-agent.md) | 단일 스펙을 TDD(RED-GREEN-REFACTOR)로 구현 후 커밋 |
| [`review-agent.md`](review-agent.md) | 리뷰 오케스트레이터 — 아키텍처·성능·스타일 리뷰어와 네이티브 `/security-review`(불가 시 `security-reviewer`)를 병렬 실행한 뒤 aggregator·issue-fixer 실행 |

---

> **`_shared/`**는 에이전트가 아님. `_shared/implementation-conventions.md`는 `developer-agent`의 코딩 규율·스택 컨벤션 매트릭스·테스트 컨벤션, `_shared/stack-profile.md`는 `project-context`의 스택 감지 매트릭스와 `stack-profile.json` 스키마, `_shared/review-orchestration.md`는 `review-agent`의 리뷰 팬아웃 계약을 담으며, 모두 `@include` 지시로 참조됨. [`AGENTS.md`](../AGENTS.md)의 불변식 9 참조.
>
> `.claude-plugin/plugin.json`의 `agents` 목록에 있는 파일만 에이전트로 로드되므로, 이 README와 `_shared/`는 에이전트 목록에 나타나지 않음 (불변식 11).
