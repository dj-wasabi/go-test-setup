package in

type IUser interface {
	GetId() string
	GetOrgId() string
	GetUsername() string
	GetPassword() string
	GetEnabled() bool
	GetRole() string
}

func (o *UserIn) GetOrgId() string {
	return o.OrgId
}

func (o *UserIn) GetUsername() string {
	return o.Username
}

func (o *UserIn) GetPassword() string {
	return o.Password
}

func (o *UserIn) GetEnabled() bool {
	return *o.Enabled
}

func (o *UserIn) GetRole() string {
	return *o.Role
}

func NewUserIn(username, password, role string, enabled bool, orgid string) *UserIn {
	return &UserIn{
		Username: username,
		Password: password,
		Enabled:  &enabled,
		Role:     &role,
		OrgId:    orgid,
	}
}
