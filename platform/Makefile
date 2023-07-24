.PHONY: cert servers integration

global:
	(cd global && sam deploy --parameter-overrides DomainName=${DOMAIN_NAME} HostedZoneId=${HOST_ZONE_ID})

servers:
	$(eval CertArn := $(shell aws cloudformation describe-stacks --region us-east-1 --stack-name sam-tsumiki-uploader-global | jq -r '.Stacks[0].Outputs[0].OutputValue'))
	$(eval ApiCertArn := $(shell aws cloudformation describe-stacks --region us-east-1 --stack-name sam-tsumiki-uploader-global | jq -r '.Stacks[0].Outputs[1].OutputValue'))
	cd servers && sam deploy --parameter-overrides CertificateArn=$(CertArn) HostedZoneId=${HOST_ZONE_ID} HostedZoneName=${HOST_ZONE_NAME} DomainName=${DOMAIN_NAME} ApiDomainName=${API_DOMAIN_NAME} ApiCertificateArn=$(ApiCertArn)

integration:
	cd integration && sam deploy
