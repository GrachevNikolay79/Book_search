package book

type Book struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
	Path string `json:"path"`
}

func NewBook(name string, size int64, path string) *Book {
	return &Book{
		ID:   "",
		Name: name,
		Size: size,
		Path: path,
	}
}
