package selfwatch

import (
	"os"
	"testing"
)

const testDbName = "test.db"

func cleanDb() {
	os.Remove(testDbName)
}

func TestNewStorage(t *testing.T) {
	_, err := NewWatchStorage(testDbName)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestCreateSchema(t *testing.T) {
	storage, err := NewWatchStorage(testDbName)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = storage.CreateSchema()

	if err != nil {
		t.Fatal(err.Error())
	}
}
