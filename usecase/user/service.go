package user

import (
	"strings"

	"github.com/david-kartopranoto/go-base/entity"
)

//Service  interface
type Service struct {
	repo Repository
}

//NewService create new use case
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

//Register register an user
func (s *Service) Register(email, password, username string) (int64, error) {
	e, err := entity.NewUser(email, password, username)
	if err != nil {
		return e.ID, err
	}
	return s.repo.Create(e)
}

//GetUser Get an user
func (s *Service) GetUser(id int64) (*entity.User, error) {
	return s.repo.Get(id)
}

//SearchUsers Search users
func (s *Service) SearchUsers(query string) ([]*entity.User, error) {
	return s.repo.Search(strings.ToLower(query))
}

//ListUsers List users
func (s *Service) ListUsers() ([]*entity.User, error) {
	return s.repo.List()
}
