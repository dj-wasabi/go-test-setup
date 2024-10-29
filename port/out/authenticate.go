package out

type AuthenticationRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type AuthenticationToken struct {
	Token string
}

type Authentication interface {
	GetUsername() string
	GetPassword() string
}

func (a *AuthenticationRequest) GetUsername() string {
	return a.Username
}

func (a *AuthenticationRequest) GetPassword() string {
	return a.Password
}
