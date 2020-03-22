<script>
  import Homepage from "./Homepage";
  import Login from "./Login";
  import request from "./request";
  import { onMount } from "svelte";

  let isLoggedIn = "";

  function checkLoggedIn() {
    const regex = /dispatch-logged-in=(?<loginStatus>\w+)/;
    const cookies = document.cookie;
    const cookieMatch = cookies.match(regex);
    isLoggedIn = cookieMatch && cookieMatch.groups.loginStatus == "true";
  }

  function logout() {
    request("/api/logout", "GET").then(checkLoggedIn);
  }

  onMount(() => {
    checkLoggedIn();
  });
</script>

<nav>
  <div class="row">
    <div class="nav-wrapper grey darken-4">
      <div class="brand-logo">pjournal</div>
      <ul class="col right">
        {#if isLoggedIn}
          <li>
            <button class="btn grey darken-2" on:click={logout}>Log Out</button>
          </li>
        {/if}
      </ul>
    </div>
  </div>
</nav>
<div class="container">
  <div class="col">
    {#if isLoggedIn}
      <Homepage />
    {:else}
      <Login on:loggedIn={checkLoggedIn} />
    {/if}
  </div>
</div>
