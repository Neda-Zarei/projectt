package util

// Map applies a function f to each element of the input slice s
// and returns a new slice with the results.
func Map[T, U any](s []T, f func(T) U) []U {
	res := make([]U, len(s))
	for i, v := range s {
		res[i] = f(v)
	}
	return res
}
