.PHONY: global storages servers integration-api web-server

storages:
	cd storages && sam deploy

global:
	(cd global && sam deploy --parameter-overrides DomainName=${DOMAIN_NAME} HostedZoneId=${HOST_ZONE_ID})

web-server:
	$(eval CertArn := $(shell aws cloudformation describe-stacks --region us-east-1 --stack-name sam-tsumiki-uploader-global | jq -r '.Stacks[0].Outputs[0].OutputValue'))
	cd web-server && sam deploy --parameter-overrides CertificateArn=$(CertArn) HostedZoneId=${HOST_ZONE_ID} DomainName=${DOMAIN_NAME}

integration-api:
	cd integration-api && sam deploy --parameter-overrides HostedZoneId=${HOST_ZONE_ID} HostedZoneName=${HOST_ZONE_NAME} DomainName=${DOMAIN_NAME} ApiDomainName=${API_DOMAIN_NAME}
