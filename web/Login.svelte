<script>
  import Auth from "./Auth";
  import { onMount, createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  const FlowState = Object.freeze({
    loading: 1,
    initial: 2,
    signUpConfirmation: 3,
    forgotPassword: 4
  });

  let formFields = {
    email: "",
    password: "",
    newPassword: "",
    confirmationCode: ""
  };

  let flowState = FlowState.loading;
  let infoMessage = "";
  let errorMessage = "";

  onMount(() => {
    Auth.currentAuthenticatedUser()
      .then(signedIn)
      .catch(err => {
        if (err === "not authenticated") {
          flowState = FlowState.initial;
        } else {
          displayError(err);
        }
      });
  });

  function msgToString(msg) {
    if (typeof msg === "string") {
      return msg;
    }
    if (msg.message && typeof msg.message === "string") {
      return msg.message;
    }
    return JSON.stringify(msg);
  }

  function clearMessages() {
    infoMessage = "";
    errorMessage = "";
  }

  function displayError(msg) {
    clearMessages();
    errorMessage = msgToString(msg);
  }

  function displayInfo(msg) {
    clearMessages();
    infoMessage = msgToString(msg);
  }

  function signedIn() {
    dispatch("signedIn");
  }

  function signUp() {
    displayInfo("Registering...");
    // TODO: Special case handling for common errors or problems (eg "@gmail.con")
    formFields.email = formFields.email;
    Auth.signUp({
      username: formFields.email,
      password: formFields.password,
      attributes: {
        email: formFields.email
      }
    })
      .then(({ user, userConfirmed, userSub }) => {
        if (!userConfirmed) {
          // User is not confirmed; we need a confirmation code
          displayInfo(
            "Please confirm by entering the confirmation code sent to your email address"
          );
          flowState = FlowState.signUpConfirmation;
          return;
        }
        signedIn();
      })
      .catch(displayError);
  }

  function resendCode() {
    displayInfo(`Resending confirmation code for ${formFields.email}`);
    Auth.resendSignUp(formFields.email).catch(displayError);
  }

  function confirmSignUp() {
    displayInfo("Confirming user...");
    Auth.confirmSignUp(formFields.email, formFields.confirmationCode)
      .then(result => {
        if (result === "SUCCESS") {
          // Set back to initial, now that user needs to sign in for the first time
          displayInfo("User confirmed! Please sign in");
          flowState = FlowState.initial;
          return;
        }
      })
      .catch(displayError);
  }

  function signIn() {
    displayInfo("Signing in...");
    Auth.signIn(formFields.email, formFields.password)
      .then(signedIn)
      .catch(err => {
        if (err.code === "UserNotConfirmedException") {
          displayInfo(
            "Please confirm by entering the confirmation code sent to your email address"
          );
          flowState = FlowState.signUpConfirmation;
        } else {
          displayError(err);
        }
      });
  }

  async function startPasswordReset() {
    if (!formFields.email) {
      displayError("Please enter your email first");
      return;
    }
    try {
      await Auth.forgotPassword(formFields.email);
      displayInfo(
        "Please enter your new password and the confirmation code sent to your email address"
      );
      flowState = FlowState.forgotPassword;
    } catch (e) {
      displayError(e);
    }
  }

  function cancelPasswordReset() {
    clearMessages();
    flowState = FlowState.initial;
  }

  async function confirmResetPassword() {
    displayInfo("Resetting password...");
    try {
      await Auth.forgotPasswordSubmit(
        formFields.email,
        formFields.confirmationCode,
        formFields.newPassword
      );
      displayInfo("Please log in with your new password");
      flowState = FlowState.initial;
    } catch (e) {
      displayError(e);
    }
  }
</script>

<style>
  .info-message {
    margin: 10px;
    padding: 5px;
    border: 1px solid rgba(0, 195, 209, 0.5);
    border-radius: 5px;
    color: antiquewhite;
  }

  .error-message {
    margin: 10px;
    padding: 5px;
    border: 1px solid rgba(209, 0, 0, 0.5);
    border-radius: 5px;
    color: rgb(250, 208, 153);
  }
</style>

<div class="row">
  <div class="col l8 offset-l2 m10 offset-m1 s12 center">
    <div class="card grey darken-4">
      {#if errorMessage}
        <div class="col s12">
          <div class="error-message">{errorMessage}</div>
        </div>
      {:else if infoMessage}
        <div class="col s12">
          <div class="info-message">{infoMessage}</div>
        </div>
      {/if}
      <div class="card-content white-text">
        <div class="row">
          <h5 class="card-title">Log in</h5>
          {#if flowState === FlowState.initial}
            <form on:submit|preventDefault={signIn}>
              <div class="form-group">
                <label>Email Address</label>
                <input
                  type="email"
                  class="white-text form-control"
                  bind:value={formFields.email} />
              </div>
              <div class="form-group">
                <label>Password</label>
                <input
                  type="button"
                  class="btn btn-flat btn-small grey darken-4 grey-text"
                  on:click={startPasswordReset}
                  value="(Forgot password?)" />
                <input
                  type="password"
                  class="white-text form-control"
                  bind:value={formFields.password} />
              </div>
              <input
                type="button"
                class="btn grey darken-2 left"
                on:click={signUp}
                value="Register" />
              <input type="submit" class="btn right" value="Log In" />
            </form>
          {:else if flowState === FlowState.signUpConfirmation}
            <form on:submit|preventDefault={confirmSignUp}>
              <div class="form-group">
                <label>Confirmation Code</label>
                <input
                  type="text"
                  class="white-text form-control"
                  bind:value={formFields.confirmationCode} />
              </div>
              <input
                type="button"
                class="btn grey darken-2 left"
                on:click={resendCode}
                value="Resend Confirmation Code" />
              <input type="submit" class="btn right" value="Confirm" />
            </form>
          {:else if flowState === FlowState.forgotPassword}
            <form on:submit|preventDefault={confirmResetPassword}>
              <div class="form-group">
                <label>New Password</label>
                <input
                  type="password"
                  class="white-text form-control"
                  bind:value={formFields.newPassword} />
              </div>
              <div class="form-group">
                <label>Confirmation Code</label>
                <input
                  type="text"
                  class="white-text form-control"
                  bind:value={formFields.confirmationCode} />
              </div>
              <div class="left">
                <input
                  type="button"
                  class="btn grey darken-2"
                  on:click={cancelPasswordReset}
                  value="Cancel Reset" />
                <input
                  type="button"
                  class="btn grey darken-2"
                  on:click={startPasswordReset}
                  value="Resend Code" />
              </div>
              <input type="submit" class="btn right" value="Confirm" />
            </form>
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>
