package selfwatch

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const DefaultConfigFname = "selfwatch.json"

type config struct {
	DbName           string
	RemoteUrl        string
	RemoteFlushDelay float64
	SyncDelay        float64
}

var defaultConfig = config{
	DbName:           "selfwatch.db",
	RemoteUrl:        "",
	RemoteFlushDelay: 60,
	SyncDelay:        60,
}

func LoadConfig(fname string) *config {
	c := defaultConfig

	if fname == "" {
		return &c
	}

	jsonBlob, err := ioutil.ReadFile(fname)
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
