package transform

type Transformer interface {
	Transform(a any) error
	Mapping(src, dst any) (err error)
	RegisterTransformer(name string, f f)
	SetTag(tag string) Transformer
	Tag() string
}

// Item for transforming
type I struct {
	Name string
	F    f
}

type shell struct {
	v any
}

type f func(*shell, string) *shell
