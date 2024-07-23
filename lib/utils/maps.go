package utils

import (
	"kato-studio/katoengine/types"
	"sort"
)

func SortIntKeys(input types.IntMap) []int {
	keys := make([]int, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func SortStrKeys(input map[string]string) []string {
	keys := make([]string, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
