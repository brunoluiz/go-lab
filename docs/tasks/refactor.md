I want the following changes in the services/todo

- Interface definitions should always be close to the callers
- Validate library should always leverage `StructCtx`
- Validate library should be injected in service constructors
- Services should always receive interfaces, never concrete classes
- `dto` should split things into file, such as `task.go` and `list.go`
- If using `errx`, you probably want to use `Wrap` instead of `Wrapf`
- After the closing statement of an `if`, it must have a breakline, except for defer functions after conditionals
