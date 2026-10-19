# 공유 지침 — 리뷰 팬아웃 오케스트레이션 계약

`review-agent`(복잡 경로 Phase 4)가 리뷰어·aggregator·issue-fixer를 부리는 계약임. `review-agent` 정의가 `@include`로 참조하며, Codex 미러에는 인라인됨.

review-agent가 정하는 값(작업 디렉토리·리뷰 대상·게이트·커밋 여부)은 이 문서 말미의 "review-agent가 정하는 것" 절에 둠 — 그 외 항목을 review-agent 정의에 다시 기술하지 않음.

## 리뷰 축과 산출 파일

4개 축을 **전부** 실행함. 축을 임의로 생략하지 않음.

| 축 | 수행 주체 | 산출 파일 |
| --- | --- | --- |
| 아키텍처 | `architecture-reviewer` | `{reviewsDir}/architecture.md` |
| 보안 | 네이티브 `/security-review` (Skill 도구) — 아래 조건이면 `security-reviewer`로 폴백 | `{reviewsDir}/security.md` |
| 성능 | `performance-reviewer` | `{reviewsDir}/performance.md` |
| 스타일 | `style-reviewer` | `{reviewsDir}/style.md` |

`{reviewsDir}` = `{workspaceDir}/reviews/`.

### 보안 축: 네이티브 우선, 에이전트 폴백

네이티브 `/security-review`는 대상 경로 인수를 받지 않고 **세션 저장소의 현재 브랜치와 origin 기본 브랜치 사이 diff**를 검토함 (2026-07-28 실측). 다음 중 하나면 네이티브 대신 `security-reviewer`를 스폰함:

- Skill 목록에 `security-review`가 없거나 호출이 실패함
- 리뷰 대상이 세션 작업 디렉토리의 저장소가 아님
- origin 원격이나 `origin/HEAD`가 없어 네이티브가 diff 범위를 계산하지 못함

네이티브를 썼으면 결과를 `{reviewsDir}/security.md`에 저장하고 첫 줄에 실제 검토 범위를 기록함 — `{base}`와 다르면 그 차이를 명시함. 어느 쪽으로 수행했는지 리포트에 남김.

## 병렬 디스패치 규칙

**아키텍처·성능·스타일 리뷰어 Agent 호출과 `/security-review` Skill 호출(폴백이면 `security-reviewer` Agent 호출)을 단일 응답에 전부 포함함.** 순차 실행 금지 — 축이 서로 독립이므로 순차 실행은 대기 시간만 늘림.

리뷰어는 `subagent_type`으로 스폰함 — `Agent(subagent_type="my-poor-ai:architecture-reviewer")`. 플러그인 설치 시 에이전트는 `my-poor-ai:` 접두사로 네임스페이스됨. 에이전트가 등록되지 않은 환경에서만 범용 에이전트에 `${CLAUDE_PLUGIN_ROOT}/agents/<name>.md`를 읽혀 같은 입력으로 실행함 — 읽은 정의 파일 본문의 변수는 치환되지 않으므로 플러그인 루트 절대 경로를 함께 전달함.

## 리뷰어 공통 입력 계약

리뷰어 에이전트 전부에 동일하게 전달함:

| 항목 | 값 |
| --- | --- |
| 리뷰 대상 | review-agent가 산출한 diff |
| 스택 프로필 | `_workspaces/stack-profile.json` (없으면 "없음" — 리뷰어가 대상에서 스택 추론) |
| 출력 경로 | 위 표의 산출 파일 경로 — "aggregator가 읽는 입력 파일이므로 반드시 작성"이라고 프롬프트에 명시함 (AGENTS.md 불변식 12) |

리뷰어가 파일을 쓰지 않고 표를 텍스트로 반환했으면 review-agent가 그 텍스트를 해당 산출 파일에 저장함.

### 축별 추가 지시

- **architecture-reviewer**: 리뷰 대상에 신규 DB 테이블·엔티티·스키마 생성이 포함되면, 마이그레이션 파일(Flyway, Liquibase, EF Migrations, Alembic 등) 존재 여부 확인을 명시적으로 지시함. 미발견 시 High 이슈로 보고하도록 함 — 스키마 변경의 배포 가능성은 아키텍처 축임 (`/my-poor-ai:code-review`와 같은 배정)

## aggregator 계약

4개 축 산출이 **모두 도착한 뒤** 호출함. 일부만으로 통합하지 않음.

| 항목 | 값 |
| --- | --- |
| 입력 | `{reviewsDir}/` 하위 4개 축 파일 (architecture·security·performance·style) |
| 스택 프로필 | `_workspaces/stack-profile.json` (있으면) |
| 출력 | `{workspaceDir}/review-report.md` — issue-fixer 분배의 입력 파일임을 프롬프트에 명시함. 파일 대신 본문이 반환되면 review-agent가 저장함 |

`review-report.md`는 **issue-fixer 입력 호환 형식**이어야 함 — 파일별 이슈 섹션을 포함하여 파일 단위 분배가 가능해야 함.

## 심각도 → 후속 조치

| 결과 | 조치 |
| --- | --- |
| Critical 0건, High 0건 | 리뷰 통과 — 후속 수정 없이 main agent에게 반환 |
| Critical 또는 High 1건 이상 | issue-fixer 단계 자동 진행 |
| Medium·Low만 존재 | 리포트에 기록하되 자동 수정 대상 아님 — main agent가 사용자에게 보고 |

## issue-fixer 계약

`review-report.md`의 파일별 이슈 섹션에서 대상 파일 목록을 추출하여 **파일 단위로 병렬 호출**함. 서로 다른 파일은 독립이므로 리뷰어와 동일하게 단일 응답에서 동시 스폰함.

| 항목 | 값 |
| --- | --- |
| 대상 파일 | 이슈가 있는 파일 1개 (에이전트 1개당 1파일) |
| 이슈 목록 | 해당 파일의 이슈만 발췌 |
| 스택 프로필 | `_workspaces/stack-profile.json` (있으면) |
| 커밋 여부 | `commit: true` |

동일 파일에 대해 issue-fixer를 2개 이상 띄우지 않음 — 편집 충돌이 발생함.

### 파일에 귀속되지 않는 이슈

일부 이슈는 **변경된 어느 파일에도 귀속되지 않음** — 대표적으로 "마이그레이션 파일 누락"처럼 *있어야 할 것이 없는* 부류임. 파일 단위 분배로는 배정 대상이 없어 그대로 두면 리포트에만 남고 수정되지 않음.

aggregator는 이런 이슈를 `review-report.md`의 **`## 파일 미귀속 이슈`** 섹션에 별도로 모음. review-agent는 각 항목에 대해 다음 중 하나를 선택함:

- **생성할 파일 경로가 특정되면** 그 경로를 대상 파일로 삼아 issue-fixer 배정 (예: 마이그레이션 파일 신규 작성)
- **특정되지 않으면** 자동 수정 대상에서 제외하고 main agent에게 명시적으로 보고 — 조용히 누락시키지 않음

## 절대 금지

- 리뷰 축 **순차 실행** (반드시 단일 응답 병렬)
- 4개 축 산출 중 일부만으로 aggregator 호출
- 폴백 조건이 아닌데 보안 축을 `security-reviewer`로 수행하거나, 네이티브 호출 실패를 숨기고 보안 축을 생략
- 오케스트레이터가 **코드 파일 직접 수정** — 수정은 issue-fixer에 위임
- 동일 파일에 issue-fixer 중복 스폰
- `_workspaces/` 루트에 직접 저장 — 반드시 `{workspaceDir}/` 하위
- 절대 경로·`~/` 사용

## review-agent가 정하는 것

| 항목 | 값 |
| --- | --- |
| `{workspaceDir}` | `_workspaces/review-{branch-slug}/` |
| 리뷰 대상 | 브랜치 diff `{base}...HEAD` |
| 사용자 게이트 | 없음 — 판정 후 자동으로 수정 단계 진행 |
| issue-fixer 커밋 | `commit: true` |
| 출력 프로토콜 | review-agent 정의의 STATUS 블록 |
| 진행 상태 기록 | main agent의 `pipeline-state.md` |

## 단독 커맨드와의 차이 (의도된 분기)

`/my-poor-ai:code-review` 단독 커맨드는 **이 계약을 따르지 않음.** 보안·정확성 축을 네이티브 `/security-review`·`/code-review`에 위임하고 스타일 축을 제거한 구성임(직접 수행은 아키텍처·성능 2축) — 대화형 단독 실행에서는 네이티브 리뷰와의 중복 회피가 우선이기 때문임.

review-agent는 무인 실행이 전제라 스타일 축까지 4축을 유지하고 보안 축만 네이티브 우선으로 둠. **드리프트가 아니라 문서화된 의도적 분기임.**
