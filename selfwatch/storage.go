package selfwatch

import (
	"database/sql"
	"fmt"
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
	log.Print("Loading database ", fname)
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
			strftime('%Y-%m-%d %H:%M:%S', created_at),
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

func (s *WatchStorage) BindRecorder(recorder *Recorder, syncDelay float64) error {
	counter := 0
	last := time.Unix(0, 0)
	var lastWindow int64

	recorder.KeyRelease = func(event Event) {
		counter += 1
		if time.Now().Sub(last).Seconds() > syncDelay || event.Window != lastWindow {
			if counter > 0 {
				log.Println("Syncing keys...", counter)
				s.WriteKeys(counter)
				counter = 0
			}

			last = time.Now()
			lastWindow = event.Window
		}
	}

	return nil
}

type DailyCount struct {
	Day   string
	Count int64
}

func (s *WatchStorage) DailyCounts(days int, newDayHour int) ([]DailyCount, error) {
	rows, err := s.db.Query(`
		select strftime('%Y-%m-%d',
			datetime(datetime(created_at, 'localtime'), ?)
		), sum(nrkeys)
		from keys where created_at > datetime('now', ?)
		group by 1;
	`, fmt.Sprintf("-%v hours", newDayHour), fmt.Sprintf("-%v days", days))

	if err != nil {
		return nil, err
	}

	out := make([]DailyCount, 0)

	defer rows.Close()
	for rows.Next() {
		var day string
		var count int64

		err = rows.Scan(&day, &count)

		if err != nil {
			return nil, err
		}

		out = append(out, DailyCount{
			day, count,
		})
	}

	return out, nil
}
