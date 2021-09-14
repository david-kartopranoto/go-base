package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	u, err := NewUser("john@doe.com", "new_password", "JohnDoe")
	assert.Nil(t, err)
	assert.Equal(t, u.Username, "JohnDoe")
	assert.NotEqual(t, u.Password, "new_password")
}

func TestValidatePassword(t *testing.T) {
	u, err := NewUser("john@doe.com", "new_password", "JohnDoe")
	err = u.ValidatePassword("new_password")
	assert.Nil(t, err)
	err = u.ValidatePassword("wrong_password")
	assert.NotNil(t, err)
}

func TestUserValidate(t *testing.T) {
	type test struct {
		email     string
		password  string
		username string
		want      error
	}

	tests := []test{
		{
			email:     "john@doe.com",
			password:  "new_password",
			username:  "JohnDoe",
			want:      nil,
		},
		{
			email:     "",
			password:  "new_password",
			username: "JohnDoe",
			want:      ErrInvalidEntity,
		},
	}
	for _, tc := range tests {

		_, err := NewUser(tc.email, tc.password, tc.username)
		assert.Equal(t, err, tc.want)
	}

}