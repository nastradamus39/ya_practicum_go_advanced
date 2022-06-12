package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	//_ "github.com/go-sql-driver/mysql"
	"github.com/nastradamus39/ya_practicum_go_advanced/internal/types"
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

func (r *DbRepository) Save(url *types.URL) (err error) {
	if r.DB == nil {
		err = errors.New("нет подключения к бд")
		return
	}

	_, err = r.DB.Exec(`BEGIN
		INSERT INTO urls (hash, uuid, url, short_url)
		VALUES ($1, $2, $3, $4);
		EXCEPTION WHEN unique_violation THEN
		-- Ignore duplicate inserts.
		END;`, url.Hash, url.UUID, url.URL, url.ShortURL)

	return err
}

func (r *DbRepository) FindByHash(hash string) (exist bool, url *types.URL, err error) {
	if r.DB == nil {
		exist = false
		url = nil
		err = errors.New("нет подключения к бд")
		return
	}

	rows, err := r.DB.QueryContext(context.Background(), "SELECT u.hash, u.uuid, u.url, u.short_url FROM urls u WHERE u.hash = ? limit ?", hash, 1)
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
	_, err := r.DB.Exec(`CREATE TABLE IF NOT EXISTS urls
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
