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

func mergeUnique(arr1, arr2 []int) []int {
	// Создание карты для отслеживания уникальных элементов
	uniqueMap := make(map[int]bool)
	for _, num := range arr1 {
		uniqueMap[num] = true
	}
	for _, num := range arr2 {
		uniqueMap[num] = true
	}

	// Создание массива с уникальными элементами
	var mergedArray []int
	for num := range uniqueMap {
		mergedArray = append(mergedArray, num)
	}
	sort.Ints(mergedArray)

	return mergedArray
}

func copyMap[K comparable, V any](original map[K]V) map[K]V {
	cpy := make(map[K]V)
	for key, value := range original {
		cpy[key] = value
	}

	return cpy
}

func copyWithAdd(arr []int, add int) []int {
	res := make([]int, 0, len(arr))
	for _, el := range arr {
		res = append(res, el+add)
	}
	return res
}
