package tools

func SliceToSet[T comparable](items []T) map[T]struct{} {
	set := make(map[T]struct{}, len(items))
	for _, v := range items {
		set[v] = struct{}{}
	}

	return set
}

func InSlice[T comparable](items []T, value T) bool {
	set := SliceToSet(items)
	if _, ok := set[value]; ok {
		return true
	}
	return false
}

func TernaryOperator[T any](exp bool, e1, e2 T) T {
	if exp {
		return e1
	}
	return e2
}

func GetMapDefault[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if value, exists := m[key]; exists {
		return value
	}
	return defaultValue
}
