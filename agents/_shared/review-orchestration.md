# 공유 지침 — 리뷰 팬아웃 오케스트레이션 계약

`review-agent`(복잡 경로 Phase 4)와 `feature-pipeline`(Phase 4~5)이 공통으로 쓰는 리뷰 계약임. 두 파이프라인은 **용도가 달라 각각 유지**하되, 리뷰어 4종·aggregator·issue-fixer를 부리는 방식은 이 문서를 단일 소스로 함.

호출자가 정하는 것은 이 문서 말미의 "호출자 책임" 절에 명시함 — 그 외 항목을 호출자가 다시 기술하지 않음.

## 리뷰 축과 산출 파일

4개 축을 **전부** 실행함. 축을 임의로 생략하지 않음.

| 축 | 에이전트 | 산출 파일 |
| --- | --- | --- |
| 아키텍처 | `architecture-reviewer` | `{reviewsDir}/architecture.md` |
| 보안 | `security-reviewer` | `{reviewsDir}/security.md` |
| 성능 | `performance-reviewer` | `{reviewsDir}/performance.md` |
| 스타일 | `style-reviewer` | `{reviewsDir}/style.md` |

`{reviewsDir}` = `{workspaceDir}/reviews/` (호출자가 `{workspaceDir}`를 정함).

## 병렬 디스패치 규칙

**4개 리뷰어 호출을 단일 응답에 전부 포함함.** 순차 실행 금지 — 축이 서로 독립이므로 순차 실행은 대기 시간만 늘림.

호출 방식은 두 가지가 동등한 계약으로 허용됨 (`AGENTS.md` 호출 방식 2가지 참조):

- `subagent_type`으로 직접 스폰 — `Agent(subagent_type="my-poor-ai:architecture-reviewer")`
- 에이전트 정의 주입 — `{팀_위치}/agents/architecture-reviewer.md를 읽고 그 지침에 따라 작업한다`

> 플러그인 설치 시 에이전트는 `my-poor-ai:` 접두사로 네임스페이스됨. 접두사가 해소되지 않는 환경에서만 bare name을 씀.

## 리뷰어 공통 입력 계약

4개 리뷰어 전부에 동일하게 전달함:

| 항목 | 값 |
| --- | --- |
| 리뷰 대상 | 호출자가 산출한 diff 또는 변경 파일 목록 |
| 스택 프로필 | `{workspaceDir}/stack-profile.json` (없으면 "없음" — 리뷰어가 대상에서 스택 추론) |
| 출력 경로 | 위 표의 산출 파일 경로 |

### 축별 추가 지시

- **security-reviewer**: 리뷰 대상에 신규 DB 테이블·엔티티·스키마 생성이 포함되면, 마이그레이션 파일(Flyway, Liquibase, EF Migrations, Alembic 등) 존재 여부 확인을 명시적으로 지시함. 마이그레이션 미발견 시 High 이슈로 보고하도록 함

## aggregator 계약

4개 산출이 **모두 도착한 뒤** 호출함. 일부만으로 통합하지 않음.

| 항목 | 값 |
| --- | --- |
| 입력 | `{reviewsDir}/` 하위 4개 `*.md` |
| 스택 프로필 | `{workspaceDir}/stack-profile.json` (있으면) |
| 출력 | `{workspaceDir}/review-report.md` |

`review-report.md`는 **issue-fixer 입력 호환 형식**이어야 함 — 파일별 이슈 섹션을 포함하여 파일 단위 분배가 가능해야 함.

## 심각도 → 후속 조치

| 결과 | 조치 |
| --- | --- |
| Critical 0건, High 0건 | 리뷰 통과 — 후속 수정 없이 호출자에게 반환 |
| Critical 또는 High 1건 이상 | issue-fixer 단계 진행 (호출자의 게이트 정책에 따름) |
| Medium·Low만 존재 | 리포트에 기록하되 자동 수정 대상 아님 — 호출자 판단 |

## issue-fixer 계약

`review-report.md`의 파일별 이슈 섹션에서 대상 파일 목록을 추출하여 **파일 단위로 병렬 호출**함. 서로 다른 파일은 독립이므로 리뷰어와 동일하게 단일 응답에서 동시 스폰함.

| 항목 | 값 |
| --- | --- |
| 대상 파일 | 이슈가 있는 파일 1개 (에이전트 1개당 1파일) |
| 이슈 목록 | 해당 파일의 이슈만 발췌 |
| 스택 프로필 | `{workspaceDir}/stack-profile.json` (있으면) |
| 커밋 여부 | 호출자가 지정 (아래 "호출자 책임" 참조) |

동일 파일에 대해 issue-fixer를 2개 이상 띄우지 않음 — 편집 충돌이 발생함.

### 파일에 귀속되지 않는 이슈

일부 이슈는 **변경된 어느 파일에도 귀속되지 않음** — 대표적으로 "마이그레이션 파일 누락"처럼 *있어야 할 것이 없는* 부류임. 파일 단위 분배로는 배정 대상이 없어 그대로 두면 리포트에만 남고 수정되지 않음.

aggregator는 이런 이슈를 `review-report.md`의 **`## 파일 미귀속 이슈`** 섹션에 별도로 모음. 호출자는 각 항목에 대해 다음 중 하나를 선택함:

- **생성할 파일 경로가 특정되면** 그 경로를 대상 파일로 삼아 issue-fixer 배정 (예: 마이그레이션 파일 신규 작성)
- **특정되지 않으면** 자동 수정 대상에서 제외하고 사용자에게 명시적으로 보고 — 조용히 누락시키지 않음

## 절대 금지

- 4개 리뷰어 **순차 실행** (반드시 단일 응답 병렬)
- 4개 산출 중 일부만으로 aggregator 호출
- 오케스트레이터가 **코드 파일 직접 수정** — 수정은 issue-fixer에 위임
- 동일 파일에 issue-fixer 중복 스폰
- `_workspaces/` 루트에 직접 저장 — 반드시 `{workspaceDir}/` 하위
- 절대 경로·`~/` 사용

## 호출자 책임 (이 문서가 정하지 않는 것)

용도가 다르므로 아래는 각 파이프라인이 정함:

| 항목 | 설명 |
| --- | --- |
| `{workspaceDir}` 명명 | `review-agent`는 `_workspaces/review-{branch-slug}/`, `feature-pipeline`은 `_workspaces/{workspaceName}/` |
| 리뷰 대상 산출 방법 | 브랜치 diff(`{base}...HEAD`) vs Phase 2~3 변경 파일 목록 |
| 사용자 게이트 | 리뷰 결과를 사용자 확인에 걸지, 자동으로 수정 단계까지 진행할지 |
| issue-fixer 커밋 | 수정 후 커밋까지 시킬지 여부 |
| 출력 프로토콜 | 호출자에게 무엇을 어떤 형식으로 반환할지 |
| 진행 상태 기록 | `plan.md` 체크박스 / `pipeline-state.md` 갱신 등 |

## 단독 커맨드와의 차이 (의도된 분기)

`/my-poor-ai:code-review` 단독 커맨드는 **이 계약을 따르지 않음.** 보안 축을 네이티브 `/security-review`에 위임하고 스타일 축을 제거한 2축 구성임 — 대화형 단독 실행에서는 네이티브 리뷰와의 중복 회피가 우선이기 때문임.

파이프라인 2종은 무인 실행이 전제라 4축 자체 수행을 유지함. **드리프트가 아니라 문서화된 의도적 분기임.**
