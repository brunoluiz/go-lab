package todo

//go:generate npx openapi-flattener -s openapi.yaml -o /tmp/todo.yaml
//go:generate oapi-codegen -config=./oapi-codegen.yaml /tmp/todo.yaml
//go:generate sqlc generate --file ./sqlc.yaml
