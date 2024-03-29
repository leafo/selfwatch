package selfwatch

import (
	"encoding/json"
	"log"
	"os"
)

const DefaultConfigFname = "~/.selfwatch/config.json"

type config struct {
	DbName           string
	RemoteUrl        string
	RemoteFlushDelay float64
	SyncDelay        float64
	NewDayHour       int
}

var defaultConfig = config{
	DbName:           "~/.selfwatch/selfwatch.db",
	RemoteUrl:        "",
	RemoteFlushDelay: 60,
	SyncDelay:        60,
	NewDayHour:       4,
}

func LoadConfig(fname string) *config {
	c := defaultConfig

	if fname == "" {
		return &c
	}

	log.Print("Loading config ", fname)
	jsonBlob, err := os.ReadFile(fname)
	if err == nil {
		err = json.Unmarshal(jsonBlob, &c)

		if err != nil {
			log.Fatal("Failed parsing config: ", fname, ": ", err.Error())
		}
	} else {
		log.Print(err.Error())
	}

	return &c
}
