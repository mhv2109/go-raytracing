package main

type MaterialType int

type Material interface {
	Scatter(Ray, HitRecord) (*Color, *Ray)
}

type material struct {
	mt     MaterialType
	albedo Color
}

const (
	// Enumeration of all Materials
	Lambertian MaterialType = iota
	SimpleDiffusion
	Metal
)

func NewMaterial(albedo Color, mt MaterialType) Material {
	return material{mt, albedo}
}

func (m material) Scatter(r Ray, hr HitRecord) (att *Color, scatt *Ray) {
	switch m.mt {
	case Lambertian, SimpleDiffusion:
		dir := hr.N.Add(diffuse(m.mt, hr))
		if dir.NearZero() {
			dir = hr.N
		}
		scatt = &Ray{hr.P, dir}
		att = &m.albedo
	case Metal:
		reflected := reflect(r.Dir.Unit(), hr.N)
		s := Ray{hr.P, reflected}
		a := m.albedo
		if s.Dir.Dot(hr.N) > 0 {
			scatt = &s
			att = &a
		}
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
