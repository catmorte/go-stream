package main

import (
	"fmt"

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

func forEachValues(i, a int) error {
	fmt.Printf("For each %v: %v\n", i, a)
	return nil
}

func main() {
	x := []int{1, 9, 9, 9, 2, 3, 4, 5, 5, 5, 5, 6, 7, 8, 9}
	originalS := stream.New(x).Sort(sortValues)
	distinctS := originalS.Distinct(distinctValues)
	filteredS := distinctS.Peek(peekValues).Skip(5).Peek(peekValues).Filter(filterValues).Peek(peekValues)
	fmt.Println("Filtered:")
	filteredS.ForEach(forEachValues)
	fmt.Println("Distinct:")
	distinctS.ForEach(forEachValues)
	fmt.Println("Original sorted:")
	originalS.ForEach(forEachValues)
}
