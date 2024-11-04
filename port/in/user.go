package in

import (
	"time"
)

type UserIn struct {
	Id        string    `json:"-"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Roles     []string  `json:"roles"`
}

type IUser interface {
	GetId() string
	GetUsername() string
	GetPassword() string
	GetEnabled() bool
	GetCreated() time.Time
	GetUpdated() time.Time
	GetRoles() []string
}

func (o *UserIn) GetId() string {
	return o.Id
}

func (o *UserIn) GetUsername() string {
	return o.Username
}

func (o *UserIn) GetPassword() string {
	return o.Password
}

func (o *UserIn) GetEnabled() bool {
	return o.Enabled
}

func (o *UserIn) GetCreated() time.Time {
	return o.CreatedAt
}

func (o *UserIn) GetUpdated() time.Time {
	return o.UpdatedAt
}

func (o *UserIn) GetRoles() []string {
	return o.Roles
}

func NewUserIn(username, password string, enabled bool, roles []string) *UserIn {
	return &UserIn{
		Username:  username,
		Password:  password,
		Enabled:   enabled,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Roles:     roles,
	}
}
