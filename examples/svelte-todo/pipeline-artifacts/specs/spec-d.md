---
complexity: simple
depends-on: spec-b
estimated-tasks: 5
---

# Spec D: TodoItem 컴포넌트 (개별 항목 토글/삭제)

## 목표

단일 todo 를 렌더링하는 컴포넌트를 구현한다. 체크박스로 완료 토글, 텍스트 표시,
삭제 버튼을 제공한다. 상위와의 통신은 콜백 prop `ontoggle(id)`, `onremove(id)` 로
단방향 처리한다.

## 구현 범위

### 변경 파일

- 생성: `src/lib/TodoItem.svelte`
- 테스트: `src/lib/TodoItem.test.ts`

### 태스크 목록

- [ ] **태스크 1: 실패하는 테스트 작성** (`src/lib/TodoItem.test.ts`)

  ```ts
  import { render, screen, fireEvent } from '@testing-library/svelte';
  import { describe, it, expect, vi } from 'vitest';
  import TodoItem from './TodoItem.svelte';
  import type { Todo } from './types';

  const todo: Todo = { id: '1', text: 'read book', completed: false };

  describe('TodoItem', () => {
    it('renders the text with an unchecked checkbox', () => {
      render(TodoItem, {
        props: { todo, ontoggle: vi.fn(), onremove: vi.fn() },
      });
      expect(screen.getByText('read book')).toBeInTheDocument();
      expect(screen.getByRole('checkbox')).not.toBeChecked();
    });

    it('reflects completed state as checked', () => {
      render(TodoItem, {
        props: {
          todo: { ...todo, completed: true },
          ontoggle: vi.fn(),
          onremove: vi.fn(),
        },
      });
      expect(screen.getByRole('checkbox')).toBeChecked();
    });

    it('calls ontoggle with the id on checkbox change', async () => {
      const ontoggle = vi.fn();
      render(TodoItem, { props: { todo, ontoggle, onremove: vi.fn() } });
      await fireEvent.click(screen.getByRole('checkbox'));
      expect(ontoggle).toHaveBeenCalledWith('1');
    });

    it('calls onremove with the id on delete button click', async () => {
      const onremove = vi.fn();
      render(TodoItem, { props: { todo, ontoggle: vi.fn(), onremove } });
      await fireEvent.click(
        screen.getByRole('button', { name: 'Delete read book' }),
      );
      expect(onremove).toHaveBeenCalledWith('1');
    });
  });
  ```

- [ ] **태스크 2: 테스트 실행 → 실패 확인**

  ```bash
  npm test -- TodoItem
  ```

  예상: FAIL — `Failed to resolve import "./TodoItem.svelte"`.

- [ ] **태스크 3: 컴포넌트 구현** (`src/lib/TodoItem.svelte`)

  ```svelte
  <script lang="ts">
    import type { Todo } from './types';

    interface Props {
      todo: Todo;
      ontoggle: (id: string) => void;
      onremove: (id: string) => void;
    }

    let { todo, ontoggle, onremove }: Props = $props();
  </script>

  <li class="todo-item" class:completed={todo.completed}>
    <input
      type="checkbox"
      checked={todo.completed}
      aria-label={`Toggle ${todo.text}`}
      onchange={() => ontoggle(todo.id)}
    />
    <span class="text">{todo.text}</span>
    <button
      type="button"
      class="delete"
      aria-label={`Delete ${todo.text}`}
      onclick={() => onremove(todo.id)}
    >
      ×
    </button>
  </li>

  <style>
    .todo-item {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.4rem 0;
    }
    .todo-item.completed .text {
      text-decoration: line-through;
      opacity: 0.6;
    }
    .text {
      flex: 1;
    }
    .delete {
      border: none;
      background: transparent;
      cursor: pointer;
      font-size: 1.1rem;
      line-height: 1;
    }
  </style>
  ```

- [ ] **태스크 4: 테스트 실행 → 통과 확인**

  ```bash
  npm test -- TodoItem
  npm run check
  ```

  예상: PASS (4 tests), check 0 errors.

- [ ] **태스크 5: 커밋**

  ```bash
  git add src/lib/TodoItem.svelte src/lib/TodoItem.test.ts
  git commit -m "feat: TodoItem component with toggle and delete"
  ```

## 완료 기준

- [ ] 텍스트/체크박스 렌더, `completed` 시 checked 및 취소선 스타일.
- [ ] 체크박스 변경 → `ontoggle(id)`, 삭제 버튼 → `onremove(id)`.
- [ ] `npm test`, `npm run check` 무오류.

## 제외 범위

- 인라인 편집, 애니메이션.
