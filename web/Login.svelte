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

  onMount(() => {
    authToken = localStorage.getItem("authToken");
  });
</script>

<div class="row">
  <div class="col l6 offset-l3 m8 offset-m2 s12 center">
    <div class="card grey darken-4">
      <div class="card-content white-text">
        <form id="login-form" on:submit={handleLogin}>
          <label for="username" class="white-text">Username</label>
          <br />
          <input
            type="text"
            class="white-text"
            name="username"
            bind:value={username} />
          <br />
          <label for="password" class="white-text">Password</label>
          <br />
          <input
            type="password"
            class="white-text"
            name="password"
            bind:value={password} />
          <br />
          <button
            class="btn waves-effect waves-light grey darken-2"
            on:click={handleLogin}>
            Log In
          </button>
          <br />
        </form>
      </div>
    </div>
  </div>
</div>
