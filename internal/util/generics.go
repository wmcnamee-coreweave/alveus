package util

func Ptr[T any](v T) *T {
	return &v
}

func MergeMapsShallow[K comparable, T any](vals ...map[K]T) map[K]T {
	result := make(map[K]T)
	for _, val := range vals {
		for k, v := range val {
			result[k] = v
		}
	}

	return result
}

func CoalesceSlices[T any](vals ...[]T) []T {
	for _, val := range vals {
		if len(val) > 0 {
			return val
		}
	}

	return nil
}

func CoalesceMaps[K comparable, T any](vals ...map[K]T) map[K]T {
	for _, val := range vals {
		if len(val) > 0 {
			return val
		}
	}

	return nil
}

func CoalescePointers[T any](vals ...*T) *T {
	for _, val := range vals {
		if val != nil {
			return val
		}
	}

	return nil
}
