package stream

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type v[V any] struct {
	value V
}

func TestNew(t *testing.T) {
	t.Run("string stream", func(t *testing.T) {
		expected := []string{"a", "b", "c", "d", "e", "f", "g"}
		s := New(expected)
		assert.ElementsMatch(t, s.Get(), expected)
	})

	t.Run("int stream", func(t *testing.T) {
		expected := []int{1, 2, 3, 4, 5, 6, 7, 8}
		s := New(expected)
		assert.ElementsMatch(t, s.Get(), expected)
	})

	t.Run("struct stream", func(t *testing.T) {
		expected := []v[int]{{value: 1}, {value: 2}, {value: 3}, {value: 4}, {value: 5}, {value: 6}, {value: 7}, {value: 8}}
		s := New(expected)
		assert.ElementsMatch(t, s.Get(), expected)
	})

	t.Run("sort peek filter and get", func(t *testing.T) {
		original := []int{5, 6, 7, 8, 1, 2, 4, 9, 3}
		s := New(original)
		expectedPeek := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		expectedGet := []int{3, 6, 9}
		actualPeek := make([]int, len(original))
		actualGet := s.Sort(func(i, a, j, b int) bool {
			return a < b
		}).Peek(func(i, a int) {
			actualPeek[i] = a
		}).Filter(func(i, a int) bool {
			return a%3 == 0
		}).Get()
		assert.ElementsMatch(t, actualGet, expectedGet)
		assert.ElementsMatch(t, actualPeek, expectedPeek)
	})

	t.Run("skip limit get", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		s := New(original)
		expected := []int{3, 4}
		actual := s.Skip(2).Limit(2).Get()
		assert.ElementsMatch(t, actual, expected)
	})

	t.Run("distinct get", func(t *testing.T) {
		original := []int{1, 1, 1, 2, 3, 3, 4, 2, 5, 5, 6, 7, 7, 8, 8, 8, 8, 9, 9}
		s := New(original)
		expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		actual := s.Distinct(func(i, a, j, b int) bool {
			return a == b
		}).Get()
		assert.ElementsMatch(t, actual, expected)
	})

	t.Run("distinct by key get", func(t *testing.T) {
		original := []int{1, 1, 1, 2, 3, 3, 4, 2, 5, 5, 6, 7, 7, 8, 8, 8, 8, 9, 9}
		s := New(original)
		expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		actual := s.DistinctByKey(func(i, a int) interface{} {
			return a
		}).Get()
		assert.ElementsMatch(t, actual, expected)
	})

	t.Run("expand get", func(t *testing.T) {
		original := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
		s := New(original)
		expected := []string{"0", "a", "1", "b", "2", "c", "3", "d", "4", "e", "5", "f", "6", "g", "7", "h", "8", "i"}
		actual := s.Expand(func(i int, a string) []string {
			return []string{strconv.Itoa(i), a}
		}).Get()
		assert.ElementsMatch(t, actual, expected)
	})

}
