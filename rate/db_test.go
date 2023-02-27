package rate

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMyDB_GetDB(t *testing.T) {
	dbPath := "./test.db"
	defer os.Remove(dbPath)
	db := NewMyDB(dbPath)
	err := db.GetDB().Ping()
	assert.NoError(t, err)
}

func TestMyDB_Close(t *testing.T) {
	dbPath := "./test.db"
	defer os.Remove(dbPath)
	db := NewMyDB(dbPath)
	db.GetDB().Close()
	err := db.GetDB().Ping()
	assert.Error(t, err)
}
