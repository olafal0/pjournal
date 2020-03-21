<script>
  import Modal from "./Modal";
  import marked from "marked";
  import hljs from "highlight.js";
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  export let id = "";
  export let content = "";
  export let createdAt = Date.now();
  export const username = "";

  let timestamp;
  let showDelete = false;
  let showDeleteModal = false;

  $: timestamp = new Date(createdAt).toLocaleString(undefined, {
    weekday: "short",
    month: "long",
    day: "numeric",
    year: "numeric",
    hour: "numeric",
    minute: "numeric"
  });

  marked.setOptions({
    gfm: true,
    highlight: function(code, language) {
      const validLanguage = hljs.getLanguage(language) ? language : "plaintext";
      return hljs.highlight(validLanguage, code).value;
    }
  });
</script>

<style>
  .aside {
    font-size: 75%;
  }
</style>

<div class="row">
  <div class="col l8 s12 offset-l2">
    <div
      class="card grey darken-4"
      on:mouseenter={() => (showDelete = true)}
      on:mouseleave={() => (showDelete = false)}>
      <div class="card-content white-text">
        {#if showDelete}
          <div
            class="btn-floating red right"
            on:click={() => (showDeleteModal = true)}>
            <i class="material-icons white-text">delete</i>
          </div>
        {/if}
        <i class="aside">{timestamp}</i>
        <p>
          {@html marked(content, { gfm: true })}
        </p>
      </div>
    </div>
  </div>
</div>

{#if showDeleteModal}
  <Modal on:close={() => (showDeleteModal = false)}>
    <div class="card grey darken-4">
      <div class="card-content white-text">
        <div>Are you sure you want to delete this post?</div>
      </div>
      <div class="card-action">
        <button
          class="btn grey darken-2"
          on:click={() => (showDeleteModal = false)}>
          Cancel
        </button>
        <button
          class="btn red darken-2 right"
          on:click={dispatch('deletePost', id)}>
          Delete
        </button>
      </div>
    </div>
  </Modal>
{/if}
