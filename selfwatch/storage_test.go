package selfwatch

import "testing"

func TestNewStorage(t *testing.T) {
	_, err := NewWatchStorage("test.db")
	if err != nil {
		t.Fatal(err.Error())
	}
}
