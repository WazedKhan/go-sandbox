package cleanup

// Resource is anything that must be closed after use: a file, a DB row set,
// an HTTP response body.
type Resource interface {
	Size() int
	Close() error
}

// Opener opens a resource by name.
type Opener func(name string) (Resource, error)

// SumSizes opens each named resource, adds up the sizes, and closes it.
//
// Rules:
//   - At most one resource may be open at any moment.
//   - Every opened resource must be closed, even when something fails.
//   - If opening fails, return the error (wrapped, with the name).
//   - If closing fails, return that error too (wrapped, with the name).
func SumSizes(open Opener, names []string) (int, error) {
	total := 0
	for _, name := range names {
		r, err := open(name)
		if err != nil {
			return 0, err
		}
		defer r.Close()
		total += r.Size()
	}
	return total, nil
}
