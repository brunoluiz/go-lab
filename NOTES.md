# Notes

These are notes I am gathering from my experiments.

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
