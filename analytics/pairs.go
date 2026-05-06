package analytics

import (
	"context"
	"log"
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

func (p *PairsIndicator) StaticCorrelation(ctx context.Context, req1, req2 projectx.BarHistoryRequest) float64 {

	history1, err := p.client.Markets.History(ctx, req1)
	if err != nil {
		log.Printf("failed history1: %v", err)
		return 0
	}

	history2, err := p.client.Markets.History(ctx, req2)
	if err != nil {
		log.Printf("failed history2: %v", err)
		return 0
	}

	if len(history1) < 2 || len(history2) < 2 {
		return 0
	}

	// Use shortest aligned length
	n := len(history1)
	if len(history2) < n {
		n = len(history2)
	}

	// Build return series
	returns1 := make([]float64, 0, n-1)
	returns2 := make([]float64, 0, n-1)

	for i := 1; i < n; i++ {
		prev1 := history1[i-1].Close
		cur1 := history1[i].Close

		prev2 := history2[i-1].Close
		cur2 := history2[i].Close

		// Prevent divide-by-zero
		if prev1 == 0 || prev2 == 0 {
			continue
		}

		// Simple returns
		r1 := (cur1 - prev1) / prev1
		r2 := (cur2 - prev2) / prev2

		returns1 = append(returns1, r1)
		returns2 = append(returns2, r2)
	}

	if len(returns1) < 2 || len(returns2) < 2 {
		return 0
	}

	return corr(returns1, returns2)
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
