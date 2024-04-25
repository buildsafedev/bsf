package strings

type parse func(string) string

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

// PreferNewSliceElements returns a  slice with elements from new slice that are not in existing slice
// Note: existing elements will not be removed over new ones.
func PreferNewSliceElements(existing []string, new []string, parseFunc parse) []string {
	newEleMap := SliceToMapParse(new, parseFunc)

	// this is the final slice we'll return
	finalEle := make([]string, 0, len(new)+len(existing))

	for _, v := range existing {
		if v == "" {
			continue
		}
		// do not add common elements from existing slice.
		if !newEleMap[parseFunc(v)] {
			finalEle = append(finalEle, v)
		}
	}

	for _, v := range new {
		if v == "" {
			continue
		}
		finalEle = append(finalEle, v)
	}

	return finalEle
}

// SliceToMapParse converts a slice of strings to a map of strings, keys are applied with parseFunc.
func SliceToMapParse(slice []string, parseFunc parse) map[string]bool {
	m := make(map[string]bool, 0)
	for _, s := range slice {
		m[parseFunc(s)] = true
	}

	return m
}

// SliceToMap converts a slice of strings to a map of strings
func SliceToMap(slice []string) map[string]bool {
	m := make(map[string]bool, 0)
	for _, s := range slice {
		m[s] = true
	}

	return m
}
