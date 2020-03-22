<script>
  import Modal from "./Modal";
  import request from "./request";
  import marked from "marked";
  import hljs from "highlight.js";
  import { createEventDispatcher, tick } from "svelte";
  import M from "materialize-css";

  const dispatch = createEventDispatcher();

  export let id = "";
  export let content = "";
  export let createdAt = Date.now();
  export let updatedAt = undefined;
  export const username = "";

  let createdTimeStr;
  let updatedTimeStr;
  let hovering = false;
  let showDeleteModal = false;
  let editing = false;
  let editedContent = content;
  let updateInProgress = false;
  let updateError = "";

  function updatePost() {
    updateInProgress = true;
    request(`/api/post/${id}`, "POST", editedContent)
      .then(() => {
        content = editedContent;
        editing = false;
      })
      .catch(error => {
        updateError = error;
      })
      .finally(() => {
        updateInProgress = false;
      });
  }

  function formatDate(d) {
    return new Date(d).toLocaleString(undefined, {
      weekday: "short",
      month: "long",
      day: "numeric",
      year: "numeric",
      hour: "numeric",
      minute: "numeric"
    });
  }

  $: createdTimeStr = formatDate(createdAt);
  $: updatedTimeStr = formatDate(updatedAt);

  $: if (editing) {
    tick().then(() =>
      M.textareaAutoResize(document.getElementById(`${id}-edit-field`))
    );
  }

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
      on:mouseenter={() => (hovering = true)}
      on:mouseleave={() => (hovering = false)}>
      <div class="card-content white-text">
        {#if hovering && !editing}
          <div class="right">
            <div class="btn red" on:click={() => (showDeleteModal = true)}>
              Delete
            </div>
            <div class="btn grey" on:click={() => (editing = true)}>Edit</div>
          </div>
        {/if}
        <i class="aside">
          {createdTimeStr}
          {#if updatedAt}edited {updatedTimeStr}{/if}
        </i>
        <p>
          {#if editing}
            <div class="card-content">
              {#if updateError}
                <div class="red-text">{updateError}</div>
              {/if}
              <textarea
                disabled={updateInProgress}
                bind:value={editedContent}
                id={`${id}-edit-field`}
                spellcheck="true"
                placeholder="How was your day?"
                class="materialize-textarea input-field s12 white-text" />
            </div>
            <div class="card-action">
              <button
                class="btn waves-effect waves-light grey darken-2"
                on:click={() => {
                  editing = false;
                }}>
                Cancel
              </button>
              <button
                class="btn waves-effect waves-light blue darken-4 right"
                on:click={updatePost}>
                Submit
              </button>
            </div>
          {:else}
            {@html marked(content, { gfm: true })}
          {/if}
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
