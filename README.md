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

This project uses a `Makefile` to streamline common tasks.

### Test

``` shell
make test
```

### Lint

``` shell
make lint
```

### Benchmark

``` shell
make bench
```

### Build (with PGO)

By default, the build target uses a platform-specific PGO profile named
`$(GOOS)-$(GOARCH).pgo` and outputs a binary to `build/$(GOOS)/$(GOARCH)/rt`.

``` shell
make build
```

You can override `GOOS` and `GOARCH` as usual if you want to cross‑compile:

``` shell
GOOS=linux GOARCH=amd64 make build
GOOS=darwin GOARCH=arm64 make build
```

### Generate a PGO profile

The `profile` target runs the renderer with a small image and writes a
CPU profile named `$(GOOS)-$(GOARCH).pgo` in the project root:

``` shell
make profile
```

### All‑in‑one

Run lint, tests, profile generation, and build in one go:

``` shell
make
```
