package analytics

func FairValue(c, r, x float64) float64 {
	return c * (1 + r*(x/360))
}
