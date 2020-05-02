package index

type record struct {
	index int
	key   int
	ok    bool
}

type records []record

func (rs records) Len() int {
	return len(rs)
}

func (rs records) Less(i, j int) bool {
	ri := rs[i]
	rj := rs[j]
	return ri.ok && rj.ok && ri.key < rj.key
}

func (rs records) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}
