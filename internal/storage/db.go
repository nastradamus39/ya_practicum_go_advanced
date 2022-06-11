package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/nastradamus39/ya_practicum_go_advanced/internal/types"
	//_ "github.com/go-sql-driver/mysql"
	//_ "github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

type DbRepository struct {
	DB  *sql.DB
	cfg *types.Config
}

func NewDbRepository(cfg *types.Config) *DbRepository {
	repo := &DbRepository{
		cfg: cfg,
		DB:  nil,
	}

	if cfg.DatabaseDsn != "" {
		db, err := sql.Open("postgres", cfg.DatabaseDsn)
		if err == nil {
			repo.DB = db
			repo.migrate()
		} else {
			fmt.Println(err)
		}
	}

	return repo
}

func (r *DbRepository) Save(url *types.URL) error {
	return nil
}

func (r *DbRepository) FindByHash(hash string) (exist bool, url *types.URL, err error) {
	if r.DB == nil {
		exist = false
		url = nil
		err = errors.New("нет подключения к бд")
		return
	}

	rows, err := r.DB.QueryContext(context.Background(), "SELECT hash, uuid, url, short_url FROM urls where hash = ? limit ?", hash, 1)
	defer rows.Close()

	if err != nil {
		exist = false
		return
	}

	url = &types.URL{}
	for rows.Next() {
		exist = true
		rows.Scan(&url.Hash, &url.UUID, &url.URL, &url.ShortURL)
	}

	return
}

func (r *DbRepository) FindByUUID(uuid string) (exist bool, urls map[string]*types.URL, err error) {
	if r.DB == nil {
		exist = false
		urls = nil
		err = errors.New("нет подключения к бд")
		return
	}

	rows, err := r.DB.QueryContext(context.Background(), "SELECT hash, uuid, url, short_url FROM urls where uuid = ?", uuid)
	defer rows.Close()

	if err != nil {
		exist = false
		return
	}

	urls = map[string]*types.URL{}
	//for rows.Next() {
	//	exist = true
	//	rows.Scan(&url.Hash, &url.UUID, &url.URL, &url.ShortURL)
	//}

	return
}

func (r *DbRepository) Ping() (err error) {
	if r.DB == nil {
		return errors.New("нет подключения к бд")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.DB.PingContext(ctx)
}

func (r *DbRepository) migrate() {
	_, err := r.DB.Exec("create table if not exists urls(hash varchar(256) null, uuid varchar(256) null, url text null, short_url varchar(256) null)")

	fmt.Println(err)
}
