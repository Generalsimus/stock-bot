package utils

func Sum[T int | uint | float64 | float32](array []T) T {
	var result T = 0
	for _, v := range array {
		result += v
	}
	return result
}
