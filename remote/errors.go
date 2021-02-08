package remote

import (
	"fmt"
	"net/http"
)

type RequestErrors struct {
	Errors []RequestError `json:"errors"`
}

type RequestError struct {
	response *http.Response
	Code     string `json:"code"`
	Status   string `json:"status"`
	Detail   string `json:"detail"`
}

func IsRequestError(err error) bool {
	_, ok := err.(*RequestError)

	return ok
}

// Returns the error response in a string form that can be more easily consumed.
func (re *RequestError) Error() string {
	c := 0
	if re.response != nil {
		c = re.response.StatusCode
	}

	return fmt.Sprintf("Error response from Panel: %s: %s (HTTP/%d)", re.Code, re.Detail, c)
}

type sftpInvalidCredentialsError struct {
}

func (ice sftpInvalidCredentialsError) Error() string {
	return "the credentials provided were invalid"
}

func IsInvalidCredentialsError(err error) bool {
	_, ok := err.(*sftpInvalidCredentialsError)

	return ok
}
