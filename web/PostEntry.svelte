<script>
  import request from "./request.ts";
  import config from "./config";
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
  let textareaStyle;

  function submitContent() {
    encryptPost(newPostContent, $localEncryptEnabled, $encryptKey)
      .then(postData => {
        return request
          .post(`${config.apiUrl}/posts/new`, postData)
          .then(postId => {
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
    if (!textareaStyle) {
      return;
    }
    // Set height to auto to shrink if above min-height
    textarea.style.height = "auto";

    // Round height to multiples of lineHeight
    let lineHeightPx;
    if (textareaStyle.lineHeight === "normal") {
      lineHeightPx = parseFloat(textareaStyle.fontSize) * 1.5;
    } else {
      lineHeightPx = parseFloat(textareaStyle.lineHeight);
    }
    let height = textarea.scrollHeight;
    height += lineHeightPx - (height % lineHeightPx);
    textarea.style.height = height + "px";
  }

  $: if (textarea) {
    textareaStyle = getComputedStyle(textarea);
  }

  onMount(() => {
    updateFieldHeight();
  });
</script>

<style>
  .pj-textarea {
    background-color: transparent;
    line-height: 1.5em;
    border: 1px solid grey;
    width: 100%;
    outline: none;
    min-height: 3em;
    resize: none;
    box-sizing: border-box;
    overflow-y: hidden;
    padding: 10px;
  }

  .card {
    padding: 10px 0px;
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
