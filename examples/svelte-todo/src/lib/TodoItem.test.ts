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
