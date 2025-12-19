package main

import "testing"

func TestSumPoint(t *testing.T) {
	p1 := Point3{1, 2, 3}

	if s := p1.Sum(); s != 6 {
		t.Fail()
	}
}

func TestNegPoint(t *testing.T) {
	p1 := Point3{1, 2, 3}

	if p2 := p1.Neg(); p2.X != -1 || p2.Y != -2 || p2.Z != -3 {
		t.Fail()
	}
}

func TestAddPoints(t *testing.T) {
	p1 := Point3{0, 0, 0}
	p2 := Point3{1, 2, 3}

	if p3 := p1.Add(p2); p3.X != 1 || p3.Y != 2 || p3.Z != 3 {
		t.Fail()
	}
}

func TestMulSPoint(t *testing.T) {
	p1 := Point3{1, 2, 3}

	if p3 := p1.MulS(2); p3.X != 2 || p3.Y != 4 || p3.Z != 6 {
		t.Fail()
	}
}

func TestMulPoints(t *testing.T) {
	p1 := Point3{1, 2, 3}
	p2 := Point3{4, 5, 6}

	if p3 := p1.Mul(p2); p3.X != 4 || p3.Y != 10 || p3.Z != 18 {
		t.Fail()
	}
}

func TestDivSPoint(t *testing.T) {
	p1 := Point3{2, 4, 6}

	if p3 := p1.DivS(2); p3.X != 1 || p3.Y != 2 || p3.Z != 3 {
		t.Fail()
	}
}

func TestDivPoint(t *testing.T) {
	p1 := Point3{2, 4, 6}
	p2 := Point3{2, 4, 6}

	if p3 := p1.Div(p2); p3.X != 1 || p3.Y != 1 || p3.Z != 1 {
		t.Fail()
	}
}

func TestVec3DivSZeroPanics(t *testing.T) {
	v := Vec3{2, 4, 6}

	assertPanics(t, "Vec3.DivS: division by zero", func() {
		_ = v.DivS(0)
	})

	assertPanics(t, "Vec3.DivS: division by zero", func() {
		_ = v.DivS(2, 0)
	})
}

func TestVec3DivZeroComponentPanics(t *testing.T) {
	v := Vec3{2, 4, 6}

	assertPanics(t, "Vec3.Div: division by zero", func() {
		_ = v.Div(Vec3{0, 1, 1})
	})

	assertPanics(t, "Vec3.Div: division by zero", func() {
		_ = v.Div(Vec3{1, 1, 1}, Vec3{1, 0, 1})
	})
}

func TestVec3DivSNormal(t *testing.T) {
	v := Vec3{2, 4, 6}

	if got := v.DivS(2); got != (Vec3{1, 2, 3}) {
		t.Fatalf("DivS(2) = %#v, want %#v", got, Vec3{1, 2, 3})
	}

	if got := v.DivS(2, 2); got != (Vec3{0.5, 1, 1.5}) {
		t.Fatalf("DivS(2,2) = %#v, want %#v", got, Vec3{0.5, 1, 1.5})
	}
}

func TestVec3DivNormal(t *testing.T) {
	v := Vec3{2, 4, 6}

	if got := v.Div(Vec3{2, 4, 6}); got != (Vec3{1, 1, 1}) {
		t.Fatalf("Div = %#v, want %#v", got, Vec3{1, 1, 1})
	}

	if got := v.Div(Vec3{2, 4, 6}, Vec3{1, 2, 3}); got != (Vec3{1, 0.5, 0.3333333333333333}) {
		t.Fatalf("Div chain = %#v, want %#v", got, Vec3{1, 0.5, 0.3333333333333333})
	}
}

func assertPanics(t *testing.T, wantMsg string, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			msg, ok := r.(string)
			if !ok {
				return
			}
			if msg != wantMsg {
				t.Fatalf("panic message = %q, want %q", msg, wantMsg)
			}
		} else {
			t.Fatalf("expected panic %q, but function returned normally", wantMsg)
		}
	}()
	f()
}

func TestLenPoint(t *testing.T) {
	p1 := Point3{1, 1, -2}

	if l := p1.Len(); l != 2.449489742783178 {
		t.Fail()
	}
}

func TestSumColor(t *testing.T) {
	c1 := Color{1, 2, 3}

	if s := c1.Sum(); s != 6 {
		t.Fail()
	}
}

func TestNegColor(t *testing.T) {
	c1 := Color{1, 2, 3}

	if c2 := c1.Neg(); c2.X != -1 || c2.Y != -2 || c2.Z != -3 {
		t.Fail()
	}
}

func TestAddColors(t *testing.T) {
	c1 := Color{0, 0, 0}
	c2 := Color{1, 2, 3}

	if c3 := c1.Add(c2); c3.X != 1 || c3.Y != 2 || c3.Z != 3 {
		t.Fail()
	}
}

func TestMulSColor(t *testing.T) {
	c1 := Color{1, 2, 3}

	if c3 := c1.MulS(2); c3.X != 2 || c3.Y != 4 || c3.Z != 6 {
		t.Fail()
	}
}

func TestMulColors(t *testing.T) {
	c1 := Color{1, 2, 3}
	c2 := Color{4, 5, 6}

	if c3 := c1.Mul(c2); c3.X != 4 || c3.Y != 10 || c3.Z != 18 {
		t.Fail()
	}
}

func TestDivSColor(t *testing.T) {
	c1 := Color{2, 4, 6}

	if c3 := c1.DivS(2); c3.X != 1 || c3.Y != 2 || c3.Z != 3 {
		t.Fail()
	}
}

func TestDivColor(t *testing.T) {
	c1 := Color{2, 4, 6}
	c2 := Color{2, 4, 6}

	if c3 := c1.Div(c2); c3.X != 1 || c3.Y != 1 || c3.Z != 1 {
		t.Fail()
	}
}

func TestLenColor(t *testing.T) {
	c1 := Color{1, 1, -2}

	if l := c1.Len(); l != 2.449489742783178 {
		t.Fail()
	}
}
