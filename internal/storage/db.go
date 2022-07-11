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

type DBRepository struct {
	DB  *sqlx.DB
	cfg *types.Config
}

func NewDBRepository(cfg *types.Config) *DBRepository {
	repo := &DBRepository{
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

func (r *DBRepository) Save(url *types.URL) (err error) {
	if r.DB == nil {
		return fmt.Errorf("%w", shortenerErrors.ErrNoDBConnection)
	}

	rows, err := r.DB.NamedQuery(
		"SELECT * FROM urls u WHERE u.`hash` = :hash LIMIT 1",
		map[string]interface{}{"hash": url.Hash},
	)
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	if rows.Err() != nil {
		return err
	}

	u := types.URL{}
	if rows.Next() && rows.StructScan(&u) == nil { // такой url есть - дубль
		return fmt.Errorf("%w", shortenerErrors.ErrURLConflict)
	}

	// Новый url - сохраняем
	_, err = r.DB.NamedExec(`INSERT INTO urls (hash, uuid, url, short_url)
		VALUES (:hash, :uuid, :url, :short_url)`, url)

	return err
}

func (r *DBRepository) SaveBatch(url []*types.URL) (err error) {
	if r.DB == nil {
		err = errors.New("нет подключения к бд")
		return
	}

	_, err = r.DB.NamedExec(`INSERT INTO urls (hash, uuid, url, short_url)
        VALUES (:hash, :uuid, :url, :short_url)`, url)

	return err
}

func (r *DBRepository) FindByHash(hash string) (exist bool, url *types.URL, err error) {
	if r.DB == nil {
		exist = false
		url = nil
		err = errors.New("нет подключения к бд")
		return
	}

	rows, err := r.DB.NamedQuery(
		"SELECT * FROM urls u WHERE u.`hash` = :hash LIMIT 1",
		map[string]interface{}{"hash": hash},
	)
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	if rows.Err() != nil {
		exist = false
		return
	}

	if rows.Next() {
		url = &types.URL{}
		err = rows.StructScan(url)
		if err != nil {
			exist = false
		}
	}

	return
}

func (r *DBRepository) FindByUUID(uuid string) (exist bool, urls map[string]*types.URL, err error) {
	if r.DB == nil {
		exist = false
		urls = nil
		err = errors.New("нет подключения к бд")
		return
	}

	rows, err := r.DB.NamedQuery(
		"SELECT hash, uuid, url, short_url FROM urls u where u.`uuid` = :uuid",
		map[string]interface{}{"uuid": uuid},
	)
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	if rows.Err() != nil {
		exist = false
		return
	}

	urls = map[string]*types.URL{}
	err = rows.StructScan(&urls)

	return
}

func (r *DBRepository) DeleteByHash(hash []string) (err error) {
	if r.DB == nil {
		err = errors.New("нет подключения к бд")
		return
	}

	_, err = r.DB.QueryContext(context.Background(), `UPDATE urls SET deletedAt = $1 WHERE hash IN ($2)`, time.Now(), hash)

	return
}

func (r *DBRepository) Ping() (err error) {
	if r.DB == nil {
		return errors.New("нет подключения к бд")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.DB.PingContext(ctx)
}

func (r *DBRepository) migrate() {
	_, err := r.DB.Exec(`CREATE TABLE IF NOT EXISTS urls
		(
			hash      varchar(256) not null,
			uuid      varchar(256) not null,
			url       text         not null,
			short_url varchar(256) not null,
    		deleted_at datetime DEFAULT NULL,
			constraint uk
				unique (hash, uuid)
		)`,
	)

	log.Println(err)
}
