package main

type Ray struct {
	Orig, Dir Vec3 // A, b
}

func (r Ray) At(t float64) Vec3 {
	return r.Orig.Add(r.Dir.MulS(t)) // (A + t*b)
}
