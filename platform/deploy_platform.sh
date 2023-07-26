#!/bin/bash

set -eu

ORG_DIR=$PWD
cd $(dirname $0)

ENV=$2
TARGET=$1

source .env.$1

if [ "${TARGET}" = "storages" ]; then
	cd storages && sam deploy
elif [ "${TARGET}" = "global" ]; then
	cd global && sam deploy --parameter-overrides DomainName=${WEB_DOMAIN_NAME} HostedZoneId=${HOSTED_ZONE_ID}
elif [ "${TARGET}" = "web-server" ]; then
	$(eval CertArn := $(shell aws cloudformation describe-stacks --region us-east-1 --stack-name sam-tsumiki-uploader-global | jq -r '.Stacks[0].Outputs[0].OutputValue'))
	cd web-server && sam deploy --parameter-overrides CertificateArn=$(CertArn) HostedZoneId=${HOSTED_ZONE_ID} DomainName=${WEB_DOMAIN_NAME}
elif [ "${TARGET}" = "integration-api" ]; then
	cd integration-api && sam deploy --parameter-overrides HostedZoneId=${HOSTED_ZONE_ID} ApiDomainName=${API_DOMAIN_NAME}
else
    echo "${TARGET} is not defined"
fi