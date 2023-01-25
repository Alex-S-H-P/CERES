package ceres


func PutInto[K any](array []K, index int, element K) []K{
    return append(array[:index],
                append([]K{element},
                    array[index:]...)...)
}
