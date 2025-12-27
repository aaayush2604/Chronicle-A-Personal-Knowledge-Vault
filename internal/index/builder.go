package index

import "chronicle/internal/entry"

func (i *Index) Build(entries []entry.KnowledgeEntry) {
	i.Reset()

	for _, e := range entries {
		words := tokenize(e.Content)
		for _, w := range words {
			i.terms[w] = append(i.terms[w], e.ID)
		}
	}
}
