<script>
  import Post from "./Post";
  import PostEntry from "./PostEntry.svelte";
  import { onMount } from "svelte";
  import request from "./request.ts";

  export let authToken = "";
  let posts = [];
  let newPostContent = "";
  let postInProgress = false;
  let postError = null;

  function loadPosts() {
    request("http://localhost:8000/posts/all", "GET", null, {
      auth: authToken
    }).then(p => {
      p.reverse();
      posts = p;
    });
  }

  onMount(loadPosts);
</script>

<div>
  <PostEntry {authToken} on:newPostTrigger={loadPosts} />
  {#each posts as { content, createdAt }}
    <Post {content} {createdAt} />
  {/each}
</div>
