<script>
  import Homepage from "./Homepage";
  import Settings from "./Settings";
  import Login from "./Login";
  import request from "./request";
  import { username } from "./store";
  import { onMount } from "svelte";

  let isLoggedIn = "";
  let showSettings = false;

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

<style>
  nav .brand-logo {
    padding: 0px 10px;
  }
</style>

<nav>
  <div class="row">
    <div class="nav-wrapper grey darken-4">
      <div class="brand-logo">pjournal</div>
      <ul class="col right">
        {#if isLoggedIn}
          <li>
            <button
              class="btn-flat white-text"
              on:click={() => (showSettings = !showSettings)}>
              Settings
            </button>
          </li>
          <li>
            <button class="btn-flat white-text" on:click={logout}>
              Log Out
            </button>
          </li>
        {/if}
      </ul>
    </div>
  </div>
</nav>
<div class="container">
  <div class="col">
    {#if isLoggedIn}
      {#if showSettings}
        <Settings on:closeSettings={() => (showSettings = false)} />
      {:else}
        <Homepage />
      {/if}
    {:else}
      <Login on:loggedIn={checkLoggedIn} />
    {/if}
  </div>
</div>
