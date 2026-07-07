---
description: >
  현재 세션의 git 저장소에서 지정한 GitHub ID 또는 이름으로 이번 주(월~오늘) commit 내역을
  날짜순으로 markdown 테이블로 출력한다. 모노레포인 경우 프로젝트별로 분리하여 표시.
  다음 표현이나 의미적으로 동등한 요청 시 자동 실행:
  "이번주 작업이력", "이번 주 커밋", "이번주에 작업한 내용", "주간 커밋 이력",
  "이번주 작업 내역", "이번 주 작업 정리", "weekly commit", "weekly commits",
  "이번주 {이름} 커밋", "{이름} 이번주 작업", "이번주 작업 보여줘",
  "/weekly-commits", "/weekly-commits {username}"
model: haiku
allowed-tools: [Read, Bash]
---

# 이번 주 커밋 요약

> **실행 원칙**: 모든 단계를 사용자에게 표시하지 않고 조용히 수행함. 최종 마크다운 결과만 출력함.

---

## Step 1: AUTHOR 결정

아래 우선순위로 `AUTHOR`를 결정함.

1. **슬래시 커맨드 인자** — `{args}`가 비어 있지 않으면 그대로 사용
2. **자연어에서 추출** — 사용자 메시지에 이름/ID가 포함되어 있으면 파싱
   - 예: "홍길동 이번주 작업이력" → `AUTHOR = "홍길동"`
   - 예: "이번주 octocat 커밋 보여줘" → `AUTHOR = "octocat"`
3. **미지정 (기본값)** — `git config user.name` 조회. 비어 있으면 `git config user.email` 사용
4. **git config도 없음** — 사용자에게 입력 요청 (이 경우에만 메시지 출력)

**검증 (필수)**: 결정된 `AUTHOR`는 셸 명령에 문자열로 치환되므로, 영숫자·한글·공백·`.` `-` `_` `@` 외의 문자(특히 `` ` `` `$` `"` `'` `;` `|` `&` `(` `)` `<` `>`)가 포함되어 있으면 실행하지 않고 사용자에게 이름을 다시 확인함. 대화에 붙여넣어진 텍스트에서 추출한 값도 예외 없이 검증함.

---

## Step 2: 저장소 확인 및 날짜 계산

```bash
git rev-parse --show-toplevel
```

실패 시 "현재 디렉토리가 git 저장소가 아님." 출력 후 종료.

성공 시 결과를 `REPO_ROOT`, 마지막 경로 세그먼트를 `REPO_NAME`으로 저장.

이번 주 월요일 날짜를 계산하여 `START_DATE`, 오늘을 `END_DATE`로 저장:

```bash
python3 -c "
from datetime import date, timedelta
today = date.today()
monday = today - timedelta(days=today.weekday())
print('START=' + monday.isoformat())
print('END=' + today.isoformat())
"
```

python3 실패 시 node 대안:

```bash
node -e "
const t = new Date(), d = t.getDay();
const m = new Date(t);
m.setDate(t.getDate() - (d === 0 ? 6 : d - 1));
console.log('START=' + m.toISOString().slice(0,10));
console.log('END=' + t.toISOString().slice(0,10));
"
```

둘 다 실패 시 Claude가 현재 날짜 기준으로 직접 계산.

---

## Step 3: 모노레포 감지

`REPO_ROOT`에서 아래를 순서대로 확인함. 하나라도 해당하면 `IS_MONOREPO=true`.

```bash
# 모노레포 마커 파일 존재 여부
ls "{REPO_ROOT}/pnpm-workspace.yaml" 2>/dev/null
ls "{REPO_ROOT}/nx.json" 2>/dev/null
ls "{REPO_ROOT}/lerna.json" 2>/dev/null
ls "{REPO_ROOT}/turbo.json" 2>/dev/null
```

```bash
# 대표 폴더 존재 여부 (packages/, apps/, libs/ 중 하나라도)
ls -d "{REPO_ROOT}/packages" 2>/dev/null
ls -d "{REPO_ROOT}/apps" 2>/dev/null
ls -d "{REPO_ROOT}/libs" 2>/dev/null
```

모노레포인 경우 `PROJECT_ROOTS`를 수집함 — 존재하는 폴더(`packages`, `apps`, `libs`) 각각의 직계 하위 디렉토리 이름 목록.

---

## Step 4: git log 조회

```bash
git -C "{REPO_ROOT}" log --all \
  --author="{AUTHOR}" \
  --after="{START_DATE}T00:00:00" \
  --before="{END_DATE}T23:59:59" \
  --no-merges \
  --format="COMMIT|%H|%h|%ad|%s" \
  --date=format:"%Y-%m-%d %H:%M" \
  --name-only
```

`--name-only`로 각 커밋이 수정한 파일 목록을 함께 가져옴 (모노레포 프로젝트 판별용).

결과가 없으면 아래 메시지만 출력하고 종료:

```
이번 주({START_DATE} ~ {END_DATE}) {AUTHOR}의 커밋이 없습니다.
```

---

## Step 5: 파싱 및 중복 제거

출력 형식:

```
COMMIT|{full_sha}|{short_sha}|{date}|{message}
{파일1}
{파일2}
             ← 빈 줄 (커밋 구분자)
COMMIT|...
```

파싱 규칙:

- `COMMIT|`으로 시작하는 줄 → 새 커밋 시작, 필드 분리
- 그 외 비어 있지 않은 줄 → 해당 커밋의 변경 파일 목록에 추가
- **중복 제거**: `full_sha`가 이미 처리된 커밋은 건너뜀 (동일 커밋이 여러 브랜치에 중복 출력되는 경우 방지)

---

## Step 6: 그룹화 및 마크다운 출력

### 일반 저장소 (`IS_MONOREPO=false`)

브랜치 구분 없이 날짜 내림차순으로 단일 테이블 출력:

```markdown
# {AUTHOR}의 이번 주 커밋 요약

> **저장소**: {REPO_NAME}
> **기간**: {START_DATE} (월) ~ {END_DATE}
> **총 커밋**: {N}개

| SHA       | 날짜             | 메시지      |
| --------- | ---------------- | ----------- |
| `abc1234` | YYYY-MM-DD HH:MM | 커밋 메시지 |
```

### 모노레포 (`IS_MONOREPO=true`)

각 커밋의 변경 파일 경로에서 프로젝트를 판별함:

- `packages/foo/src/index.ts` → 프로젝트 `foo`
- `apps/web/pages/index.tsx` → 프로젝트 `web`
- 위 규칙에 해당하지 않는 파일(루트 설정 파일 등) → `(root)`

한 커밋이 여러 프로젝트에 걸쳐 있으면 관련 프로젝트 모두에 표시 (SHA 기준 중복 처리와 별개 — 같은 SHA가 두 프로젝트 섹션에 모두 나올 수 있음).

커밋 수 내림차순으로 프로젝트 섹션 정렬:

```markdown
# {AUTHOR}의 이번 주 커밋 요약

> **저장소**: {REPO_NAME} (모노레포)
> **기간**: {START_DATE} (월) ~ {END_DATE}
> **총 커밋**: {N}개 | 관련 프로젝트: {P}개

## {프로젝트명1}

| SHA       | 날짜             | 메시지      |
| --------- | ---------------- | ----------- |
| `abc1234` | YYYY-MM-DD HH:MM | 커밋 메시지 |

## {프로젝트명2}

| SHA       | 날짜             | 메시지      |
| --------- | ---------------- | ----------- |
| `def5678` | YYYY-MM-DD HH:MM | 커밋 메시지 |
```

---

## Step 7: 커밋 메시지 분석 및 정리 요약 출력

Step 6 출력 직후, 수집된 전체 커밋 메시지를 Claude가 직접 분석하여 정리 요약을 출력함.

### 분석 규칙

- **의미 중복 제거**: 동일하거나 실질적으로 같은 작업을 가리키는 메시지는 하나로 합침
  - 예: "fix: 로그인 버그 수정", "fix: 로그인 오류 재수정" → "로그인 버그 수정"
- **관련 작업 묶기**: 같은 기능/컴포넌트에 대한 여러 커밋은 하나의 항목으로 통합함
  - 예: "feat: 결제 UI 추가", "style: 결제 UI 스타일 조정", "fix: 결제 UI 버그" → "결제 UI 구현"
- **Conventional Commits 접두어 제거**: `feat:`, `fix:`, `chore:`, `refactor:` 등 타입 접두어를 제거하고 핵심 내용만 남김
- **날짜·SHA 미포함**: 정리 요약에는 날짜와 SHA를 표시하지 않음

### 출력 형식

**일반 저장소**:

```markdown
---

## 이번 주 작업 정리

- 작업 항목 1
- 작업 항목 2
- 작업 항목 3
```

**모노레포**:

```markdown
---

## 이번 주 작업 정리

### {프로젝트명1}

- 작업 항목 1
- 작업 항목 2

### {프로젝트명2}

- 작업 항목 1
```

프로젝트 순서는 Step 6와 동일(커밋 수 내림차순). 항목은 중요도·완결성 기준으로 위에서 아래로 정렬.

---

## 에러 처리

| 상황              | 처리                                                                       |
| ----------------- | -------------------------------------------------------------------------- |
| git 저장소 아님   | 에러 메시지 출력 후 종료                                                   |
| 커밋 없음         | "커밋 없음" 한 줄 출력. `git shortlog -sn --all`로 후보 이름이 있으면 제안 |
| python3/node 없음 | Claude가 날짜 직접 계산                                                    |
