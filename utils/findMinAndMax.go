package utils

func FindMinAndMax[T int | uint | float64 | int64 | float32](a []T) (min T, max T) {
	min = a[0]
	max = a[0]
	for _, value := range a {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}
