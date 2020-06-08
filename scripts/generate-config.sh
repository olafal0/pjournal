#!/bin/bash

set -e
source scripts/aws-tools.sh

# Check dependencies
source scripts/depcheck.sh
depcheck jq
depcheck aws

if [[ -z "$1" ]]; then
  export DEPLOY_ENV=development
else
  export DEPLOY_ENV=$1
fi

prefix="pjournal-$DEPLOY_ENV"

BACKEND_API_ID=$(getStackOutput $prefix BackendAPIId)
USER_POOL_ID=$(getStackOutput $prefix UserPoolId)
WEB_CLIENT_ID=$(getStackOutput $prefix WebClientId)

echo "export default {
  Auth: {
    region: 'us-east-1',
    userPoolId: '$USER_POOL_ID',
    userPoolWebClientId: '$WEB_CLIENT_ID',
  },
  apiUrl: 'https://$BACKEND_API_ID.execute-api.us-east-1.amazonaws.com/api',
}
" > web/config.js
