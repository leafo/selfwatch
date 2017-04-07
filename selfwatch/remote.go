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
	MaxId int64 `json:"max_id"`
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

	var r maxRows
	log.Print(string(body))
	err = json.Unmarshal(body, &r)

	if err != nil {
		return err
	}

	log.Print(r.MaxId)

	return nil
}

func (s *RemoteSync) SendRows() error {
	return nil
}
