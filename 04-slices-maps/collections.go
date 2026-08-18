package collections

import "strings"

// TODO: implement each function below.
// Run `go test ./04-slices-maps/...` to check your work.

// Reverse returns a NEW slice with the elements of s in reverse order.
// Do not modify s itself.
func Reverse(s []int) []int {
	for i := 0; i < len(s)/2; i++ {

	}
	return nil
}

// Dedup returns a new slice with duplicate values removed, preserving the
// order of first occurrence. Use a map[int]bool as a "set" — Go has no
// built-in Set type like JS's `Set`.
func Dedup(slice []int) []int {
	seen := make(map[int]bool)
	result := []int{}
	for _, value := range slice {
		// 判断 在 map是否存在 存在，不做。不存在 放map，并且放result里面一份
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

// WordCount splits s on whitespace and returns how many times each word
// appears. Remember: a nil map can be read but not written to — initialize
// with make(map[string]int) first.
// Hint: you'll need to add "strings" to an import block, then use
// strings.Fields(s) to split on whitespace.
func WordCount(s string) map[string]int {
	// strings.Fields(s) 切出来是 []string
	words := strings.Fields(s)
	counts := make(map[string]int)
	for _, word := range words {
		counts[word]++
	}
	return counts
}
