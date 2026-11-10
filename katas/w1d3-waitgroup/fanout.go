package fanout

// FetchAll calls fetch for every id concurrently and returns the results
// in the same order as ids.
func FetchAll(ids []int, fetch func(id int) string) []string {
	var results []string
	for _, id := range ids {
		go func() {
			results = append(results, fetch(id))
		}()
	}
	return results
}

// CountParallel starts n goroutines that each add 1 to a shared counter,
// and returns the final count. It must always return n.
func CountParallel(n int) int {
	count := 0
	done := make(chan bool)
	for i := 0; i < n; i++ {
		go func() {
			count++
			done <- true
		}()
	}
	for i := 0; i < n; i++ {
		<-done
	}
	return count
}
