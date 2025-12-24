# go-raytracing

Go implementation for [Ray Tracing in One Weekend](https://raytracing.github.io/books/RayTracingInOneWeekend.html).

## Requirements

- [Go](https://go.dev/)

## Working with Development Container

### Requirements

- [Docker](https://www.docker.com/)
- [Dev Container CLI](https://github.com/devcontainers/cli)

### Usage

Start a "remote" development container:

``` shell
devcontainer up --workspace-folder ${PROJECT_ROOT}
```

Start a shell:

``` shell
devcontainer exec --workspace-folder ${PROJECT_ROOT} /bin/bash
```

You can then run Test, Run, and Build commands.

## Test, Run, and Build

Test:

``` shell
go test ./...
```

Run:

``` shell
go run ./... > <output file> # use --help for all arguments
```

Build:

``` shell
go build -o go-rt-cli -pgo {mac|linux}.pgo ./...
```
