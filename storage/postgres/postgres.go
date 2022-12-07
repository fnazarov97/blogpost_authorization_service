package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sqlx.DB
}

func InitDB(auth string) (p *Postgres, err error) {
	db, err := sqlx.Connect("postgres", auth)
	if err != nil {
		return p, err
	}
	return &Postgres{
		DB: db,
	}, nil
}
