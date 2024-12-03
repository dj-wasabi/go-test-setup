package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_new_user(t *testing.T) {
	username := "testuser"
	password := "secure_password"
	enabled := true
	role := "admin"
	orgid := "orgid123456"

	user, err := NewUser(username, password, role, enabled, orgid)

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
	orgid := "orgid123456"

	_, err := NewUser("", password, role, enabled, orgid)
	assert.Equal(t, err.Error(), "The field 'username' is required.")
}

func Test_new_user_without_username_incorrect(t *testing.T) {
	username := "test_user"
	password := "securepassword"
	enabled := true
	role := "admin"
	orgid := "orgid123456"

	_, err := NewUser(username, password, role, enabled, orgid)
	assert.Equal(t, err.Error(), "Only alphabetical and numerical characters are allowed.")
}

func Test_new_user_without_password(t *testing.T) {
	username := "testuser"
	enabled := true
	role := "admin"
	orgid := "orgid123456"

	_, err := NewUser(username, "", role, enabled, orgid)
	assert.Equal(t, err.Error(), "The field 'password' is required.")
}

func Test_new_user_not_minimum_characters(t *testing.T) {
	username := "test"
	password := "securepassword"
	enabled := true
	role := "admin"
	orgid := "orgid123456"

	_, err := NewUser(username, password, role, enabled, orgid)
	assert.Equal(t, err.Error(), "The 'username' field needs a minimum amount of 6 characters.")
}

func Test_new_user_more_than_maximum_characters(t *testing.T) {
	username := "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest"
	password := "securepassword"
	enabled := true
	role := "admin"
	orgid := "orgid123456"

	_, err := NewUser(username, password, role, enabled, orgid)
	assert.Equal(t, err.Error(), "The 'username' field has a maximum amount of 64 characters.")
}

func Test_new_user_funcs(t *testing.T) {
	username := "testuser"
	password := "secure_password"
	enabled := true
	role := "admin"
	orgid := "orgid123456"

	user, _ := NewUser(username, password, role, enabled, orgid)

	assert.Equal(t, user.GetUsername(), username, "Username should be equal.")
	assert.Equal(t, user.GetPassword(), password, "Password should be equal.")
	assert.Equal(t, user.GetEnabled(), enabled, "Enabled should be equal.")
	assert.Equal(t, user.GetRole(), "admin", "Is it part of a role?")
}
