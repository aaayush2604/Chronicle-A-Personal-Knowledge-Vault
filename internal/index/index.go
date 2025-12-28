package index

import "fmt"

type Index struct {
	terms map[string][]int
}

func New() *Index {
	return &Index{
		terms: make(map[string][]int),
	}
}

func (i *Index) Search(term string) []int {
	return i.terms[term]
}

func (i *Index) Reset() {
	i.terms = make(map[string][]int)
}

func (i *Index) PrintIndex() {
	for w, c := range i.terms {
		fmt.Printf("[%s]:%d\n", w, c)
	}
}
