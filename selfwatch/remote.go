package selfwatch

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type RemoteSync struct {
	Url     string
	Storage *WatchStorage
}

type maxRows struct {
	MaxId int64 `json:"max_id"`
}

func (s *RemoteSync) GetLastRowId() (int64, error) {
	resp, err := http.Get(s.Url)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return 0, err
	}

	var r maxRows
	log.Print(string(body))
	err = json.Unmarshal(body, &r)

	if err != nil {
		return 0, err
	}

	return r.MaxId, nil
}

func (s *RemoteSync) SendRows(rows [][]interface{}) error {
	payload, err := json.Marshal(rows)

	if err != nil {
		return err
	}

	res, err := http.Post(s.Url, "application/json", bytes.NewReader(payload))

	if err != nil {
		return err
	}

	res.Body.Close()

	return nil
}

func (s *RemoteSync) FlushKeys() error {
	maxRowId, err := s.GetLastRowId()
	rows, err := s.Storage.SerializeRecentKeyCounts(maxRowId)

	if err != nil {
		return err
	}

	return s.SendRows(rows)
}
