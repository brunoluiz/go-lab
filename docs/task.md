You are Golang product engineering implementing the `services/todo` and this task is specifically to replace the Key Value database with Postgres, through bob as an ORM.

- You must abide by docs/golang-guidelines.md
- You must read the following docs from bob:
  - <https://bob.stephenafamo.com/docs/code-generation/intro>
  - <https://bob.stephenafamo.com/docs/code-generation/usage>
  - <https://bob.stephenafamo.com/docs/code-generation/relationships>
  - <https://bob.stephenafamo.com/docs/code-generation/factories>
  - <https://bob.stephenafamo.com/docs/code-generation/queries>
  - <https://bob.stephenafamo.com/docs/code-generation/enums>
- You must change the repository to use `internal/database/bob` generated code
- Repository models and bob models must be different, so keep them separate
- In theory the other layers (grpc and service) must be kept intact, unless there were breaking changes on the repository layer
- Do not make any assumptions without asking first.
