package slice

func ToAny[T any](list []T) []any {
	result := make([]any, len(list))

	for i := range list {
		result[i] = list[i]
	}

	return result
}

func Safe[T any](data []T, index int) (val T) {
	if len(data) > index {
		return data[index]
	}
	return
}
