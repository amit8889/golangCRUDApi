package storage

import "github.com/amit8889/golangCRUDApi/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudent(id any) (types.Student, error)
}
