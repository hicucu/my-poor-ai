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
