<script>
  import Login from "./Login";
  import Spinner from "./Spinner";
  import { onMount } from "svelte";
  import request from "./request.ts";
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  let username = "";
  let password = "";
  let inProgress = false;
  let loginError = "";

  function handleLogin(event) {
    event.preventDefault();
    inProgress = true;
    request("/api/login", "POST", {
      username,
      password
    })
      .then(() => {
        dispatch("loggedIn");
      })
      .catch(error => {
        loginError = error;
      })
      .finally(() => {
        inProgress = false;
      });
  }

  function handleRegistration(event) {
    event.preventDefault();
    request("/api/register", "POST", {
      username,
      password
    })
      .then(() => {
        dispatch("loggedIn");
      })
      .catch(error => {
        loginError = error;
      })
      .finally(() => {
        inProgress = false;
      });
  }
</script>

<div class="row">
  <div class="col l6 offset-l3 m8 offset-m2 s12 center">
    <div class="card grey darken-4">
      <div class="card-content white-text">
        <div class="row">
          <form id="login-form" on:submit={handleLogin}>
            <label for="username" class="white-text">Username</label>
            <input
              type="text"
              class="white-text"
              name="username"
              bind:value={username} />
            <label for="password" class="white-text">Password</label>
            <input
              type="password"
              class="white-text"
              name="password"
              bind:value={password} />
            <input type="submit" class="hide" />
            <button
              class="btn waves-effect waves-light grey darken-2 right"
              on:click={handleLogin}>
              Log In
            </button>
            <button
              class="btn waves-effect waves-light grey darken-2 left"
              on:click={handleRegistration}>
              Register
            </button>
          </form>
        </div>
        {#if inProgress}
          <Spinner />
        {:else if loginError}
          <div class="row red-text">{loginError}</div>
        {/if}
      </div>
    </div>
  </div>
</div>
