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
