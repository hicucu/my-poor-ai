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
