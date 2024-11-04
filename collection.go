package main

import (
	"sort"
)

func getSortedIntKeys(m map[int]map[rune]int) []int {
	keys := make([]int, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	return keys
}

func getSortedPositionKeys(m map[Position]Message) []Position {
	keys := make([]Position, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].line < keys[j].line ||
			(keys[i].line == keys[j].line && keys[i].pos < keys[j].pos)
	})

	return keys
}
