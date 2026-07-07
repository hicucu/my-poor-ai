<script lang="ts">
  import type { Filter } from './types';

  interface Props {
    filter: Filter;
    onfilter: (filter: Filter) => void;
  }

  let { filter, onfilter }: Props = $props();

  const options: { value: Filter; label: string }[] = [
    { value: 'all', label: 'All' },
    { value: 'active', label: 'Active' },
    { value: 'completed', label: 'Completed' },
  ];
</script>

<div class="todo-filter" role="group" aria-label="Filter todos">
  {#each options as option (option.value)}
    <button
      type="button"
      class:selected={filter === option.value}
      aria-pressed={filter === option.value}
      onclick={() => onfilter(option.value)}
    >
      {option.label}
    </button>
  {/each}
</div>

<style>
  .todo-filter {
    display: flex;
    gap: 0.25rem;
  }
  .todo-filter button.selected {
    font-weight: 700;
    text-decoration: underline;
  }
</style>
