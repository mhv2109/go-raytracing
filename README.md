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

## Project Layout

- `build/` — build output (gitignored)
- `tmp/` — temporary intermediate artifacts: benchmark results, profiling data, test builds, plans, etc. (gitignored)
- `*.pgo` — platform-specific PGO profiles (e.g. `linux-amd64.pgo`)

## Test, Run, and Build

Test:

``` shell
go test ./...
```

Run:

``` shell
go run ./... > <output file> # use --help for all arguments
```

Build (with PGO):

``` shell
go build -o build/rt -pgo linux-amd64.pgo ./...
```
