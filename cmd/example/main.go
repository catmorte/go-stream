package main

import (
	"fmt"
	"time"

	"github.com/catmorte/go-streams/pkg/stream"
)

func distinctValues(_, a, _, b int) bool {
	return a == b
}

func sortValues(_, a, _, b int) bool {
	return b < a
}

func peekValues(i, a int) {
	fmt.Printf("%v: %v\n", i, a)
}

func filterValues(i, a int) bool {
	return a%3 == 0
}

func forEachValues(i, a int) error {
	time.Sleep(time.Second * 5)
	fmt.Printf("%v - %v", i, a)
	return nil
}

func main() {
	x := []int{1, 9, 9, 9, 2, 3, 4, 5, 5, 5, 5, 6, 7, 8, 9}
	start := time.Now()
	result := stream.New(x).Sort(sortValues).Distinct(distinctValues).Peek(peekValues).Skip(5).Peek(peekValues).Filter(filterValues).Peek(peekValues).ForEach(forEachValues)
	end := time.Now()
	fmt.Println(result)
	fmt.Println(end.Sub(start).Seconds())
}
