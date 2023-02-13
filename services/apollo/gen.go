package apollo

//go:generate npx openapi-flattener -s openapi.yaml -o /tmp/apollo.yaml
//go:generate oapi-codegen -config=./oapi-codegen.yaml /tmp/apollo.yaml
//go:generate sqlc generate --file ./sqlc.yaml
