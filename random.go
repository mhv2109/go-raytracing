package main

import "math/rand"

/*
We see a marginal performance increase by generating random numbers concurrently
vs. calling rand.Float64 directly.
*/

const RandomChSize = 4096

var RandomCh chan float64

func init() {
	RandomCh = make(chan float64, RandomChSize)
	go func() {
		for {
			RandomCh <- rand.Float64()
		}
	}()
}
