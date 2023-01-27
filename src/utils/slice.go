package utils

/*
Returns a slice of all elements in both slices.

If no elements are foud, returns an empty slice.

No matter how many times the elements are found, only one is sent.
*/
func Intersect[K comparable](sl1, sl2 []K) []K {
	return long_intersect(sl1, sl2)
}

func long_intersect[K comparable](sl1, sl2 []K) []K {
	var n int // min(len) > len(intersection)
	n = len(sl1)
	if len(sl2) < n {
		n = len(sl2)
	}

	var arr = make([]K, 0, n)
	var m map[K]int8 = make(map[K]int8)

	for _, el := range sl1 {
		m[el] = 0
	}

	for _, el := range sl2 {
		if _, ok := m[el]; ok {
			arr = append(arr, el)
		}
	}

	return arr
}
