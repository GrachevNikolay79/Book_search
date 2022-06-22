package app

import (
	"book_search/internal/book"
	"book_search/internal/config"
	"book_search/internal/utils"
	postgres "book_search/pkg/database/postgresql"
	"container/list"
	"context"
	"fmt"
	"io/fs"
	"log"
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
	var b book.Book
	pushdbActiv = true
	for queue.Len() > 0 {
		bq := queue.Front()
		b = bq.Value.(book.Book)
		//all entries have already been checked for extension
		a.insertBook(&b)

		lock.Lock()
		queue.Remove(bq)
		lock.Unlock()
	}
	log.Println(b.Path)
	//log.Println(b.Name)
	pushdbActiv = false
}

// обойдем все директории и сформируем очередь из файлов с книжками
// go around all the directories and form a queue of entries from files with books
func (a *App) visitAllSubDirs(path string) {
	err := filepath.WalkDir(path,
		func(lpath string, lfile fs.DirEntry, err error) error {
			if !lfile.Type().IsDir() {
				llpath := strings.TrimRight(strings.Replace(lpath, lfile.Name(), "", -1), "/")
				extension := filepath.Ext(lfile.Name())

				//check ext and add to queue
				if a.cfg.Ext[extension] {
					sha256, size, err := utils.CalcFileSHA256(lpath)
					if err != nil {
						log.Panicln("Calc sha256:", err)
						return nil
					}
					if err != nil {
						log.Println(err)
					}
					lock.Lock()
					queue.PushBack(book.Book{
						ID:   sha256,
						Name: lfile.Name(),
						Size: size,
						Path: llpath,
						Ext:  extension})
					lock.Unlock()
				}

				if queue.Len() > 5 && !pushdbActiv {
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
		path varchar(1024),
		ext varchar(8))`

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
			(id,name, length, path, ext) 
			VALUES ($1, $2, $3, $4, $5) 
			ON CONFLICT(id) do UPDATE
			SET name   = excluded.name,
				length = excluded.length,
				path   = excluded.path,
				ext    = excluded.ext
		RETURNING id;`,
		b.ID, b.Name, b.Size, b.Path, b.Ext)

	var id string
	err = row.Scan(&id)
	if err != nil {
		log.Printf("Unable to INSERT: %v\n", err)
		log.Println(b)
	}
}
