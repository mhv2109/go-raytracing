package main

import "testing"

func TestCameraCoordsCoverageAndBounds(t *testing.T) {
	cam := NewCamera(4, 3, 1, 1, 1,
		Point3{0, 0, 0},
		Point3{0, 0, -1},
		Vec3{0, 1, 0},
		90,
		0,
		1,
	)

	seen := make(map[Coords]bool)
	count := 0

	for c := range cam.coords() {
		if c.i < 0 || c.i >= cam.ImageWidth() {
			t.Fatalf("coord i out of range: %d", c.i)
		}
		if c.j < 0 || c.j >= cam.ImageHeight() {
			t.Fatalf("coord j out of range: %d", c.j)
		}

		if seen[c] {
			t.Fatalf("duplicate coord: %+v", c)
		}
		seen[c] = true
		count++
	}

	want := cam.ImageWidth() * cam.ImageHeight()
	if count != want {
		t.Fatalf("coords count = %d, want %d", count, want)
	}

	// Optional sanity check on ordering: j should start at height-1 and decrease
	prevJ := cam.ImageHeight()
	first := true
	for j := cam.ImageHeight() - 1; j >= 0; j-- {
		for i := 0; i < cam.ImageWidth(); i++ {
			c := Coords{i, j}
			if !seen[c] {
				t.Fatalf("expected coord %+v to be present", c)
			}
			if first {
				first = false
			} else if j > prevJ {
				t.Fatalf("coords not in expected top-to-bottom order: j=%d after j=%d", j, prevJ)
			}
			prevJ = j
		}
	}
}
