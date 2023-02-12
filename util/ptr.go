package util

func Ptr[T any](x T) *T {
	var y T = x
	return &y
}

func Copy[T any](x *T) *T {
	if x == nil {
		return nil
	}

	return Ptr(*x)
}
