package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/leafo/selfwatch/selfwatch"
)

var (
	configFname string
	debugOutput bool
)

func init() {
	flag.StringVar(&configFname, "config", selfwatch.DefaultConfigFname, "Path to json config file")
	flag.BoolVar(&debugOutput, "dump", false, "Print extra debug information")
}

func main() {
	flag.Parse()
	config := selfwatch.LoadConfig(configFname)

	storage, err := selfwatch.NewWatchStorage(config.DbName)
	if err != nil {
		log.Fatal(err.Error())
	}

	if !storage.SchemaExists() {
		storage.CreateSchema()
	}

	command := flag.Arg(0)
	if command == "" {
		command = "start"
	}

	switch command {
	case "summary":
		out, err := storage.DailyCounts(7, config.NewDayHour)
		if err != nil {
			log.Fatal(err.Error())
		}

		for _, row := range out {
			fmt.Println(row.Day, "\t", row.Count)
		}

	case "status":
		out, err := storage.DailyCounts(7, config.NewDayHour)
		if err != nil {
			log.Fatal(err.Error())
		}

		dateKey := func(t time.Time) string {
			return fmt.Sprintf("%v-%v-%v",
				t.Year(),
				fmt.Sprintf("%02d", t.Month()),
				fmt.Sprintf("%02d", t.Day()))
		}

		byDay := map[string]int64{}
		for _, row := range out {
			byDay[row.Day] = row.Count
		}

		today := time.Now().Add(time.Hour * time.Duration(-config.NewDayHour))
		yesterday := today.Add(time.Hour * -24)

		todaysCount := byDay[dateKey(today)]
		yesterdaysCount := byDay[dateKey(yesterday)]

		delta := todaysCount - yesterdaysCount

		p := ""
		if delta >= 0 {
			p = "+"
		}

		fmt.Print(todaysCount, " (", p, delta, ")\n")

	case "start":
		recorder := selfwatch.NewRecorder()
		storage.BindRecorder(recorder, config.SyncDelay)

		if config.RemoteUrl != "" {
			remote := selfwatch.RemoteSync{
				Url:     config.RemoteUrl,
				Storage: storage,
			}

			if config.RemoteFlushDelay > 0 {
				remote.FlushEvery(config.RemoteFlushDelay)
			}
		}

		recorder.Bind()
	default:
		log.Fatal("Unknown command:", command)
	}

}
