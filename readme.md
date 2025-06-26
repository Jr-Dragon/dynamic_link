# Dynamic Link

[![Go](https://github.com/Jr-Dragon/dynamic_link/actions/workflows/workflow.yml/badge.svg)](https://github.com/Jr-Dragon/dynamic_link/actions/workflows/workflow.yml)

A simple URL shortener.

## Development

### Getting Start

Dynamic Link is built with [DevContainer](https://containers.dev/).

- VSCode: https://code.visualstudio.com/docs/devcontainers/containers
- Jetbrains IDEs: https://www.jetbrains.com/help/idea/connect-to-devcontainer.html

```shell
(devcontainer) > make init
```

### Project Structure

- `cmd`: The entrypoint of the application.
  - `http_server/main.go`: The HTTP server entrypoint.
- `api`: The API definitions
  - `{servie_name}`: The API scope of a service.
    - `{version}`: Define the API route here.
  - `internal/response`: Useful response components.
- `internal`: The main package of the repository, contains most of the business logic.
  - `biz`: The main business logic
  - `data`: The connections, such as Redis
  - `server`: The server implementations
  - `library`: The common libraries used in the project

### Code Generating

We use [wire](https://github.com/google/wire) to generate the dependency injection code.

```shell
(devcontainer) > make generate
```

### Testing

```shell
(devcontainer) > make test
```