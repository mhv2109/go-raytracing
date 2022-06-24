package main

import "math"

// TODO: make Point3 and Color distinct types without redeclaring each method
type Point3 = Vec3
type Color = Vec3

type Vec3 struct {
	X, Y, Z float64
}

func (v Vec3) Sum() float64 {
	return v.X + v.Y + v.Z
}

func (v Vec3) Neg() Vec3 {
	return v.MulS(-1)
}

func (v Vec3) Add(o Vec3) Vec3 {
	return Vec3{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}

func (v Vec3) MulS(t float64) Vec3 {
	return Vec3{v.X * t, v.Y * t, v.Z * t}
}

func (v Vec3) Mul(o Vec3) Vec3 {
	return Vec3{v.X * o.X, v.Y * o.Y, v.Z * o.Z}
}

func (v Vec3) DivS(t float64) Vec3 {
	return v.MulS(1 / t)
}

func (v Vec3) Div(o Vec3) Vec3 {
	return Vec3{v.X / o.X, v.Y / o.Y, v.Z / o.Z}
}

func (v Vec3) Len() float64 {
	return math.Sqrt(v.LenSq())
}

func (v Vec3) LenSq() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vec3) Dot(o Vec3) float64 {
	return v.Mul(o).Sum()
}

func (v Vec3) Cross(o Vec3) Vec3 {
	return Vec3{
		v.Y*o.Z - v.Z*o.Y,
		v.Z*o.X - v.X*o.Z,
		v.X*o.Y - v.Y*o.Z,
	}
}

func (v Vec3) Unit() Vec3 {
	return v.DivS(v.Len())
}

// TODO: are these necessary?

func Add(a, b Vec3) Vec3 {
	return a.Add(b)
}

func Sub(a, b Vec3) Vec3 {
	return a.Add(b.Neg())
}

func MulS(a Vec3, t float64) Vec3 {
	return a.MulS(t)
}

func Mul(a, b Vec3) Vec3 {
	return a.Mul(b)
}

func DivS(a Vec3, t float64) Vec3 {
	return a.DivS(t)
}

func Div(a, b Vec3) Vec3 {
	return a.Div(b)
}

func Dot(a, b Vec3) float64 {
	return a.Dot(b)
}

func Cross(a, b Vec3) Vec3 {
	return a.Cross(b)
}

func Unit(a Vec3) Vec3 {
	return a.Unit()
}
