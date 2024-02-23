package table

// Row represents one line in the table.
type Row[T any] struct{
    Id string
    Data T
    Renderer func(T) []string
}

func (r Row[T]) Render() []string {
    return r.Renderer(r.Data)
}
