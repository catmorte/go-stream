package stream

import (
	"sort"

	"golang.org/x/sync/errgroup"
)

type (
	SortFunc[V any]                               func(i int, a V, j int, b V) bool
	FilterFunc[V any]                             func(i int, a V) bool
	DoFunc[V any]                                 func(i int, a V) error
	ExpandFunc[V any]                             func(i int, a V) []V
	PeekFunc[V any]                               func(i int, a V)
	EqFunc[V any]                                 func(i int, a V, j int, b V) bool
	DoChunkFunc[V any]                            func(from, to int, chunk []V) error
	ExtractKeyFunc[V any]                         func(i int, a V) interface{}
	ExtractComparableKeyFunc[V any, K comparable] func(i int, a V) K
	WrapValueFunc[V any, NV any]                  func(i int, value V) []NV
	MergeValuesFunc[V any, W any, VW any]         func(okV bool, v V, okW bool, w W) []VW
	stream[V any]                                 struct {
		values []V
		chain  []func(*[]V)
	}

	Stream[V any] interface {
		First() (V, bool)
		FirstBy(checkValues FilterFunc[V]) (int, V, bool)
		Last() (V, bool)
		LastBy(checkValues FilterFunc[V]) (int, V, bool)
		Count() int
		AllMatch(checkValues FilterFunc[V]) bool
		AnyMatch(checkValues FilterFunc[V]) bool
		NoneMatch(checkValues FilterFunc[V]) bool
		ForEach(do DoFunc[V]) error
		ForEachAsync(do DoFunc[V]) error
		ForEachChunk(chunkSize int, do DoChunkFunc[V]) error
		ForEachChunkAsync(chunkSize int, do DoChunkFunc[V]) error
		Get() []V
		Reverse() Stream[V]
		Sort(sortValues SortFunc[V]) Stream[V]
		Filter(checkValues FilterFunc[V]) Stream[V]
		Peek(peekValues PeekFunc[V]) Stream[V]
		Expand(remap ExpandFunc[V]) Stream[V]
		Limit(int) Stream[V]
		Skip(int) Stream[V]
		Distinct(compareValues EqFunc[V]) Stream[V]
		DistinctByKey(distinctKey ExtractKeyFunc[V]) Stream[V]
	}
)

func Join[V any, W any, VW any, K comparable](sLeft Stream[V], extractKeyLeft ExtractComparableKeyFunc[V, K], sRight Stream[W], extractKeyRight ExtractComparableKeyFunc[W, K], merge MergeValuesFunc[V, W, VW]) Stream[VW] {
	rightMap := map[K]W{}
	rightFound := map[K]struct{}{}
	sRight.ForEach(func(i int, a W) error {
		key := extractKeyRight(i, a)
		rightFound[key] = struct{}{}
		rightMap[key] = a
		return nil
	})

	newValues := []VW{}

	sLeft.ForEach(func(i int, a V) error {
		key := extractKeyLeft(i, a)
		right, okRight := rightMap[key]
		newValues = append(newValues, merge(true, a, okRight, right)...)
		if okRight {
			delete(rightFound, key)
		}
		return nil
	})

	for k := range rightFound {
		var defaultV V
		newValues = append(newValues, merge(false, defaultV, true, rightMap[k])...)
	}
	return newStream(newValues)
}

func Wrap[V any, NV any](s Stream[V], wrapValue WrapValueFunc[V, NV]) Stream[NV] {
	values := s.Get()
	newValues := []NV{}
	for i, v := range values {
		newValues = append(newValues, wrapValue(i, v)...)
	}
	return newStream(newValues)
}

func New[V any](slice ...V) Stream[V] {
	return newStream(slice)
}

func newStream[V any](slice []V) Stream[V] {
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

func (p stream[V]) ForEachChunk(chunkSize int, do DoChunkFunc[V]) error {
	values := p.callChain()
	for i := 0; i < len(values); i += chunkSize {
		end := i + chunkSize
		if end > len(values) {
			end = len(values)
		}
		err := do(i, end, values[i:end])
		if err != nil {
			return err
		}
	}
	return nil
}

func (p stream[V]) ForEachChunkAsync(chunkSize int, do DoChunkFunc[V]) error {
	values := p.callChain()
	g := new(errgroup.Group)
	for i := 0; i < len(values); i += chunkSize {
		start := i
		end := i + chunkSize
		if end > len(values) {
			end = len(values)
		}
		subvalues := values[start:end]
		g.Go(func() error { return do(start, end, subvalues) })
	}
	return g.Wait()
}

func (p stream[V]) Count() int {
	values := p.callChain()
	return len(values)
}

func (p stream[V]) Peek(peekValues PeekFunc[V]) Stream[V] {
	return p.wrapWithFunc(func(values *[]V) {
		for i, v := range *values {
			peekValues(i, v)
		}
	})
}

func (p stream[V]) Expand(expandFunc ExpandFunc[V]) Stream[V] {
	return p.wrapWithFunc(func(values *[]V) {
		var newValues []V
		for i, v := range *values {
			newValues = append(newValues, expandFunc(i, v)...)
		}
		*values = newValues
	})
}

func (p stream[V]) Sort(sortValues SortFunc[V]) Stream[V] {
	return p.wrapWithFunc(func(values *[]V) {
		sort.Slice(*values, func(i, j int) bool {
			return sortValues(i, (*values)[i], j, (*values)[j])
		})
	})
}

func (p stream[V]) LastBy(checkValues FilterFunc[V]) (int, V, bool) {
	values := p.callChain()
	for i := len(values) - 1; i >= 0; i-- {
		v := values[i]
		if checkValues(i, v) {
			return i, v, true
		}
	}
	var defaultValue V
	return 0, defaultValue, false
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

func (p stream[V]) wrapWithFunc(f func(values *[]V)) Stream[V] {
	chain := append(p.chain, f)
	return stream[V]{
		values: p.values,
		chain:  chain,
	}
}

func (p stream[V]) Reverse() Stream[V] {
	return p.wrapWithFunc(func(values *[]V) {
		for i, j := 0, len(*values)-1; i < j; i, j = i+1, j-1 {
			(*values)[i], (*values)[j] = (*values)[j], (*values)[i]
		}
	})
}

func (p stream[V]) Filter(checkValues FilterFunc[V]) Stream[V] {
	return p.wrapWithFunc(func(values *[]V) {
		var newValues []V
		for i, v := range *values {
			if checkValues(i, v) {
				newValues = append(newValues, v)
			}
		}
		*values = newValues
	})
}

func (p stream[V]) Limit(limit int) Stream[V] {
	return p.wrapWithFunc(func(values *[]V) {
		length := len(*values)
		if limit >= length {
			limit = length
		}
		*values = (*values)[:limit]
	})
}

func (p stream[V]) Skip(skip int) Stream[V] {
	return p.wrapWithFunc(func(values *[]V) {
		length := len(*values)
		if skip < length {
			*values = (*values)[skip:]
		} else {
			*values = []V{}
		}
	})
}

func (p stream[V]) DistinctByKey(distinctKey ExtractKeyFunc[V]) Stream[V] {
	return p.wrapWithFunc(func(values *[]V) {
		var unique []V
		keys := map[interface{}]struct{}{}
		for i, v := range *values {
			key := distinctKey(i, v)
			if _, ok := keys[key]; ok {
				continue
			}
			keys[key] = struct{}{}
			unique = append(unique, v)
		}
		*values = unique
	})
}
func (p stream[V]) Distinct(eqValues EqFunc[V]) Stream[V] {
	return p.wrapWithFunc(func(values *[]V) {
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
