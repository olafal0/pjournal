import { createKey } from './cryption';

// persistentStore implements the svelte/stores contract, but persists values to
// local storage automatically.
function persistentStore(key, setHook?: Function) {
  let subscribers = [];
  let _value = localStorage.getItem(key) || '';

  return {
    get() {
      return _value;
    },

    subscribe(listener) {
      listener(_value);
      subscribers.push(listener);
      return () => {
        subscribers = subscribers.filter(sub => sub !== listener);
      };
    },

    async set(newValue) {
      if (setHook) {
        _value = await setHook(newValue);
      } else {
        _value = newValue;
      }
      localStorage.setItem(key, _value);
      subscribers.forEach(s => s(_value));
    }
  };
}

export const username = persistentStore('pjournal-username');
export const localEncryptEnabled = persistentStore('pjournal-encrypt-enabled');

// encryptKey can be set with a password, and will automatically derive an
// AES-CBC JSON Web Key using PBKDF2 and store it. Accessing this value will return
// a stringified JWK object, not the password that it's set with.
export const encryptKey = persistentStore('pjournal-local-key', async pw => {
  if (!localEncryptEnabled.get || !pw) {
    return '';
  }
  const exportedKey = await createKey(pw, username.get());

  return JSON.stringify(exportedKey);
});
