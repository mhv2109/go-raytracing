# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go implementation of [Ray Tracing in One Weekend](https://raytracing.github.io/books/RayTracingInOneWeekend.html). This is a single-package raytracer that renders scenes by simulating light physics including reflection, refraction, and diffusion.

## Commands

**Test:**
```shell
go test ./...
```

**Run single test:**
```shell
go test -run TestName
```

**Benchmark:**
```shell
go test -bench=. -benchmem
```

**Run:**
```shell
go run ./... > output.ppm  # outputs PPM image format
go run ./... --help        # see all available arguments
```

**Build (with PGO):**
```shell
go build -o go-rt-cli -pgo default.pgo ./...
```

**CPU Profiling:**
```shell
go run ./... -cpuprofile cpu.prof > output.ppm
```

## Architecture

### Core Data Types

- **Vec3**: Foundation type used as both `Point3` (3D positions) and `Color` (RGB values). Provides vector math operations (add, sub, mul, dot, cross, etc.)
- **Ray**: Defined by origin and direction vectors, represents light paths through the scene
- **RGB**: Final pixel color after gamma correction and clamping

### Physics and Materials (material.go)

The raytracer implements three material types following the `Material` interface:

1. **Diffusion**: Matte/rough surfaces with two diffusion models:
   - Lambertian (default): Physically accurate diffuse reflection
   - SimpleDiffusion: Simpler approximation (use `--simple` flag)

2. **Metal**: Reflective surfaces with optional fuzziness parameter for imperfect mirrors

3. **Dielectric**: Glass-like materials with refraction using Snell's Law and Schlick's approximation for reflectance

Each material implements `Scatter(Ray, HitRecord, *Color, *Ray) bool` to calculate how rays interact with surfaces.

### Scene Objects (hittable.go)

- **Hittable** interface: Any object that can be intersected by a ray
- **Sphere**: Currently the only shape, defined by center point, radius, and material
- **Hittables**: Collection that manages multiple objects and finds closest ray intersection

The `Hit()` method uses quadratic equation solving to determine ray-sphere intersections.

### Rendering Pipeline (camera.go)

1. **Camera** generates rays from a viewpoint through each pixel with configurable:
   - Field of view (vfov)
   - Depth of field (aperture, focus distance)
   - Anti-aliasing (samples per pixel)
   - Ray bounce depth for reflections/refractions

2. **Rendering** is parallelized using `ParallelMap` (iter.go):
   - Generates pixel coordinates as an iterator sequence
   - Maps each coordinate to RGB via parallel workers
   - Worker count controlled by `--jobs` flag (default: 4×CPU count)

3. **Ray tracing** (`rayColor`) iteratively bounces rays through the scene:
   - Stops at max depth or when ray absorbed
   - Accumulates color attenuation from each material scatter
   - Returns background gradient if no objects hit

### Parallelization (iter.go)

`ParallelMap[T, V]` is a generic parallel iterator implementation:
- Pulls from input sequence, spawns worker goroutines up to chunk size
- Workers process items through function `f` and send results to buffered channel
- Maintains ordering via channel sequencing
- Context-aware cancellation propagates downstream

### Scene Construction (main.go)

`randomScene()` creates the demo scene:
- Large ground sphere
- Three featured spheres (glass, diffuse, metal)
- Hundreds of randomly placed/colored small spheres
- Materials distributed: 80% diffuse, 15% metal, 5% glass

## Performance Considerations

- **BVH Acceleration**: Uses Bounding Volume Hierarchy to reduce ray-object intersection tests from O(n) to O(log n). Provides ~4-5x speedup for scenes with 100+ objects.
- **PGO (Profile-Guided Optimization)**: The `default.pgo` file contains profiling data to guide compiler optimizations. Rebuild after significant changes to hot paths.
- **Memory reuse**: `Hittables.Clear()` explicitly nils out slice elements before truncating to allow GC to reclaim memory
- **Ray bounce depth**: Higher `--depth` increases realism but impacts performance exponentially
- **Samples per pixel**: More samples reduce noise but scale linearly with render time
- **Job parallelism**: Default of 4×CPU cores determined through benchmarking

## BVH Implementation

The raytracer uses a Bounding Volume Hierarchy (BVH) to accelerate ray-scene intersection:

- **Structure**: Binary tree where each node has an axis-aligned bounding box (AABB)
- **Construction**: Top-down recursive partitioning using largest-extent axis strategy
- **Traversal**: Early rejection via fast ray-AABB tests before expensive ray-sphere tests
- **Files**: `aabb.go` (AABB + slab method), `bvh.go` (tree construction and traversal)
- **Performance**: Reduces ~100 sphere tests per ray to ~15-25 AABB tests + 1-2 sphere tests

## Output Format

Generates PPM (Portable Pixmap) P3 format:
- ASCII text format with RGB values
- Header: `P3`, width/height, max color value (255)
- Gamma correction (sqrt) applied during RGB conversion
- Redirect stdout to `.ppm` file, viewable in most image viewers
