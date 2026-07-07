import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect, vi } from 'vitest';
import TodoFilter from './TodoFilter.svelte';

describe('TodoFilter', () => {
  it('renders the three filter options', () => {
    render(TodoFilter, { props: { filter: 'all', onfilter: vi.fn() } });
    expect(screen.getByRole('button', { name: 'All' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Active' })).toBeInTheDocument();
    expect(
      screen.getByRole('button', { name: 'Completed' }),
    ).toBeInTheDocument();
  });

  it('marks the active filter as pressed', () => {
    render(TodoFilter, { props: { filter: 'active', onfilter: vi.fn() } });
    expect(screen.getByRole('button', { name: 'Active' })).toHaveAttribute(
      'aria-pressed',
      'true',
    );
    expect(screen.getByRole('button', { name: 'All' })).toHaveAttribute(
      'aria-pressed',
      'false',
    );
  });

  it('calls onfilter with the chosen value', async () => {
    const onfilter = vi.fn();
    render(TodoFilter, { props: { filter: 'all', onfilter } });
    await fireEvent.click(screen.getByRole('button', { name: 'Completed' }));
    expect(onfilter).toHaveBeenCalledWith('completed');
  });
});
