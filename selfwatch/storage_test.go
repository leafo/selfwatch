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

	exists, err := storage.SchemaExists()
	if err != nil {
		t.Fatal(err.Error())
	}
	if exists {
		t.Fatal("Schema exists when it shouldn't yet.")
	}
	storage.CreateSchema()
	exists, err = storage.SchemaExists()
	if err != nil {
		t.Fatal(err.Error())
	}
	if !exists {
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

	if err = storage.WriteKeys(5); err != nil {
		t.Fatal(err.Error())
	}
	if err = storage.WriteKeys(2); err != nil {
		t.Fatal(err.Error())
	}
	if err = storage.WriteKeys(6); err != nil {
		t.Fatal(err.Error())
	}

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
