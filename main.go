package main

import (
	"log"
	"time"

	"github.com/leafo/selfwatch/selfwatch"
)

func main() {
	storage, err := selfwatch.NewWatchStorage("keys.db")
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
		if time.Now().Sub(last).Seconds() > 60 || event.Window != lastWindow {
			if counter > 0 {
				log.Println("Flushing....", counter)
				storage.WriteKeys(counter)
				counter = 0
			}

			last = time.Now()
			lastWindow = event.Window
		}
	}

	// flush example
	// remote := selfwatch.RemoteSync{
	// 	Url:     "http://localhost/put.php",
	// 	Storage: storage,
	// }

	// remote.FlushKeys()

	recorder.Bind()
}
