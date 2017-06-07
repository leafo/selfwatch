package main

import (
	"flag"
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

	counter := 0
	last := time.Unix(0, 0)
	var lastWindow int64

	recorder := selfwatch.NewRecorder()
	recorder.KeyRelease = func(event selfwatch.Event) {
		counter += 1
		if time.Now().Sub(last).Seconds() > config.SyncDelay || event.Window != lastWindow {
			if counter > 0 {
				log.Println("Syncing keys...", counter)
				storage.WriteKeys(counter)
				counter = 0
			}

			last = time.Now()
			lastWindow = event.Window
		}
	}

	if config.RemoteUrl != "" {
		// flush example
		remote := selfwatch.RemoteSync{
			Url:     config.RemoteUrl,
			Storage: storage,
		}

		if config.RemoteFlushDelay > 0 {
			remote.FlushEvery(config.RemoteFlushDelay)
		}
	}

	recorder.Bind()
}
