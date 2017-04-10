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

func TestInsertKeys(t *testing.T) {
	cleanDb()

	storage, err := NewWatchStorage(testDbName)
	if err != nil {
		t.Fatal(err.Error())
	}

	storage.CreateSchema()

	storage.WriteKeys(5)
	storage.WriteKeys(2)
	storage.WriteKeys(6)

	rows, err := storage.db.Query(`select count(*) from keys;`)
	rows.Next()

	var count int
	err = rows.Scan(&count)

	if err != nil {
		t.Fatal(err.Error())
	}

	if count != 3 {
		t.Fatal("Expected three rows")
	}

}
