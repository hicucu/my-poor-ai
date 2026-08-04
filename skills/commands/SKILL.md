---
name: commands
description: my-poor-ai 커맨드 목록 표시 및 진입점. /my-poor-ai:{커맨드} 형태로 세부 커맨드 실행.
---

# my-poor-ai 커맨드

사용 가능한 목록. 전부 `/my-poor-ai:{이름}` 형태로 실행 가능함.

> 커맨드와 스킬은 같은 형식임 — 모두 `skills/{이름}/SKILL.md`에 위치함. 차이는 **누가 시작하는가**뿐임: 아래 "스킬"은 요청 내용이 맞으면 Claude가 자동 로드하고, "커맨드"는 사용자가 호출함. 프로젝트 외부에 부작용이 있는 셋업 3종은 `disable-model-invocation`으로 자동 발동이 차단됨.

## 스킬 (자연어 트리거)

| 스킬                           | 트리거                                 | 역할                                 |
| ------------------------------ | -------------------------------------- | ------------------------------------ |
| `feature-pipeline`             | "기능 추가해줘", "엔드포인트 만들어줘" | 스택 무관 5단계 기능 개발 파이프라인 |
| `sync-docs-from-diff`          | "docs sync", "문서 동기화"             | 브랜치 diff → README·docs 동기화     |
| `socratic-plan-review`         | 복잡 플랜 검증 요청                    | 산파술로 숨겨진 가정 표면화          |

## 커맨드

| 커맨드                            | 역할                                                                    |
| --------------------------------- | ----------------------------------------------------------------------- |
| `/my-poor-ai`                       | 진입점 — 요청을 알맞은 파이프라인(DEBUG/SIMPLE/FULL)으로 라우팅          |
| `/my-poor-ai:code-review`           | 아키텍처·성능 리뷰 + 네이티브 `/security-review` 통합, `REVIEW.md` 생성  |
| `/my-poor-ai:detect-stack`          | 프로젝트 스택 감지 → `_workspaces/stack-profile.json`                    |
| `/my-poor-ai:git-resume`            | 과거 commit 기반 작업 맥락 복원                                         |
| `/my-poor-ai:session-manager`       | 로컬 Claude 세션 목록 조회·이름 변경·삭제                               |
| `/my-poor-ai:roles`                 | 역할 프리셋 카탈로그 — 역할명으로 스킬 번들 진입                        |

### 셋업 (사용자 호출 전용 — 프로젝트 외부에 기록함)

| 커맨드                            | 역할                                                                    |
| --------------------------------- | ----------------------------------------------------------------------- |
| `/my-poor-ai:setup`                 | SessionStart 훅을 `~/.claude/settings.json`에 등록                       |
| `/my-poor-ai:codex-setup`           | 에이전트를 `~/.codex/config.toml`에 등록                                 |
