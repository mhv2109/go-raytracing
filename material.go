package main

import "math"

type MaterialType int

// Material describes object + ray interactions. See ch 9.
type Material interface {
	Scatter(Ray, HitRecord) (*Color, *Ray)
}

type material struct {
	mt       MaterialType
	fuzz, ir float64
	albedo   Color
}

const (
	// Enumeration of all Materials
	Lambertian MaterialType = iota
	SimpleDiffusion
	Metal
	Dielectric
)

type MaterialOpt func(*material)

// MetalFuzz applies randomness to reflection of rays from Metal objects.
// See 9.6.
func MetalFuzz(fuzz float64) MaterialOpt {
	return func(m *material) {
		m.fuzz = fuzz
	}
}

// DielectricIndexOfRefraction defines the refractive index (eta prime) in Snell's
// Law equation. See 10.2.
func DielectricIndexOfRefraction(ir float64) MaterialOpt {
	return func(m *material) {
		m.ir = ir
	}
}

func NewMaterial(albedo Color, mt MaterialType, opts ...MaterialOpt) Material {
	m := material{mt: mt, albedo: albedo}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

func (m material) Scatter(r Ray, hr HitRecord) (att *Color, scatt *Ray) {
	switch m.mt {
	case Lambertian, SimpleDiffusion: // see 9.3
		dir := hr.N.Add(diffuse(m.mt, hr))
		if dir.NearZero() {
			dir = hr.N
		}
		scatt = &Ray{hr.P, dir}
		att = &m.albedo
	case Metal: // see 9.4
		reflected := reflect(r.Dir.Unit(), hr.N)
		s := Ray{hr.P, reflected.Add(RandomVec3InUnitSphere().MulS(m.fuzz))} // fuzziness introduced in 9.6
		a := m.albedo
		if s.Dir.Dot(hr.N) > 0 {
			scatt = &s
			att = &a
		}
	case Dielectric:
		att = &m.albedo

		var ratio float64
		if hr.F {
			ratio = 1.0 / m.ir
		} else {
			ratio = m.ir
		}

		udir := r.Dir.Unit()
		cosT := math.Min(udir.Neg().Dot(hr.N), 1)
		sinT := math.Sqrt(1 - cosT*cosT)

		var dir Vec3
		if ratio*sinT > 1 {
			// cannot refract
			dir = reflect(udir, hr.N)
		} else {
			// can refract
			dir = refract(udir, hr.N, ratio)
		}

		scatt = &Ray{hr.P, dir}
	default:
		panic("unexpected MaterialType")
	}
	return
}

func diffuse(mt MaterialType, hr HitRecord) (vec Vec3) {
	switch mt {
	case Lambertian:
		vec = RandomUnitVec3()
	case SimpleDiffusion:
		r := RandomVec3InUnitSphere()
		if r.Dot(hr.N) < 0 {
			r = r.Neg()
		}
		vec = r
	default:
		panic("unexpected MaterialType")
	}
	return
}

func reflect(v, n Vec3) Vec3 {
	return v.Sub(n.MulS(2).MulS(v.Dot(n)))
}

func refract(uv, n Vec3, eta float64) Vec3 {
	var (
		cosTheta = math.Min(uv.Neg().Dot(n), 1.0)
		perp     = n.MulS(cosTheta).Add(uv).MulS(eta)
		par      = n.MulS(-math.Sqrt(math.Abs(1 - perp.LenSq())))
	)
	return perp.Add(par)
}
