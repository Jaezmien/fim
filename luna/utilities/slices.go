package utilities

// A variant of slices.IndexFunc, but allowing a custom starting index
// Setting `start` into a non-negative value will skip the according amount of elements
func IndexFunc[S ~[]E, E any](s S, f func(E) bool, start int) int {
	if start == 0 {
		start = -1
	}

	for i := range s {
		if i < start {
			continue
		}

		if f(s[i]) {
			return i
		}
	}

	return -1
}
