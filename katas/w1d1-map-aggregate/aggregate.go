package aggregate

type Entry struct {
	key   string
	value int
}

// Aggregate sums the values for each key across all maps in the list.
//
// Example: [{"m1": 20}, {"m2": 30}, {"m1": 20}] -> {"m1": 40, "m2": 30}
func Aggregate(list []map[string]int) map[string]int {
	bucket := make(map[string]int)

	for _, currentMap := range list {
		for key, value := range currentMap {
			bucket[key] += value
		}
	}
	return bucket
}

func AggregateEntries(list []Entry) map[string]int {
	bucket := make(map[string]int)
	for _, e := range list {
		bucket[e.key] += e.value
	}
	return bucket
}
