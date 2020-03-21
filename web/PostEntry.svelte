<script>
  import request from "./request.ts";
  import { createEventDispatcher, tick } from "svelte";
  import M from "materialize-css";

  const dispatch = createEventDispatcher();

  export let authToken = "";
  let newPostContent = "";
  let postInProgress = false;
  let postError = null;

  function handleNewPost() {
    // Submit new post and reload the page
    postInProgress = true;
    request("/api/posts/new", "POST", newPostContent, {
      auth: authToken
    })
      .then(postId => {
        dispatch("newPostTrigger", { postId });
        newPostContent = "";
        tick().then(() =>
          M.textareaAutoResize(document.getElementById("post-text-field"))
        );
      })
      .catch(err => {
        postError = err;
      })
      .finally(() => {
        postInProgress = false;
      });
  }
</script>

<div class="row">
  <div class="card grey darken-4 ">
    <div class="card-content">
      {#if postError}
        <div class="red-text">{postError}</div>
      {/if}
      <textarea
        disabled={postInProgress}
        bind:value={newPostContent}
        id="post-text-field"
        spellcheck="true"
        placeholder="How was your day?"
        class="materialize-textarea input-field s12 white-text" />
    </div>
    <div class="card-action">
      <button
        class="btn waves-effect waves-light grey darken-2"
        on:click={handleNewPost}>
        Submit
      </button>
    </div>
  </div>
</div>
