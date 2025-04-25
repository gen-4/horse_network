package utils

func Map[T any, R any](arr []T, f func(T) R) []R {
	appliedArr := make([]R, len(arr))
	for i, value := range arr {
		appliedArr[i] = f(value)
	}
	return appliedArr
}

func RemoveFromSliceById[T any](slice []T, equalityFun func(a T) bool) []T {
	var index int = -1
	for i, item := range slice {
		if equalityFun(item) {
			index = i
			break
		}
	}

	if index == -1 {
		return slice
	}

	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}
