---
complexity: simple
depends-on: spec-d
estimated-tasks: 5
---

# Spec E: TodoList 컴포넌트 (필터된 목록 렌더)

## 목표

전달받은 todo 배열을 `TodoItem` 반복으로 렌더링하는 컨테이너 컴포넌트를 구현한다.
목록이 비었을 때는 안내 메시지를 표시한다. 토글/삭제 콜백은 그대로 하위로 전달한다.

> spec-d(TodoItem)에 의존한다.

## 구현 범위

### 변경 파일

- 생성: `src/lib/TodoList.svelte`
- 테스트: `src/lib/TodoList.test.ts`

### 태스크 목록

- [ ] **태스크 1: 실패하는 테스트 작성** (`src/lib/TodoList.test.ts`)

  ```ts
  import { render, screen } from '@testing-library/svelte';
  import { describe, it, expect, vi } from 'vitest';
  import TodoList from './TodoList.svelte';
  import type { Todo } from './types';

  const todos: Todo[] = [
    { id: '1', text: 'alpha', completed: false },
    { id: '2', text: 'beta', completed: true },
  ];

  describe('TodoList', () => {
    it('renders one list item per todo', () => {
      render(TodoList, {
        props: { todos, ontoggle: vi.fn(), onremove: vi.fn() },
      });
      expect(screen.getAllByRole('listitem')).toHaveLength(2);
      expect(screen.getByText('alpha')).toBeInTheDocument();
      expect(screen.getByText('beta')).toBeInTheDocument();
    });

    it('shows an empty message when there are no todos', () => {
      render(TodoList, {
        props: { todos: [], ontoggle: vi.fn(), onremove: vi.fn() },
      });
      expect(screen.getByText('No todos yet.')).toBeInTheDocument();
      expect(screen.queryByRole('listitem')).toBeNull();
    });
  });
  ```

- [ ] **태스크 2: 테스트 실행 → 실패 확인**

  ```bash
  npm test -- TodoList
  ```

  예상: FAIL — `Failed to resolve import "./TodoList.svelte"`.

- [ ] **태스크 3: 컴포넌트 구현** (`src/lib/TodoList.svelte`)

  ```svelte
  <script lang="ts">
    import type { Todo } from './types';
    import TodoItem from './TodoItem.svelte';

    interface Props {
      todos: Todo[];
      ontoggle: (id: string) => void;
      onremove: (id: string) => void;
    }

    let { todos, ontoggle, onremove }: Props = $props();
  </script>

  {#if todos.length === 0}
    <p class="empty">No todos yet.</p>
  {:else}
    <ul class="todo-list">
      {#each todos as todo (todo.id)}
        <TodoItem {todo} {ontoggle} {onremove} />
      {/each}
    </ul>
  {/if}

  <style>
    .todo-list {
      list-style: none;
      margin: 0;
      padding: 0;
    }
    .empty {
      color: #888;
      text-align: center;
    }
  </style>
  ```

- [ ] **태스크 4: 테스트 실행 → 통과 확인**

  ```bash
  npm test -- TodoList
  npm run check
  ```

  예상: PASS (2 tests), check 0 errors.

- [ ] **태스크 5: 커밋**

  ```bash
  git add src/lib/TodoList.svelte src/lib/TodoList.test.ts
  git commit -m "feat: TodoList component rendering items with empty state"
  ```

## 완료 기준

- [ ] 각 todo 당 `TodoItem` 1개 렌더 (key = id).
- [ ] 빈 목록일 때 "No todos yet." 표시.
- [ ] 토글/삭제 콜백이 하위로 전달됨.
- [ ] `npm test`, `npm run check` 무오류.

## 제외 범위

- 정렬, 페이지네이션, 가상 스크롤.
