package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_new_user(t *testing.T) {
	username := "test_user"
	password := "secure_password"
	enabled := true
	role := "admin"

	user, err := NewUser(username, password, role, enabled)

	if err != nil {
		t.Errorf("Unexpected error creating user: %v", err)
	}

	if user.Username != username || user.Password != password || user.Enabled != enabled || user.Role != role {
		t.Errorf("User details don't match expected values")
	}
}

func Test_new_user_without_username(t *testing.T) {
	password := "securepassword"
	enabled := true
	role := "admin"

	_, err := NewUser("", password, role, enabled)
	assert.Nil(t, err, "Username is not provided")
}

func Test_new_user_without_password(t *testing.T) {
	username := "test_user"
	enabled := true
	role := "admin"

	_, err := NewUser(username, "", role, enabled)
	assert.Nil(t, err, "Password is not provided")
}

func Test_new_user_funcs(t *testing.T) {
	username := "test_user"
	password := "secure_password"
	enabled := true
	role := "admin"

	user, _ := NewUser(username, password, role, enabled)

	assert.Equal(t, user.GetUsername(), username, "Username should be equal.")
	assert.Equal(t, user.GetPassword(), password, "Password should be equal.")
	assert.Equal(t, user.GetEnabled(), enabled, "Enabled should be equal.")
	assert.Equal(t, user.GetRole(), "admin", "Is it part of a role?")
}
