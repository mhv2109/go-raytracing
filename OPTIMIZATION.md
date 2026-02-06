# Optimization Summary

Go ray tracer performance work. Scene: ~480 spheres, 2560x1440, 100 samples, depth 50.

## Results

| # | Optimization | Status | Impact |
|---|---|---|---|
| 1 | BVH acceleration structure | Done | Massive — O(n) to O(log n) intersection |
| 3 | Non-variadic Vec3 ops | Done | ~1.2% improvement |
| 6 | Unchecked `Unit()` | Done | ~1.7% improvement |
| 4 | Fold `MulS` chains | Done | Within noise |
| 7 | `*Hittables` pointer params | Done | Within noise |
| 8 | Manual `Pow(x,5)` | Done | Within noise |
| 2 | Buffered I/O | Reverted | No improvement — `/dev/null` writes already near-free |
| 5 | Per-goroutine RNG | Reverted | No improvement — Go 1.24 rand already lock-free per-goroutine |
| 9 | Batch pixel formatting | Reverted | No improvement — formatting is <0.08% of total render time |

## Key Takeaways

**BVH dominates everything.** The acceleration structure reduced intersection from linear scan to logarithmic traversal across ~480 objects. All micro-optimizations combined (~3%) are eclipsed by this algorithmic change.

**I/O is irrelevant.** Both buffered writes and batch formatting showed no real improvement — rendering completely dominates wall-clock time.

**Go 1.24 rand is already fast.** Per-goroutine `*rand.Rand` via `sync.Pool` measured within noise of global rand, which is already lock-free.

## Implemented Changes

- **`bvh.go`**: AABB slab-method intersection, `BVHNode` with random-axis midpoint split, `BoundingBox()` on `Hittable` interface
- **`geometry.go`**: Non-variadic `Add`/`Sub`/`MulS`/`DivS`/`Mul`; unchecked `Unit()` via `MulS(1/Len())`
- **`material.go`**: Folded `n.MulS(2).MulS(v.Dot(n))` into single `MulS(2*v.Dot(n))`; manual `x*x*x*x*x` replacing `math.Pow`

## Not Yet Attempted

| Optimization | Impact | Effort |
|---|---|---|
| SIMD / struct-of-arrays Vec3 | Potentially high | High — architecture change |
