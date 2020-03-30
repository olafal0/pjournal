<script>
  import request from "./request.ts";
  import { localEncryptEnabled, encryptKey } from "./store";
  import { encryptPost } from "./cryption";
  import { createEventDispatcher, onMount } from "svelte";

  const dispatch = createEventDispatcher();

  export let newPostContent = "";
  export let postInProgress = false;
  export let postError = null;
  export let update = false;
  export let id = "";

  let textarea;

  function submitContent() {
    encryptPost(newPostContent, $localEncryptEnabled, $encryptKey)
      .then(postData => {
        return request("/api/posts/new", "POST", postData).then(postId => {
          dispatch("newPostTrigger", { postId });
          newPostContent = "";
        });
      })
      .catch(err => {
        postError = err;
      })
      .finally(() => {
        postInProgress = false;
        updateFieldHeight();
      });
  }

  function updateFieldHeight() {
    // Set height to auto to shrink if above min-height
    textarea.style.height = "auto";

    // Round height to multiples of lineHeight
    let lineHeightPx = parseFloat(getComputedStyle(textarea).lineHeight);
    let height = textarea.scrollHeight;
    height += lineHeightPx - (height % lineHeightPx);
    textarea.style.height = height + "px";
  }

  onMount(() => {
    updateFieldHeight();
  });
</script>

<style>
  .pj-textarea {
    background-color: transparent;
    border: none;
    border-bottom: 1px solid grey;
    width: 100%;
    outline: none;
    min-height: 3em;
    resize: none;
    box-sizing: border-box;
    overflow-y: hidden;
  }

  .pj-textarea:focus {
    border-bottom: 1px solid hotpink;
  }
</style>

<div class="card grey darken-4">
  <div class="card-content">
    {#if postError}
      <div class="red-text">{postError}</div>
    {/if}
    <textarea
      disabled={postInProgress}
      bind:this={textarea}
      bind:value={newPostContent}
      on:input={updateFieldHeight}
      id="post-text-field{id}"
      spellcheck="true"
      placeholder="How was your day?"
      class="pj-textarea input-field white-text" />
  </div>
  <div class="card-action">
    {#if update}
      <button
        class="btn grey darken-2"
        on:click={() => dispatch('cancelEdit')}
        on:click={updateFieldHeight}>
        Cancel
      </button>
      <button
        class="btn blue darken-4 right"
        on:click={() => dispatch('updatePost')}
        on:click={updateFieldHeight}>
        Submit
      </button>
    {:else}
      <button class="btn blue darken-4" on:click={submitContent}>Submit</button>
    {/if}
  </div>
</div>
