package entity

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//User data
type User struct {
	ID        int64
	Email     string
	Password  string
	Username  string
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}

//NewUser create a new user
func NewUser(email, password, username string) (*User, error) {
	u := &User{
		Email:     email,
		Username:  username,
		CreatedAt: sql.NullTime{Time: time.Now()},
	}
	pwd, err := generatePassword(password)
	if err != nil {
		return nil, err
	}
	u.Password = pwd
	err = u.Validate()
	if err != nil {
		return nil, ErrInvalidEntity
	}
	return u, nil
}

//Validate validate data
func (u *User) Validate() error {
	if u.Email == "" || u.Username == "" || u.Password == "" {
		return ErrInvalidEntity
	}

	return nil
}

//ValidatePassword validate user password
func (u *User) ValidatePassword(p string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	if err != nil {
		return err
	}

	return nil
}

func generatePassword(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), 10)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
