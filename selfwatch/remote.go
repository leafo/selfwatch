package selfwatch

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type RemoteSync struct {
	Url string
}

type maxRows struct {
	max_id int64
}

func (s *RemoteSync) GetLastRowId() error {
	resp, err := http.Get(s.Url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	max := maxRows{}
	json.Unmarshal(body, &max)

	log.Print(max)

	return nil
}

func (s *RemoteSync) SendRows() error {
	return nil
}
