# Svelte Todo 앱 설계

> 상태: **승인됨** (무인 CI 실행 — 오케스트레이터가 사용자 대리 검토 수행)
>
> **검토 노트 (오케스트레이터, 2026-07-07):** 접근법(Svelte 5 runes + 컴포넌트 분리)과 파일 목록은 요구사항 7개 항목을 모두 커버하며 과설계 없이 적절함. 초안에서 "자동화 테스트 프레임워크 도입은 범위 밖"이라 되어 있던 부분은 TDD 필수 방침 및 최종 `npm test` 통과 요건과 충돌하여 **수정**: Vitest + Testing Library를 Task 1 셋업에 포함하고, 상태 모듈/컴포넌트별 테스트 파일과 성공 기준에 `npm test`를 추가함. 이 외 설계는 그대로 승인. Planning 단계로 진행.

## 요구사항 요약

라이트한 로컬 저장형 Todo 리스트 웹앱을 Svelte로 구축한다.

**핵심 기능**

- Todo 추가
- 완료/미완료 토글
- Todo 삭제
- 필터: All / Active / Completed
- 완료 항목 일괄 삭제 (Clear completed)
- 남은(미완료) 항목 수 표시
- localStorage 영속화

**경계 (이번 범위 아님 / YAGNI)**

- 백엔드 / 서버 / 인증
- 다중 사용자 / 동기화
- Todo 수정(인라인 편집) — 요구사항에 없음
- 마감일, 우선순위, 카테고리, 드래그 정렬
- 라우팅 (단일 페이지로 충분)
- 테마 전환, i18n

**제약**

- 신규 프로젝트, 기존 코드/패턴 없음.
- 환경: Node v22.22, npm 10.9. 사용 가능 버전 — Svelte 5.56, Vite 8.1.
- Svelte 5가 기본이므로 runes(`$state`, `$derived`, `$effect`) 기반으로 구현하는 것이 관용적.

**성공 기준**

- 위 7개 기능이 브라우저에서 동작.
- 새로고침 후에도 목록/완료 상태가 유지됨.
- `npm run dev` / `npm run build` / `npm run check` 무오류.

## 선택한 접근법: 반응형 상태 모듈 + 컴포넌트 분리 (Svelte 5 Runes)

Vite + Svelte 5 + TypeScript 스캐폴드 위에, 상태 로직을 단일 반응형 모듈(`.svelte.ts`)로 분리하고 UI를 작은 컴포넌트로 나눈다. 뷰(컴포넌트)와 도메인 로직(상태 모듈)을 분리해 테스트/유지보수가 쉽고, Svelte 5 runes로 파생 값(남은 개수, 필터링된 목록)과 localStorage 동기화를 선언적으로 처리한다.

**아키텍처**

- **상태 모듈 `src/lib/todos.svelte.ts`**: `$state`로 todos 배열과 현재 필터를 보유. `add / toggle / remove / clearCompleted / setFilter` 조작 함수 제공. `$derived`로 `filtered`(현재 필터에 맞는 목록)와 `remaining`(미완료 개수) 계산. `$effect`로 todos 변경 시 localStorage에 저장, 초기값은 localStorage에서 로드(JSON 파싱 실패 시 빈 배열로 폴백).
- **타입 `src/lib/types.ts`**: `Todo { id: string; text: string; completed: boolean }`, `Filter = 'all' | 'active' | 'completed'`.
- **컴포넌트**
  - `App.svelte`: 레이아웃 컨테이너. 상태 모듈을 소비.
  - `TodoInput.svelte`: 새 todo 입력 + Enter/버튼으로 추가. 빈 문자열은 무시.
  - `TodoList.svelte`: `filtered` 목록을 렌더, `TodoItem` 반복.
  - `TodoItem.svelte`: 체크박스(토글) + 텍스트 + 삭제 버튼.
  - `TodoFilter.svelte`: All/Active/Completed 선택 버튼.
  - `TodoFooter.svelte`: 남은 개수 표시 + Clear completed 버튼.
- **데이터 흐름**: 컴포넌트는 상태 모듈의 함수를 호출(단방향). 상태 변경 → runes 반응성 → UI 자동 갱신 → `$effect`가 localStorage 반영.

**장점**

- 도메인 로직 한 곳 집중 → 가독성/테스트 용이.
- Svelte 5 관용적. 파생 값·영속화가 선언적이라 버그 표면 적음.
- 컴포넌트가 작아 개별 수정 쉬움.

**단점**

- 작은 앱 치고 파일 수가 다소 많음(초기 오버헤드 약간).

**추천 여부:** 권장 — 요구사항 규모에 과하지 않으면서 확장/검증에 유리하고 Svelte 5 관례에 부합.

## 변경 파일 목록

| 파일 | 작업 | 역할 |
| ---- | ---- | ---- |
| (프로젝트 루트 스캐폴드) | create | **Task 1: 초기 셋업** — Vite `svelte-ts` 템플릿 생성, `npm install`, Svelte 5 확인 |
| `package.json` | create | 의존성/스크립트 (dev, build, preview, check) |
| `vite.config.ts` | create | Vite + `@sveltejs/vite-plugin-svelte` 설정 (스캐폴드 기본) |
| `tsconfig.json` / `tsconfig.node.json` | create | TS 설정 (스캐폴드 기본) |
| `svelte.config.js` | create | Svelte 설정 (스캐폴드 기본) |
| `index.html` | create/modify | 앱 엔트리, 타이틀 조정 |
| `src/main.ts` | create/modify | 앱 마운트 |
| `src/app.css` | create/modify | 전역 스타일 (스캐폴드 정리) |
| `src/lib/types.ts` | create | `Todo`, `Filter` 타입 |
| `src/lib/todos.svelte.ts` | create | 반응형 상태 모듈 + localStorage 동기화 |
| `src/App.svelte` | create/modify | 레이아웃 컨테이너 |
| `src/lib/TodoInput.svelte` | create | 새 todo 입력/추가 |
| `src/lib/TodoList.svelte` | create | 필터된 목록 렌더 |
| `src/lib/TodoItem.svelte` | create | 개별 항목(토글/삭제) |
| `src/lib/TodoFilter.svelte` | create | All/Active/Completed 필터 |
| `src/lib/TodoFooter.svelte` | create | 남은 개수 + Clear completed |
| `.gitignore` | create | `node_modules`, `dist` 무시 |
| `vitest.config.ts` (또는 `vite.config.ts` 내 `test` 필드) | create | Vitest + `@testing-library/svelte` 설정 |
| `src/lib/todos.svelte.test.ts` | create | 상태 모듈 단위테스트 (add/toggle/remove/clearCompleted/filter/persist) |
| `src/lib/*.test.ts` (컴포넌트별) | create | 각 컴포넌트 렌더/동작 테스트 |

> 스캐폴드 생성 방식: `npm create vite@latest . -- --template svelte-ts` (현재 디렉토리에 생성). 생성 후 데모 컴포넌트(`Counter.svelte` 등) 제거하고 위 구조로 교체.
> **테스트 도구:** Vitest + `@testing-library/svelte` + `jsdom`을 Task 1에서 함께 설치. `package.json`에 `"test": "vitest run"` 스크립트 추가. 모든 태스크는 TDD(RED-GREEN-REFACTOR)로 진행하며 `npm test`가 최종적으로 무오류 통과해야 함.

## 제외 범위

- 인라인 편집, 마감일/우선순위, 정렬, 검색.
- 백엔드/동기화/인증.
- 라우터, 상태관리 라이브러리(runes로 충분).
- E2E 테스트(Playwright 등) — 단위/컴포넌트 테스트로 충분한 범위이므로 이번엔 도입하지 않음.

## 성공 기준

- [ ] `npm install` 및 `npm run dev` 정상 기동.
- [ ] 입력창에서 todo 추가 가능(빈 값 무시).
- [ ] 체크박스로 완료/미완료 토글.
- [ ] 개별 todo 삭제.
- [ ] All / Active / Completed 필터 정상 동작.
- [ ] Clear completed로 완료 항목 일괄 삭제.
- [ ] 남은(미완료) 개수 정확히 표시.
- [ ] 새로고침 후 목록/완료/필터 상태 유지(localStorage).
- [ ] `npm run check`(svelte-check), `npm test`(vitest), `npm run build` 모두 무오류.

## 검토된 대안

### 접근법 B: 단일 App.svelte 모놀리식

**요약:** 상태와 모든 UI를 `App.svelte` 하나에 담고 localStorage도 그 안에서 처리.

**탈락 이유:** 파일이 적어 초기엔 빠르나, 로직과 뷰가 뒤섞여 가독성/재사용/검증이 나빠짐. 요구사항이 7개 기능으로 이미 여러 관심사를 포함해 최소한의 분리가 유리. 학습·유지보수 관점에서 접근법 A가 우세.

### 접근법 C: 클래식 Svelte store(`writable`) 사용

**요약:** runes 대신 `svelte/store`의 `writable`로 상태를 만들고 `$store` 구독으로 소비, `store.subscribe`로 localStorage 저장.

**탈락 이유:** Svelte 5 신규 프로젝트에서 runes가 권장 방식이며 파생 값·영속화가 `$derived`/`$effect`로 더 간결. 클래식 store는 하위호환은 좋지만 이 프로젝트엔 이점이 없어 관용성 낮은 선택이 됨. 접근법 A가 더 자연스러움.
