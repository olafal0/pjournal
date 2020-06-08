import Auth from '@aws-amplify/auth';

import config from './config';

Auth.configure({
  Auth: config.Auth,
});

export default Auth;
