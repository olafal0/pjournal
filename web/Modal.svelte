<script>
  // Copied from https://svelte.dev/examples#modal
  import { createEventDispatcher, onDestroy } from "svelte";

  const dispatch = createEventDispatcher();

  let modal;

  const handle_keydown = e => {
    if (e.key === "Escape") {
      dispatch("close");
      return;
    }

    if (e.key === "Tab") {
      // trap focus
      const nodes = modal.querySelectorAll("*");
      const tabbable = Array.from(nodes).filter(n => n.tabIndex >= 0);

      let index = tabbable.indexOf(document.activeElement);
      if (index === -1 && e.shiftKey) index = 0;

      index += tabbable.length + (e.shiftKey ? -1 : 1);
      index %= tabbable.length;

      tabbable[index].focus();
      e.preventDefault();
    }
  };

  const previously_focused =
    typeof document !== "undefined" && document.activeElement;

  if (previously_focused) {
    onDestroy(() => {
      previously_focused.focus();
    });
  }
</script>

<style>
  .modal-bg {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.5);
    z-index: 1;
    backdrop-filter: blur(5px);
  }

  .modal-fg {
    position: absolute;
    left: 50%;
    top: 50%;
    width: calc(100vw - 4em);
    max-width: 32em;
    max-height: calc(100vh - 4em);
    overflow: auto;
    transform: translate(-50%, -50%);
    padding: 1em;
    border-radius: 0.2em;
    z-index: 0;
  }
</style>

<svelte:window on:keydown={handle_keydown} />

<div class="modal-bg" on:click={() => dispatch('close')}>

  <div
    class="modal-fg"
    role="dialog"
    aria-modal="true"
    bind:this={modal}
    on:click|stopPropagation>
    <slot />
  </div>
</div>
