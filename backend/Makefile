.DEFAULT_GOAL = build

.PHONY: build clean gen

gen: catch-all/gen/openapi/openapi.gen.go

build: catch-all/gen/openapi/openapi.gen.go
	sam build

clean:
	- rm -rf .aws-sam
	- rm -rf catch-all/gen/

catch-all/gen/openapi/openapi.gen.go:
	mkdir -p catch-all/gen/openapi/
	oapi-codegen -package openapi -generate "types" openapi.yaml > catch-all/gen/openapi/types.gen.go
	oapi-codegen -package openapi -generate "client" openapi.yaml > catch-all/gen/openapi/client.gen.go
	oapi-codegen -package openapi -generate "server" openapi.yaml > catch-all/gen/openapi/server.gen.go
	oapi-codegen -package openapi -generate "spec" openapi.yaml > catch-all/gen/openapi/spec.gen.go
