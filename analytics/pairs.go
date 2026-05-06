package analytics

import (
	"context"
	"math"

	"github.com/ryanrmg/projectx-api"
)

type PairsIndicator struct {
	client *projectx.ProjectXClient
	window int
}

type rolling struct {
	data []float64
	size int
}

func newRolling(n int) *rolling { return &rolling{size: n} }

func (r *rolling) add(v float64) {
	if len(r.data) == r.size {
		r.data = r.data[1:]
	}
	r.data = append(r.data, v)
}

func (r *rolling) full() bool { return len(r.data) == r.size }

func corr(x, y []float64) float64 {
	n := len(x)
	var mx, my float64
	for i := 0; i < n; i++ {
		mx += x[i]
		my += y[i]
	}
	mx /= float64(n)
	my /= float64(n)

	var sx, sy, sxy float64
	for i := 0; i < n; i++ {
		dx := x[i] - mx
		dy := y[i] - my
		sx += dx * dx
		sy += dy * dy
		sxy += dx * dy
	}

	return sxy / math.Sqrt(sx*sy)
}

func (p *PairsIndicator) Correlation(ctx context.Context, symbolId1, symbolId2 string) <-chan float64 {

	out := make(chan float64, 10)

	go func() {
		defer close(out)

		stream := p.client.Realtime.TradesStream()

		var last1, last2 float64
		r1 := newRolling(p.window)
		r2 := newRolling(p.window)

		for {
			select {
			case <-ctx.Done():
				return

			case t := <-stream:

				switch t.SymbolId {

				case symbolId1:
					if last1 != 0 {
						ret := math.Log(t.Price / last1)
						r1.add(ret)
					}
					last1 = t.Price

				case symbolId2:
					if last2 != 0 {
						ret := math.Log(t.Price / last2)
						r2.add(ret)
					}
					last2 = t.Price
				}

				if r1.full() && r2.full() {
					c := corr(r1.data, r2.data)

					select {
					case out <- c:
					default:
					}
				}
			}
		}
	}()

	return out
}
