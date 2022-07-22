package slices

func PickSpaced[T any](items []T, maxCount int) []T {
	if len(items) <= maxCount {
		return items
	}

	out := make([]T, 0, maxCount)
	nth := len(items) / maxCount
	for i, v := range items {
		if i == 0 || i == len(items)-1 {
			out = append(out, v)
		} else if i%nth == 0 && len(out) < maxCount-1 {
			out = append(out, v)
		}
	}
	return out
}

func Map[T any](items []T, fn func(t T) T) []T {
	out := make([]T, 0, len(items))
	for _, v := range items {
		out = append(out, fn(v))
	}
	return out
}
