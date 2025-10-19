<h1 align="center">
  Go Lab
</h1>

<p align="center">
  🧪 Laboratory for my Go and infrastructure experiments
</p>

## Requirements

- [mise](https://mise.jdx.dev/) (version manager for all CLIs used)

## Installation

1. **Install mise**: Follow the instructions on the [mise website](https://mise.jdx.dev/getting-started.html).
2. **Install project dependencies**:

   ```bash
   mise install
   ```

3. **Install lefthook pre-commit hooks**: `lefthook install`.

## Main Commands

- `make run service=<> cmd=<>`: Runs a specific service's command with hot-reloading (e.g., `make run service=hello-world cmd=api`).
- `make format`: Formats Go code using `gofumpt` and `goimports`, and other files with `prettier`.
- `make lint`: Lints Go code using `golangci-lint`.
- `make test`: Runs all Go tests.
- `make scan`: Scans the codebase for vulnerabilities, misconfigurations, and licenses using Trivy.
- `make docker-all service=<service_name> cmd=<command_name> docker_registry=<> docker_user=<> docker_password=<>`: Performs a full Docker workflow: login, build, sign, and scan images.

## Wish list

- Format and linting should only take changed files
- Linter and tests might need manual saving after running, as the cache action doesn't handle it nicely
