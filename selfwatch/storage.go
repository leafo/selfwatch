package selfwatch

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var keysSchema = `
CREATE TABLE keys (
	id INTEGER NOT NULL,
	created_at DATETIME,
	nrkeys INTEGER,
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

func (s *WatchStorage) SchemaExists() bool {
	rows, err := s.db.Query(`SELECT 1 FROM sqlite_master WHERE type='table' AND name='keys';`)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		return true
	}
	return false
}

func (s *WatchStorage) WriteKeys(keys int) error {
	tx, err := s.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer tx.Commit()

	stmt, err := tx.Prepare("insert into keys(created_at, nrkeys) values(?, ?)")

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(time.Now(), keys)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

type rowTuple struct {
	id         int
	created_at string
	nrkeys     string
}

func (s *WatchStorage) KeyCountsAfterId(id int64) ([]rowTuple, error) {
	rows, err := s.db.Query(`select
			id,
			strftime('%d-%m-%Y %H:%M:%S', created_at),
			nrkeys
		from keys where id > ?
		order by id asc;`, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var out []rowTuple

	for rows.Next() {
		var row rowTuple

		err = rows.Scan(&row.id, &row.created_at, &row.nrkeys)
		if err != nil {
			return nil, err
		}
		out = append(out, row)
	}

	return out, nil
}

func (s *WatchStorage) SerializeRecentKeyCounts(id int64) ([][]interface{}, error) {
	tuples, err := s.KeyCountsAfterId(id)
	if err != nil {
		return nil, err
	}

	var out [][]interface{}

	for _, tup := range tuples {
		var flat []interface{}
		flat = append(flat, tup.id)
		flat = append(flat, tup.created_at)
		flat = append(flat, tup.nrkeys)
		out = append(out, flat)
	}

	return out, nil
}
