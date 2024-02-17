package radars

//go:generate oapi-codegen -config=./oapi-codegen.yaml openapi/main.yaml
//go:generate sqlc generate --file ./sqlc.yaml
