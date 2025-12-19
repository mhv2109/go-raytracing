package main

import (
	"testing"
)

func TestParallelMapFullConsumption(t *testing.T) {
	input := func(yield func(int) bool) {
		for i := 0; i < 5; i++ {
			if !yield(i) {
				return
			}
		}
	}

	seq := ParallelMap[int, int](input, func(v int) int { return v * 2 }, 2)

	var got []int
	for v := range seq {
		got = append(got, v)
	}

	want := []int{0, 2, 4, 6, 8}
	if len(got) != len(want) {
		t.Fatalf("len(got)=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d, want %d", i, got[i], want[i])
		}
	}
}

func TestParallelMapEarlyStop(t *testing.T) {
	input := func(yield func(int) bool) {
		for i := 0; i < 1000; i++ {
			if !yield(i) {
				return
			}
		}
	}

	seq := ParallelMap[int, int](input, func(v int) int { return v * 2 }, 8)

	count := 0
	for range seq {
		count++
		if count == 3 {
			break
		}
	}

	if count != 3 {
		t.Fatalf("expected to consume 3 elements, got %d", count)
	}
}
