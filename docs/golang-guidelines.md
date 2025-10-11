# General Golang Guidelines

## Standards

- Default ports for external services:
  - `3000`: http
  - `4000`: grpc/connect
  - `9090`: operational (metrics, health, pprof)
- `encoding/json/v2` should be used instead of `encoding/json`
- Do not use `fmt.Println`
- Apps must handle `SIGTERM` and `SIGINT` gracefully
- Packages must never be plural

## Libraries

The below libraries must be used for the respective purposes and alternatives should not be added

- Logging: `log/slog`
- HTTP Client: `github.com/go-resty/resty` (when integrating with HTTP APIs that do not have a Go SDK)

## Best practices

- When implementing a method, it should always support `context.Context` as the first parameter
- Do not use `http.DefaultClient`: always create a custom `*http.Client` with sensible timeouts
- When logging, always pass the `context.Context` via the `*Context` variants in `slog`
- Always inject dependencies via constructors, even for things such as loggers.

## Service structure

- New services must be placed under `services/{name}`
- Not every service needs to have all the below layers, but the structure should be followed as much as possible

```
│   # Main entrypoints per type. The following are some common ones
├── cmd
│   ├── connect
│   ├── grpc
│   ├── http
│   ├── cron
│   ├── consumer
│   └── cli
│
└── internal
    │ # Storage layer (eg: SQL, NoSQL, in-memory, etc)
    │ # NOTE: model differs than dto because it might have storage specific tags/details
    ├── database
    │   ├── repository
    │   ├── model
    │   │ # (Only reqyuired if using SQL)
    │   └── migrations
    │
    │ # GRPC handlers
    ├── grpc
    │
    │ # HTTP handlers and related packages
    ├── http
    │   ├── middleware
    │   │   └── etc.go
    │   └── handler
    │       ├── hello.go
    │       └── todo.go
    │
    │ # Used to convert handlers requests (grpc/http/cli) <=> business domain (service)
    ├── dto
    │   ├── hello.go
    │   └── todo.go
    │
    │ # Service containing the business logic
    └── service
        ├── hello
        └── todo
```
