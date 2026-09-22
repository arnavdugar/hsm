package orderedset

type OrderedSet[Value comparable] struct {
	set    map[Value]struct{}
	values []Value
}

func Create[Value comparable]() OrderedSet[Value] {
	return OrderedSet[Value]{
		set:    map[Value]struct{}{},
		values: []Value{},
	}
}

func (orderedset *OrderedSet[Value]) Add(value Value) {
	_, ok := orderedset.set[value]
	if ok {
		return
	}
	if orderedset.set == nil {
		orderedset.set = map[Value]struct{}{}
	}
	orderedset.set[value] = struct{}{}
	orderedset.values = append(orderedset.values, value)
}

func (orderedset OrderedSet[Value]) Contains(value Value) bool {
	_, ok := orderedset.set[value]
	return ok
}

func (orderedset OrderedSet[Value]) Values() []Value {
	return orderedset.values
}
