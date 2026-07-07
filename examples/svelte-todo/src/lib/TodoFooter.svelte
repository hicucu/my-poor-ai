<script lang="ts">
  interface Props {
    remaining: number;
    onclear: () => void;
  }

  let { remaining, onclear }: Props = $props();
</script>

<footer class="todo-footer">
  <span class="count">
    {remaining} {remaining === 1 ? 'item' : 'items'} left
  </span>
  <button type="button" onclick={onclear}>Clear completed</button>
</footer>

<style>
  .todo-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 0.75rem;
    font-size: 0.9rem;
    color: #555;
  }
</style>
