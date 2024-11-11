package util

import (
	"unique"
)

// Filter filters the given slice based on the given predicate.
func Filter[T any](data []T, predicate func(T) bool) []T {
	filtered := make([]T, 0)
	for _, item := range data {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// Unique returns a new slice containing only the unique elements of the given slice.
func Unique[T comparable](data []T) []T {
	seen := make(map[unique.Handle[T]]struct{})
	uniqueData := make([]T, 0)
	for _, item := range data {
		handle := unique.Make(item)
		if _, ok := seen[handle]; !ok {
			seen[handle] = struct{}{}
			uniqueData = append(uniqueData, item)
		}
	}
	return uniqueData
}

// Map returns a new slice containing the results of applying the given mapper function to each element of the given slice.
func Map[T, U any](data []T, mapper func(T) U) []U {
	mapped := make([]U, len(data))
	for i, item := range data {
		mapped[i] = mapper(item)
	}
	return mapped
}
