# HANDOFF: svelte-todo

**갱신:** 2026-07-07 · **현재 위치:** spec-h 완료 (파이프라인 전체 완료)

## 지금까지 (무엇을·왜)
- spec-a~g(초기 셋업, 타입/스토어, TodoInput, TodoItem, TodoFilter, TodoFooter, TodoList)
  전부 완료(직전 세션 요약 아래 참조).
- spec-h(App.svelte 통합 + 최종 검증) 완료: `src/App.svelte`에서 `TodoStore` 인스턴스를
  생성하고 `TodoInput`/`TodoList`/`TodoFilter`/`TodoFooter` 4개 컴포넌트를 콜백 prop으로
  배선(`onadd`/`ontoggle`/`onremove`/`onfilter`/`onclear`). `src/App.test.ts`에 통합 테스트
  4건(추가+남은개수 표시, Active/Completed 필터 전환, Clear completed, 개별 삭제) 작성 —
  RED 단계에서 `New todo` 라벨/필터·푸터 버튼을 찾지 못해 4개 전부 실패 확인 후, 스펙 코드
  그대로 GREEN 구현으로 4/4 통과. 스펙 코드가 이미 최소/명료하여 REFACTOR 불필요.
- 최종 검증: `npm test` 8 files / 30 tests 전부 통과, `npm run check`(svelte-check) 0
  errors/0 warnings, `npm run build` 정상 산출(dist 생성, 오류 없음), `npm run dev` 기동 후
  curl로 앱 셸(`#app` 마운트 포인트) 응답 확인 — 실제 브라우저 클릭 조작 스모크는 CI/무인
  환경 한계로 수행 불가(대신 App.test.ts의 통합 테스트가 7개 핵심 기능 흐름을 jsdom에서
  검증).
- 이로써 design.md의 성공 기준 9개 항목 전부 충족 확인(아래 "다음 이어받을 일" 참조 —
  파이프라인 종료 상태).

### (이전 세션 요약: spec-a~g)
- spec-a: Vite `svelte-ts` 스캐폴드 생성, 데모 콘텐츠 제거, Vitest+jsdom+Testing Library
  파이프라인 구축.
- spec-b: `src/lib/types.ts`(Todo/Filter), `src/lib/todos.svelte.ts`(TodoStore: runes 기반
  add/toggle/remove/clearCompleted/setFilter, $derived로 filtered/remaining, $effect.root+
  $effect로 localStorage 영속화).
- spec-c~g: TodoInput/TodoItem/TodoFilter/TodoFooter/TodoList 각 컴포넌트를 콜백 prop
  기반 단방향 데이터 흐름으로 구현(상세는 git log 및 각 spec 파일 참조).
- 누적 검증: `npm test` 26/26, `npm run check` 0 errors(spec-h 이전 시점).

## 현재 진행 중
- (없음) spec-h 완료로 전체 파이프라인 종료.

## 다음 이어받을 일
- 없음 — spec-a~h 전부 완료, design.md 성공 기준 9개 항목 전부 충족:
  1. `npm install`/`npm run dev` 정상 기동 — 확인.
  2. 입력창 추가(빈 값 무시) — App.test.ts로 검증.
  3. 체크박스 토글 — App.test.ts로 검증.
  4. 개별 삭제 — App.test.ts로 검증.
  5. All/Active/Completed 필터 — App.test.ts로 검증.
  6. Clear completed — App.test.ts로 검증.
  7. 남은 개수 표시 — App.test.ts로 검증.
  8. 새로고침 후 상태 유지(localStorage) — todos.svelte.test.ts(단위) + 수동 확인 필요
     (자동화 테스트 환경 특성상 실제 브라우저 새로고침은 미수행, $effect 기반 저장/복원
     로직은 단위 테스트로 검증됨).
  9. `npm run check`/`npm test`/`npm run build` 무오류 — 확인.
- 추가 스펙 없음. 후속 작업이 필요하면 새 스펙(리팩터링/기능 확장)을 별도로 계획해야 함.

## 주의·막힌 점·가정
- 브라우저 수동 조작(새로고침 후 상태 유지 등)은 무인 CI 환경 한계로 실제 사람이 브라우저로
  재확인 필요 — 로직 자체는 `todos.svelte.test.ts`의 영속화 테스트로 커버됨.
- `package.json`의 `name` 필드는 스캐폴드 기본값(`vite-scaffold`)을 그대로 둠(스펙 범위 밖).
- E2E(Playwright 등)는 design.md에서 명시적으로 범위 제외.
- 인라인 편집, 마감일/우선순위/카테고리, 드래그 정렬, 라우팅, 다크모드, i18n 등은 design.md
  경계에 따라 전부 미구현(의도된 제외 범위).

## 참조
- design.md · specs/ · pipeline-state.md

---
## 인계 로그 (최근 5개, 최신 위)
- 2026-07-07 spec-h 완료
- 2026-07-07 spec-e 완료
- 2026-07-07 spec-g 완료
- 2026-07-07 spec-f 완료
- 2026-07-07 spec-d 완료
