# Agents

[English](README.md) | **한국어**

단일 책임 원칙을 따르는 서브에이전트 정의 24개, 각각 명확한 입출력 계약을 가짐. 각 에이전트 파일의 YAML frontmatter(`name`, `description`, `model`, `tools`)가 정본이며, 이 문서는 빠른 색인 용도임. 호출 순서·Phase 다이어그램·전체 I/O 계약은 저장소 루트의 [`AGENTS.md`](../AGENTS.md) 참조.

이 에이전트들은 사용자가 직접 호출하지 않음 — 스킬 오케스트레이터(`feature-pipeline`, `generate-claude-instructions`, `sync-docs-from-diff`)가 호출하거나, `using-my-poor-ai` 복잡 경로(FULL)에서 main agent가 `subagent_type`으로 직접 스폰함.

## 공통 인프라 (1개)

| 에이전트 | 역할 |
| --- | --- |
| [`project-context.md`](project-context.md) | 기능 개발 시작 전 프로젝트 구조·스택·컨벤션·최근 커밋 캡처. 24시간 캐시. |

## docs-suite — `generate-claude-instructions` 서브그룹 (5개)

Phase 1에서 앞 4개를 병렬 호출하고, Phase 2에서 composer가 그 산출물을 `CLAUDE.md`로 합성함.

| 에이전트 | 역할 |
| --- | --- |
| [`dev-principles.md`](dev-principles.md) | `DEVELOPMENT.md` 작성 — SOLID, TDD, 클린 코드, 보안, 성능 원칙 |
| [`language-guidelines.md`](language-guidelines.md) | `LANGUAGE_GUIDELINES.md` 작성 — 감지된 언어/프레임워크별 섹션 구성 |
| [`ai-behavior.md`](ai-behavior.md) | `AI_BEHAVIOR.md` 작성 — 응답 형식, 워크플로우, 도구 사용, 자기 검증 원칙 |
| [`commit-convention.md`](commit-convention.md) | `COMMIT_CONVENTION.md` 작성 — 프로젝트 commitlint 설정 또는 Conventional Commits 기반 |
| [`claude-md-composer.md`](claude-md-composer.md) | 위 4개 문서를 읽고 포인터 형태의 간결한 `CLAUDE.md`로 합성 |

## docs-suite — `sync-docs-from-diff` 서브그룹 (5개)

브랜치 diff를 분석해 문서 패치를 제안(직접 수정 금지)하고, 사용자 승인 후 반영 결과를 검증함.

| 에이전트 | 역할 |
| --- | --- |
| [`change-analyzer.md`](change-analyzer.md) | 베이스 브랜치~HEAD 간 git diff·커밋을 구조화된 변경 분석 보고서로 정리 |
| [`readme-updater.md`](readme-updater.md) | 사용자 대면 동작이 바뀌었을 때 루트 `README.md` 갱신 제안 |
| [`docs-updater.md`](docs-updater.md) | `./docs/` 하위 가이드·튜토리얼·API 레퍼런스·아키텍처 문서 갱신 제안 |
| [`inline-doc-updater.md`](inline-doc-updater.md) | 변경 파일 근처 인라인 문서(컴포넌트 README, 모듈 노트) 갱신 제안 |
| [`doc-sync-validator.md`](doc-sync-validator.md) | 패치 적용 후 누락·표현 불일치·링크/시그니처 불일치 여부 최종 검증 |

## feature-pipeline 그룹 (9개)

`feature-pipeline` 스킬의 5단계 파이프라인: 계획 → 구현 → 테스트 → 리뷰 → 수정.

| 에이전트 | 역할 |
| --- | --- |
| [`feature-planner.md`](feature-planner.md) | 스택 감지 후 기능 요청을 파일 단위 계획(`stack-profile.json`, `plan.md`, `file-manifest.json`)으로 분해 |
| [`file-developer.md`](file-developer.md) | 명세에 따라 단일 파일 구현/수정, 언어·프레임워크 무관 |
| [`test-writer.md`](test-writer.md) | 비즈니스 로직 파일 단위테스트 작성, 감지된 테스트 프레임워크 자동 선택 |
| [`architecture-reviewer.md`](architecture-reviewer.md) | 레이어 위반·결합도·SRP·추상화 관점 리뷰 |
| [`security-reviewer.md`](security-reviewer.md) | OWASP Top 10·인증/인가·입력 검증·시크릿 노출 관점 리뷰 |
| [`performance-reviewer.md`](performance-reviewer.md) | N+1 쿼리·동기 블로킹·메모리 누수·캐싱 관점 리뷰 |
| [`style-reviewer.md`](style-reviewer.md) | 네이밍·중복·매직 넘버·언어 관용 규칙 관점 리뷰 |
| [`review-aggregator.md`](review-aggregator.md) | 4개 리뷰 결과를 파일 단위 `review-report.md`로 통합 |
| [`issue-fixer.md`](issue-fixer.md) | 스택 컨벤션에 맞춰 단일 파일의 리뷰 이슈 수정 |

## subagent-driven 플로우 그룹 (4개)

`using-my-poor-ai` 복잡 경로(FULL)에서 main agent가 `subagent_type`으로 직접 스폰함. 리뷰 단계는 위 feature-pipeline 그룹의 리뷰어 4종·review-aggregator·issue-fixer를 재사용함.

| 에이전트 | 역할 |
| --- | --- |
| [`brainstorming-agent.md`](brainstorming-agent.md) | 요구사항 분석 후 설계안 2~3개 비교, 사용자 승인용 `design.md` 작성 |
| [`planning-agent.md`](planning-agent.md) | 승인된 `design.md`를 TDD 태스크 스펙(`specs/*.md`, `file-manifest.json`)으로 분해 |
| [`developer-agent.md`](developer-agent.md) | 단일 스펙을 TDD(RED-GREEN-REFACTOR)로 구현 후 커밋 |
| [`review-agent.md`](review-agent.md) | 리뷰 오케스트레이터 — 리뷰어 4종·aggregator·issue-fixer를 직접 스폰하는 독립 플로우 |
