import App from './App';

// Load regenerator-runtime to allow transpilation of async Svelte functions
import 'regenerator-runtime/runtime';

const app = new App({
  target: document.getElementById('svelte-app'),
});

if ('serviceWorker' in navigator) {
  window.addEventListener('load', function () {
    navigator.serviceWorker.register('/serviceWorker.js').then(
      function (registration) {
        console.log(
          'ServiceWorker registration successful with scope: ',
          registration.scope
        );
      },
      function (err) {
        console.log('ServiceWorker registration failed: ', err);
      }
    );
  });
}
