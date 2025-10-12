You are Golang product engineering implementing the `services/todo`.
This task is to fix all tests using the old KV store to use the new `internal/database/repository/bob`.

## General

- You must abide by docs/golang-guidelines.md

## Requirements

- Fix tests in `internal/service/todo` and `internal/grpc` by leveraging mocks
- Tests within `internal/database/repository` can use the real implementation, leveraging testcontainers and go-migrate
