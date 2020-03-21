interface options {
  auth?: string;
  contentType?: string;
}

const defaultOptions = {
  auth: null,
  contentType: 'application/json'
};

export default function request(
  url: string,
  method?: string,
  formData?: object | string,
  options?: options
) {
  url = url;
  method = method || 'GET';
  formData = formData || null;
  let { contentType, auth } = Object.assign({}, defaultOptions, options);

  const req = new XMLHttpRequest();
  req.open(method, url);
  req.setRequestHeader('Content-Type', contentType);
  if (auth) {
    req.setRequestHeader('Authorization', auth);
  }

  return new Promise((resolve, reject) => {
    req.addEventListener('load', () => {
      if (req.status >= 400) {
        reject(req.statusText);
        return;
      }
      const contentType = req.getResponseHeader('Content-Type');
      if (contentType == 'application/json') {
        resolve(JSON.parse(req.response));
      }
      resolve(req.response);
    });
    req.addEventListener('error', () => reject(req.statusText));
    if (contentType == 'application/json' && formData) {
      formData = JSON.stringify(formData);
    } else {
      formData = formData as string;
    }
    req.send(formData);
  });
}
