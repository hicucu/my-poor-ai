---
complexity: simple
depends-on: spec-b
estimated-tasks: 5
---

# Spec F: TodoFilter 컴포넌트 (All / Active / Completed)

## 목표

All / Active / Completed 세 개의 필터 버튼을 렌더링하고, 현재 선택된 필터를
시각/접근성(aria-pressed) 상으로 표시한다. 클릭 시 콜백 prop `onfilter(filter)`
로 선택 값을 전달한다.

## 구현 범위

### 변경 파일

- 생성: `src/lib/TodoFilter.svelte`
- 테스트: `src/lib/TodoFilter.test.ts`

### 태스크 목록

- [ ] **태스크 1: 실패하는 테스트 작성** (`src/lib/TodoFilter.test.ts`)

  ```ts
  import { render, screen, fireEvent } from '@testing-library/svelte';
  import { describe, it, expect, vi } from 'vitest';
  import TodoFilter from './TodoFilter.svelte';

  describe('TodoFilter', () => {
    it('renders the three filter options', () => {
      render(TodoFilter, { props: { filter: 'all', onfilter: vi.fn() } });
      expect(screen.getByRole('button', { name: 'All' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Active' })).toBeInTheDocument();
      expect(
        screen.getByRole('button', { name: 'Completed' }),
      ).toBeInTheDocument();
    });

    it('marks the active filter as pressed', () => {
      render(TodoFilter, { props: { filter: 'active', onfilter: vi.fn() } });
      expect(screen.getByRole('button', { name: 'Active' })).toHaveAttribute(
        'aria-pressed',
        'true',
      );
      expect(screen.getByRole('button', { name: 'All' })).toHaveAttribute(
        'aria-pressed',
        'false',
      );
    });

    it('calls onfilter with the chosen value', async () => {
      const onfilter = vi.fn();
      render(TodoFilter, { props: { filter: 'all', onfilter } });
      await fireEvent.click(screen.getByRole('button', { name: 'Completed' }));
      expect(onfilter).toHaveBeenCalledWith('completed');
    });
  });
  ```

- [ ] **태스크 2: 테스트 실행 → 실패 확인**

  ```bash
  npm test -- TodoFilter
  ```

  예상: FAIL — `Failed to resolve import "./TodoFilter.svelte"`.

- [ ] **태스크 3: 컴포넌트 구현** (`src/lib/TodoFilter.svelte`)

  ```svelte
  <script lang="ts">
    import type { Filter } from './types';

    interface Props {
      filter: Filter;
      onfilter: (filter: Filter) => void;
    }

    let { filter, onfilter }: Props = $props();

    const options: { value: Filter; label: string }[] = [
      { value: 'all', label: 'All' },
      { value: 'active', label: 'Active' },
      { value: 'completed', label: 'Completed' },
    ];
  </script>

  <div class="todo-filter" role="group" aria-label="Filter todos">
    {#each options as option (option.value)}
      <button
        type="button"
        class:selected={filter === option.value}
        aria-pressed={filter === option.value}
        onclick={() => onfilter(option.value)}
      >
        {option.label}
      </button>
    {/each}
  </div>

  <style>
    .todo-filter {
      display: flex;
      gap: 0.25rem;
    }
    .todo-filter button.selected {
      font-weight: 700;
      text-decoration: underline;
    }
  </style>
  ```

- [ ] **태스크 4: 테스트 실행 → 통과 확인**

  ```bash
  npm test -- TodoFilter
  npm run check
  ```

  예상: PASS (3 tests), check 0 errors.

- [ ] **태스크 5: 커밋**

  ```bash
  git add src/lib/TodoFilter.svelte src/lib/TodoFilter.test.ts
  git commit -m "feat: TodoFilter component with All/Active/Completed"
  ```

## 완료 기준

- [ ] 세 버튼 렌더, 선택된 필터가 `aria-pressed="true"` 및 강조 표시.
- [ ] 클릭 시 `onfilter(value)` 호출.
- [ ] `npm test`, `npm run check` 무오류.

## 제외 범위

- 커스텀 필터, URL 동기화.
