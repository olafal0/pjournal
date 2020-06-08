import App from './App';

// Load regenerator-runtime to allow transpilation of async Svelte functions
import 'regenerator-runtime/runtime';

const app = new App({
  target: document.getElementById('svelte-app'),
});
