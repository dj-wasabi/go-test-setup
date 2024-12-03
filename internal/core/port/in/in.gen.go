// Package in provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package in

import (
	"time"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// AuthenticatePostRequest Authenticate
type AuthenticatePostRequest struct {
	Password string `json:"password" validate:"required,min=6,max=256"`
	Username string `json:"username" validate:"required,min=6,max=64,alphanum"`
}

// AuthenticatePostResponse Authenticate
type AuthenticatePostResponse struct {
	Token string `json:"token" validate:"required"`
}

// Error Default error response
type Error struct {
	Message string `json:"message"`
}

// MetricsGetResponse Metrics for Prometheus response
type MetricsGetResponse = string

// Organisation Return of the Organisation object
type Organisation struct {
	Admins      []string  `json:"admins"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description" validate:"required,min=1,max=256"`
	Enabled     bool      `json:"enabled"`
	Fqdn        string    `json:"fqdn"`
	Id          string    `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=6,max=64,alphanum"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OrganisationIn Creating the Organisation object
type OrganisationIn struct {
	Admins      []string `json:"admins"`
	Description string   `json:"description" validate:"required,min=1,max=256"`
	Enabled     bool     `json:"enabled"`
	Fqdn        string   `json:"fqdn" validate:"required,fqdn"`
	Name        string   `json:"name" validate:"required,min=6,max=64,alphanum"`
}

// Organisations Overview of all organisations
type Organisations = []Organisation

// User Return of the User object
type User struct {
	CreatedAt time.Time `json:"created_at"`
	Enabled   bool      `json:"enabled"`
	Id        string    `json:"id"`
	OrgId     string    `json:"org_id"`
	Password  string    `json:"password" validate:"required,min=6,max=256"`
	Role      string    `json:"role"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username" validate:"required,min=6,max=64,alphanum"`
}

// UserIn Object to create the User
type UserIn struct {
	Enabled  bool   `json:"enabled,omitempty"`
	OrgId    string `json:"org_id" validate:"required"`
	Password string `json:"password" validate:"required,min=6,max=256"`
	Role     string `json:"role,omitempty"`
	Username string `json:"username" validate:"required,min=6,max=64,alphanum"`
}

// UserNoPassword Return of the User object
type UserNoPassword struct {
	CreatedAt time.Time `json:"created_at"`
	Enabled   bool      `json:"enabled"`
	Id        string    `json:"id" validate:"required"`
	OrgId     string    `json:"org_id" validate:"required"`
	Role      string    `json:"role"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username" validate:"required,min=6,max=64,alphanum"`
}

// GetAllOrganisationsParams defines parameters for GetAllOrganisations.
type GetAllOrganisationsParams struct {
	// Limit Size of the page, maximum is 100, default is 25
	Limit *int32 `form:"limit,omitempty" json:"limit,omitempty"`

	// Page The page to return
	Page *int32 `form:"page,omitempty" json:"page,omitempty"`
}

// AuthenticateLoginJSONRequestBody defines body for AuthenticateLogin for application/json ContentType.
type AuthenticateLoginJSONRequestBody = AuthenticatePostRequest

// CreateOrganisationJSONRequestBody defines body for CreateOrganisation for application/json ContentType.
type CreateOrganisationJSONRequestBody = OrganisationIn

// UserCreateJSONRequestBody defines body for UserCreate for application/json ContentType.
type UserCreateJSONRequestBody = UserIn
