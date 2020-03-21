<style>
  .icon-spinner {
    animation-duration: 1s;
    transition-property: transform;
    animation-name: infinite-spin;
    animation-iteration-count: infinite;
    animation-timing-function: linear;
  }
  @keyframes infinite-spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(359deg);
    }
  }
</style>

<i class="material-icons icon-spinner medium">rotate_right</i>
