package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	username = "testuser"
	password = "1Secure_password!"
	enabled  = true
	role     = "admin"
	orgid    = "orgid123456"
)

func Test_new_user(t *testing.T) {
	user, err := NewUser(username, password, role, enabled, orgid)

	if err != nil {
		t.Errorf("Unexpected error creating user: %v", err)
	}

	if user.Username != username || user.Password != password || user.Enabled != enabled || user.Role != role {
		t.Errorf("User details don't match expected values")
	}
}

func Test_new_user_without_username(t *testing.T) {
	_, err := NewUser("", password, role, enabled, orgid)
	assert.Equal(t, err.Error(), "The field 'username' is required.")
}

func Test_new_user_username_incorrect(t *testing.T) {
	username := "test_user"

	_, err := NewUser(username, password, role, enabled, orgid)
	assert.Equal(t, err.Error(), "Only alphabetical and numerical characters are allowed.")
}

func Test_new_user_without_password(t *testing.T) {
	_, err := NewUser(username, "", role, enabled, orgid)
	assert.Equal(t, err.Error(), "The field 'password' is required.")
}

func Test_new_user_username_not_minimum_characters(t *testing.T) {
	username := "test"

	_, err := NewUser(username, password, role, enabled, orgid)
	assert.Equal(t, err.Error(), "The 'username' field needs a minimum amount of 6 characters.")
}

func Test_new_user_username_more_than_maximum_characters(t *testing.T) {
	username := "testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest"

	_, err := NewUser(username, password, role, enabled, orgid)
	assert.Equal(t, err.Error(), "The 'username' field has a maximum amount of 64 characters.")
}

func Test_new_user_password_needs_uppercase(t *testing.T) {
	password := "1secure_password!"

	_, err := NewUser(username, password, role, enabled, orgid)
	assert.Equal(t, err.Error(), "The password needs to have at least 1 uppercase, lowercase, number and any of the following characters: !\"#$%&'()*+,\\-./:;<=>?@[\\\\]^_`{|}~]")
}

func Test_new_user_password_needs_lowercase(t *testing.T) {
	password := "1SECURE_PASSWORD!"

	_, err := NewUser(username, password, role, enabled, orgid)
	assert.Equal(t, err.Error(), "The password needs to have at least 1 uppercase, lowercase, number and any of the following characters: !\"#$%&'()*+,\\-./:;<=>?@[\\\\]^_`{|}~]")
}

func Test_new_user_password_needs_number(t *testing.T) {
	password := "secure_password!"

	_, err := NewUser(username, password, role, enabled, orgid)
	assert.Equal(t, err.Error(), "The password needs to have at least 1 uppercase, lowercase, number and any of the following characters: !\"#$%&'()*+,\\-./:;<=>?@[\\\\]^_`{|}~]")
}

func Test_new_user_password_needs_character(t *testing.T) {
	password := "1securepassword"

	_, err := NewUser(username, password, role, enabled, orgid)
	assert.Equal(t, err.Error(), "The password needs to have at least 1 uppercase, lowercase, number and any of the following characters: !\"#$%&'()*+,\\-./:;<=>?@[\\\\]^_`{|}~]")
}

func Test_new_user_without_roles(t *testing.T) {
	_, err := NewUser(username, password, "", enabled, orgid)
	assert.Equal(t, err.Error(), "Only 'admin', 'write' or 'readonly' are allowed.")
}

func Test_new_user_incorrect_role(t *testing.T) {
	role := "pizza"

	_, err := NewUser(username, password, role, enabled, orgid)
	assert.Equal(t, err.Error(), "Only 'admin', 'write' or 'readonly' are allowed.")
}

func Test_new_user_funcs(t *testing.T) {

	user, _ := NewUser(username, password, role, enabled, orgid)

	assert.Equal(t, user.GetUsername(), username, "Username should be equal.")
	assert.Equal(t, user.GetPassword(), password, "Password should be equal.")
	assert.Equal(t, user.GetEnabled(), enabled, "Enabled should be equal.")
	assert.Equal(t, user.GetRole(), "admin", "Is it part of a role?")
}
