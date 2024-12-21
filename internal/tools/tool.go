package tools

// SliceToSet 切片转Set
func SliceToSet[T comparable](items []T) map[T]struct{} {
	set := make(map[T]struct{}, len(items))
	for _, v := range items {
		set[v] = struct{}{}
	}

	return set
}

// InSlice 元素是否存在
func InSlice[T comparable](items []T, value T) bool {
	set := SliceToSet(items)
	if _, ok := set[value]; ok {
		return true
	}
	return false
}
