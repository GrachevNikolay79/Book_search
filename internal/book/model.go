package book

type Book struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
	Path string `json:"path"`
	Ext  string `json:"ext"`
}

func NewBook(name string, size int64, path string, ext string) *Book {
	return &Book{
		ID:   "",
		Name: name,
		Size: size,
		Path: path,
		Ext:  ext,
	}
}
