export default function request(
  url: string,
  method?: string,
  formData?: object | string,
  contentType?: string
) {
  url = url;
  method = method || 'GET';
  formData = formData || null;
  contentType = contentType || 'application/json';

  const req = new XMLHttpRequest();
  req.open(method, url);
  req.setRequestHeader('Content-Type', contentType);

  return new Promise((resolve, reject) => {
    req.addEventListener('load', () => {
      if (req.status >= 400) {
        reject(req.responseText);
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
