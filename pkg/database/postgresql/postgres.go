package postgres

import (
	"book_search/internal/config"
	"book_search/internal/utils"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// type Client interface {
// 	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
// 	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
// 	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
// 	Begin(ctx context.Context) (pgx.Tx, error)
// }

func NewPool(ctx context.Context, maxAttempts int, sc *config.Config) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		sc.PgSQL.Username,
		sc.PgSQL.Password,
		sc.PgSQL.Host,
		sc.PgSQL.Port,
		sc.PgSQL.Database)
	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("error do with tries postgresql")
	}

	return pool, nil
}
