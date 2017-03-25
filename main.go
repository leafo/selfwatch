package main

import (
	"log"

	"github.com/leafo/selfwatch/selfwatch"
)

func main() {
	storage, err := selfwatch.NewWatchStorage("keys.db")
	if err != nil {
		log.Fatal(err.Error())
	}

	storage.CreateSchema()

	recorder := selfwatch.NewRecorder()

	recorder.KeyRelease = func(code int32) {
		storage.WriteKeys(1)
	}

	recorder.Bind()
}
