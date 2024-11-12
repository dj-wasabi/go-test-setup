package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_new_user(t *testing.T) {
	username := "test_user"
	password := "secure_password"
	enabled := true
	roles := []string{"admin"}

	user, err := NewUser(username, password, enabled, roles)

	if err != nil {
		t.Errorf("Unexpected error creating user: %v", err)
	}

	if user.Username != username || user.Password != password || user.Enabled != enabled || user.Roles[0] != roles[0] {
		t.Errorf("User details don't match expected values")
	}
}

func Test_new_user_without_username(t *testing.T) {
	password := "securepassword"
	enabled := true
	roles := []string{"admin"}

	_, err := NewUser("", password, enabled, roles)
	assert.Nil(t, err, "Username is not provided")
}

func Test_new_user_without_password(t *testing.T) {
	username := "test_user"
	enabled := true
	roles := []string{"admin"}

	_, err := NewUser(username, "", enabled, roles)
	assert.Nil(t, err, "Password is not provided")
}

func Test_new_user_funcs(t *testing.T) {
	username := "test_user"
	password := "secure_password"
	enabled := true
	roles := []string{"admin"}

	user, _ := NewUser(username, password, enabled, roles)

	assert.Equal(t, user.GetUsername(), username, "Username should be equal.")
	assert.Equal(t, user.GetPassword(), password, "Password should be equal.")
	assert.Equal(t, user.GetEnabled(), enabled, "Enabled should be equal.")
	assert.Contains(t, user.GetRoles(), "admin", "Is it part of a role?")
}
