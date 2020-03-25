<script>
  import Homepage from "./Homepage";
  import Login from "./Login";
  import request from "./request";
  import { localEncryptEnabled, encryptKey } from "./store";
  import { onMount, createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  let localEncryptionEnabled = $localEncryptEnabled == "true";
  let localEncryptPass = "";

  function savePrefs() {
    localEncryptEnabled.set(localEncryptionEnabled ? "true" : "false");
    encryptKey.set(localEncryptPass);
    dispatch("closeSettings");
  }

  function closeSettings() {
    localEncryptionEnabled = $localEncryptEnabled == "true";
    dispatch("closeSettings");
  }
</script>

<div class="container">
  <div class="col">
    <div class="row">
      <div class="card grey darken-4 white-text">
        <div class="card-content white-text">
          <button
            class="btn-flat material-icons white-text left"
            style="padding-left:0px;"
            on:click={closeSettings}>
            arrow_back
          </button>
          <h1 class="small">Settings</h1>
          <p>
            Local encryption enables you to save and edit journal entries
            without the actual content of your posts ever leaving your machine
            in unencrypted form. All posts are encrypted with a key derived from
            your local encryption password before they are sent to the server.
            The key itself is never sent to the server, and only the encrypted
            content of your posts are stored on the server.
          </p>
          <br />
          <p>Please use a different password than your account password.</p>
          <br />
          <label class="white-text">
            <input type="checkbox" bind:checked={localEncryptionEnabled} />
            <span>Enable local encryption</span>
          </label>
          <form id="settings-form" on:submit|preventDefault={savePrefs}>
            <label class="white-text">
              <input name="dummy-password" type="password" class="hide" />
              <input
                class="white-text"
                type="password"
                id="encryption-pass"
                name="encryption-pass"
                bind:value={localEncryptPass}
                disabled={!localEncryptionEnabled} />
              <span>Local encryption password</span>
            </label>
            <input type="submit" class="hide" />
          </form>
        </div>
        <div class="card-action">
          <button class="btn grey darken-2" on:click={savePrefs}>Save</button>
        </div>
      </div>
    </div>
  </div>
</div>
