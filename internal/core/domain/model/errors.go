package model

type Error struct {
	// Code    string `json:"code"`
	Message string `json:"error"`
}

var ErrorCodes = map[string]string{
	"UNKNOWN": "Unkown error message",
	"ORG0001": "Duplicate organisation",

	"USR0001": "Duplicate username",
	"USR0002": "Invalid username/password combination",
	"USR0003": "Error while decoding user object",

	"AUTH001": "Error while validating the token.",
	"AUTH002": "Error while updating the token, try again.",
	"AUTH003": "Invalid token/user combination.",
}

func GetError(error string) *Error {
	message := &Error{
		Message: ErrorCodes[error],
	}
	return message
}
