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
	cleanDb()

	_, err := NewWatchStorage(testDbName)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestCreateSchema(t *testing.T) {
	cleanDb()

	storage, err := NewWatchStorage(testDbName)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = storage.CreateSchema()

	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestSchemaExists(t *testing.T) {
	cleanDb()

	storage, err := NewWatchStorage(testDbName)
	if err != nil {
		t.Fatal(err.Error())
	}

	if storage.SchemaExists() {
		t.Fatal("Schema exists when it shouldn't yet.")
	}
	storage.CreateSchema()
	if !storage.SchemaExists() {
		t.Fatal("Schema does not exist but should.")
	}
}
