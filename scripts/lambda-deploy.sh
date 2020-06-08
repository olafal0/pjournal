#!/bin/bash

bucket=$1
pkg=$2

if [ -z "$pkg" ]; then
    echo "Missing pkg variable"
    exit 1
fi

mkdir -p tmp/bin
mkdir -p tmp/dist

WORK_DIR=$(pwd)

pkgname=$(cut -f 2 -d '/' <<< $pkg)

# build lambda
cd ./$pkg
GOOS=linux GOARCH=amd64 go build -o ${WORK_DIR}/tmp/bin/${pkgname}
cd -

# zip handlers
zip -j tmp/dist/${pkgname}.zip tmp/bin/${pkgname}

# Upload to S3
aws s3 cp tmp/dist/${pkgname}.zip s3://$bucket/lambdas/${pkgname}.zip
