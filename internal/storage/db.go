package storage

import (
	"context"
	"errors"
	"fmt"
	shortenerErrors "github.com/nastradamus39/ya_practicum_go_advanced/internal/errors"
	"log"
	"time"

	//_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/types"
)

type DbRepository struct {
	DB  *sqlx.DB
	cfg *types.Config
}

func NewDbRepository(cfg *types.Config) *DbRepository {
	repo := &DbRepository{
		cfg: cfg,
		DB:  nil,
	}

	if cfg.DatabaseDsn != "" {
		db, err := sqlx.Open("postgres", cfg.DatabaseDsn) // mysql || postgres
		if err == nil {
			repo.DB = db
			repo.migrate()
		} else {
			log.Println(err)
		}
	}

	return repo
}

func (r *DbRepository) Save(url *types.URL) (err error) {
	if r.DB == nil {
		return fmt.Errorf("%w", shortenerErrors.NoDbConnection)
	}

	rows, err := r.DB.QueryContext(context.Background(), "SELECT * FROM urls where `hash` = $1", url.Hash)

	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		return err
	}

	if rows.Next() {
		return fmt.Errorf("%w", shortenerErrors.UrlConflict)
	}

	_, err = r.DB.NamedExec(`INSERT INTO urls (hash, uuid, url, short_url)
		VALUES (:hash, :uuid, :url, :short_url)`, url)

	return err
}

func (r *DbRepository) SaveBatch(url []*types.URL) (err error) {
	if r.DB == nil {
		err = errors.New("нет подключения к бд")
		return
	}

	_, err = r.DB.NamedExec(`INSERT INTO urls (hash, uuid, url, short_url)
        VALUES (:hash, :uuid, :url, :short_url)`, url)

	return err
}

func (r *DbRepository) FindByHash(hash string) (exist bool, url *types.URL, err error) {
	if r.DB == nil {
		exist = false
		url = nil
		err = errors.New("нет подключения к бд")
		return
	}

	rows, err := r.DB.QueryContext(context.Background(), "SELECT u.hash, u.uuid, u.url, u.short_url FROM urls u WHERE u.hash = $1 limit $2", hash, 1)

	if rows != nil {
		defer rows.Close()
	}

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

	rows, err := r.DB.QueryContext(context.Background(), "SELECT hash, uuid, url, short_url FROM urls where uuid = $1", uuid)
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
	_, err := r.DB.Queryx(`CREATE TABLE IF NOT EXISTS urls
		(
			hash      varchar(256) not null,
			uuid      varchar(256) not null,
			url       text         not null,
			short_url varchar(256) not null,
			constraint uk
				unique (hash, uuid)
		)`,
	)

	fmt.Println(err)
}
