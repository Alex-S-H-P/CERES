package utils

import (
	"fmt"
)

func DeepCopy[K comparable, V any](m map[K]V) map[K]V {
	n := make(map[K]V)
	for k, v := range m {
		n[k] = v
	}
	return n
}

/*
Creates a copy of the map1, adds into it the element of map2
*/
func Merge[K comparable, V any](m1, m2 map[K]V) map[K]V {
	m := DeepCopy(m1)

	for k2, v2 := range m2 {
		m[k2] = v2
	}

	return m
}

func MinKey[K comparable, V any](m map[K]V, keyInf func(K, K) bool) K {
	if m == nil {
		return *new(K)
	}

	var minKey *K
	for k := range m {
		if minKey == nil {
			minKey = new(K)
			*minKey = k
			continue
		}
		if keyInf(k, *minKey) {
			minKey = &k
		}

	}
	return *minKey
}

func ReverseMap[K, V comparable](m map[K]V) (map[V]K, error) {
	M := make(map[V]K)
	for k, v := range m {
		if _, ok := M[v]; ok {
			return nil, fmt.Errorf("could not reverse map. Value %v was present twice", v)
		}
		M[v] = k
	}
	return M, nil
}

type number interface {
	int | float64 | float32 | complex128 | complex64 | uint | int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64
}

func SumMap[K comparable, N number](m map[K]N) N {
	var rslt = new(N)
	for _, n := range m {
		*rslt += n
	}
	return *rslt
}
