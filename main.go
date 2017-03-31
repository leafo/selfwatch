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

	recorder := selfwatch.NewRecorder()
	recorder.KeyRelease = func(code int32) {
		if time.Now().Sub(last).Seconds() > 60 {
			storage.WriteKeys(counter)
			counter = 0
			last = time.Now()
		} else {
			counter += 1
		}

	}

	recorder.Bind()
}
