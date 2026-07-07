---
description: 프로젝트를 분석하여 기술 스택과 실제 존재하는 파일/디렉토리를 기반으로 최적화된 .claudeignore를 생성하거나 기존 파일에 누락 항목을 병합
argument-hint: "[경로] [check|reset] (기본값: 현재 디렉토리, 기존 파일은 병합 모드)"
---

# generate-claudeignore

지정 경로(기본: 현재 디렉토리)를 스캔하여 프로젝트에 최적화된 `.claudeignore` 파일을 생성함.
기존 파일이 있으면 덮어쓰지 않고 누락 항목만 병합함.

## 사용법

```
/generate-claudeignore                    현재 디렉토리 분석 → .claudeignore 생성 또는 병합
/generate-claudeignore <경로>             지정 경로 분석 → <경로>/.claudeignore 생성 또는 병합
/generate-claudeignore check              현재 디렉토리 분석, 추가될 항목 미리보기 (파일 미수정)
/generate-claudeignore <경로> check       지정 경로 분석, 미리보기
/generate-claudeignore reset              기존 .claudeignore 백업 후 전체 재생성
/generate-claudeignore <경로> reset       지정 경로 .claudeignore 전체 재생성
```

---

## Phase 0: 입력 인자 해석

### 0-1. 대상 경로 확정

`$ARGUMENTS`를 분석함.

- `check`·`reset`이 아닌 첫 번째 위치 인수가 있으면 대상 경로로 사용
- 없으면 현재 작업 디렉토리 (`.`)

대상 경로가 존재하지 않으면 아래를 출력하고 종료:

```
오류: 경로를 찾을 수 없습니다 — <경로>
  해결: 경로가 올바른지 확인하거나, 경로 없이 실행하면 현재 디렉토리를 사용합니다.
```

### 0-2. 실행 모드 결정

| 인수    | 모드           | 설명                                               |
| ------- | -------------- | -------------------------------------------------- |
| (없음)  | 생성 또는 병합 | 기존 파일 없으면 생성, 있으면 누락 항목만 추가     |
| `check` | 미리보기       | 파일 수정 없이 추가될 항목만 출력                  |
| `reset` | 강제 재생성    | 기존 파일을 `.claudeignore.bak`으로 백업 후 재생성 |

### 0-3. 기존 .claudeignore 확인

대상 경로에서 `.claudeignore` 존재 여부를 확인함.

| 상황                  | 동작                                                      |
| --------------------- | --------------------------------------------------------- |
| 파일 없음             | 새로 생성 (전체 내용)                                     |
| 파일 있음 + 기본 모드 | 기존 파일 읽기 → 누락 항목만 식별 → 병합 섹션을 끝에 추가 |
| 파일 있음 + `reset`   | 기존 파일을 `.claudeignore.bak`으로 복사 후 전체 재생성   |
| `check` 모드          | 파일 수정 없이 Phase 4 미리보기만 출력                    |

`reset` 시 `.claudeignore.bak`이 이미 존재하면 `.claudeignore.bak2`, `.claudeignore.bak3` 형식으로 넘버링함.

---

## Phase 1: 프로젝트 구조 스캔

### 1-1. 기술 스택 마커 감지

대상 경로에서 아래 마커 파일 존재 여부를 확인하고 해당 패턴을 수집함.
복수 마커가 발견되면 모두 수집함.

| 마커 파일                                               | 감지 스택      | 수집할 패턴                                                                                                                                       |
| ------------------------------------------------------- | -------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| `package.json`                                          | Node.js (기본) | `node_modules/`, `.npm/`, `.yarn/`, `.pnp.*`, `.pnp.js`                                                                                           |
| `package.json` + `next.config.*`                        | Next.js        | `.next/`, `out/`                                                                                                                                  |
| `package.json` + `nuxt.config.*`                        | Nuxt.js        | `.nuxt/`, `.output/`                                                                                                                              |
| `package.json` + `vite.config.*`                        | Vite           | `dist/`                                                                                                                                           |
| `package.json` + `svelte.config.*`                      | SvelteKit      | `.svelte-kit/`, `build/`                                                                                                                          |
| `package.json` + `remix.config.*`                       | Remix          | `build/`, `.cache/`                                                                                                                               |
| `package.json` + `angular.json`                         | Angular        | `dist/`, `.angular/`                                                                                                                              |
| `requirements.txt` 또는 `Pipfile` 또는 `pyproject.toml` | Python         | `__pycache__/`, `*.pyc`, `*.pyo`, `.venv/`, `venv/`, `env/`, `.env/`, `.pytest_cache/`, `*.egg-info/`, `.mypy_cache/`, `.ruff_cache/`, `.pytype/` |
| `Cargo.toml`                                            | Rust           | `target/`                                                                                                                                         |
| `go.mod`                                                | Go             | `vendor/`                                                                                                                                         |
| `*.csproj` 또는 `*.sln` 또는 `*.fsproj`                 | .NET           | `bin/`, `obj/`, `.vs/`                                                                                                                            |
| `pom.xml`                                               | Maven/Java     | `target/`                                                                                                                                         |
| `build.gradle` 또는 `build.gradle.kts`                  | Gradle/JVM     | `build/`, `.gradle/`                                                                                                                              |
| `composer.json`                                         | PHP            | `vendor/`, `storage/logs/`, `bootstrap/cache/`                                                                                                    |
| `Gemfile`                                               | Ruby/Rails     | `vendor/bundle/`, `.bundle/`, `log/`                                                                                                              |
| `pubspec.yaml`                                          | Flutter/Dart   | `.dart_tool/`, `build/`, `.flutter-plugins`                                                                                                       |
| `mix.exs`                                               | Elixir         | `_build/`, `deps/`, `.elixir_ls/`                                                                                                                 |

### 1-2. 실제 디렉토리 존재 확인

마커 감지와 무관하게 대상 경로에서 실제 존재 여부를 직접 확인하여 추가 수집함.

**디렉토리 확인:**

| 실제 존재 디렉토리 | 수집할 패턴                                       |
| ------------------ | ------------------------------------------------- |
| `node_modules/`    | `node_modules/`                                   |
| `dist/`            | `dist/`                                           |
| `build/`           | `build/`                                          |
| `out/`             | `out/`                                            |
| `.next/`           | `.next/`                                          |
| `.nuxt/`           | `.nuxt/`                                          |
| `coverage/`        | `coverage/`                                       |
| `.nyc_output/`     | `.nyc_output/`                                    |
| `_workspaces/`      | `_workspaces/`                                     |
| `logs/`            | `logs/`                                           |
| `tmp/`             | `tmp/`                                            |
| `temp/`            | `temp/`                                           |
| `.cache/`          | `.cache/`                                         |
| `graphify-out/`    | `graphify-out/*`, `!graphify-out/GRAPH_REPORT.md` |
| `target/`          | `target/`                                         |
| `__pycache__/`     | `__pycache__/`, `*.pyc`                           |

**파일 패턴 확인:**

| 실제 존재 파일 패턴                          | 수집할 패턴                     |
| -------------------------------------------- | ------------------------------- |
| `.env` 또는 `.env.*` 파일                    | `.env`, `.env.*`                |
| `*.log` 파일                                 | `*.log`                         |
| `*.sqlite` 또는 `*.sqlite3` 또는 `*.db` 파일 | `*.sqlite`, `*.sqlite3`, `*.db` |
| `*.pem` 또는 `*.key` 파일                    | `*.pem`, `*.key`, `*.cert`      |
| `*.rdb` 파일                                 | `*.rdb`                         |

---

## Phase 2: 콘텐츠 조립

수집된 패턴을 아래 섹션 순서로 조립함.

**조립 규칙:**

- 이미 기존 `.claudeignore`에 있는 패턴은 중복 추가하지 않음.
- 각 섹션은 해당 패턴이 하나라도 있을 때만 출력함.
- 스택별 패턴이 기본 섹션과 겹치면 기본 섹션에 합산하고 스택 섹션은 비워둠.

**중복 판정 알고리즘**: 패턴을 정규화(trailing `/` 제거, `**` → `*` 단순화, 대소문자 무시)한 뒤 집합(set) 비교로 중복 여부 판정함. 예: 기존 `node_modules/`와 신규 `node_modules/**`는 동일 패턴으로 간주해 신규 항목을 추가하지 않음.

### 2-1. 신규 생성 또는 reset 모드 — 전체 파일 구조

```
# .claudeignore
# Claude Code가 컨텍스트에서 제외할 파일/디렉토리
# 생성: generate-claudeignore (<YYYY-MM-DD>)
# 문법: .gitignore 동일

# ── 의존성 ──────────────────────────────────────────────────────────────────
<감지된 의존성 패턴>

# ── 빌드·컴파일 산출물 ──────────────────────────────────────────────────────
<감지된 빌드 산출물 패턴>

# ── 민감 정보 ────────────────────────────────────────────────────────────────
.env
.env.*
*.pem
*.key
*.cert
*.p12
*.pfx
secrets.json
credentials.json
service-account*.json

# ── 대용량 바이너리·미디어 ──────────────────────────────────────────────────
*.jpg
*.jpeg
*.png
*.gif
*.webp
*.ico
*.mp4
*.mp3
*.wav
*.ogg
*.pdf
*.docx
*.xlsx
*.pptx
*.zip
*.tar.gz
*.tar.bz2
*.rar
*.7z

# ── 로그·임시 파일 ──────────────────────────────────────────────────────────
*.log
logs/
tmp/
temp/
.cache/
*.tmp
*.swp
*.bak

# ── 테스트 커버리지·리포트 ──────────────────────────────────────────────────
coverage/
.nyc_output/
test-results/
htmlcov/
.coverage
junit*.xml

# ── 데이터베이스·스토리지 ────────────────────────────────────────────────────
*.sqlite
*.sqlite3
*.db
*.rdb

# ── OS 시스템 파일 ──────────────────────────────────────────────────────────
.DS_Store
Thumbs.db
desktop.ini
ehthumbs.db

# ── IDE·에디터 ───────────────────────────────────────────────────────────────
.idea/
*.suo
*.user
*.userprefs

# ── 스택별 추가 항목 (위 섹션에 미포함 패턴만) ──────────────────────────────
<Phase 1-1에서 수집된 패턴 중 위에 없는 것>

# ── 실제 존재 확인된 추가 항목 (위 섹션에 미포함 패턴만) ────────────────────
<Phase 1-2에서 수집된 패턴 중 위에 없는 것>
```

해당 패턴이 없는 섹션(스택별 추가, 실제 존재 확인)은 출력하지 않음.

### 2-2. 병합 모드 — 기존 파일 끝에 추가

```

# ── 자동 병합 항목 (generate-claudeignore <YYYY-MM-DD>) ─────────────────────
<기존에 없는 항목만>
```

추가할 항목이 하나도 없으면 파일 수정 없이 Phase 4로 이동함.

---

## Phase 3: 파일 쓰기

**신규 생성:**

- Phase 2-1에서 조립된 전체 내용을 `<대상경로>/.claudeignore`에 작성함.

**reset 모드:**

1. 기존 `.claudeignore`를 `.claudeignore.bak`(이미 존재 시 넘버링)으로 복사
2. Phase 2-1 전체 내용으로 새로 작성

**병합 모드:**

- 기존 파일 끝에 빈 줄 1개를 삽입한 후 Phase 2-2 섹션을 추가함.

**check 모드:**

- 파일 수정 없이 Phase 4로 이동함.

---

## Phase 4: 결과 보고

아래 형식으로 출력함.

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  .claudeignore <동작> — <대상 경로>
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  감지된 스택     : <스택 목록 (없으면 "(없음)")>
  마커 파일       : <감지된 마커 파일 목록>
  추가된 패턴 수  : <N>개
  출력 파일       : <경로>/.claudeignore

  추가된 주요 항목:
    <추가된 패턴 목록 (최대 10개, 초과 시 "외 N개" 표기)>

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

참고:
  - SVG(*.svg)는 텍스트 형식이므로 기본 제외 안 함. 대용량 아이콘 번들이면 직접 추가.
  - Lock 파일(package-lock.json 등)은 버전 충돌 디버깅에 유용해 기본 제외 안 함.
  - .env.example 등 예시 파일은 제외 대상 아님 (설정 문서화 역할).
```

`<동작>` 자리:

| 상황           | 표기                                         |
| -------------- | -------------------------------------------- |
| 신규 생성      | `생성 완료`                                  |
| 병합 완료      | `병합 완료`                                  |
| 추가 항목 없음 | `이미 최신 상태`                             |
| check 모드     | `미리보기`                                   |
| reset 모드     | `전체 재생성 완료 (백업: .claudeignore.bak)` |

---

## 에러 핸들링

| 상황                                     | 대응                                                                        |
| ---------------------------------------- | --------------------------------------------------------------------------- |
| 대상 경로 미존재                         | 에러 메시지 출력 후 종료                                                    |
| `.claudeignore` 쓰기 권한 없음           | 에러 메시지 + 경로 확인 안내 후 종료                                        |
| 마커 파일 전혀 없음                      | 스택별 섹션 생략, 공통 기본 항목만 생성 + "스택을 감지하지 못함" 안내 |
| `reset` 시 `.claudeignore.bak` 이미 존재 | `.claudeignore.bak2`, `.bak3` 넘버링으로 충돌 방지                          |

---

## 절대 금지

- 대상 경로 외의 파일 수정 (`.gitignore`, `CLAUDE.md`, `package.json` 등 건드리지 않음)
- 절대 경로·`~/` 사용 (모든 경로는 CWD 기준 상대)
- `.env.example`, `*.example`, `*.sample` 패턴 추가 (예시 파일은 제외 대상 아님)
- `*.svg` 기본 추가 (텍스트 형식, 선택적 제외)
- 마커 미발견 패턴을 임의 추가 (Phase 1 탐지 결과에만 기반할 것)
