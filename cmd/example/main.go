package main

import (
	"fmt"
	"strings"

	"github.com/catmorte/go-streams/pkg/stream"
)

func distinctValues(_, a, _, b int) bool {
	return a == b
}

func sortValues(_, a, _, b int) bool {
	return b < a
}

func peekValues(i, a int) {
	fmt.Printf("Peek %v: %v\n", i, a)
}

func filterValues(i, a int) bool {
	return a%3 == 0
}

func forEachValues[V any](i int, a V) error {
	fmt.Printf("For each %v: %v\n", i, a)
	return nil
}

func main() {
	x := []int{1, 9, 9, 9, 2, 3, 4, 5, 5, 5, 5, 6, 7, 8, 9}
	originalS := stream.New(x).Sort(sortValues)
	distinctS := originalS.Distinct(distinctValues)
	filteredS := distinctS.Peek(peekValues).Skip(5).Peek(peekValues).Filter(filterValues).Peek(peekValues)
	mappedS := stream.Wrap(originalS, func(i int, v int) []string {
		return []string{fmt.Sprintf("index:%v", i), fmt.Sprintf("value:%v", v)}
	})
	fmt.Println("Filtered:")
	filteredS.ForEach(forEachValues[int])
	fmt.Println("Distinct:")
	distinctS.ForEach(forEachValues[int])
	fmt.Println("Original sorted:")
	originalS.ForEach(forEachValues[int])
	fmt.Println("Map:")
	mappedS.ForEach(forEachValues[string])
	fmt.Println(mappedS.FirstBy(func(i int, a string) bool {
		return strings.HasPrefix(a, "index:3")
	}))

	originalS.ForEachChunk(6, func(from, to int, chunk []int) error {
		fmt.Printf("From: %v    To: %v   Chunk: %v\n", from, to, chunk)
		return nil
	} )
	fmt.Println(originalS.Get())
	fmt.Println(originalS.Count())
}
