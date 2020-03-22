<script>
  import Post from "./Post";
  import PostEntry from "./PostEntry.svelte";
  import { onMount } from "svelte";
  import request from "./request.ts";

  let posts = [];
  let newPostContent = "";
  let postInProgress = false;
  let postError = null;

  function loadPosts() {
    request("/api/posts/all", "GET").then(p => {
      p.reverse();
      posts = p;
    });
  }

  function deletePost(event) {
    const id = event.detail;
    request(`/api/post/${id}`, "DELETE")
      .then(loadPosts)
      .catch(error => {
        postError = error;
      });
  }

  onMount(loadPosts);
</script>

<div>
  <PostEntry on:newPostTrigger={loadPosts} />
  {#each posts as post (post.id)}
    <Post {...post} on:deletePost={deletePost} />
  {/each}
</div>
