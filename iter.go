package main

import (
	"iter"
	"sync"
)

func ParallelMap[T, V any](
	seq iter.Seq[T],
	f func(T) V,
	chunksize int,
) iter.Seq[V] {
	return func(yield func(V) bool) {
		var (
			next, stop = iter.Pull(seq)
			buf        = make(chan chan V, chunksize)

			wg sync.WaitGroup
		)

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer stop()
			defer close(buf)

			for {
				v, ok := next()
				if !ok {
					break
				}

				// spawn a goroutine for each element and put the results in the shared buffer
				ch := make(chan V)
				buf <- ch
				go func() {
					defer close(ch)
					ch <- f(v)
				}()
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			// read each result as they become available and yield to caller
			for ch := range buf {
				if !yield(<-ch) {
					return
				}
			}
		}()

		wg.Wait()
	}
}
