package sqllite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/amit8889/golangCRUDApi/internal/config"
	"github.com/amit8889/golangCRUDApi/internal/types"
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

func (s *Sqlite) GetStudent(id any) (types.Student, error) {
	smt, err := s.DB.Prepare("SELECT * FROM student WHERE id=?")
	if err != nil {
		return types.Student{}, err
	}
	defer smt.Close()
	var student types.Student
	err = smt.QueryRow(id).Scan(&student.ID, &student.Name, &student.Email, &student.Age)
	if err != nil {
		return types.Student{}, err
	}
	return student, nil
}
func (s *Sqlite) GetAllStudents() ([]types.Student, error) {
	smt, err := s.DB.Prepare("SELECT * FROM student")
	if err != nil {
		return nil, err
	}
	defer smt.Close()
	var students []types.Student
	rows, err := smt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var student types.Student
		err = rows.Scan(&student.ID, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}
func (s *Sqlite) DeleteStudent(id any) error {
	smt, err := s.DB.Prepare("DELETE FROM student WHERE id=?")
	if err != nil {
		return err
	}
	defer smt.Close()
	fmt.Println("=====id==", id)
	// Execute the delete query
	result, err := smt.Exec(id) // assuming 1 is the id you want to delete
	if err != nil {
		return err
	}

	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected")
	} else {
		fmt.Println("Row deleted successfully")
	}

	return nil

}
func (s *Sqlite) UpdateStudent(id any, student types.Student) error {
	// Prepare the update query
	smt, err := s.DB.Prepare("UPDATE student SET name=?, email=?, age=? WHERE id = ?")
	if err != nil {
		return err
	}
	defer smt.Close()

	// Execute the update query
	result, err := smt.Exec(student.Name, student.Email, student.Age, id)
	if err != nil {
		return err
	}

	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, return an error
	if rowsAffected == 0 {
		return errors.New("no rows updated, possibly invalid ID")
	}

	fmt.Println("Row updated successfully")
	return nil
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
