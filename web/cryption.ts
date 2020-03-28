export async function createKey(pw, salt) {
  const enc = new TextEncoder();
  const pwEncoded = enc.encode(pw);
  salt = enc.encode(salt);

  const cryptoPw = await crypto.subtle.importKey(
    'raw',
    pwEncoded,
    'PBKDF2',
    false,
    ['deriveKey']
  );
  const key = await crypto.subtle.deriveKey(
    {
      name: 'PBKDF2',
      salt,
      iterations: 10000,
      hash: 'SHA-256'
    },
    cryptoPw,
    {
      name: 'AES-CBC',
      length: 256
    },
    true, // must be able to export to save locally
    ['encrypt', 'decrypt']
  );

  // export as a JWK to save to local storage
  const exportedKey = await crypto.subtle.exportKey('jwk', key);

  return exportedKey;
}

async function encrypt(exportedKey, content) {
  const enc = new TextEncoder();
  const dec = new TextDecoder();

  const iv = crypto.getRandomValues(new Uint8Array(16));

  const key = await crypto.subtle.importKey(
    'jwk',
    exportedKey,
    'AES-CBC',
    true,
    ['encrypt', 'decrypt']
  );

  const ciphertext = await crypto.subtle.encrypt(
    {
      name: 'AES-CBC',
      iv
    },
    key,
    enc.encode(content)
  );

  return { ciphertext, iv };
}

async function decrypt(exportedKey, ciphertext, iv) {
  const dec = new TextDecoder();

  const key = await crypto.subtle.importKey(
    'jwk',
    exportedKey,
    'AES-CBC',
    true,
    ['encrypt', 'decrypt']
  );

  const plaintext = await crypto.subtle.decrypt(
    { name: 'AES-CBC', iv },
    key,
    ciphertext
  );

  return dec.decode(plaintext);
}

export async function encryptPost(content, encryptEnabled, key) {
  if (encryptEnabled == 'true') {
    const { ciphertext, iv } = await encrypt(JSON.parse(key), content);
    const ciphertextView = new Uint8Array(ciphertext);

    return {
      content: String.fromCodePoint(...ciphertextView),
      iv: String.fromCodePoint(...iv),
      encrypted: true
    };
  }

  return { content };
}

export async function decryptPost(
  { content, encrypted, iv },
  encryptEnabled,
  key
) {
  if (!encrypted || !encryptEnabled) {
    return content;
  }

  // convert unicode strings into arrays of code point values
  const contentBytes = [];
  for (const codePoint of content) {
    contentBytes.push(codePoint.codePointAt(0));
  }
  const ivBytes = [];
  for (const codePoint of iv) {
    ivBytes.push(codePoint.codePointAt(0));
  }

  let plaintext: string;
  try {
    plaintext = await decrypt(
      JSON.parse(key),
      Uint8Array.from(contentBytes),
      Uint8Array.from(ivBytes)
    );
  } catch (err) {
    plaintext = `Decryption failed: ${err}\n\n${content}`;
  }

  return plaintext;
}
