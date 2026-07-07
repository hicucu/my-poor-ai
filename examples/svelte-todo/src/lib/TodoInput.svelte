<script lang="ts">
  interface Props {
    onadd: (text: string) => void;
  }

  let { onadd }: Props = $props();
  let value = $state('');

  function submit() {
    const text = value.trim();
    if (!text) return;
    onadd(text);
    value = '';
  }

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Enter') submit();
  }
</script>

<div class="todo-input">
  <input
    type="text"
    placeholder="What needs to be done?"
    aria-label="New todo"
    bind:value
    onkeydown={handleKeydown}
  />
  <button type="button" onclick={submit}>Add</button>
</div>
