<script>
  import Homepage from "./Homepage";
  import Settings from "./Settings";
  import Login from "./Login";
  import Auth from "./Auth";
  import request from "./request";
  import { username, localEncryptEnabled, encryptKey } from "./store";
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
      $localEncryptEnabled = false;
      $encryptKey = null;
      user = null;
    });
  }
</script>

<nav>
  <h1>pjournal</h1>
  <div class="right">
    {#if user}
      <button
        class="borderless"
        on:click={() => (showSettings = !showSettings)}>
        Settings
      </button>
      <button class="borderless" on:click={logout}>Log Out</button>
    {/if}
  </div>
</nav>
{#if user}
  {#if showSettings}
    <Settings on:closeSettings={() => (showSettings = false)} />
  {:else}
    <Homepage />
  {/if}
{:else}
  <Login on:signedIn={signedIn} />
{/if}
