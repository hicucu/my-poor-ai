# 코드 리뷰 보고서

## 적용 컨텍스트

- 스택: Svelte 5 (runes) / 클라이언트 전용 SPA (TypeScript, Vite, localStorage 영속화)
- 브랜치: svelte-todo vs master
- 검토 파일 수: 소스 6개 (전체 diff 약 19~27개) / 이슈 등재 파일 8개
- 검토자: Architecture / Security / Performance / Style
- 비고: `_workspaces/stack-profile.json` 부재 — 스택 정보는 각 reviewer 헤더에서 추출

## 요약

| 카테고리     | Critical | High  | Medium | Low    |
| ------------ | -------- | ----- | ------ | ------ |
| Architecture | 0        | 0     | 3      | 2      |
| Security     | 0        | 0     | 0      | 5      |
| Performance  | 0        | 0     | 1      | 3      |
| Style        | 0        | 1     | 1      | 5      |
| **합계\***   | **0**    | **1** | **5**  | **15** |

\* 동일 위치 중복 이슈(localStorage load 스키마 미검증 = Style · Architecture · Security 3중 지적, App.svelte 스토어 정리 = Architecture · Performance 2중 지적)는 각 카테고리 카운트에는 그대로 반영하되, 아래 "파일별 이슈"에서는 병합해 1건으로 표기함.

## 우선순위 (Critical / High만)

| #   | 카테고리                        | 파일:라인                    | 문제                            | 조치                                   |
| --- | ------------------------------- | ---------------------------- | ------------------------------- | -------------------------------------- |
| 1   | Style · Architecture · Security | src/lib/todos.svelte.ts:16-17 | localStorage load 값 스키마 미검증 (`as Todo[]` 무검증 단언) | `isTodo` 타입 가드 + `parsed.todos.filter(isTodo)` |

## 파일별 이슈

### src/lib/todos.svelte.ts

**[High / Style · Architecture · Security]** localStorage load 값 스키마 미검증 (line 16-17)

- 근거: `Array.isArray(parsed.todos)`만 확인하고 각 원소가 `{id, text, completed}` 스키마를 만족하는지 검사 없이 `as Todo[]`로 단언. 손상/오염된 항목이 `$state`로 유입되면 `todo.text` undefined 전파, `t.id === id` 비교 및 `{#each ... (todo.id)}` 키 로직 오작동 가능
- 수정 방향(채택): `typeof item?.id === 'string' && typeof item?.text === 'string' && typeof item?.completed === 'boolean'` 형태의 타입 가드 `isTodo`를 만들어 `parsed.todos.filter(isTodo)`로 필터링 후 사용
- 참고 권장: 실패 항목 제외 또는 전체 빈 배열 폴백 (Security), zod 등 스키마 검증 라이브러리 대안 (Style/Security)

**[Medium / Architecture]** 도메인 상태 관리와 영속화 관심사 혼재 (line 12-25, 52-58)

- 근거: add/toggle/remove/filter 도메인 로직과 localStorage load/save·JSON 직렬화가 한 클래스에 혼재, 저장소 추상화 인터페이스 부재로 교체·모킹 곤란
- 수정 방향: `Storage` 인터페이스(`load(): Persisted` / `save(state): void`) 정의 후 `localStorage` 구현을 생성자 주입, 순수 도메인 로직과 영속화 어댑터 분리

**[Medium / Style]** `catch { ... }` 원인 정보 소실 (line 21-23)

- 근거: JSON 파싱 오류·localStorage 접근 차단 등 모든 예외를 구분 없이 빈 상태로 폴백해 디버깅 곤란
- 수정 방향: `catch (err)`로 받아 개발 모드 한정 `console.warn('failed to load persisted todos', err)` 로깅, 또는 무시 사유를 WHY 주석으로 명시

**[Low / Architecture]** 브라우저 전역(`localStorage`·`crypto.randomUUID()`) 직접 의존 (line 14, 62)

- 수정 방향: ID 생성기를 주입 가능한 의존성으로 분리(선택), 저장소 추상화 도입 시 함께 처리

**[Low / Architecture]** 도메인/저장 표현 미분리 (line 6-9, 52-56)

- 근거: `Persisted`가 도메인 `Todo`를 그대로 직렬화 → 향후 도메인 변경 시 저장 포맷 마이그레이션 부담. 현 규모에선 수용 가능
- 수정 방향: 저장 전용 DTO 매핑 계층 도입을 향후 확장 시 고려

**[Low / Performance]** 필터 변경 시 전체 todos 재직렬화 (line 53-58)

- 근거: 영속화 `$effect`가 `todos`와 `filter`를 한 effect에서 직렬화 → `setFilter`만 해도 전체 `todos`를 `JSON.stringify`로 재기록
- 수정 방향: 저장 대상 분리(todos 전용 + filter 전용 effect) 또는 쓰기 debounce 적용

**[Low / Performance]** `filtered`·`remaining` 독립 순회 중복 (line 35-45, 47)

- 근거: `active` 필터 시 동일한 `!t.completed` 스캔을 두 번 수행. 소규모에선 무시 가능
- 수정 방향: 현 규모에선 조치 불필요, 대규모 확장 시 활성 개수를 1회 순회로 공유 계산

**[Low / Performance]** `filtered` 참조 변경 시 새 배열 생성 (line 36-44) — 정보성

- 근거: keyed `{#each}`가 재조정을 처리하므로 실질 리렌더는 변경 항목에 국한. 병목 아님
- 수정 방향: 조치 불필요(현재 keyed each 구성 올바름)

**[Low / Style]** 콜백 파라미터명 `t` 의미 불명확 (line 33, 35, 41, 71, 77, 81)

- 수정 방향: `t` → `todo`로 리네이밍 (예: `this.todos.filter((todo) => !todo.completed)`)

### src/App.svelte

**[Medium / Architecture · Performance]** 스토어 `destroy()` 미호출 → `$effect.root` 누수 (line 8; 관련 todos.svelte.ts:52-59)

- 근거: `new TodoStore()`가 `$effect.root`로 영속화 이펙트를 등록하지만 컴포넌트 언마운트 시 `store.destroy()`가 호출되지 않아 이펙트 루트 미정리. 실사용 단일 루트 App에선 사실상 무해하나, 테스트 반복(`new TodoStore()` per test) 시 root effect 누적 누수
- 수정 방향: Svelte `onDestroy(() => store.destroy())`로 스토어 수명주기를 컴포넌트 수명주기에 결선
- 참고 권장: 테스트에서 `afterEach`로 `store.destroy()` 호출 (Performance)

### src/lib/TodoInput.svelte

**[Low / Security]** 입력 길이 제한 없음 (line 9-13, `submit`)

- 근거: `text.trim()`만 수행하고 최대 길이 제한이 없어 초장문 반복 추가 시 localStorage 용량(오리진당 5~10MB) 소진 → 저장 실패/앱 비정상(가용성 저하)
- 수정 방향: `maxlength` 속성 및 `submit()` 내 길이 상한 검증(예: 500자) 추가, 초과 시 사용자 안내

### index.html

**[Low / Security]** Content-Security-Policy 미설정 (line 1-13)

- 근거: CSP meta/헤더 부재로 향후 XSS 벡터 발생 시 방어 계층 없음(심층 방어 누락). 현재 코드엔 위험 sink 없음
- 수정 방향: `index.html` meta 태그 또는 정적 호스팅 설정으로 `default-src 'self'` 수준 CSP 추가

### src/main.ts

**[Low / Style]** 세미콜론 미사용 (line 1-9)

- 근거: 프로젝트 다른 `.ts`/`.svelte`는 세미콜론 사용, Vite 스캐폴드 기본 파일만 미사용 → 스타일 혼용
- 수정 방향: 프로젝트 컨벤션(세미콜론 사용)에 맞춰 통일

### svelte.config.js

**[Low / Style]** 세미콜론 미사용 (line 2, `export default {}`)

- 수정 방향: 프로젝트 컨벤션에 맞춰 통일

### src/lib/todos.svelte.test.ts

**[Low / Style]** 매직 스트링 `'test-todos'` 반복 (line 15, 22, 24, 33, 42, 49, 58, 63, 74, 84, 94)

- 수정 방향: 파일 상단에 `const TEST_KEY = 'test-todos';` 상수 추출 후 재사용

### src/lib/TodoFilter.svelte

**[Low / Style]** `filter === option.value` 중복 평가 (line 406-407)

- 근거: `class:selected`와 `aria-pressed`에서 각각 재평가 (우선순위 낮음)
- 수정 방향: `{#each}` 내부에서 `const isSelected = filter === option.value;` 지역 변수 추출

### _workspaces/svelte-todo/HANDOFF.md

**[Low / Security]** 내부 작업 로그가 저장소 이력에 포함 (전체)

- 근거: 파이프라인 내부 메모(작업 경로·진행상황)가 커밋에 포함 → 공개 저장소 배포 시 내부 정보 노출 가능. 보안 결함은 아님
- 수정 방향: 배포/공개 전 `.gitignore`에 `_workspaces/` 추가 여부 검토

### (범위 외 참고 이슈)

**[Low / Security]** 사용자 입력 평문 localStorage 영구 저장 (src/lib/todos.svelte.ts, 영속화 `$effect`)

- 근거: Todo 텍스트에 실수로 민감정보 입력 시 암호화 없이 평문 저장, 동일 브라우저 프로필 타 스크립트·물리적 접근자 열람 가능. 순수 로컬 앱 특성상 통상 설계
- 수정 방향: 앱 성격상 필수 아님. 필요 시 "민감 정보 입력 금지" 안내 또는 저장 전 마스킹/암호화 고려 (문서화 수준 권장)

## 판정

**NEEDS_FIXES**: Critical 0건, High 1건 (localStorage load 스키마 미검증 — Style·Architecture·Security 3중 지적). 해당 High 이슈는 반드시 수정 필요. Medium 5건 / Low 15건은 수정 시 함께 정리 권장하나 대부분 소규모 클라이언트 앱 특성상 실사용 무해 수준.

---

## Appendix: Architecture 원본

# Architecture Review

**검토 범위**: 신규 소스 6개 파일 (App.svelte, lib/todos.svelte.ts, lib/types.ts, lib/Todo{Input,Item,List,Filter,Footer}.svelte, main.ts) / Svelte 5 (runes) + Vite + TypeScript, 컴포넌트 기반 단방향 데이터 흐름 + 클래스 스토어(상태/영속화)

## 총평

전반적으로 아키텍처 품질이 양호함. `App.svelte`를 합성 루트로 두고 `TodoStore`(상태·도메인) 하나만 인스턴스화한 뒤 프레젠테이션 컴포넌트에는 콜백 prop(`onadd`/`ontoggle`/`onremove`/`onfilter`/`onclear`)만 전달하는 단방향 흐름이 일관됨. 의존성 그래프는 `types.ts`(리프) ← `todos.svelte.ts` ← 컴포넌트 ← `App.svelte`로 트리 구조이며 **순환 의존성 없음**. 파일·함수 길이 모두 임계치 이내로 크기 기반 SRP 위반 없음. Critical/High 급 결함은 없으며, 아래는 스토어의 관심사 혼재와 리소스 수명주기 관련 개선 항목임.

| 심각도 | 파일명:라인 | 문제 | 권장 수정 |
| ------ | ----------- | ---- | --------- |
| Medium | src/lib/todos.svelte.ts:12-25, 52-58 | 도메인 상태 관리(add/toggle/remove/filter)와 영속화 관심사(localStorage load/save·JSON 직렬화)가 한 클래스에 혼재. 저장소 추상화 인터페이스가 없어 스토리지 교체·모킹이 어려움 (관심사 분리·추상화 누락) | `Storage` 인터페이스(예: `load(): Persisted` / `save(state): void`)를 정의하고 `localStorage` 구현을 생성자 주입. 순수 도메인 로직과 영속화 어댑터 분리 |
| Medium | src/App.svelte:8 | `new TodoStore()`로 생성한 스토어의 `destroy()`가 어디서도 호출되지 않음. 스토어가 `$effect.root`로 영속화 이펙트를 등록하지만(todos.svelte.ts:52) 컴포넌트 언마운트 시 정리되지 않아 이펙트 루트 누수 (수명주기 계약 미이행·캡슐화) | Svelte `onDestroy(() => store.destroy())`로 스토어 수명주기를 컴포넌트 수명주기에 결선 |
| Medium | src/lib/todos.svelte.ts:17 | `load()`가 `parsed.todos`의 배열 여부만 확인하고 각 항목을 `as Todo[]`로 무검증 캐스팅. 외부(localStorage) 데이터 구조를 신뢰하여 필드 누락/오염된 항목이 도메인으로 유입 가능 (경계에서의 검증 누락) | 항목 단위로 `id`/`text`/`completed` 형태를 검증하거나 스키마 파서(예: 좁히기 가드)를 두어 유효 항목만 통과 |
| Low | src/lib/todos.svelte.ts:14, 62 | `localStorage`·`crypto.randomUUID()` 등 브라우저 전역에 직접 의존. 위 저장소 추상화 항목과 함께 ID 생성기도 주입 가능하게 하면 테스트·환경 이식성이 향상됨 (구체 구현 직접 의존) | ID 생성기를 주입 가능한 의존성으로 분리(선택). 저장소 추상화 도입 시 함께 처리 |
| Low | src/lib/todos.svelte.ts:6-9, 52-56 | 영속화 스키마(`Persisted`)가 도메인 `Todo`를 그대로 직렬화. 도메인 모델과 저장 포맷이 결합되어 향후 도메인 변경 시 저장 포맷 마이그레이션 부담 발생 (도메인/저장 표현 미분리). 현 규모에선 수용 가능 | 저장 전용 DTO로 매핑 계층을 두는 것을 향후 확장 시 고려 |

## 참고 (범위 외·긍정 확인)

- **도메인 노출**: API 레이어가 없는 프론트엔드 앱으로 `Todo`가 뷰 모델을 겸함. 응답 DTO 분리 이슈 해당 없음.
- **강결합**: 프레젠테이션 컴포넌트 간 상호 의존 없음. 콜백 prop 기반으로 느슨하게 결합되어 재사용·단위 테스트 용이. `App.svelte`가 4개 컴포넌트+스토어를 아는 것은 합성 루트 역할상 정상.
- **레이어 위반**: 컴포넌트가 `localStorage`나 스토어 내부 상태를 직접 조작하지 않고 노출된 메서드만 호출. UI/상태 경계 준수.

## Appendix: Security 원본

# Security Review

**검토 범위**: 신규 프로젝트 초기 커밋 (master...HEAD) — 약 27개 파일 (설정/문서 제외 소스 파일 약 12개) / Svelte 5 + TypeScript + Vite, 백엔드 없음(순수 클라이언트 앱, localStorage 영속화) / 민감 데이터 처리 여부: 아니오 (사용자가 자유 텍스트로 입력하는 todo 항목 외 개인정보·인증정보 없음)

이 diff는 신규 백엔드/DB 스키마가 없는 클라이언트 전용 Todo 앱의 최초 스캐폴딩입니다. 서버 API, 인증, DB가 없어 SQL/NoSQL Injection, 인증/인가, CSRF, 감사 로그, DB 마이그레이션 항목은 해당 없음(N/A)으로 확인했습니다. `{@html}`, `innerHTML`, `eval`, `new Function`, `document.write` 등 위험 API는 코드 전체에서 검색했으나 발견되지 않았습니다(`src/lib/TodoItem.svelte` 등 모든 사용자 입력 렌더링은 Svelte의 `{expression}` 자동 이스케이프를 사용하여 표시됨 — XSS 벡터 없음).

아래는 그럼에도 확인이 필요한 항목입니다.

| 심각도 | 파일명:라인 | 취약점 | 근거 | 권장 수정 |
| ------ | ----------- | ------ | ---- | --------- |
| Low | src/lib/todos.svelte.ts:912-925 (`load` 함수) | localStorage 역직렬화 데이터의 스키마 미검증 | `parsed.todos`가 배열인지만 확인(`Array.isArray`)하고 각 원소의 `id`/`text`/`completed` 타입·존재 여부는 검증하지 않은 채 `Todo[]`로 캐스팅함. 사용자가 DevTools로 `localStorage`를 직접 조작하거나 향후 다른 스크립트/확장 프로그램이 같은 오리진에 값을 기록하면 잘못된 형태의 객체가 `$state`에 그대로 로드됨. 현재는 Svelte의 자동 이스케이프 덕분에 XSS로 이어지지는 않으나, `todo.id`가 문자열이 아닐 경우 `toggle`/`remove`의 `t.id === id` 비교나 `{#each todos as todo (todo.id)}` 키 로직이 예기치 않게 동작할 수 있음 | 로드 시 각 항목에 대해 `typeof id === 'string' && typeof text === 'string' && typeof completed === 'boolean'`을 만족하는지 검증(zod 등 스키마 검증 라이브러리 또는 수동 가드)하고, 하나라도 실패하면 해당 항목을 제외하거나 전체를 빈 배열로 폴백 |
| Low | src/lib/TodoInput.svelte:9-13 (`submit` 함수) | 입력 길이 제한 없음 | `text.trim()`만 수행하고 최대 길이 제한이 없어 임의로 매우 긴 문자열을 반복 추가하면 `localStorage`(오리진당 통상 5~10MB) 용량을 소진해 이후 저장이 실패하거나 앱이 비정상 동작할 수 있음(클라이언트 측 가용성 저하) | `maxlength` 속성 및 `submit()` 내 길이 상한 검증(예: 500자) 추가, 초과 시 사용자에게 안내 |
| Low | index.html:1-13 | Content-Security-Policy 미설정 | `<meta http-equiv="Content-Security-Policy">` 또는 서버 응답 헤더가 없어 향후 의존성 추가나 코드 변경으로 XSS 벡터가 생길 경우 방어 계층이 없음. 현재 코드 자체에는 위험 sink가 없어 즉각적인 위험은 아니나 심층 방어(defense-in-depth) 관점에서 누락 | 정적 호스팅 설정 또는 `index.html` meta 태그로 `default-src 'self'` 수준의 CSP 추가 권장 |
| Low | src/lib/todos.svelte.ts:952-959 (영속화 `$effect`) | 사용자 입력이 평문으로 브라우저 저장소에 영구 저장 | Todo 텍스트를 사용자가 자유롭게 입력할 수 있어 실수로 비밀번호·개인정보 등 민감 정보를 입력하면 암호화 없이 `localStorage`에 평문 저장되고 같은 브라우저 프로필의 다른 스크립트나 물리적 접근자가 열람 가능. 순수 로컬 Todo 앱 특성상 통상적인 설계이나, 민감 정보 입력 가능성이 있다면 문서화가 필요 | 앱 성격상 필수는 아니나, 필요 시 사용자에게 "민감 정보 입력 금지" 안내 문구 추가 또는 저장 전 마스킹/암호화 고려 |
| Low | _workspaces/svelte-todo/HANDOFF.md (전체, diff 상단) | 내부 작업 로그가 저장소 이력에 포함 | 보안 결함은 아니지만 파이프라인 내부 메모(작업 경로, 진행상황)가 커밋에 포함되어 공개 저장소로 배포 시 불필요한 내부 정보 노출 가능성 | 배포/공개 전 `.gitignore`에 `_workspaces/` 추가 여부 검토 |

## 확인했으나 이슈 없음으로 판단한 항목

- **XSS**: 모든 사용자 입력(`todo.text` 등)은 Svelte `{expression}` 중괄호 보간으로만 렌더링되며 `{@html}`은 어디에도 사용되지 않음 — 자동 이스케이프 적용됨.
- **Injection(SQL/NoSQL/Command)**: 백엔드/DB가 없어 해당 없음.
- **인증/인가, CSRF**: 서버 통신이 전혀 없는 순수 클라이언트 앱이라 해당 없음.
- **하드코딩 시크릿**: `package.json`, `vite.config.ts`, `svelte.config.js` 등 전 설정 파일에 API 키·토큰·연결 문자열 없음.
- **감사 로그 누락**: 서버 측 데이터 변경이 없으므로 규정상 감사 로그 요건 해당 없음(민감 도메인 아님).
- **DB 마이그레이션 누락**: 신규 테이블/스키마 생성 없음(클라이언트 localStorage 기반), 해당 없음.
- **에러 처리 시 정보 노출**: `todos.svelte.ts`의 `load()`는 JSON 파싱 실패 시 콘솔 로그나 스택 트레이스 노출 없이 조용히 빈 상태로 폴백 — 적절함.

## Appendix: Performance 원본

# Performance Review

**검토 범위**: 6개 소스 파일 (src/lib/todos.svelte.ts, App.svelte, TodoList.svelte, TodoItem.svelte, TodoFilter.svelte, TodoInput.svelte) / Svelte 5 runes + Vite + TypeScript, localStorage 영속화, 클라이언트 전용 SPA (DB/네트워크 I/O 없음)

## 요약

클라이언트 전용 소규모 Todo 앱으로 DB/네트워크 I/O, N+1, 인덱스 이슈는 해당 없음. `{#each}`는 두 곳 모두 안정적 key(`todo.id`, `option.value`)를 사용해 리스트 재조정이 올바르다. 반응성 그래프도 단순·정상. Critical/High 병목은 없으며, 아래는 정확성보다는 자원 효율·누수 예방 관점의 경미한 개선 사항이다.

| 심각도 | 파일명:라인 | 병목 원인 | 예상 영향 | 권장 수정 |
| ------ | ----------- | --------- | --------- | --------- |
| Medium | src/lib/todos.svelte.ts:52-59 | `$effect.root`로 생성한 localStorage 영속화 effect가 컴포넌트 생명주기에 묶이지 않음. `App.svelte`는 `store.destroy()`를 호출하지 않아 App 언마운트/HMR/테스트 반복 시 effect가 해제되지 않음 | 실사용(단일 루트 App)에서는 앱 수명과 동일해 사실상 무해하나, `todos.svelte.test.ts`에서 `new TodoStore()`가 매 테스트마다 해제되지 않는 root effect를 누적 생성 → 테스트 스위트 내 반응성 리스너 누수 | App 언마운트 시 `store.destroy()` 호출(`$effect(() => store.destroy)` 또는 `onDestroy`)하거나, 테스트에서 `afterEach`로 destroy 호출하여 root effect를 정리 |
| Low | src/lib/todos.svelte.ts:53-58 | 영속화 `$effect`가 `todos`와 `filter`를 한 effect에서 함께 직렬화. 필터만 변경(`setFilter`)해도 전체 `todos` 배열을 매번 `JSON.stringify`로 재직렬화·재기록 | 필터 전환 시 불필요한 전체 todos 직렬화/localStorage 쓰기 1회. 목록이 작아 영향 미미하나 항목 수에 비례해 낭비 증가 | 저장 대상을 분리(todos 전용 effect + filter 전용 effect)하거나, 목록이 커질 것을 대비해 쓰기 debounce 적용 |
| Low | src/lib/todos.svelte.ts:35-45,47 | `filtered`(`$derived.by`)와 `remaining`(`$derived`)이 각각 `this.todos`를 독립적으로 `filter` 순회. `active` 필터 시 동일한 `!t.completed` 스캔을 두 번 수행 | todos 변경마다 배열을 최대 2회 전체 순회. 소규모 목록에서는 무시 가능한 마이크로 비용 | 현 규모에선 조치 불필요. 대규모로 확장 시 활성 개수를 한 번의 순회로 공유 계산하는 방식 고려 |
| Low | src/lib/todos.svelte.ts:36-44 | `filtered`가 todos 참조가 바뀔 때마다 새 배열을 생성(`all` 케이스는 원본 반환). TodoList에 매 변경 시 새 배열 전달 | 키 기반 `{#each}`가 재조정을 처리하므로 실질 리렌더는 변경된 항목에 국한됨. 정상 동작이며 병목 아님 | 조치 불필요(정보성). 현재 keyed each 구성이 올바름 |

## 정상 확인 항목 (문제 없음)

- **`{#each}` key**: `TodoList.svelte:19`(`todo.id`), `TodoFilter.svelte:16`(`option.value`) 모두 안정적 key 사용 — 리스트 재정렬/삭제 시 DOM 노드 오재사용·과다 재생성 없음.
- **루프 내 I/O**: 반복문 내 DB/네트워크/localStorage 호출 없음.
- **불필요한 렌더링**: props 기반 단방향 흐름, 컴포넌트 분할이 세밀해 토글/삭제 시 변경 항목만 갱신됨. 큰 목록 가상화는 Todo 앱 도메인상 불필요.
- **이벤트 리스너/타이머 누수**: `setInterval`/수동 `addEventListener`/구독 없음. `onkeydown`/`onclick` 등은 Svelte가 언마운트 시 자동 정리.
- **동기 블로킹**: 무거운 동기 연산 없음. `crypto.randomUUID()`는 경량.

## 결론

Critical/High 성능 이슈 없음. 실제 배포 사용(단일 루트 App) 기준으로는 모두 무해 수준이며, 가장 실질적인 개선은 `$effect.root` 정리(Medium) — 특히 테스트 반복 실행 시 누수 예방 목적이다.

## Appendix: Style 원본

# Code Style Review

**검토 범위**: 19 개 파일 (신규 프로젝트 초기 커밋, package-lock.json 제외) / Svelte 5 (runes) + TypeScript + Vite

주요 검토 대상: `src/App.svelte`, `src/lib/todos.svelte.ts`, `src/lib/TodoInput.svelte`, `src/lib/TodoItem.svelte`,
`src/lib/TodoList.svelte`, `src/lib/TodoFilter.svelte`, `src/lib/TodoFooter.svelte`, `src/lib/types.ts`,
`src/main.ts`, `svelte.config.js`, 각 `*.test.ts`.

Svelte 5 runes(`$props`, `$state`, `$derived`, `$derived.by`, `$effect.root`/`$effect`) 사용은 전반적으로
관용적이며 콜백 prop 기반 단방향 데이터 흐름도 잘 지켜지고 있음. 발견된 이슈는 대부분 Medium/Low 수준.

| 심각도 | 파일명:라인 | 위반 항목 | 권장 수정 |
| ------ | ----------- | --------- | --------- |
| High | src/lib/todos.svelte.ts:16 | localStorage에서 읽은 값을 각 항목의 형태(shape) 검증 없이 `parsed.todos as Todo[]`로 단언. `Array.isArray` 체크만 하고 배열 내부 각 요소가 `{id, text, completed}` 스키마를 만족하는지 확인하지 않아, 손상되거나 임의 형태의 객체가 `Todo`로 취급되어 런타임에 `undefined` 필드가 컴포넌트로 전파될 수 있음(예: `todo.text`가 `undefined`이면 `toUpperCase` 등 호출 시 타입 안전성 붕괴) | 배열 요소마다 `typeof item?.id === 'string' && typeof item?.text === 'string' && typeof item?.completed === 'boolean'` 형태의 타입 가드 함수(`isTodo`)를 만들어 `parsed.todos.filter(isTodo)`로 검증 후 사용 |
| Medium | src/lib/todos.svelte.ts:21-23 | `catch { ... }` 블록이 원인 구분 없이 모든 예외(예: JSON 파싱 오류, localStorage 접근 차단 등)를 동일하게 빈 상태로 폴백. 의도된 설계일 수 있으나 원인 정보가 완전히 소실되어 디버깅이 어려움 | `catch (err)`로 받아 개발 모드에서만 `console.warn('failed to load persisted todos', err)` 등 최소 로깅 추가, 또는 왜 무시해도 되는지 주석(WHY)으로 명시 |
| Low | src/main.ts:1-9 | 세미콜론 미사용. 프로젝트의 다른 모든 `.ts`/`.svelte` 파일(App.svelte, todos.svelte.ts, TodoInput.svelte 등)은 세미콜론을 사용하는데, Vite 스캐폴드 기본 파일만 세미콜론 없이 남아있어 파일 간 스타일이 혼용됨 | 프로젝트 컨벤션(세미콜론 사용)에 맞춰 통일 |
| Low | svelte.config.js:2 | 세미콜론 미사용(`export default {}`), 위와 동일한 스캐폴드 잔존 스타일 불일치 | 프로젝트 컨벤션에 맞춰 통일 |
| Low | src/lib/todos.svelte.ts:33,35,41,71,77,81 | 콜백 파라미터명 `t` (todo 객체)가 의미가 불명확. 파일 전체가 `Todo` 타입을 다루므로 `todo`로 통일하면 가독성 향상 | `t` → `todo`로 리네이밍 (예: `this.todos.filter((todo) => !todo.completed)`) |
| Low | src/lib/todos.svelte.test.ts:15,22,24,33,42,49,58,63,74,84,94 | 문자열 리터럴 `'test-todos'`가 테스트 전체에서 반복(각 `it` 블록마다 `new TodoStore('test-todos')`). DRY 위반은 아니지만 매직 스트링 성격 | 파일 상단에 `const TEST_KEY = 'test-todos';` 상수로 추출 후 재사용 |
| Low | src/lib/TodoFilter.svelte:406-407 | `filter === option.value` 비교식이 `class:selected`와 `aria-pressed`에서 각각 중복 평가됨 | `{#each}` 블록 내부에서 `const isSelected = filter === option.value;` 같은 지역 변수로 추출(단, Svelte 템플릿 표현식 특성상 우선순위는 낮음) |

## 참고 (이슈로 등재하지 않은 확인 사항)
- runes 사용(`$state`, `$derived`, `$derived.by`, `$props`, `$effect.root`/`$effect`)은 Svelte 5 관용 패턴에 부합하며 위반 없음.
- 컴포넌트별 `interface Props` 명명 패턴이 모든 `.svelte` 파일에서 일관되게 사용됨(양호).
- 함수 길이(30줄 이하), 파라미터 개수(4개 이하), 중첩 깊이 모두 기준 내.
- `any` 타입 사용 없음. `unknown` 대체가 필요한 지점 없음.
- 주석 처리된 코드나 WHAT 설명형 주석 없음.
- `TodoStore.destroy()`가 정의되어 있으나 `App.svelte`에서 호출되지 않는 점은 컴포넌트 라이프사이클/리소스 정리 이슈로, 스타일이 아닌 아키텍처 리뷰 영역으로 판단하여 본 리포트에서는 제외함.
