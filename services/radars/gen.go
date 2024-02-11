package radars

//go:generate oapi-codegen -config=./oapi-codegen.yaml openapi/main.yaml
//go:generate sqlc generate --file ./sqlc.yaml
//go:generate killall gopls # otherwise vim will not detect changes properly
