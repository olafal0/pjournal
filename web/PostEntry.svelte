<script>
  import request from "./request.ts";
  import { localEncryptEnabled, encryptKey } from "./store";
  import { encryptPost } from "./cryption";
  import { createEventDispatcher, tick, onMount } from "svelte";
  import M from "materialize-css";

  const dispatch = createEventDispatcher();

  export let newPostContent = "";
  export let postInProgress = false;
  export let postError = null;
  export let update = false;
  export let id = "";

  function submitContent() {
    encryptPost(newPostContent, $localEncryptEnabled, $encryptKey)
      .then(postData => {
        return request("/api/posts/new", "POST", postData).then(postId => {
          dispatch("newPostTrigger", { postId });
          newPostContent = "";
          tick().then(() =>
            M.textareaAutoResize(document.getElementById("post-text-field"))
          );
        });
      })
      .catch(err => {
        postError = err;
      })
      .finally(() => {
        postInProgress = false;
      });
  }

  onMount(() => {
    tick().then(() =>
      M.textareaAutoResize(document.getElementById(`post-text-field${id}`))
    );
  });
</script>

<div class="card grey darken-4">
  <div class="card-content">
    {#if postError}
      <div class="red-text">{postError}</div>
    {/if}
    <textarea
      disabled={postInProgress}
      bind:value={newPostContent}
      id="post-text-field{id}"
      spellcheck="true"
      placeholder="How was your day?"
      class="materialize-textarea input-field white-text" />
  </div>
  <div class="card-action">
    {#if update}
      <button
        class="btn waves-effect waves-light grey darken-2"
        on:click={() => dispatch('cancelEdit')}>
        Cancel
      </button>
      <button
        class="btn waves-effect waves-light blue darken-4 right"
        on:click={() => dispatch('updatePost')}>
        Submit
      </button>
    {:else}
      <button
        class="btn waves-effect waves-light blue darken-4"
        on:click={submitContent}>
        Submit
      </button>
    {/if}
  </div>
</div>
