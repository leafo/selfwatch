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
	expandedFname, err := expandHomePath(fname)
	if err != nil {
		return nil, err
	}
	fname = expandedFname
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

type HourlyCount struct {
	Hour  string
	Count int64
}

type WeeklyHourlyCount struct {
	Day   string `json:"day"`
	Hour  int    `json:"hour"`
	Count int64  `json:"count"`
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

func (s *WatchStorage) HourlyCounts(hours int, dayOffset int) ([]HourlyCount, error) {
	// dayOffset: 0 = current period, 1 = previous 24h, etc.
	startHours := hours + (dayOffset * 24)

	var rows *sql.Rows
	var err error

	if dayOffset == 0 {
		// Current period: no upper bound needed
		rows, err = s.db.Query(`
			select strftime('%Y-%m-%d %H',
				datetime(created_at, 'localtime')
			), sum(nrkeys)
			from keys
			where created_at > datetime('now', 'localtime', ?)
			group by 1
			order by 1;
		`, fmt.Sprintf("-%v hours", startHours))
	} else {
		// Historical period: need both bounds
		endHours := dayOffset * 24
		rows, err = s.db.Query(`
			select strftime('%Y-%m-%d %H',
				datetime(created_at, 'localtime')
			), sum(nrkeys)
			from keys
			where created_at > datetime('now', 'localtime', ?)
			  and created_at <= datetime('now', 'localtime', ?)
			group by 1
			order by 1;
		`, fmt.Sprintf("-%v hours", startHours), fmt.Sprintf("-%v hours", endHours))
	}

	if err != nil {
		return nil, err
	}

	out := make([]HourlyCount, 0)

	defer rows.Close()
	for rows.Next() {
		var hour string
		var count int64

		err = rows.Scan(&hour, &count)

		if err != nil {
			return nil, err
		}

		out = append(out, HourlyCount{
			hour, count,
		})
	}

	return out, nil
}

func (s *WatchStorage) HourlyCountsForDate(date string) ([]HourlyCount, error) {
	// date format: "2024-12-10"
	rows, err := s.db.Query(`
		select strftime('%Y-%m-%d %H',
			datetime(created_at, 'localtime')
		), sum(nrkeys)
		from keys
		where date(created_at, 'localtime') = ?
		group by 1
		order by 1;
	`, date)

	if err != nil {
		return nil, err
	}

	out := make([]HourlyCount, 0)

	defer rows.Close()
	for rows.Next() {
		var hour string
		var count int64

		err = rows.Scan(&hour, &count)

		if err != nil {
			return nil, err
		}

		out = append(out, HourlyCount{
			hour, count,
		})
	}

	return out, nil
}

func (s *WatchStorage) YearlyCounts(year int, newDayHour int) ([]DailyCount, error) {
	startDate := fmt.Sprintf("%d-01-01", year)
	endDate := fmt.Sprintf("%d-12-31", year)

	rows, err := s.db.Query(`
		select strftime('%Y-%m-%d',
			datetime(datetime(created_at, 'localtime'), ?)
		), sum(nrkeys)
		from keys
		where date(datetime(created_at, 'localtime'), ?) between ? and ?
		group by 1
		order by 1;
	`, fmt.Sprintf("-%v hours", newDayHour), fmt.Sprintf("-%v hours", newDayHour), startDate, endDate)

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

func (s *WatchStorage) WeeklyHourlyGrid() ([]WeeklyHourlyCount, error) {
	rows, err := s.db.Query(`
		select
			strftime('%Y-%m-%d', datetime(created_at, 'localtime')) as day,
			cast(strftime('%H', datetime(created_at, 'localtime')) as integer) as hour,
			sum(nrkeys)
		from keys
		where created_at > datetime('now', '-7 days')
		group by 1, 2
		order by 1, 2;
	`)

	if err != nil {
		return nil, err
	}

	out := make([]WeeklyHourlyCount, 0)

	defer rows.Close()
	for rows.Next() {
		var day string
		var hour int
		var count int64

		err = rows.Scan(&day, &hour, &count)

		if err != nil {
			return nil, err
		}

		out = append(out, WeeklyHourlyCount{
			Day:   day,
			Hour:  hour,
			Count: count,
		})
	}

	return out, nil
}
