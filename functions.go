package transform

// Input type and output type is same.
func F1[T any](f func(T, string) T) f {
	return func(s *shell, p string) *shell {
		return &shell{f(s.v.(T), p)}
	}
}

// Input type and output type is different.
// This can be used for mapping.
func F2[T, U any](f func(T, string) U) f {
	return func(s *shell, p string) *shell {
		return &shell{f(s.v.(T), p)}
	}
}
