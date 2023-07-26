#!/bin/bash

echo $1 $2

local_directory=$1

bucket_name=$(aws cloudformation describe-stacks --region ap-northeast-1 --stack-name sam-tsumiki-uploader-platform-storages | jq -r '.Stacks[0].Outputs[0].OutputValue')
distribution=$(aws cloudformation describe-stacks --region ap-northeast-1 --stack-name sam-tsumiki-uploader-platform-web-server | jq -r '.Stacks[0].Outputs[0].OutputValue')

echo ${distribution}
aws s3 sync "$local_directory" "s3://${bucket_name}/" --delete
aws cloudfront create-invalidation --distribution-id ${distribution} --paths /index.html /service-worker.js
