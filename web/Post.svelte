<script>
  import Modal from "./Modal";
  import PostEntry from "./PostEntry";
  import request from "./request";
  import config from "./config";
  import { encryptKey, localEncryptEnabled } from "./store";
  import { encryptPost, decryptPost } from "./cryption";
  import marked from "marked";
  import { createEventDispatcher, onMount } from "svelte";

  const dispatch = createEventDispatcher();

  export let id = "";
  export let content = "";
  export let createdAt = Date.now();
  export let updatedAt = undefined;
  export let encrypted = false;
  export let iv = "";

  let displayContent = "";
  let createdTimeStr;
  let updatedTimeStr;
  let hovering = false;
  let showDeleteModal = false;
  let editing = false;
  let editedContent = displayContent;
  let updateInProgress = false;
  let updateError = "";

  function updatePost() {
    updateInProgress = true;
    encryptPost(editedContent, $localEncryptEnabled, $encryptKey)
      .then(encryptedPost => {
        displayContent = editedContent;
        encrypted = encryptedPost.encrypted;
        editing = false;
        const encryptedFullPost = {
          id,
          createdAt,
          updatedAt,
          encrypted,
          iv,
          ...encryptedPost
        };
        return request.post(`${config.apiUrl}/post/${id}`, encryptedFullPost);
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

  $: editedContent = displayContent;

  onMount(() => {
    marked.setOptions({
      gfm: true
    });

    decryptPost(
      {
        content,
        encrypted,
        iv
      },
      $localEncryptEnabled == "true",
      $encryptKey
    ).then(c => {
      displayContent = c;
    });
  });
</script>

<style>
  .aside {
    font-size: 75%;
  }
  :global(.markdown-content p) {
    padding-top: 0.5em;
    padding-bottom: 0.5em;
  }
</style>

<div class="row">
  <div class="col l8 s12 offset-l2">
    {#if editing}
      <PostEntry
        bind:newPostContent={editedContent}
        {id}
        postInProgress={updateInProgress}
        on:updatePost={updatePost}
        update="true"
        on:cancelEdit={() => {
          editing = false;
        }} />
    {:else}
      <div
        class="card grey darken-4"
        on:mouseenter={() => (hovering = true)}
        on:mouseleave={() => (hovering = false)}>
        <div class="card-content white-text">
          {#if encrypted}
            <div class="material-icons white-text right" title="Encrypted">
              lock
            </div>
          {/if}
          {#if hovering}
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
          <div class="markdown-content">
            {@html marked(displayContent, { gfm: true })}
          </div>
        </div>
      </div>
    {/if}
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
