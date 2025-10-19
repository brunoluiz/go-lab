We will be working together in the `services/todo`.

Ensure you are following docs/golang-guidelines.md

1. You must adjust `lib/app/errors.go` to use `github.com/samber/oops`
2. These errors should be used within the repository layer and replace any current ones
3. Any error middleware in the `handler/*` package should be handled correctly (eg: NotFound shouldnt be a 500, it should be a 404, same for gRPC and so on for other codes)
