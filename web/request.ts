import Auth from './Auth';

export default {
  get: async (url, params?: Object) => {
    if (params) {
      const searchParams = new URLSearchParams();
      Object.entries(params).forEach(([k, v]) => searchParams.append(k, v));
      url = `${url}?${searchParams.toString()}`;
    }
    const session = await Auth.currentSession();
    const response = await fetch(url, {
      method: 'GET',
      mode: 'cors',
      headers: {
        Authorization: session.getIdToken().getJwtToken(),
      },
    });
    if (response.status >= 300) {
      throw `${response.status}: ${response.statusText}`;
    }
    return response.json();
  },

  post: async (url, data) => {
    const session = await Auth.currentSession();
    const response = await fetch(url, {
      method: 'POST',
      mode: 'cors',
      headers: {
        Authorization: session.getIdToken().getJwtToken(),
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });
    if (response.status >= 300) {
      throw `${response.status}: ${response.statusText}`;
    }
    return response.json();
  },

  delete: async (url) => {
    const session = await Auth.currentSession();
    const response = await fetch(url, {
      method: 'DELETE',
      mode: 'cors',
      headers: {
        Authorization: session.getIdToken().getJwtToken(),
      },
    });
    if (response.status >= 300) {
      throw `${response.status}: ${response.statusText}`;
    }
    return response.json();
  },
};
