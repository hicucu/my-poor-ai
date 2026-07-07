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
