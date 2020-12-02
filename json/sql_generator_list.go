package json

type sqlGenListEl struct {
	v    SqlGenerator
	prev *sqlGenListEl
}

// sqlGenList is an immutable singly linked list for generators
type sqlGenList struct {
	tail *sqlGenListEl
}

func (l *sqlGenList) append(v SqlGenerator) *sqlGenList {
	return &sqlGenList{
		tail: &sqlGenListEl{v, l.tail},
	}
}

// iterate over list from the start by using a stack based approach (because we only store a prev pointer)
func (l *sqlGenList) iterate(f func(v SqlGenerator)) {
	sqlGenListIterateEl(l.tail, f)
}

func sqlGenListIterateEl(el *sqlGenListEl, f func(v SqlGenerator)) {
	if el == nil {
		return
	}
	if el.prev != nil {
		sqlGenListIterateEl(el.prev, f)
	}
	f(el.v)
}
