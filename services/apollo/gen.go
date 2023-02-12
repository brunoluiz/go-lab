//go:build +gen

package gen

//go:generate npx openapi-flattener -s openapi.yaml -o /tmp/apollo.yaml
//go:generate oapi-codegen -config=./oapi-codegen.yaml /tmp/apollo.yaml
