package teamtailorgo

import (
	"fmt"
)

type ErrorStruct struct {
	StatusCode int
	Message    string
}

func (e ErrorStruct) Error() string {
	return fmt.Sprintf("Error %v: %v", e.StatusCode, e.Message)
}

func UnauthorizedError(status int) error {
	return ErrorStruct{
		status,
		"Your token is not authorized",
	}
}
