---
complexity: simple
depends-on: spec-b
estimated-tasks: 5
---

# Spec G: TodoFooter 컴포넌트 (남은 개수 + Clear completed)

## 목표

미완료(남은) 항목 개수를 표시하고(단수/복수 라벨 처리), Clear completed 버튼으로
완료 항목 일괄 삭제를 요청하는 컴포넌트를 구현한다. 상위와의 통신은 prop
`remaining`(수)과 콜백 `onclear()` 로 처리한다.

## 구현 범위

### 변경 파일

- 생성: `src/lib/TodoFooter.svelte`
- 테스트: `src/lib/TodoFooter.test.ts`

### 태스크 목록

- [ ] **태스크 1: 실패하는 테스트 작성** (`src/lib/TodoFooter.test.ts`)

  ```ts
  import { render, screen, fireEvent } from '@testing-library/svelte';
  import { describe, it, expect, vi } from 'vitest';
  import TodoFooter from './TodoFooter.svelte';

  describe('TodoFooter', () => {
    it('shows singular label for exactly one remaining item', () => {
      render(TodoFooter, { props: { remaining: 1, onclear: vi.fn() } });
      expect(screen.getByText('1 item left')).toBeInTheDocument();
    });

    it('shows plural label for zero or many remaining items', () => {
      render(TodoFooter, { props: { remaining: 3, onclear: vi.fn() } });
      expect(screen.getByText('3 items left')).toBeInTheDocument();
    });

    it('calls onclear when Clear completed is clicked', async () => {
      const onclear = vi.fn();
      render(TodoFooter, { props: { remaining: 0, onclear } });
      await fireEvent.click(
        screen.getByRole('button', { name: 'Clear completed' }),
      );
      expect(onclear).toHaveBeenCalledTimes(1);
    });
  });
  ```

- [ ] **태스크 2: 테스트 실행 → 실패 확인**

  ```bash
  npm test -- TodoFooter
  ```

  예상: FAIL — `Failed to resolve import "./TodoFooter.svelte"`.

- [ ] **태스크 3: 컴포넌트 구현** (`src/lib/TodoFooter.svelte`)

  ```svelte
  <script lang="ts">
    interface Props {
      remaining: number;
      onclear: () => void;
    }

    let { remaining, onclear }: Props = $props();
  </script>

  <footer class="todo-footer">
    <span class="count">
      {remaining} {remaining === 1 ? 'item' : 'items'} left
    </span>
    <button type="button" onclick={onclear}>Clear completed</button>
  </footer>

  <style>
    .todo-footer {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-top: 0.75rem;
      font-size: 0.9rem;
      color: #555;
    }
  </style>
  ```

  참고: `{remaining} {...} left` 사이 공백이 텍스트로 렌더되어 `1 item left`,
  `3 items left` 문자열이 그대로 매칭된다.

- [ ] **태스크 4: 테스트 실행 → 통과 확인**

  ```bash
  npm test -- TodoFooter
  npm run check
  ```

  예상: PASS (3 tests), check 0 errors.

- [ ] **태스크 5: 커밋**

  ```bash
  git add src/lib/TodoFooter.svelte src/lib/TodoFooter.test.ts
  git commit -m "feat: TodoFooter with remaining count and clear-completed"
  ```

## 완료 기준

- [ ] `remaining === 1` 이면 "1 item left", 그 외 "N items left".
- [ ] Clear completed 클릭 시 `onclear()` 호출.
- [ ] `npm test`, `npm run check` 무오류.

## 제외 범위

- 완료 개수 표시, 진행률 바.
