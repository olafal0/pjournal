<script>
  import Login from "./Login";
  import { onMount } from "svelte";
  import request from "./request.ts";
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  let authToken = "";
  let username = "";
  let password = "";

  function handleLogin(event) {
    event.preventDefault();
    request("http://localhost:8000/login", "POST", {
      Username: username,
      Password: password
    }).then(token => {
      dispatch("loggedIn", token);
    });
  }

  function handleRegistration(event) {
    event.preventDefault();
    request("http://localhost:8000/register", "POST", {
      Username: username,
      Password: password
    }).then(token => {
      dispatch("loggedIn", token);
    });
  }

  onMount(() => {
    authToken = localStorage.getItem("authToken");
  });
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
          </form>
        </div>
        <div class="row">
          <button
            class="btn waves-effect waves-light grey darken-2 col s6 m3"
            on:click={handleRegistration}>
            Register
          </button>
          <button
            class="btn waves-effect waves-light grey darken-2 col s6 m3 right"
            on:click={handleLogin}>
            Log In
          </button>
        </div>
      </div>
    </div>
  </div>
</div>
