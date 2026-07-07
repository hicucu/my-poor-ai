<script lang="ts">
  import type { Todo } from './types';
  import TodoItem from './TodoItem.svelte';

  interface Props {
    todos: Todo[];
    ontoggle: (id: string) => void;
    onremove: (id: string) => void;
  }

  let { todos, ontoggle, onremove }: Props = $props();
</script>

{#if todos.length === 0}
  <p class="empty">No todos yet.</p>
{:else}
  <ul class="todo-list">
    {#each todos as todo (todo.id)}
      <TodoItem {todo} {ontoggle} {onremove} />
    {/each}
  </ul>
{/if}

<style>
  .todo-list {
    list-style: none;
    margin: 0;
    padding: 0;
  }
  .empty {
    color: #888;
    text-align: center;
  }
</style>
