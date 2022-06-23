package book

type Book struct {
	ID     string `json:"id"`
	SHA256 string `json:"SHA256"`
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	Path   string `json:"path"`
	Ext    string `json:"ext"`
}

func NewBook(name string, size int64, path string, ext string) *Book {
	return &Book{
		ID:     "",
		SHA256: "",
		Name:   name,
		Size:   size,
		Path:   path,
		Ext:    ext,
	}
}
