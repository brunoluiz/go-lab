# Notes

These are notes I am gathering from my experiments.

# 2024-02-18

- As the handlers start to grow, perhaps is better to keep each method in separate files, making it easier to fuzzy search. It might reach a point though where there are too many files.

# 2024-02-16

- Seems `oapi-codegen` does not like schemas with nested $ref calls... Had to place it back to `openapi/main.yaml` again

# 2024-02-15

- Continued playing with UberFX, but now it is in a better shape
- Potentially should finish the basic CRUD APIs (will struggle a bit with OpenAPI) and/or implement metrics to it (use otel and push into docker-compose)

# 2024-02-14

- Converted the project to use `uber/fx` dependency injection. It is actually not that bad, besides the fact that it needs some `fx.go` files spread around the project

# 2024-02-11

- Organised openapi is partial files: seems there is a problem with partial files since nested items are created with anonymous structs. This makes it slightly annoying to map as output
- First use of the WithTx wrapper

# 2024-02-07

- Trying to refresh my mind on how to `sqlc` and do transactions

# 2024-02-04

- Added `tools.go` to pin dependency versions
- Changing things to Tech Radar, as it will probably be a more palpable example
  - Requires routes and UI
  - New routes would be:
    - `GET /api/v1/radars`: get all radars available
    - `POST /api/v1/radars`: create a new radar
    - `PUT /api/v1/radar/{}`: update existing radar settings (eg: quadradants, name, labels)
    - `PUT /api/v1/radar/{}/items`: update full items list (should be in a separate table)
    - `GET /api/v1/radar/{}`: fetch radar + items
- While renaming openapi schema, I've noticed that using the same name for the request and response body leads to conflicts. I suffixed all the responses with `Out` to sort this.

## 2023-02-14

- Created a `app.Err*`: is going to be useful to standardise errors across application
  - See `./core/app/errors_test.go` to see some interesting tests around how `errors` package can be leveraged for type assertion
  - See `xmiddlewares.ErrorHandler` to see how it can be handled through interface assert
- Adding `x*` in front of packages that are "wrapped" makes life easier for imports
- `sed -i '' 's/apollo/todo/g' **/**/*.go` to rename all packages in Mac
- `sqlc` output generates a massive interface, but better than creating lots of small packages, as otherwise you will have lots of "duplicate" models
  - Problem 1: very big interface to mock
  - Problem 2: `sql.Err*` leaks into application/service layers

## 2023-02-11

- Using `godotenv` and `envconfig` instead of `urfave` or `viper` because this is not a CLI and the configurations will always be loaded through env configs.
- Using `gin` because of the amount of extras it includes when compared to `net/http`.
- Using `oapi-codegen` to play around with OpenAPI3 generated code
  - Problem 1: it needs YAML flattening to work well, otherwise the code is can be very verbose
  - Problem 2: no default values are set, which could be useful for errors, for example
