package conversion

func Convert(amount float64, rateFrom, rateTo float64) float64 {
	return amount * (rateTo / rateFrom)
}
