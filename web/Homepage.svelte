<script>
  import Post from "./Post";
  import config from "./config";
  import PostEntry from "./PostEntry.svelte";
  import { onMount } from "svelte";
  import request from "./request.ts";

  let posts = [];
  let newPostContent = "";
  let postInProgress = false;
  let postError = null;

  // loadPosts loads the initial page of posts only
  async function loadPosts() {
    const response = await request.get(`${config.apiUrl}/posts?n=25`);
    posts = response.posts;
  }

  function deletePost(event) {
    const id = event.detail;
    request
      .delete(`${config.apiUrl}/post/${id}`)
      .then(loadPosts)
      .catch(error => {
        postError = error;
      });
  }

  onMount(loadPosts);
</script>

<div class="container">
  <PostEntry on:newPostTrigger={loadPosts} />
  {#each posts as post (post.id)}
    <Post {...post} on:deletePost={deletePost} />
  {/each}
</div>
