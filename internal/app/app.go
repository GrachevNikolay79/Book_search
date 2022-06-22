package app

import (
	"book_search/internal/book"
	"book_search/internal/config"
	"book_search/internal/utils"
	postgres "book_search/pkg/database/postgresql"
	"container/list"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type App struct {
	cfg  *config.Config
	pool *pgxpool.Pool
}

var queue = list.New()
var pushdbActiv = false

var lock sync.Mutex

func NewApp(cfg *config.Config) App {
	pool, err := postgres.NewPool(context.Background(), 5, cfg)
	if err != nil {
		log.Fatal(err)
	}

	return App{
		cfg:  cfg,
		pool: pool,
	}
}

func (a *App) ShutDown() {
	a.pool.Close()
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
	for pushdbActiv {
		time.Sleep(1 * time.Second)
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
			a.insertBook(&book)
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

func (a *App) InitDatabase() {
	conn, err := a.pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v\n", err)
	}
	defer conn.Release()

	sql := `
		CREATE TABLE IF NOT EXISTS public.TEMP_BOOK 
		(id varchar(64) primary key, 
		name varchar(256), 
		length bigint, 
		path varchar(1024))`

	row := conn.QueryRow(context.Background(), sql)
	_ = row
}

func (a *App) insertBook(b *book.Book) {
	conn, err := a.pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v\n", err)
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`INSERT INTO TEMP_BOOK 
			(id,name, length, path) 
			VALUES ($1, $2, $3, $4) 
			ON CONFLICT(id) do UPDATE
			SET name   = excluded.name,
				length = excluded.length,
				path   = excluded.path
		RETURNING id;`,
		b.ID, b.Name, b.Size, b.Path)

	var id string
	err = row.Scan(&id)
	if err != nil {
		log.Printf("Unable to INSERT: %v\n", err)
		log.Println(b)
	}
}
