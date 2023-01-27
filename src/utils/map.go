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

func MinKey[K comparable, V any](m map[K]V, keyInf func(K, K) bool) K {
	var minKey *K
	for k, _ := range m {
		if minKey == nil {
			minKey = &k
		} else {
			if keyInf(k, *minKey) {
				minKey = &k
			}
		}
	}
	return *minKey
}

func ReverseMap[K, V comparable](m map[K]V) (map[V]K, error) {
	M := make(map[V]K)
	for k, v := range m {
		if _, ok := M[v]; ok {
			return nil, fmt.Errorf("Could not reverse map. Value %v was present twice", v)
		}
		M[v] = k
	}
	return M, nil
}
