# General Golang Guidelines

## Standards

- `mise` is used for tooling
- Default ports for external services:
  - `3000`: http
  - `4000`: grpc/connect
  - `9090`: operational (metrics, health, pprof)

## Libraries

The below libraries must be used for the respective purposes and alternatives should not be added

- Entrypoint flags and environment variables management: `github.com/alecthomas/kong`
- Logging: `log/slog`
- Errors: `github.com/samber/oops`
- HTTP REST: `github.com/go-resty/resty`
- JSON: `encoding/json/v2`
- YAML: `github.com/yaml/go-yaml`
- Validation: `github.com/go-playground/validator`
- Test (assertions): `github.com/stretchr/testify` (might change in the future, is just quite convenient for now... could be using `matryer/is` in the future)
- Test (containers): `github.com/testcontainers/testcontainers-go`
- Test (mocks): `github.com/uber-go/mock`
- SQL (query): `github.com/stephenafamo/bob`
- SQL (migration): `github.com/golang-migrate/migrate`
- Health checks: `github.com/hellofresh/health-go`

## Best practices

- Apps must handle `SIGTERM` and `SIGINT` gracefully using `signal.NotifyContext`
- Packages must never be plural
- When implementing a method, it should always support `context.Context` as the first parameter
- Always inject dependencies via constructors, even for things such as loggers.
- When input arguments are not used, replace the name with `_`, for example `(_ context.Context, input Bla)`
- When `go build` is trigged, always delete after finishing any tests
- Migrations must not be done on application (human operator will do separate), unless in tests
- Services should always receive interfaces, never concrete classes
- When creating a new service, leverage `lib/app` so the app contains some strong defaults

## Error handling

- Always wrap errors on adapters such as repositories or external clients using `lib/errx.Err*` errors
- The service layer it should never re-wrap calls from adapters, but they must wrap errors from third-party libraries or others happening within its own layer
- All handlers must handle errors coming from the service layer and convert them into appropriate responses (eg: HTTP status codes) using a middleware. An example is `./lib/handler/connectrpc/interceptor`

## Testing

- Table-driven tests must be used when it makes sense. If so, create a `prepare` and `assert` methods for all scenarios to keep a common structure
- Tests must be done so they can run using `t.Parallel()`, with exception to integration tests which should be on best-effort basis
- Mocks should be created as a `mock` package within the package being tested
- Mock generation must be done using `go:generate` and `mockgen`

## Observability

1. When using `lib/app`, it will always expose metrics, debug probes and health checks on port `9090` + give a health check registry

## Flow of data

```
$request --> handler --> $dto --> service --> $model --> repository --> $db_model --> database
```

1. User command via HTTP/GRPC/CLI
2. It gets converted into a DTO (data transfer object) and sent in the correct service
3. Service executes business logic, potentially using repositories to fetch/store data
4. When using repositories, data is converted from DTO to database/model and vice-versa
5. Within the database/repository, it will do the required storage operations and convert to the required structures (eg: bob/model)

## Folder and package structure

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
│   # Scripts related to this specific service/cmds (eg: clients)
├── script
│
└── internal
    │ # Storage layer (eg: SQL, NoSQL, in-memory, etc)
    │ # NOTE: model differs than dto because it might have storage specific tags/details
    ├── database
    │   ├── repository
    │   ├── model # used by calees, should be implementation (kv/sql) agnostic
    │   │ # (Only reqyuired if using SQL)
    │   └── migrations
    │
    │ # Handlers for incoming commands (like ports)
    ├── handler
    │   │
    │   │ # GRPC handlers
    │   ├── grpc
    │   │
    │   │ # HTTP handlers and related packages
    │   └── http
    │       ├── middleware
    │       │   └── etc.go
    │       └── handler
    │           ├── hello.go
    │           └── todo.go
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
