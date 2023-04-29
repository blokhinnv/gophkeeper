// package slices contains utils fuctions for slices.
package slices

// Remove returns a new slice of strings that contains all elements of s
// except for the ones that are equal to x.
func Remove(s []string, x string) []string {
	newS := make([]string, 0)
	for _, v := range s {
		if v != x {
			newS = append(newS, v)
		}
	}
	return newS
}
