# Notes

These are notes I am gathering from my experiments.

# 2025-10-11

- I am trialling adding mono-repository tooling via `monogo`, a separate project... main challenge seems to be not the build steps anymore, but linting or testing, since at the moment it runs for all and counts on the cache solely to skip
- I am trialling `opencode` to create a service from scratch based on guidelines (see docs/golang-guidelines.md)
- `core` needs to be re-thought
- Perhaps `go.work` can help on creating mono repositories, but `monogo` must support

# 2024-04-30

- Some work has been done in another repo (argocd-lab), but I got a basic deployment working in a k3d cluster
- The same k3d cluster has a few deployments running, but the setup is not perfect + require some manual setup
- Building with `ko.build` although simple is very constrained... I will need to create a CLI to allow me build this in a more suitable way (see `bob`)

# 2024-04-17

- Added tests using test-containers... Seems to work ok
- Need to understand a bit how to improve the tests when wanting to leverage fx... Seems a bit dirty
- `httpjson` package (needs to be renamed) allows easy step-by-step integration test

# 2024-02-46

- Added all missing radar items + radar endpoints and adapted all payloads do use the data response object
- Error middlewares probably need some love, as they are doing some map[string]string magic
- Probably the support to update quadrants here is not super important for now, as it would require thinking if a rename means deleting all previous ones or not

- Next items: add tests using `test-containers`

# 2024-02-24

- Seems `kin-openapi@v0.122.0` broke the path's parsing, making quite hard to check the API Envelope suggested on the previous iteration to be checked. Need to open a GH issue surrounding it
- Although the API Envelope idea works, it might be a bit confusing for client users... question is: do I care at this stage? I guess this can potentially be used til later notice. More specific modules that do not comply with it can have its own payloads.

# 2024-02-19

- The nesting `oapi-codegen` anonymous struct issue is very annoying... It is not only on $ref, but on any nesting. Considering seriously using `kin` to generate the schema and sort this with Golang code instead.
- I realised that potentially I need to think a bit on how the API design will fit here. Probably JSONApi is an overkill, but something like JSend <https://github.com/omniti-labs/jsend> seems okay-ish.
  - It seems I can declare a common "API Envelope" and use that in my responses. **I just need to check how this looks like for clients / Swagger UI**, as I think it is not a good practice.

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
