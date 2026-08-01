package terminal

import (
	"strconv"
	"strings"
)

func isIntList(l string) (bool, []int) {
	parts := strings.Split(l, ",")
	var temp []int
	var i int
	var err error
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return false, []int{}
		}
		if i, err = strconv.Atoi(p); err != nil {
			return false, []int{}
		}
		temp = append(temp, i)
	}
	return true, temp
}
