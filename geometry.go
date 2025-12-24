package main

import (
	"math"
	"math/rand"
)

// TODO: make Point3 and Color distinct types without redeclaring each method
type Point3 = Vec3
type Color = Vec3

type Vec3 struct {
	X, Y, Z float64
}

func (v Vec3) Abs() Vec3 {
	return Vec3{math.Abs(v.X), math.Abs(v.Y), math.Abs(v.Z)}
}

func (v Vec3) Sum() float64 {
	return v.X + v.Y + v.Z
}

func (v Vec3) Neg() Vec3 {
	return Vec3{-v.X, -v.Y, -v.Z}
}

func (v Vec3) Add(o Vec3, others ...Vec3) (res Vec3) {
	res = v.add(o)
	for _, o := range others {
		res = res.add(o)
	}
	return
}

func (v Vec3) add(o Vec3) Vec3 {
	return Vec3{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}

func (v Vec3) AddS(t float64, others ...float64) (res Vec3) {
	res = v.addS(t)
	for _, o := range others {
		res = res.addS(o)
	}
	return
}

func (v Vec3) addS(t float64) Vec3 {
	return Vec3{v.X + t, v.Y + t, v.Z + t}
}

func (v Vec3) Sub(o Vec3, others ...Vec3) (res Vec3) {
	res = v.sub(o)
	for _, o := range others {
		res = res.sub(o)
	}
	return
}

func (v Vec3) sub(o Vec3) Vec3 {
	return Vec3{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}

func (v Vec3) SubS(t float64, others ...float64) (res Vec3) {
	res = v.subS(t)
	for _, o := range others {
		res = res.subS(o)
	}
	return
}

func (v Vec3) subS(t float64) Vec3 {
	return Vec3{v.X - t, v.Y - t, v.Z - t}
}

func (v Vec3) MulS(t float64, others ...float64) (res Vec3) {
	res = v.mulS(t)
	for _, o := range others {
		res = res.mulS(o)
	}
	return
}

func (v Vec3) mulS(t float64) Vec3 {
	return Vec3{v.X * t, v.Y * t, v.Z * t}
}

func (v Vec3) Mul(o Vec3, others ...Vec3) (res Vec3) {
	res = v.mul(o)
	for _, o := range others {
		res = res.mul(o)
	}
	return
}

func (v Vec3) mul(o Vec3) Vec3 {
	return Vec3{v.X * o.X, v.Y * o.Y, v.Z * o.Z}
}

func (v Vec3) DivS(t float64, others ...float64) (res Vec3) {
	res = v.divS(t)
	for _, o := range others {
		res = res.divS(o)
	}
	return
}

func (v Vec3) divS(t float64) Vec3 {
	if t == 0 {
		panic("Vec3.DivS: division by zero")
	}
	return Vec3{v.X / t, v.Y / t, v.Z / t}
}

func (v Vec3) Div(o Vec3, others ...Vec3) (res Vec3) {
	res = v.div(o)
	for _, o := range others {
		res = res.div(o)
	}
	return
}

func (v Vec3) div(o Vec3) Vec3 {
	if o.X == 0 || o.Y == 0 || o.Z == 0 {
		panic("Vec3.Div: division by zero")
	}
	return Vec3{v.X / o.X, v.Y / o.Y, v.Z / o.Z}
}

func (v Vec3) Len() float64 {
	return math.Sqrt(v.LenSq())
}

func (v Vec3) LenSq() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vec3) Sq() Vec3 {
	return Vec3{v.X * v.X, v.Y * v.Y, v.Z * v.Z}
}

func (v Vec3) Dot(o Vec3) float64 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

func (v Vec3) Cross(o Vec3) Vec3 {
	return Vec3{
		v.Y*o.Z - v.Z*o.Y,
		v.Z*o.X - v.X*o.Z,
		v.X*o.Y - v.Y*o.X,
	}
}

func (v Vec3) Unit() Vec3 {
	l := v.Len()
	return v.DivS(l)
}

func (v Vec3) NearZero() bool {
	const s = 1e-8
	v = v.Abs()
	return v.X < s && v.Y < s && v.Z < s
}

type RGB struct {
	R, G, B int
}

func (v Vec3) RGB(scale float64) RGB {
	return v.rgb(v.DivS(scale))
}

func (Vec3) rgb(v Vec3) RGB {
	return RGB{
		R: int(255.999 * v.clamp(math.Sqrt(v.X), 0.0, 0.999)),
		G: int(255.999 * v.clamp(math.Sqrt(v.Y), 0.0, 0.999)),
		B: int(255.999 * v.clamp(math.Sqrt(v.Z), 0.0, 0.999)),
	}
}

func (Vec3) clamp(x, min, max float64) float64 {
	if x < min {
		return min
	} else if x > max {
		return max
	}
	return x
}

func RandomVec3(min, max float64) Vec3 {
	r1 := rand.Float64()
	r2 := rand.Float64()
	r3 := rand.Float64()
	scale := max - min
	return Vec3{
		min + r1*scale,
		min + r2*scale,
		min + r3*scale,
	}
}

func RandomVec3InUnitSphere() Vec3 {
	for {
		// Rejection sampling in cube [-1,1]^3, identical distribution to the
		// previous implementation but avoids an intermediate Vec3 before the
		// length-squared check.
		x := -1 + 2*rand.Float64()
		y := -1 + 2*rand.Float64()
		z := -1 + 2*rand.Float64()
		if x*x+y*y+z*z < 1 {
			return Vec3{x, y, z}
		}
	}
}

func RandomUnitVec3() Vec3 {
	return RandomVec3InUnitSphere().Unit()
}

type Ray struct {
	Orig, Dir Vec3 // A, b
}

func (r Ray) At(t float64) Vec3 {
	return r.Orig.Add(r.Dir.MulS(t)) // (A + t*b)
}

type AABB struct {
	Min, Max Point3
}

func NewAABB(a, b Point3) AABB {
	return AABB{Min: a, Max: b}
}

func (b AABB) Hit(r Ray, tMin, tMax float64) bool {
	for a := 0; a < 3; a++ {
		var origin, direction, minA, maxA float64
		switch a {
		case 0:
			origin, direction = r.Orig.X, r.Dir.X
			minA, maxA = b.Min.X, b.Max.X
		case 1:
			origin, direction = r.Orig.Y, r.Dir.Y
			minA, maxA = b.Min.Y, b.Max.Y
		default:
			origin, direction = r.Orig.Z, r.Dir.Z
			minA, maxA = b.Min.Z, b.Max.Z
		}

		invD := 1.0 / direction
		t0 := (minA - origin) * invD
		t1 := (maxA - origin) * invD
		if invD < 0 {
			t0, t1 = t1, t0
		}

		if t0 > tMin {
			tMin = t0
		}
		if t1 < tMax {
			tMax = t1
		}
		if tMax <= tMin {
			return false
		}
	}
	return true
}
