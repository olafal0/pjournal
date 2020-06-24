<script>
  import Homepage from "./Homepage";
  import Login from "./Login";
  import { localEncryptEnabled, encryptKey } from "./store";
  import { onMount, createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  let settings = {
    encryption: $localEncryptEnabled,
    encryptPass: ""
  };

  function savePrefs() {
    localEncryptEnabled.set(settings.encryption);
    encryptKey.set(settings.encryptPass);
    dispatch("closeSettings");
  }

  function closeSettings() {
    settings.encryption = $localEncryptEnabled;
    dispatch("closeSettings");
  }
</script>

<div class="container">
  <div class="card padded">
    <button
      class="borderless material-icons left"
      style="padding-left:0px;"
      on:click={closeSettings}>
      arrow_back
    </button>
    <h2>Settings</h2>
    <p>
      Local encryption enables you to save and edit journal entries without the
      actual content of your posts ever leaving your machine in unencrypted
      form. All posts are encrypted with a key derived from your local
      encryption password before they are sent to the server. The key itself is
      never sent to the server, and only the encrypted content of your posts are
      stored on the server.
    </p>
    <p>Please use a different password than your account password.</p>
    <br />
    <form class="nopad" on:submit|preventDefault={savePrefs}>
      <label class="right">
        Enable local encryption
        <input type="checkbox" bind:checked={settings.encryption} />
      </label>
      <input name="dummy-password" type="password" class="hide" />
      <label>
        Local encryption password
        <input
          type="password"
          name="encryption-pass"
          bind:value={settings.encryptPass}
          disabled={!settings.encryption} />
      </label>
      <input type="submit" value="Save" />
    </form>
  </div>
</div>
