package selfwatch

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var keysSchema = `
CREATE TABLE keys (
	id INTEGER NOT NULL,
	created_at DATETIME,
	text BLOB NOT NULL,
	started DATETIME NOT NULL,
	nrkeys INTEGER,
	keys BLOB,
	PRIMARY KEY (id)
);
CREATE INDEX ix_keys_nrkeys ON keys (nrkeys);
CREATE INDEX ix_keys_created_at ON keys (created_at);
`

type WatchStorage struct {
	fname string
	db    *sql.DB
}

func NewWatchStorage(fname string) (*WatchStorage, error) {
	db, err := sql.Open("sqlite3", fname)

	if err != nil {
		return nil, err
	}

	return &WatchStorage{
		fname: fname,
		db:    db,
	}, nil
}

func (s *WatchStorage) CreateSchema() error {
	_, err := s.db.Exec(keysSchema)
	if err != nil {
		log.Fatal(err.Error())
	}
	return nil
}
