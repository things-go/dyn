package http

import (
	"fmt"
	"net/http"
)

type ErrorReply struct {
	Code   int
	Body   []byte
	Header http.Header
}

func (e *ErrorReply) Error() string {
	return fmt.Sprintf("Invoke: Status Code: %d, Status Text: %s", e.Code, http.StatusText(e.Code))
}
