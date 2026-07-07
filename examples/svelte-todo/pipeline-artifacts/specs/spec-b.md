---
complexity: complex
depends-on: spec-a
estimated-tasks: 6
---

# Spec B: 타입 + 반응형 상태 모듈 (localStorage 영속화)

## 목표

도메인 타입(`Todo`, `Filter`)과 Svelte 5 runes 기반 상태 모듈 `TodoStore` 를
구현한다. `add / toggle / remove / clearCompleted / setFilter` 조작, `$derived`
로 `filtered`(현재 필터 목록)와 `remaining`(미완료 개수)를 계산하고, `$effect`
로 todos/filter 변경 시 localStorage 에 저장, 생성 시 localStorage 에서 복원한다
(JSON 파싱 실패 시 빈 상태 폴백).

## 구현 범위

### 변경 파일

- 생성: `src/lib/types.ts`
- 생성: `src/lib/todos.svelte.ts`
- 테스트: `src/lib/todos.svelte.test.ts`

### 태스크 목록

- [ ] **태스크 1: 실패하는 테스트 작성** (`src/lib/todos.svelte.test.ts`)

  ```ts
  import { describe, it, expect, beforeEach } from 'vitest';
  import { flushSync } from 'svelte';
  import { TodoStore } from './todos.svelte';

  beforeEach(() => localStorage.clear());

  describe('TodoStore', () => {
    it('starts empty with the "all" filter', () => {
      const store = new TodoStore('test-todos');
      expect(store.todos).toEqual([]);
      expect(store.filter).toBe('all');
      expect(store.remaining).toBe(0);
    });

    it('adds a trimmed todo and ignores blank input', () => {
      const store = new TodoStore('test-todos');
      store.add('  write tests  ');
      store.add('   ');
      expect(store.todos).toHaveLength(1);
      expect(store.todos[0].text).toBe('write tests');
      expect(store.todos[0].completed).toBe(false);
      expect(typeof store.todos[0].id).toBe('string');
    });

    it('toggles completion', () => {
      const store = new TodoStore('test-todos');
      store.add('a');
      const id = store.todos[0].id;
      store.toggle(id);
      expect(store.todos[0].completed).toBe(true);
      store.toggle(id);
      expect(store.todos[0].completed).toBe(false);
    });

    it('removes a todo by id', () => {
      const store = new TodoStore('test-todos');
      store.add('a');
      store.add('b');
      const id = store.todos[0].id;
      store.remove(id);
      expect(store.todos).toHaveLength(1);
      expect(store.todos[0].text).toBe('b');
    });

    it('computes remaining as the count of active todos', () => {
      const store = new TodoStore('test-todos');
      store.add('a');
      store.add('b');
      store.toggle(store.todos[0].id);
      expect(store.remaining).toBe(1);
    });

    it('filters todos by status', () => {
      const store = new TodoStore('test-todos');
      store.add('a');
      store.add('b');
      store.toggle(store.todos[0].id);
      store.setFilter('active');
      expect(store.filtered.map((t) => t.text)).toEqual(['b']);
      store.setFilter('completed');
      expect(store.filtered.map((t) => t.text)).toEqual(['a']);
      store.setFilter('all');
      expect(store.filtered).toHaveLength(2);
    });

    it('clears completed todos', () => {
      const store = new TodoStore('test-todos');
      store.add('a');
      store.add('b');
      store.toggle(store.todos[0].id);
      store.clearCompleted();
      expect(store.todos.map((t) => t.text)).toEqual(['b']);
    });

    it('persists todos and filter to localStorage', () => {
      const store = new TodoStore('test-todos');
      store.add('persist me');
      store.setFilter('active');
      flushSync();
      const raw = localStorage.getItem('test-todos');
      expect(raw).not.toBeNull();
      const saved = JSON.parse(raw!);
      expect(saved.todos).toHaveLength(1);
      expect(saved.filter).toBe('active');
    });

    it('restores state from localStorage on construction', () => {
      localStorage.setItem(
        'test-todos',
        JSON.stringify({
          todos: [{ id: '1', text: 'restored', completed: true }],
          filter: 'completed',
        }),
      );
      const store = new TodoStore('test-todos');
      expect(store.todos).toEqual([
        { id: '1', text: 'restored', completed: true },
      ]);
      expect(store.filter).toBe('completed');
    });

    it('falls back to empty state on corrupt localStorage', () => {
      localStorage.setItem('test-todos', '{ not valid json');
      const store = new TodoStore('test-todos');
      expect(store.todos).toEqual([]);
      expect(store.filter).toBe('all');
    });
  });
  ```

- [ ] **태스크 2: 테스트 실행 → 실패 확인**

  ```bash
  npm test
  ```

  예상: FAIL — `Failed to resolve import "./todos.svelte"` (모듈 미존재).

- [ ] **태스크 3: 타입 정의** (`src/lib/types.ts`)

  ```ts
  export interface Todo {
    id: string;
    text: string;
    completed: boolean;
  }

  export type Filter = 'all' | 'active' | 'completed';
  ```

- [ ] **태스크 4: 상태 모듈 구현** (`src/lib/todos.svelte.ts`)

  ```ts
  import type { Todo, Filter } from './types';

  const STORAGE_KEY = 'svelte-todos';
  const FILTERS: readonly Filter[] = ['all', 'active', 'completed'];

  interface Persisted {
    todos: Todo[];
    filter: Filter;
  }

  function load(key: string): Persisted {
    try {
      const raw = localStorage.getItem(key);
      if (!raw) return { todos: [], filter: 'all' };
      const parsed = JSON.parse(raw);
      const todos = Array.isArray(parsed?.todos)
        ? (parsed.todos as Todo[])
        : [];
      const filter: Filter = FILTERS.includes(parsed?.filter)
        ? (parsed.filter as Filter)
        : 'all';
      return { todos, filter };
    } catch {
      return { todos: [], filter: 'all' };
    }
  }

  export class TodoStore {
    todos = $state<Todo[]>([]);
    filter = $state<Filter>('all');

    filtered = $derived.by(() => {
      switch (this.filter) {
        case 'active':
          return this.todos.filter((t) => !t.completed);
        case 'completed':
          return this.todos.filter((t) => t.completed);
        default:
          return this.todos;
      }
    });

    remaining = $derived(this.todos.filter((t) => !t.completed).length);

    #key: string;
    #stop: () => void;

    constructor(key: string = STORAGE_KEY) {
      this.#key = key;
      const initial = load(key);
      this.todos = initial.todos;
      this.filter = initial.filter;
      this.#stop = $effect.root(() => {
        $effect(() => {
          localStorage.setItem(
            this.#key,
            JSON.stringify({ todos: this.todos, filter: this.filter }),
          );
        });
      });
    }

    add(text: string): void {
      const trimmed = text.trim();
      if (!trimmed) return;
      this.todos = [
        ...this.todos,
        { id: crypto.randomUUID(), text: trimmed, completed: false },
      ];
    }

    toggle(id: string): void {
      this.todos = this.todos.map((t) =>
        t.id === id ? { ...t, completed: !t.completed } : t,
      );
    }

    remove(id: string): void {
      this.todos = this.todos.filter((t) => t.id !== id);
    }

    clearCompleted(): void {
      this.todos = this.todos.filter((t) => !t.completed);
    }

    setFilter(filter: Filter): void {
      this.filter = filter;
    }

    destroy(): void {
      this.#stop();
    }
  }
  ```

  참고: `$effect` 는 컴포넌트 밖(모듈/클래스)에서 실행되므로 `$effect.root` 로
  감싸 반응 컨텍스트를 만든다. 테스트는 `flushSync()` 로 이펙트를 강제 flush 해
  localStorage 저장을 검증한다.

- [ ] **태스크 5: 테스트 실행 → 통과 확인**

  ```bash
  npm test
  npm run check
  ```

  예상: PASS (todos.svelte.test.ts 전체 통과), check 0 errors.

- [ ] **태스크 6: 커밋**

  ```bash
  git add src/lib/types.ts src/lib/todos.svelte.ts src/lib/todos.svelte.test.ts
  git commit -m "feat: reactive todo store with localStorage persistence"
  ```

## 완료 기준

- [ ] `Todo`, `Filter` 타입 정의됨.
- [ ] `add`(빈 문자열 무시/trim), `toggle`, `remove`, `clearCompleted`, `setFilter` 동작.
- [ ] `filtered`, `remaining` 파생값이 정확.
- [ ] todos/filter 변경 시 localStorage 저장, 생성 시 복원, 파싱 실패 시 폴백.
- [ ] `npm test`, `npm run check` 무오류.

## 제외 범위

- UI 컴포넌트 (spec-c 이후).
- 싱글턴 export — 상태 인스턴스는 `App.svelte`(spec-h)에서 `new TodoStore()` 로 생성.
