<script lang="ts">
  import { onDestroy } from 'svelte';
  import { TodoStore } from './lib/todos.svelte';
  import TodoInput from './lib/TodoInput.svelte';
  import TodoList from './lib/TodoList.svelte';
  import TodoFilter from './lib/TodoFilter.svelte';
  import TodoFooter from './lib/TodoFooter.svelte';

  const store = new TodoStore();

  onDestroy(() => store.destroy());
</script>

<main class="app">
  <h1>Todos</h1>
  <TodoInput onadd={(text) => store.add(text)} />
  <TodoList
    todos={store.filtered}
    ontoggle={(id) => store.toggle(id)}
    onremove={(id) => store.remove(id)}
  />
  <TodoFilter filter={store.filter} onfilter={(f) => store.setFilter(f)} />
  <TodoFooter remaining={store.remaining} onclear={() => store.clearCompleted()} />
</main>

<style>
  .app {
    max-width: 480px;
    margin: 2rem auto;
    padding: 0 1rem;
    font-family: system-ui, sans-serif;
  }
  h1 {
    font-size: 1.5rem;
  }
</style>
