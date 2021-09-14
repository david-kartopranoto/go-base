package user

import (
	"github.com/david-kartopranoto/go-base/entity"
)

//Reader interface
type Reader interface {
	Get(id int64) (*entity.User, error)
	Search(query string) ([]*entity.User, error)
	List() ([]*entity.User, error)
}

//Writer user writer
type Writer interface {
	Create(e *entity.User) (int64, error)
	Update(e *entity.User) error
	Delete(id int64) error
}

//Repository interface
type Repository interface {
	Reader
	Writer
}

//UseCase interface
type UseCase interface {
	GetUser(id int64) (*entity.User, error)
	SearchUsers(query string) ([]*entity.User, error)
	ListUsers() ([]*entity.User, error)
	Register(email, password, username string) (int64, error)
}
