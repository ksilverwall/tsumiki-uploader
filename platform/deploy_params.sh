#!/bin/bash

echo $1 $2

local_directory=$1

bucket_name=$(aws cloudformation describe-stacks --region ap-northeast-1 --stack-name sam-tsumiki-uploader-platform-storages | jq -r '.Stacks[0].Outputs[0].OutputValue')

aws s3 sync "$local_directory" "s3://${bucket_name}/" --delete