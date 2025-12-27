package index

import "sort"

func Rank(ids []int) []int {
	count := make(map[int]int)
	for _, id := range ids {
		count[id]++
	}

	unique := make([]int, 0, len(count))
	for id := range count {
		unique = append(unique, id)
	}

	sort.Slice(unique, func(i, j int) bool {
		return count[unique[i]] > count[unique[j]]
	})

	return unique
}
