package utilities

import "golang.org/x/exp/maps"

// GetDistinct values of a slice
func GetDistinct[K comparable](slice []K) []K {
	unique := map[K]struct{}{}

	for _, item := range slice {
		unique[item] = struct{}{}
	}

	return maps.Keys(unique)
}
