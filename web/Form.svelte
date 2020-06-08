<script>
  import { onMount, createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  export let formFields = {};

  let formData = {};

  export function getFormData() {
    return formData;
  }
</script>

<form on:submit|preventDefault={dispatch('submit', formData)}>
  {#each Object.entries(formFields) as [formKey, formItem]}
    <div class="form-group">
      <label for={formKey}>{formItem.label}</label>
      {#if formItem.type === 'email'}
        <input
          type="email"
          class="white-text form-control"
          id={formKey}
          bind:value={formData[formKey]} />
      {:else if formItem.type === 'password'}
        <input
          type="password"
          class="white-text form-control"
          id={formKey}
          bind:value={formData[formKey]} />
      {:else if formItem.type === 'text'}
        <input
          type="text"
          class="white-text form-control"
          id={formKey}
          bind:value={formData[formKey]} />
      {/if}
    </div>
  {/each}
  <slot />
</form>
