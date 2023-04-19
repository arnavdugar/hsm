package orderedset

type OrderedSet[V comparable] struct {
	Values []V
	set    map[V]struct{}
}

func Create[V comparable]() OrderedSet[V] {
	return OrderedSet[V]{
		Values: []V{},
		set:    map[V]struct{}{},
	}
}

func (orderedset *OrderedSet[V]) Add(value V) {
	_, ok := orderedset.set[value]
	if ok {
		return
	}
	orderedset.set[value] = struct{}{}
	orderedset.Values = append(orderedset.Values, value)
}
