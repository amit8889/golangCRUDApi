package sqllite

import (
	"database/sql"

	"github.com/amit8889/golangCRUDApi/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	DB *sql.DB
}

// CreateStudent implements storage.Storage.
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	smt, err := s.DB.Prepare("INSERT INTO student (name,email,age) VALUES(?,?,?)")
	if err != nil {
		return 0, err
	}
	defer smt.Close()

	res, err := smt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func New(cfg *config.Config) (*Sqlite, error) {
	// Open the SQLite3 database
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	// Corrected SQL syntax: CREATE TABLE IF NOT EXISTS
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS student(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		age INTEGER
	)
	`)
	if err != nil {
		return nil, err
	}

	return &Sqlite{
		DB: db,
	}, nil
}
