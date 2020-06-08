#!/bin/bash

# Provides some aws-related tools, and automates setup of awscli in CI environments.

aws configure set default.region us-east-1
if [[ -z $GITLAB_CI ]]; then
    echo "Running locally"
    if [[ -z "$1" ]]; then
        export DEPLOY_ENV=development
    else
        export DEPLOY_ENV=$1
    fi
else
    aws configure set aws_access_key_id $AWS_ACCESS_KEY_ID
    aws configure set aws_secret_access_key $AWS_SECRET_ACCESS_KEY
    echo "Configured aws cli"
    if [[ -z "$1" ]]; then
        export DEPLOY_ENV=$CI_COMMIT_REF_NAME
    else
        export DEPLOY_ENV=$1
    fi
    echo "Set DEPLOY_ENV=$DEPLOY_ENV"
fi

function getResourceId {
    aws cloudformation describe-stack-resource \
        --stack-name $1 \
        --logical-resource-id $2 \
        --query 'StackResourceDetail.PhysicalResourceId' --output text
}

function getStackOutput {
    STACK_NAME=$1
    OUTPUT_NAME=$2
    aws cloudformation describe-stacks --stack-name $STACK_NAME --output text \
        --query "Stacks[0].Outputs[?OutputKey==\`$OUTPUT_NAME\`][OutputValue]"
}

function getLatestVersion {
    aws s3api list-object-versions --bucket $1 --prefix $2 \
    --query 'Versions[0].VersionId' --output text
}

function deployStack {
    # Disable exiting on error so that we can display custom error output
    set +e
    aws cloudformation deploy --template-file $1 --stack-name $2 --no-fail-on-empty-changeset $3
    ERR=$?

    # Re-enable exit on error
    set -e

    if [[ $ERR -ne 0 ]]; then
        echo "Deployment of $2 failed with error $ERR! Events:" >&2
        aws cloudformation describe-stack-events --stack-name $2 --max-items 10 \
        | jq --compact-output \
        '.StackEvents | reverse | .[] | {
            id: .LogicalResourceId,
            status: .ResourceStatus,
            why: .ResourceStatusReason
        }' >&2
        exit 1
    fi
}

function deployLambdaCode {
    bucket=$1
    pkgname=$2
    fname=$3

    if [[ -z "$fname" ]]; then
        fname=$pkgname
    fi

    mkdir -p tmp/bin
    mkdir -p tmp/dist

    WORK_DIR=$(pwd)

    # build lambda
    cd ./$pkgname
    GOOS=linux GOARCH=amd64 go build -o ${WORK_DIR}/tmp/bin/${fname}
    cd $WORK_DIR

    # zip handlers
    zip -j tmp/dist/${fname}.zip tmp/bin/${fname}

    # Upload to S3
    aws s3 cp tmp/dist/${fname}.zip s3://$bucket/lambdas/${fname}.zip
}
