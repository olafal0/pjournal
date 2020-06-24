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
  export let encrypted = $localEncryptEnabled;

  let textarea;
  let textareaStyle;

  function submitContent() {
    encryptPost(newPostContent, $localEncryptEnabled && encrypted, $encryptKey)
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
    width: 100%;
    outline: none;
    min-height: 3em;
    resize: none;
    box-sizing: border-box;
    overflow-y: hidden;
    padding: 10px;
    margin-bottom: 15px;
  }
</style>

<div class="card">
  <form on:submit|preventDefault>
    {#if postError}
      <div class="error-message">{postError}</div>
    {/if}
    <textarea
      disabled={postInProgress}
      bind:this={textarea}
      bind:value={newPostContent}
      on:input={updateFieldHeight}
      id="post-text-field{id}"
      spellcheck="true"
      placeholder="How was your day?"
      class="pj-textarea" />
    {#if update}
      <button
        on:click={() => dispatch('cancelEdit')}
        on:click={updateFieldHeight}>
        Cancel
      </button>
      <button
        class="primary right"
        on:click={() => dispatch('updatePost')}
        on:click={updateFieldHeight}>
        Submit
      </button>
    {:else}
      <button class="primary" on:click={submitContent}>Submit</button>
    {/if}
    <label class="right">
      <input type="checkbox" bind:checked={encrypted} />
      <span>Encrypt</span>
    </label>
  </form>
</div>
