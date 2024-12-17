package model

var ErrorCodes = map[string]string{
	"UNKNOWN": "Unkown error message",
	"ORG0001": "Duplicate organisation",

	"USR0001": "Duplicate username",
	"USR0002": "Invalid username/password combination",
	"USR0003": "Error while decoding user object",
	"USR0004": "User not found",
	"USR0005": "Invalid User Object ID",
	"USR0006": "Access denied",

	"AUTH001": "Error while validating the token.",
	"AUTH002": "Error while updating the token, try again.",
	"AUTH003": "Invalid token/user combination.",
	"AUTH004": "Role does not have permissions to access this endpoint.",
}

func GetError(error, log_id string) *Error {
	message := &Error{
		Error: ErrorCodes[error],
		Logid: log_id,
	}
	return message
}

func NewError(error string) *Error {
	return &Error{
		Error: error,
	}
}
