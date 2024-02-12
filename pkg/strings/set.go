package strings

// SliceToSet converts a slice of strings to a set of strings
func SliceToSet(s []string) []string {
	set := make(map[string]bool)
	for _, v := range s {
		set[v] = true
	}

	result := make([]string, 0, len(set))
	for k := range set {
		result = append(result, k)
	}

	return result
}

// SliceToMap converts a slice of strings to a map of strings
func SliceToMap(slice []string) map[string]bool {
	m := make(map[string]bool, 0)
	for _, s := range slice {
		m[s] = true
	}

	return m
}
