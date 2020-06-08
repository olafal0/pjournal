<script>
  import Homepage from "./Homepage";
  import Settings from "./Settings";
  import Login from "./Login";
  import Auth from "./Auth";
  import request from "./request";
  import { username } from "./store";
  import { onMount } from "svelte";

  let user = null;
  let showSettings = false;

  function signedIn() {
    Auth.currentAuthenticatedUser()
      .then(userData => {
        user = userData;
      })
      .catch(console.error);
  }

  function logout() {
    Auth.signOut().then(() => {
      user = null;
    });
  }
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
        {#if user}
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
    {#if user}
      {#if showSettings}
        <Settings on:closeSettings={() => (showSettings = false)} />
      {:else}
        <Homepage />
      {/if}
    {:else}
      <Login on:signedIn={signedIn} />
    {/if}
  </div>
</div>
