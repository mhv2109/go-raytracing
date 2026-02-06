# AGENTS.md

## Project

Go implementation of [Ray Tracing in One Weekend](https://raytracing.github.io/books/RayTracingInOneWeekend.html).
Single-package CLI ray tracer outputting PPM images.

- **Language**: Go 1.24
- **Module**: `github.com/mhv2109/RayTracing`
- **Structure**: Flat `main` package — no sub-packages

## Files

| File | Purpose |
|------|---------|
| `main.go` | CLI entry point, flag parsing, PPM output, profiling |
| `geometry.go` | `Vec3`, `Point3`, `Color`, `Ray`, `RGB`, random vector utils |
| `hittable.go` | `Hittable` interface, `HitRecord`, `Sphere`, `Hittables` |
| `material.go` | `Material` interface, `Metal`, `Dielectric`, `Diffusion` |
| `camera.go` | `Camera`, rendering pipeline (`Render` returns `iter.Seq[RGB]`) |
| `iter.go` | `ParallelMap` generic concurrent mapping utility |
| `benchmark_test.go` | `BenchmarkRender` |

## CLI Flags

`-width`, `-height`, `-samples`, `-depth`, `-jobs`, `-simpleDiff`, `-cpuprofile`, `-o`

## Commands

```bash
# Test
go test ./...

# Lint
golangci-lint run

# Format
gofmt -w .

# Build (with PGO)
go build -pgo linux-amd64.pgo -o build/rt ./...

# Run
go run ./... -o output.ppm
go run ./... -width 400 -height 225 -samples 100 -depth 50 -jobs 8 -o out.ppm

# Profile
go run ./... -cpuprofile cpu.prof -o /dev/null

## PGO profiles and GOOS/GOARCH

PGO profiles are platform-specific and named by target OS/architecture:

- `linux-amd64.pgo`
- `darwin-arm64.pgo`

**Guidelines:**

- Use the profile that matches your build target (`GOOS`, `GOARCH`).
- If no matching profile exists, build without `-pgo`.

### Examples

```bash
# Linux / amd64
GOOS=linux GOARCH=amd64 go build \
  -pgo linux-amd64.pgo \
  -o build/rt-linux-amd64 \
  ./...

# macOS / arm64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build \
  -pgo darwin-arm64.pgo \
  -o build/rt-darwin-arm64 \
  ./...

# If there is no profile for the target, skip -pgo
GOOS=linux GOARCH=amd64 go build -o build/rt-linux-amd64 ./...
```

## PGO profiles and GOOS/GOARCH

PGO profiles are platform-specific and named by target OS/architecture:

- `linux-amd64.pgo`
- `darwin-arm64.pgo`

**Guidelines:**

- Use the profile that matches your build target (`GOOS`, `GOARCH`).
- If no matching profile exists, build without `-pgo`.

### Examples

```bash
# Linux / amd64
GOOS=linux GOARCH=amd64 go build \
  -pgo linux-amd64.pgo \
  -o build/rt-linux-amd64 \
  ./...

# macOS / arm64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build \
  -pgo darwin-arm64.pgo \
  -o build/rt-darwin-arm64 \
  ./...

# If there is no profile for the target, skip -pgo
GOOS=linux GOARCH=amd64 go build -o build/rt-linux-amd64 ./...
```
```

## Directories

- `build/` — build output (gitignored)
- `tmp/` — temporary intermediate artifacts: benchmark results, profiling data, test builds, plans, etc. (gitignored, contents excluded except `.gitignore`)

## Task Completion

After any code change, run:

1. `gofmt -w .`
2. `golangci-lint run`
3. `go test ./...`
4. `go mod tidy` (if dependencies changed)

## Style

- Standard Go conventions, `gofmt` formatted
- Exported: PascalCase. Unexported: camelCase
- Value receivers for immutable types (`Vec3`, `Ray`), pointer receivers for mutable ones (`Hittables`)
- Functional options pattern for material constructors (`MetalOpt`, `DielectricOpt`, `DiffusionOpt`)
- `DiffusionType` enum: `Lambertian`, `SimpleDiffusion`
- Table-driven tests, benchmarks in `benchmark_test.go`
- Iterator-based rendering with `iter.Seq[RGB]`
- `ParallelMap[I, O]` for concurrent work distribution
- PGO profile: `linux-amd64.pgo`
