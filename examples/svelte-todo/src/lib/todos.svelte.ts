import type { Todo, Filter } from './types';

const STORAGE_KEY = 'svelte-todos';
const FILTERS: readonly Filter[] = ['all', 'active', 'completed'];

interface Persisted {
  todos: Todo[];
  filter: Filter;
}

function isTodo(item: unknown): item is Todo {
  return (
    typeof item === 'object' &&
    item !== null &&
    typeof (item as Record<string, unknown>).id === 'string' &&
    typeof (item as Record<string, unknown>).text === 'string' &&
    typeof (item as Record<string, unknown>).completed === 'boolean'
  );
}

function load(key: string): Persisted {
  try {
    const raw = localStorage.getItem(key);
    if (!raw) return { todos: [], filter: 'all' };
    const parsed = JSON.parse(raw);
    const todos = Array.isArray(parsed?.todos) ? parsed.todos.filter(isTodo) : [];
    const filter: Filter = FILTERS.includes(parsed?.filter)
      ? (parsed.filter as Filter)
      : 'all';
    return { todos, filter };
  } catch (err) {
    if (import.meta.env.DEV) {
      console.warn('failed to load persisted todos', err);
    }
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
