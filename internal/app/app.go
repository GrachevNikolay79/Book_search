package app

import (
	"book_search/internal/book"
	"book_search/internal/config"
	"book_search/internal/utils"
	"container/list"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type App struct {
	cfg *config.Config
}

var queue = list.New()
var pushdbActiv = false

var lock sync.Mutex

func NewApp(cfg *config.Config) App {
	return App{
		cfg: cfg,
	}
}

func (a *App) ShutDown() {

}

func (a *App) Run() {

	fmt.Println("===========")
	for _, v := range a.cfg.Paths {
		a.visitAllSubDirs(v)
		break
	}

	if queue.Len() > 0 && !pushdbActiv {
		fmt.Println("------------------------------------")
		pushdbActiv = true
		go a.pushToDB()
	}
}

// Поместим данные из очереди в базу
func (a *App) pushToDB() {
	pushdbActiv = true
	for queue.Len() > 0 {
		b := queue.Front()
		book := b.Value.(book.Book)
		extension := filepath.Ext(book.Name)
		if a.cfg.Ext[extension] {
			fmt.Println(extension, book.ID)
		}

		lock.Lock()
		queue.Remove(b)
		lock.Unlock()
	}
	pushdbActiv = false
}

// обойдем все директории и сформируем очередь из файлов с книжками
func (a *App) visitAllSubDirs(path string) {
	err := filepath.Walk(path,
		func(lpath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				p := strings.Replace(lpath, info.Name(), "", -1)
				p = strings.TrimRight(p, "/")

				sha256, err := utils.CalcFileSHA256(lpath)
				if err != nil {
					sha256 = ""
				}

				lock.Lock()
				queue.PushBack(book.Book{
					ID:   sha256,
					Name: info.Name(),
					Size: info.Size(),
					Path: p})
				lock.Unlock()

				if queue.Len() > 3 && !pushdbActiv {
					pushdbActiv = true
					go a.pushToDB()
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}
