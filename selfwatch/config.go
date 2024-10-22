package selfwatch

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
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

func expandHomePath(path string) (string, error) {
	if path[:2] == "~/" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, path[2:]), nil
	}
	return path, nil
}

func LoadConfig(fname string) *config {
	c := defaultConfig

	if fname == "" {
		return &c
	}

	expandedFname, err := expandHomePath(fname)
	if err != nil {
		log.Fatal("Failed to expand config path: ", err.Error())
	}
	fname = expandedFname

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
