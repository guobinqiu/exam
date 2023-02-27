package rate

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)

type MyDB struct {
	db   *sql.DB
	mu   *sync.Mutex
	once *sync.Once
	path string
}

func NewMyDB(path string) *MyDB {
	return &MyDB{
		mu:   &sync.Mutex{},
		once: &sync.Once{},
		path: path,
	}
}

func (s *MyDB) GetDB() *sql.DB {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		db, err := sql.Open("sqlite3", s.path)
		if err != nil {
			panic(fmt.Sprintf("Init database connection error: %s", err.Error()))
		}
		s.db = db
	}
	return s.db
}

func (s *MyDB) Close() {
	if s.db != nil {
		s.once.Do(func() {
			s.db.Close()
			s.db = nil
		})
	}
}
