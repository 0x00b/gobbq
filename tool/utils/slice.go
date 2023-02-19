package utils

func DeleteSlice[T comparable](a []T, t T) []T {
	idx := 0
	for i, v := range a {
		if v != t {
			if i != idx {
				a[idx] = t
			}
			idx++
		}
	}
	return a[:idx]
}
