---
complexity: simple
depends-on: spec-b
estimated-tasks: 5
---

# Spec C: TodoInput 컴포넌트 (새 todo 입력/추가)

## 목표

새 todo 텍스트를 입력받아 Enter 키 또는 Add 버튼으로 추가하는 컴포넌트를
구현한다. 빈/공백 문자열은 무시하고, 추가 후 입력창을 비운다. 상위와의 통신은
콜백 prop `onadd(text)` 로 단방향 처리한다.

## 구현 범위

### 변경 파일

- 생성: `src/lib/TodoInput.svelte`
- 테스트: `src/lib/TodoInput.test.ts`

### 태스크 목록

- [ ] **태스크 1: 실패하는 테스트 작성** (`src/lib/TodoInput.test.ts`)

  ```ts
  import { render, screen, fireEvent } from '@testing-library/svelte';
  import { describe, it, expect, vi } from 'vitest';
  import TodoInput from './TodoInput.svelte';

  describe('TodoInput', () => {
    it('calls onadd with trimmed text and clears the field on Add click', async () => {
      const onadd = vi.fn();
      render(TodoInput, { props: { onadd } });
      const input = screen.getByLabelText('New todo') as HTMLInputElement;
      await fireEvent.input(input, { target: { value: '  buy milk  ' } });
      await fireEvent.click(screen.getByRole('button', { name: 'Add' }));
      expect(onadd).toHaveBeenCalledWith('buy milk');
      expect(input.value).toBe('');
    });

    it('ignores empty / whitespace-only input', async () => {
      const onadd = vi.fn();
      render(TodoInput, { props: { onadd } });
      const input = screen.getByLabelText('New todo');
      await fireEvent.input(input, { target: { value: '   ' } });
      await fireEvent.click(screen.getByRole('button', { name: 'Add' }));
      expect(onadd).not.toHaveBeenCalled();
    });

    it('adds on Enter key', async () => {
      const onadd = vi.fn();
      render(TodoInput, { props: { onadd } });
      const input = screen.getByLabelText('New todo');
      await fireEvent.input(input, { target: { value: 'task' } });
      await fireEvent.keyDown(input, { key: 'Enter' });
      expect(onadd).toHaveBeenCalledWith('task');
    });
  });
  ```

- [ ] **태스크 2: 테스트 실행 → 실패 확인**

  ```bash
  npm test -- TodoInput
  ```

  예상: FAIL — `Failed to resolve import "./TodoInput.svelte"`.

- [ ] **태스크 3: 컴포넌트 구현** (`src/lib/TodoInput.svelte`)

  ```svelte
  <script lang="ts">
    interface Props {
      onadd: (text: string) => void;
    }

    let { onadd }: Props = $props();
    let value = $state('');

    function submit() {
      const text = value.trim();
      if (!text) return;
      onadd(text);
      value = '';
    }

    function handleKeydown(event: KeyboardEvent) {
      if (event.key === 'Enter') submit();
    }
  </script>

  <div class="todo-input">
    <input
      type="text"
      placeholder="What needs to be done?"
      aria-label="New todo"
      bind:value
      onkeydown={handleKeydown}
    />
    <button type="button" onclick={submit}>Add</button>
  </div>
  ```

- [ ] **태스크 4: 테스트 실행 → 통과 확인**

  ```bash
  npm test -- TodoInput
  npm run check
  ```

  예상: PASS (3 tests), check 0 errors.

- [ ] **태스크 5: 커밋**

  ```bash
  git add src/lib/TodoInput.svelte src/lib/TodoInput.test.ts
  git commit -m "feat: TodoInput component with add-on-enter and blank guard"
  ```

## 완료 기준

- [ ] Add 버튼/Enter 로 trim 된 텍스트가 `onadd` 로 전달됨.
- [ ] 빈/공백 입력은 무시됨.
- [ ] 추가 후 입력창이 비워짐.
- [ ] `npm test`, `npm run check` 무오류.

## 제외 범위

- 인라인 편집, 자동완성, 최대 길이 제한.
