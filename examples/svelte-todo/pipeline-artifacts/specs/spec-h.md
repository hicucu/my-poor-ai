---
complexity: complex
depends-on: spec-c, spec-e, spec-f, spec-g
estimated-tasks: 6
---

# Spec H: App.svelte 통합 + 최종 검증

## 목표

`TodoStore` 를 생성해 모든 컴포넌트(`TodoInput`, `TodoList`, `TodoFilter`,
`TodoFooter`)를 배선하고, 상태 모듈의 조작 함수를 콜백으로 연결하는 최상위
레이아웃을 완성한다. 통합 테스트로 추가/토글/필터/일괄 삭제/남은 개수 흐름을
검증하고, 성공 기준(dev/build/check/test) 전체를 최종 확인한다.

> spec-c(TodoInput), spec-e(TodoList, spec-d 포함), spec-f(TodoFilter),
> spec-g(TodoFooter), spec-b(TodoStore) 에 의존한다.

## 구현 범위

### 변경 파일

- 수정: `src/App.svelte` (spec-a 의 최소 버전을 통합 버전으로 교체)
- 테스트: `src/App.test.ts`
- (선택) 수정: `src/app.css` (필요 시 레이아웃 여백 미세 조정만)

### 태스크 목록

- [ ] **태스크 1: 실패하는 통합 테스트 작성** (`src/App.test.ts`)

  ```ts
  import { render, screen, fireEvent } from '@testing-library/svelte';
  import { describe, it, expect, beforeEach } from 'vitest';
  import App from './App.svelte';

  beforeEach(() => localStorage.clear());

  async function addTodo(text: string) {
    const input = screen.getByLabelText('New todo');
    await fireEvent.input(input, { target: { value: text } });
    await fireEvent.click(screen.getByRole('button', { name: 'Add' }));
  }

  describe('App integration', () => {
    it('adds a todo and shows the remaining count', async () => {
      render(App);
      await addTodo('first task');
      expect(screen.getByText('first task')).toBeInTheDocument();
      expect(screen.getByText('1 item left')).toBeInTheDocument();
    });

    it('filters between active and completed todos', async () => {
      render(App);
      await addTodo('task one');
      await addTodo('task two');
      await fireEvent.click(
        screen.getByRole('checkbox', { name: 'Toggle task one' }),
      );

      await fireEvent.click(screen.getByRole('button', { name: 'Active' }));
      expect(screen.queryByText('task one')).toBeNull();
      expect(screen.getByText('task two')).toBeInTheDocument();

      await fireEvent.click(screen.getByRole('button', { name: 'Completed' }));
      expect(screen.getByText('task one')).toBeInTheDocument();
      expect(screen.queryByText('task two')).toBeNull();
    });

    it('clears completed todos', async () => {
      render(App);
      await addTodo('keep me');
      await addTodo('remove me');
      await fireEvent.click(
        screen.getByRole('checkbox', { name: 'Toggle remove me' }),
      );
      await fireEvent.click(
        screen.getByRole('button', { name: 'Clear completed' }),
      );
      expect(screen.getByText('keep me')).toBeInTheDocument();
      expect(screen.queryByText('remove me')).toBeNull();
    });

    it('removes an individual todo', async () => {
      render(App);
      await addTodo('temporary');
      await fireEvent.click(
        screen.getByRole('button', { name: 'Delete temporary' }),
      );
      expect(screen.queryByText('temporary')).toBeNull();
      expect(screen.getByText('No todos yet.')).toBeInTheDocument();
    });
  });
  ```

- [ ] **태스크 2: 테스트 실행 → 실패 확인**

  ```bash
  npm test -- App
  ```

  예상: FAIL — App 이 아직 컴포넌트를 배선하지 않아 입력/필터/푸터 요소를 찾지 못함
  (예: `Unable to find a label with the text of: New todo`).

- [ ] **태스크 3: App 통합 구현** (`src/App.svelte`)

  ```svelte
  <script lang="ts">
    import { TodoStore } from './lib/todos.svelte';
    import TodoInput from './lib/TodoInput.svelte';
    import TodoList from './lib/TodoList.svelte';
    import TodoFilter from './lib/TodoFilter.svelte';
    import TodoFooter from './lib/TodoFooter.svelte';

    const store = new TodoStore();
  </script>

  <main class="app">
    <h1>Todos</h1>
    <TodoInput onadd={(text) => store.add(text)} />
    <TodoList
      todos={store.filtered}
      ontoggle={(id) => store.toggle(id)}
      onremove={(id) => store.remove(id)}
    />
    <TodoFilter filter={store.filter} onfilter={(f) => store.setFilter(f)} />
    <TodoFooter remaining={store.remaining} onclear={() => store.clearCompleted()} />
  </main>

  <style>
    .app {
      max-width: 480px;
      margin: 2rem auto;
      padding: 0 1rem;
      font-family: system-ui, sans-serif;
    }
    h1 {
      font-size: 1.5rem;
    }
  </style>
  ```

- [ ] **태스크 4: 통합 테스트 통과 확인**

  ```bash
  npm test -- App
  ```

  예상: PASS (4 tests).

- [ ] **태스크 5: 전체 성공 기준 최종 검증**

  ```bash
  npm test          # 모든 테스트 파일 통과
  npm run check     # svelte-check 0 errors
  npm run build     # dist 생성, 오류 없음
  npm run dev       # 브라우저에서 7개 기능 + 새로고침 후 상태 유지 수동 확인 후 종료
  ```

  예상:
  - `npm test`: 전체 테스트(상태 모듈 + 5개 컴포넌트 + App 통합 + smoke) 통과.
  - `npm run check` / `npm run build`: 무오류.
  - 브라우저: 추가/토글/삭제/필터/Clear completed/남은 개수/localStorage 영속화 동작.

- [ ] **태스크 6: 커밋**

  ```bash
  git add src/App.svelte src/App.test.ts src/app.css
  git commit -m "feat: wire App integrating store and all todo components"
  ```

## 완료 기준

- [ ] `App.svelte` 가 `TodoStore` 를 생성해 4개 컴포넌트를 콜백으로 배선.
- [ ] 추가/토글/삭제/필터/Clear completed/남은 개수 흐름이 통합 테스트로 검증됨.
- [ ] 새로고침 후 목록/완료/필터 상태 유지(localStorage) — 수동 확인.
- [ ] `npm test`, `npm run check`, `npm run build` 모두 무오류.

## 제외 범위

- 반응형(모바일) 정교화, 다크모드 토글, 애니메이션.
- 라우팅, 백엔드 연동.
