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
    request("/api/posts/all", "GET", null, {
      auth: authToken
    }).then(p => {
      p.reverse();
      posts = p;
    });
  }

  function deletePost(event) {
    const id = event.detail;
    request(`/api/post/${id}`, "DELETE", null, {
      auth: authToken
    })
      .then(loadPosts)
      .catch(error => {
        postError = error;
      });
  }

  onMount(loadPosts);
</script>

<div>
  <PostEntry {authToken} on:newPostTrigger={loadPosts} />
  {#each posts as post (post.id)}
    <Post {...post} on:deletePost={deletePost} />
  {/each}
</div>
