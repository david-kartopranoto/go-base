package user

import (
	"database/sql"
	"testing"
	"time"

	"github.com/david-kartopranoto/go-base/entity"
	"github.com/david-kartopranoto/go-base/usecase/user/mock"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
)

func newFixtureUser() *entity.User {
	return &entity.User{
		ID:        123456,
		Email:     "john@doe.com",
		Password:  "abc123",
		Username:  "JohnDoe",
		CreatedAt: sql.NullTime{Time: time.Now()},
	}
}

func Test_Register(t *testing.T) {
	u := newFixtureUser()

	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mock := mock.NewMockRepository(ctl)
	mock.EXPECT().Create(gomock.Any()).Return(int64(u.ID), nil)

	m := NewService(mock)
	newID, err := m.Register(u.Email, u.Password, u.Username)
	assert.Nil(t, err)
	assert.Equal(t, u.ID, newID)
	assert.False(t, u.CreatedAt.Time.IsZero())
	assert.True(t, u.UpdatedAt.Time.IsZero())
}

func Test_Get(t *testing.T) {
	u := newFixtureUser()

	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mock := mock.NewMockRepository(ctl)
	mock.EXPECT().Get(gomock.Any()).Return(u, nil)

	m := NewService(mock)
	_, err := m.GetUser(u.ID)
	assert.Nil(t, err)
}

func Test_Search(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mock := mock.NewMockRepository(ctl)
	mock.EXPECT().Search(gomock.Any()).Return(nil, nil)

	m := NewService(mock)
	_, err := m.SearchUsers("search")
	assert.Nil(t, err)
}

func Test_List(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mock := mock.NewMockRepository(ctl)
	mock.EXPECT().List().Return(nil, nil)

	m := NewService(mock)
	_, err := m.ListUsers()

	assert.Nil(t, err)
}
