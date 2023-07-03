.DEFAULT_GOAL = build

GEN_FILE_PATH=catch-all/gen/openapi

.PHONY: build clean gen

gen: ${GEN_FILE_PATH}/server.gen.go ${GEN_FILE_PATH}/types.gen.go

build:
	sam build

clean:
	- rm -rf .aws-sam
	- rm -rf catch-all/gen/

${GEN_FILE_PATH}/server.gen.go: openapi.yaml
	mkdir -p catch-all/gen/openapi/
	oapi-codegen -package openapi -generate "server" openapi.yaml > ${GEN_FILE_PATH}/server.gen.go

${GEN_FILE_PATH}/types.gen.go: openapi.yaml
	mkdir -p catch-all/gen/openapi/
	oapi-codegen -package openapi -generate "types" openapi.yaml > ${GEN_FILE_PATH}/types.gen.go
