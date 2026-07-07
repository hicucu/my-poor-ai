import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect, beforeEach } from 'vitest';
import App from './App.svelte';

beforeEach(() => localStorage.clear());

async function addTodo(text: string) {
  const input = screen.getByLabelText('New todo');
  await fireEvent.input(input, { target: { value: text } });
  await fireEvent.click(screen.getByRole('button', { name: 'Add' }));
}

describe('App integration', () => {
  it('adds a todo and shows the remaining count', async () => {
    render(App);
    await addTodo('first task');
    expect(screen.getByText('first task')).toBeInTheDocument();
    expect(screen.getByText('1 item left')).toBeInTheDocument();
  });

  it('filters between active and completed todos', async () => {
    render(App);
    await addTodo('task one');
    await addTodo('task two');
    await fireEvent.click(
      screen.getByRole('checkbox', { name: 'Toggle task one' }),
    );

    await fireEvent.click(screen.getByRole('button', { name: 'Active' }));
    expect(screen.queryByText('task one')).toBeNull();
    expect(screen.getByText('task two')).toBeInTheDocument();

    await fireEvent.click(screen.getByRole('button', { name: 'Completed' }));
    expect(screen.getByText('task one')).toBeInTheDocument();
    expect(screen.queryByText('task two')).toBeNull();
  });

  it('clears completed todos', async () => {
    render(App);
    await addTodo('keep me');
    await addTodo('remove me');
    await fireEvent.click(
      screen.getByRole('checkbox', { name: 'Toggle remove me' }),
    );
    await fireEvent.click(
      screen.getByRole('button', { name: 'Clear completed' }),
    );
    expect(screen.getByText('keep me')).toBeInTheDocument();
    expect(screen.queryByText('remove me')).toBeNull();
  });

  it('removes an individual todo', async () => {
    render(App);
    await addTodo('temporary');
    await fireEvent.click(
      screen.getByRole('button', { name: 'Delete temporary' }),
    );
    expect(screen.queryByText('temporary')).toBeNull();
    expect(screen.getByText('No todos yet.')).toBeInTheDocument();
  });
});
