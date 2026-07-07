<script lang="ts">
  import type { Todo } from './types';

  interface Props {
    todo: Todo;
    ontoggle: (id: string) => void;
    onremove: (id: string) => void;
  }

  let { todo, ontoggle, onremove }: Props = $props();
</script>

<li class="todo-item" class:completed={todo.completed}>
  <input
    type="checkbox"
    checked={todo.completed}
    aria-label={`Toggle ${todo.text}`}
    onchange={() => ontoggle(todo.id)}
  />
  <span class="text">{todo.text}</span>
  <button
    type="button"
    class="delete"
    aria-label={`Delete ${todo.text}`}
    onclick={() => onremove(todo.id)}
  >
    ×
  </button>
</li>

<style>
  .todo-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.4rem 0;
  }
  .todo-item.completed .text {
    text-decoration: line-through;
    opacity: 0.6;
  }
  .text {
    flex: 1;
  }
  .delete {
    border: none;
    background: transparent;
    cursor: pointer;
    font-size: 1.1rem;
    line-height: 1;
  }
</style>
