package main

import (
	"math"
	"math/rand"
)

var (
	_ Material = (*Metal)(nil)
	_ Material = (*Dielectric)(nil)
	_ Material = (*Diffusion)(nil)
)

// Material describes object + ray interactions. See ch 9.
type Material interface {
	Scatter(Ray, HitRecord, *Color, *Ray) bool
}

type material struct {
	albedo Color
}

type Metal struct {
	m    material
	fuzz float64
}

type MetalOpt func(*Metal)

// Fuzz applies randomness to reflection of rays from Metal objects.
// See 9.6.
func Fuzz(fuzz float64) MetalOpt {
	return func(m *Metal) {
		m.fuzz = fuzz
	}
}

func NewMetal(albedo Color, opts ...MetalOpt) Metal {
	m := Metal{m: material{albedo: albedo}}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

// Scatter - see 9.4.
func (m Metal) Scatter(r Ray, hr HitRecord, att *Color, scatt *Ray) (ok bool) {
	reflected := reflect(r.Dir.Unit(), hr.N)
	s := Ray{hr.P, reflected.Add(RandomVec3InUnitSphere().MulS(m.fuzz))} // fuzziness introduced in 9.6
	a := m.m.albedo
	if s.Dir.Dot(hr.N) > 0 {
		*scatt = s
		*att = a
		ok = true
	}
	return
}

type Dielectric struct {
	m  material
	ir float64
}

type DielectricOpt func(*Dielectric)

// IndexOfRefraction defines the refractive index (eta prime) in Snell's
// Law equation. See 10.2.
func IndexOfRefraction(ir float64) DielectricOpt {
	return func(d *Dielectric) {
		d.ir = ir
	}
}

func NewDielectric(albedo Color, opts ...DielectricOpt) Dielectric {
	d := Dielectric{m: material{albedo: albedo}, ir: 1.0}
	for _, opt := range opts {
		opt(&d)
	}
	return d
}

func (d Dielectric) Scatter(r Ray, hr HitRecord, att *Color, scatt *Ray) (ok bool) {
	*att = d.m.albedo

	var ratio float64
	if hr.F {
		ratio = 1.0 / d.ir
	} else {
		ratio = d.ir
	}

	udir := r.Dir.Unit()
	cosT := math.Min(udir.Neg().Dot(hr.N), 1)
	sinT := math.Sqrt(1 - cosT*cosT)

	var dir Vec3
	if ratio*sinT > 1 || d.reflectance(cosT, ratio) > rand.Float64() {
		// cannot refract
		dir = reflect(udir, hr.N)
	} else {
		// can refract
		dir = refract(udir, hr.N, ratio)
	}

	*scatt = Ray{hr.P, dir}
	ok = true
	return
}

// reflectance implements Schlick's approximation for reflectance. See 10.4.
func (d Dielectric) reflectance(cos, ratio float64) (refl float64) {
	r0 := (1 - ratio) / (1 + ratio)
	r0 = r0 * r0
	x := 1 - cos
	x2 := x * x
	refl = r0 + (1-r0)*x2*x2*x
	return
}

type DiffusionType int

const (
	Lambertian DiffusionType = iota
	SimpleDiffusion
)

type Diffusion struct {
	m  material
	dt DiffusionType
}

type DiffusionOpt func(*Diffusion)

func WithDiffusionType(dt DiffusionType) DiffusionOpt {
	return func(d *Diffusion) {
		d.dt = dt
	}
}

func NewDiffusion(albedo Color, opts ...DiffusionOpt) Diffusion {
	d := Diffusion{m: material{albedo: albedo}} // default DiffusionType of 0 value (Lambertian)
	for _, opt := range opts {
		opt(&d)
	}
	return d
}

// Scatter - see 9.3.
func (d Diffusion) Scatter(r Ray, hr HitRecord, att *Color, scatt *Ray) (ok bool) {
	dir := hr.N.Add(d.diffuse(hr))
	if dir.NearZero() {
		dir = hr.N
	}
	*scatt = Ray{hr.P, dir}
	*att = d.m.albedo
	ok = true
	return
}

func (d Diffusion) diffuse(hr HitRecord) (vec Vec3) {
	switch d.dt {
	case Lambertian:
		vec = RandomUnitVec3()
	case SimpleDiffusion:
		r := RandomVec3InUnitSphere()
		if r.Dot(hr.N) < 0 {
			r = r.Neg()
		}
		vec = r
	default:
		panic("unexpected DiffusionType")
	}
	return
}

func reflect(v, n Vec3) Vec3 {
	return v.Sub(n.MulS(2 * v.Dot(n)))
}

func refract(uv, n Vec3, eta float64) Vec3 {
	var (
		cosTheta = math.Min(uv.Neg().Dot(n), 1.0)
		perp     = n.MulS(cosTheta).Add(uv).MulS(eta)
		par      = n.MulS(-math.Sqrt(math.Abs(1 - perp.LenSq())))
	)
	return perp.Add(par)
}
