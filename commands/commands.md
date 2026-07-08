---
description: my-poor-ai 커맨드 목록 표시 및 진입점. /my-poor-ai:{커맨드} 형태로 세부 커맨드 실행.
---

# my-poor-ai 커맨드

사용 가능한 커맨드 목록. `/my-poor-ai:{커맨드}` 형태로 실행함.

## 스킬 (자연어 트리거)

| 스킬                           | 트리거                                 | 역할                                 |
| ------------------------------ | -------------------------------------- | ------------------------------------ |
| `feature-pipeline`             | "기능 추가해줘", "엔드포인트 만들어줘" | 스택 무관 5단계 기능 개발 파이프라인 |
| `generate-claude-instructions` | "CLAUDE.md 만들어줘", "지침서 생성"    | CLAUDE.md + 참조 문서 4종 생성       |
| `sync-docs-from-diff`          | "docs sync", "문서 동기화"             | 브랜치 diff → README·docs 동기화     |
| `socratic-plan-review`         | 복잡 플랜 검증 요청                    | 산파술로 숨겨진 가정 표면화          |

## 커맨드

| 커맨드                            | 역할                                                                    |
| --------------------------------- | ----------------------------------------------------------------------- |
| `/my-poor-ai:code-review`           | 4-전문 병렬 코드 리뷰 (Architecture/Security/Performance/Style)         |
| `/my-poor-ai:detect-stack`          | 프로젝트 스택 감지 → `_workspaces/stack-profile.json`                    |
| `/my-poor-ai:git-resume`            | 과거 commit 기반 작업 맥락 복원                                         |
| `/my-poor-ai:graphify-setup`        | 코드 그래프 도구 설치 (graphifyy 또는 codegraph)                        |
| `/my-poor-ai:generate-claudeignore` | `.claudeignore` 자동 생성·병합                                          |
| `/my-poor-ai:session-manager`       | 로컬 Claude 세션 목록 조회·이름 변경·삭제                               |
| `/my-poor-ai:weekly-commits`        | 이번 주 커밋 요약                                                       |
| `/my-poor-ai:roles`                 | 역할 프리셋 카탈로그 — 역할명으로 스킬 번들 진입                        |
