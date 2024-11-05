package model

type Error struct {
	// Code    string `json:"code"`
	Message string `json:"error"`
}

var ErrorCodes = map[string]string{
	"UNKNOWN": "Unkown error message",
	"ORG0001": "Duplicate organisation",
	"USR0001": "Duplicate username",
	"USR0002": "Incorrect username/password combination",
}

func GetError(error string) *Error {
	message := &Error{
		Message: ErrorCodes[error],
	}
	return message
}
