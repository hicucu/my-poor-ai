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
