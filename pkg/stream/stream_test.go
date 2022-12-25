package stream

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type v[V any] struct {
	value V
}

func TestWrap(t *testing.T) {
	original := []int{1, 2, 3, 4, 5, 6, 7, 8}
	expected := []string{
		"Index: 0",
		"Value: 1",
		"Index: 1",
		"Value: 2",
		"Index: 2",
		"Value: 3",
		"Index: 3",
		"Value: 4",
		"Index: 4",
		"Value: 5",
		"Index: 5",
		"Value: 6",
		"Index: 6",
		"Value: 7",
		"Index: 7",
		"Value: 8",
	}
	s := New(original...)
	newS := Wrap(s, func(i int, value int) []string {
		return []string{
			fmt.Sprintf("Index: %v", i),
			fmt.Sprintf("Value: %v", value),
		}
	})

	assert.Equal(t, expected, newS.Get())
}

func TestJoin(t *testing.T) {
	type city struct {
		name        string
		countryCode string
	}
	type country struct {
		code string
		name string
	}
	type cityView struct {
		name        string
		countryName string
	}
	extractCityJoinKey := func(i int, c city) string {
		return c.countryCode
	}
	extractCountryJoinKey := func(i int, c country) string {
		return c.code
	}
	cities := []city{
		{name: "Gomel", countryCode: "by"},
		{name: "Minsk", countryCode: "by"},
		{name: "London", countryCode: "gb"},
		{name: "Istanbul", countryCode: "tr"},
	}
	countries := []country{
		{name: "Belarus", code: "by"},
		{name: "Turkiye", code: "tr"},
		{name: "China", code: "cn"},
	}
	sCities := New(cities...)
	sCountries := New(countries...)
	actualCityView := Join(sCities, extractCityJoinKey, sCountries, extractCountryJoinKey, func(okV bool, v city, okW bool, w country) []cityView {
		cityName := v.name
		countryName := w.name
		if !okV {
			cityName = "unknown city"
		}
		if !okW {
			countryName = "unknown country"
		}
		return []cityView{{name: cityName, countryName: countryName}}
	}).Get()

	expectedCityView := []cityView{
		{name: "Gomel", countryName: "Belarus"},
		{name: "Minsk", countryName: "Belarus"},
		{name: "London", countryName: "unknown country"},
		{name: "Istanbul", countryName: "Turkiye"},
		{name: "unknown city", countryName: "China"},
	}

	assert.Equal(t, expectedCityView, actualCityView)
}

func TestBasics(t *testing.T) {
	t.Run("string stream", func(t *testing.T) {
		expected := []string{"a", "b", "c", "d", "e", "f", "g"}
		s := New(expected...)
		assert.Equal(t, expected, s.Get())
	})

	t.Run("int stream", func(t *testing.T) {
		expected := []int{1, 2, 3, 4, 5, 6, 7, 8}
		s := New(expected...)
		assert.Equal(t, expected, s.Get())
	})

	t.Run("struct stream", func(t *testing.T) {
		expected := []v[int]{{value: 1}, {value: 2}, {value: 3}, {value: 4}, {value: 5}, {value: 6}, {value: 7}, {value: 8}}
		s := New(expected...)
		assert.Equal(t, expected, s.Get())
	})

	t.Run("sort peek filter and get", func(t *testing.T) {
		original := []int{5, 6, 7, 8, 1, 2, 4, 9, 3}
		s := New(original...)
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
		assert.Equal(t, expectedGet, actualGet)
		assert.Equal(t, expectedPeek, actualPeek)
	})

	t.Run("skip limit get", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		s := New(original...)
		expected := []int{3, 4}
		actual := s.Skip(2).Limit(2).Get()
		assert.Equal(t, expected, actual)
	})

	t.Run("distinct get", func(t *testing.T) {
		original := []int{1, 1, 1, 2, 3, 3, 4, 2, 5, 5, 6, 7, 7, 8, 8, 8, 8, 9, 9}
		s := New(original...)
		expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		actual := s.Distinct(func(i, a, j, b int) bool {
			return a == b
		}).Get()
		assert.Equal(t, expected, actual)
	})

	t.Run("distinct by key get", func(t *testing.T) {
		original := []int{1, 1, 1, 2, 3, 3, 4, 2, 5, 5, 6, 7, 7, 8, 8, 8, 8, 9, 9}
		s := New(original...)
		expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		actual := s.DistinctByKey(func(i, a int) interface{} {
			return a
		}).Get()
		assert.Equal(t, expected, actual)
	})

	t.Run("expand get", func(t *testing.T) {
		original := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
		s := New(original...)
		expected := []string{"0", "a", "1", "b", "2", "c", "3", "d", "4", "e", "5", "f", "6", "g", "7", "h", "8", "i"}
		actual := s.Expand(func(i int, a string) []string {
			return []string{strconv.Itoa(i), a}
		}).Get()
		assert.Equal(t, expected, actual)
	})

	t.Run("last and first", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		s := New(original...)
		first, okFirst := s.First()
		last, okLast := s.Last()
		assert.Equal(t, 1, first)
		assert.Equal(t, true, okFirst)
		assert.Equal(t, 9, last)
		assert.Equal(t, true, okLast)
	})

	t.Run("last by and first by", func(t *testing.T) {
		original := []string{"-", "-", "#", "-", "-", "#", "-", "-", "#", "-"}
		s := New(original...)
		findHash := func(i int, v string) bool {
			return v == "#"
		}
		findStar := func(i int, v string) bool {
			return v == "*"
		}
		firstHashIndex, firstHash, okFirstHash := s.FirstBy(findHash)
		lastHashIndex, lastHash, okLastHash := s.LastBy(findHash)
		firstStarIndex, firstStar, okFirstStar := s.FirstBy(findStar)
		lastStarIndex, lastStar, okLastStar := s.LastBy(findStar)
		assert.Equal(t, 2, firstHashIndex)
		assert.Equal(t, "#", firstHash)
		assert.Equal(t, true, okFirstHash)
		assert.Equal(t, 8, lastHashIndex)
		assert.Equal(t, "#", lastHash)
		assert.Equal(t, true, okLastHash)
		assert.Equal(t, 0, firstStarIndex)
		assert.Equal(t, "", firstStar)
		assert.Equal(t, false, okFirstStar)
		assert.Equal(t, 0, lastStarIndex)
		assert.Equal(t, "", lastStar)
		assert.Equal(t, false, okLastStar)
	})

	t.Run("skip limit count", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		s := New(original...)
		originalCount := s.Count()
		skipCount := s.Skip(2).Count()
		skipLimitCount := s.Skip(2).Limit(2).Count()
		assert.Equal(t, 9, originalCount)
		assert.Equal(t, 7, skipCount)
		assert.Equal(t, 2, skipLimitCount)
	})

	t.Run("all match", func(t *testing.T) {
		originalMixed := []int{1, 2, -3, 4, -5, 6, 7, 8, 9}
		sMixed := New(originalMixed...)
		originalPositive := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		sPositive := New(originalPositive...)
		checkPositive := func(i int, v int) bool {
			return v > 0
		}
		assert.Equal(t, false, sMixed.AllMatch(checkPositive))
		assert.Equal(t, true, sPositive.AllMatch(checkPositive))
	})

	t.Run("any match", func(t *testing.T) {
		originalNegative := []int{-1, -2, -3, -4, -5, -6, -7, -8, -9}
		sNegative := New(originalNegative...)
		originalMixed := []int{1, 2, -3, 4, -5, 6, 7, 8, 9}
		sMixed := New(originalMixed...)
		originalPositive := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		sPositive := New(originalPositive...)
		checkPositive := func(i int, v int) bool {
			return v > 0
		}
		assert.Equal(t, false, sNegative.AnyMatch(checkPositive))
		assert.Equal(t, true, sMixed.AnyMatch(checkPositive))
		assert.Equal(t, true, sPositive.AnyMatch(checkPositive))
	})

	t.Run("none match", func(t *testing.T) {
		originalNegative := []int{-1, -2, -3, -4, -5, -6, -7, -8, -9}
		sNegative := New(originalNegative...)
		originalMixed := []int{1, 2, -3, 4, -5, 6, 7, 8, 9}
		sMixed := New(originalMixed...)
		originalPositive := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		sPositive := New(originalPositive...)
		checkPositive := func(i int, v int) bool {
			return v > 0
		}
		assert.Equal(t, true, sNegative.NoneMatch(checkPositive))
		assert.Equal(t, false, sMixed.NoneMatch(checkPositive))
		assert.Equal(t, false, sPositive.NoneMatch(checkPositive))
	})

	t.Run("for each", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		s := New(original...)
		actual := []int{}
		s.ForEach(func(i, a int) error {
			actual = append(actual, a)
			return nil
		})
		assert.Equal(t, original, actual)
	})

	t.Run("for each async", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		s := New(original...)
		actual := []int{}
		lock := new(sync.Mutex)
		s.ForEachAsync(func(i, a int) error {
			lock.Lock()
			actual = append(actual, a)
			lock.Unlock()
			return nil
		})
		assert.ElementsMatch(t, original, actual)
	})

	t.Run("for each chunk", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		s := New(original...)
		actualChunks := [][]int{}
		s.ForEachChunk(3, func(from, to int, a []int) error {
			actualChunks = append(actualChunks, a)
			return nil
		})
		expectedChunks := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
		assert.Equal(t, expectedChunks, actualChunks)
	})

	t.Run("for each chunk async", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		s := New(original...)

		actualChunks := [][]int{}
		lock := new(sync.Mutex)
		err := s.ForEachChunkAsync(3, func(from, to int, a []int) error {
			lock.Lock()
			actualChunks = append(actualChunks, a)
			lock.Unlock()
			return nil
		})
		if err != nil {
			panic(err)
		}
		expectedChunks := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
		assert.ElementsMatch(t, expectedChunks, actualChunks)
	})

	t.Run("reverse get", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		expected := []int{9, 8, 7, 6, 5, 4, 3, 2, 1}
		s := New(original...)
		assert.Equal(t, expected, s.Reverse().Get())
	})
}
