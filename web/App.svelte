<script>
  import Homepage from "./Homepage";
  import Login from "./Login";
  import { onMount } from "svelte";

  let authToken = "";

  function loggedIn(event) {
    localStorage.setItem("authToken", event.detail);
    authToken = event.detail;
  }

  function logout() {
    localStorage.removeItem("authToken");
    authToken = "";
  }

  onMount(() => {
    authToken = localStorage.getItem("authToken");
  });
</script>

<nav>
  <div class="row">
    <div class="nav-wrapper grey darken-4">
      <div class="brand-logo">pjournal</div>
      <ul class="col right">
        {#if authToken}
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
    {#if authToken}
      <Homepage {authToken} />
    {:else}
      <Login on:loggedIn={loggedIn} />
    {/if}
  </div>
</div>
