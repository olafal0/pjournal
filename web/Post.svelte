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
    encryptPost(editedContent, $localEncryptEnabled && encrypted, $encryptKey)
      .then(encryptedPost => {
        displayContent = editedContent;
        encrypted = encryptedPost.encrypted;
        iv = encryptedPost.iv;
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
      $localEncryptEnabled,
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

{#if editing}
  <PostEntry
    bind:newPostContent={editedContent}
    bind:encrypted
    {id}
    postInProgress={updateInProgress}
    on:updatePost={updatePost}
    update="true"
    on:cancelEdit={() => {
      editing = false;
    }} />
{:else}
  <div
    class="card padded"
    on:mouseenter={() => (hovering = true)}
    on:mouseleave={() => (hovering = false)}>
    {#if encrypted}
      <div class="material-icons right" title="Encrypted">lock</div>
    {/if}
    {#if hovering}
      <div class="right">
        <button class="danger" on:click={() => (showDeleteModal = true)}>
          Delete
        </button>
        <button on:click={() => (editing = true)}>Edit</button>
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
{/if}

{#if showDeleteModal}
  <Modal on:close={() => (showDeleteModal = false)}>
    <div class="card padded">
      <div>Are you sure you want to delete this post?</div>
      <button on:click={() => (showDeleteModal = false)}>Cancel</button>
      <button class="danger right" on:click={dispatch('deletePost', id)}>
        Delete
      </button>
    </div>
  </Modal>
{/if}
