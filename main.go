package main

import (
	"log"

	"github.com/leafo/selfwatch/selfwatch"
)

func main() {
	recorder := selfwatch.NewRecorder()
	recorder.Bind()

	storage, err := selfwatch.NewWatchStorage("keys.db")
	if err != nil {
		log.Fatal(err.Error())
	}

	storage.CreateSchema()
}
