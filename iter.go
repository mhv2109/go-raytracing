package main

import (
	"context"
	"iter"
	"sync"
)

func ParallelMap[T, V any](
	seq iter.Seq[T],
	f func(T) V,
	chunksize int,
) iter.Seq[V] {
	if chunksize <= 0 {
		chunksize = 1
	}

	return func(yield func(V) bool) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		next, stop := iter.Pull(seq)
		defer stop()

		buf := make(chan chan V, chunksize)

		var wg sync.WaitGroup

		// Producer: pulls from seq, starts workers, and feeds their result channels into buf.
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(buf)

			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				v, ok := next()
				if !ok {
					return
				}

				ch := make(chan V, 1)

				select {
				case <-ctx.Done():
					close(ch)
					return
				case buf <- ch:
				}

				wg.Add(1)
				go func(val T, out chan V) {
					defer wg.Done()
					defer close(out)

					select {
					case <-ctx.Done():
						return
					case out <- f(val):
					}
				}(v, ch)
			}
		}()

		// Consumer: reads workers' result channels from buf and yields values downstream.
		for ch := range buf {
			var v V
			ok := true

			select {
			case <-ctx.Done():
				ok = false
			case v, ok = <-ch:
			}

			if !ok {
				continue
			}

			if !yield(v) {
				cancel()
				break
			}
		}

		wg.Wait()
	}
}
