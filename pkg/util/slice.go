package util

// Generic filter function
func Filter[T any](data []T, predicate func(T) bool) []T {
	filtered := make([]T, 0)
	for _, item := range data {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
