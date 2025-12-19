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
	return v.MulS(-1)
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
	return v.Add(Vec3{t, t, t})
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
	return v.AddS(-t)
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
	return v.Sq().Sum()
}

func (v Vec3) Sq() Vec3 {
	return Vec3{v.X * v.X, v.Y * v.Y, v.Z * v.Z}
}

func (v Vec3) Dot(o Vec3) float64 {
	return v.Mul(o).Sum()
}

func (v Vec3) Cross(o Vec3) Vec3 {
	return Vec3{
		v.Y*o.Z - v.Z*o.Y,
		v.Z*o.X - v.X*o.Z,
		v.X*o.Y - v.Y*o.X,
	}
}

func (v Vec3) Unit() Vec3 {
	return v.DivS(v.Len())
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
	f := func() float64 {
		return min + rand.Float64()*(max-min)
	}
	return Vec3{f(), f(), f()}
}

func RandomVec3InUnitSphere() Vec3 {
	for {
		p := RandomVec3(-1, 1)
		if p.LenSq() < 1 {
			return p
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
