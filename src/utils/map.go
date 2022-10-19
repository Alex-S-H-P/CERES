package utils

func DeepCopy[K comparable, V any](m map[K]V) map[K]V{
    n := make(map[K]V)
    for k,v := range m {n[k]=v}
    return n
}

func MinKey[K comparable, V any](m map[K]V, keyInf func (K, K) bool) K {
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
