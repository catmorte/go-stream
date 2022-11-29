package stream

import (
	"sort"

	"golang.org/x/sync/errgroup"
)

type (
	SortFunc[V any]   func(i int, a V, j int, b V) bool
	FilterFunc[V any] func(i int, a V) bool
	DoFunc[V any]     func(i int, a V) error
	PeekFunc[V any]   func(i int, a V)
	EqFunc[V any]     func(i int, a V, j int, b V) bool
	stream[V any]     struct {
		values []V
		chain  []func(*[]V)
	}

	Stream[V any] interface {
		First() (V, bool)
		FirstBy(checkValues FilterFunc[V]) (int, V, bool)
		Last() (V, bool)
		Count() int
		AllMatch(checkValues FilterFunc[V]) bool
		AnyMatch(checkValues FilterFunc[V]) bool
		NoneMatch(checkValues FilterFunc[V]) bool
		ForEach(do DoFunc[V]) error
		ForEachAsync(do DoFunc[V]) error
		Get() []V
		Sort(sortValues SortFunc[V]) Stream[V]
		Filter(checkValues FilterFunc[V]) Stream[V]
		Peek(peekValues PeekFunc[V]) Stream[V]
		Limit(int) Stream[V]
		Skip(int) Stream[V]
		Distinct(compareValues EqFunc[V]) Stream[V]
	}
)

func Map[V any, NV any](s Stream[V], mapValue func(i int, value V) []NV) Stream[NV] {
	values := s.Get()
	newValues := []NV{}
	for i, v := range values {
		newValues = append(newValues, mapValue(i, v)...)
	}
	return stream[NV]{
		values: newValues,
	}
}

func New[S ~[]V, V any](slice S) Stream[V] {
	return stream[V]{
		values: slice,
	}
}

func (p stream[V]) callChain() []V {
	newValues := append([]V{}, p.values...)
	for _, f := range p.chain {
		f(&newValues)
	}
	return newValues
}

func (p stream[V]) First() (V, bool) {
	length := p.Count()
	if length == 0 {
		var defaultValue V
		return defaultValue, false
	}
	return p.values[0], true
}

func (p stream[V]) Last() (V, bool) {
	length := p.Count()
	if length == 0 {
		var defaultValue V
		return defaultValue, false
	}
	return p.values[length-1], true
}

func (p stream[V]) Count() int {
	values := p.callChain()
	return len(values)
}

func (p stream[V]) Peek(peekValues PeekFunc[V]) Stream[V] {
	chain := append(p.chain, func(values *[]V) {
		for i, v := range *values {
			peekValues(i, v)
		}
	})
	return stream[V]{
		values: p.values,
		chain:  chain,
	}
}

func (p stream[V]) Sort(sortValues SortFunc[V]) Stream[V] {
	chain := append(p.chain, func(values *[]V) {
		sort.Slice(*values, func(i, j int) bool {
			return sortValues(i, (*values)[i], j, (*values)[j])
		})
	})
	return stream[V]{
		values: p.values,
		chain:  chain,
	}
}

func (p stream[V]) FirstBy(checkValues FilterFunc[V]) (int, V, bool) {
	values := p.callChain()
	for i, v := range values {
		if checkValues(i, v) {
			return i, v, true
		}
	}
	var defaultValue V
	return 0, defaultValue, false
}

func (p stream[V]) Filter(checkValues FilterFunc[V]) Stream[V] {
	chain := append(p.chain, func(values *[]V) {
		var newValues []V
		for i, v := range *values {
			if checkValues(i, v) {
				newValues = append(newValues, v)
			}
		}
		*values = newValues
	})
	return stream[V]{
		values: p.values,
		chain:  chain,
	}
}

func (p stream[V]) Limit(limit int) Stream[V] {
	chain := append(p.chain, func(values *[]V) {
		length := len(*values)
		if limit >= length {
			limit = length
		}
		*values = (*values)[:limit]
	})
	return stream[V]{
		values: p.values,
		chain:  chain,
	}
}

func (p stream[V]) Skip(skip int) Stream[V] {
	chain := append(p.chain, func(values *[]V) {
		length := len(*values)
		if skip < length {
			*values = (*values)[skip:]
		} else {
			*values = []V{}
		}
	})
	return stream[V]{
		values: p.values,
		chain:  chain,
	}
}

func (p stream[V]) Distinct(eqValues EqFunc[V]) Stream[V] {
	chain := append(p.chain, func(values *[]V) {
		var unique []V
	loop:
		for j, v := range *values {
			for i, u := range unique {
				if eqValues(j, v, i, u) {
					unique[i] = v
					continue loop
				}
			}
			unique = append(unique, v)
		}

		*values = unique
	})
	return stream[V]{
		values: p.values,
		chain:  chain,
	}
}

func (p stream[V]) AllMatch(checkValues FilterFunc[V]) bool {
	values := p.callChain()
	for i, v := range values {
		if !checkValues(i, v) {
			return false
		}
	}
	return true
}

func (p stream[V]) AnyMatch(checkValues FilterFunc[V]) bool {
	values := p.callChain()
	for i, v := range values {
		if checkValues(i, v) {
			return true
		}
	}
	return false
}

func (p stream[V]) NoneMatch(checkValues FilterFunc[V]) bool {
	values := p.callChain()
	for i, v := range values {
		if checkValues(i, v) {
			return false
		}
	}
	return true
}

func (p stream[V]) ForEach(do DoFunc[V]) error {
	values := p.callChain()
	for i, v := range values {
		err := do(i, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p stream[V]) ForEachAsync(do DoFunc[V]) error {
	values := p.callChain()
	g := new(errgroup.Group)
	for i, v := range values {
		i := i
		v := v
		g.Go(func() error { return do(i, v) })
	}
	return g.Wait()
}

func (p stream[V]) Get() []V {
	return p.callChain()
}
