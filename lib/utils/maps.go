package utils

import (
	"fmt"
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

func SortIntMap(_map types.IntMap) types.IntMap {
	result := types.IntMap{}
	keys := SortIntKeys(_map)
	Debug("Sorted Keys")
	Print(fmt.Sprint(keys))
	for _, key := range keys {
		value := _map[key]
		result[key] = value
	}
	return result
}

func SortStrKeys(input map[string]string) []string {
	keys := make([]string, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func SortStrMap(_map map[string]string) map[string]string {
	result := map[string]string{}
	keys := SortStrKeys(_map)
	Debug("Sorted Keys")
	Print(fmt.Sprint(keys))
	for _, key := range keys {
		value := _map[key]
		result[key] = value
	}
	return result
}
