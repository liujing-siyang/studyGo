package errors

import (
	"github.com/go-kratos/kratos/v2/errors"
	"fmt"
	stdhttp "net/http"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// HTTPError is an HTTP error.
type HTTPError struct {
	Code          int                 `json:"code"`
	ErrorsMessage map[string][]string `json:"errors_message"`
}

func NewHttpError(code int, filed string, detail string) *HTTPError{
	return &HTTPError{
		Code:          code,
		ErrorsMessage: map[string][]string{
			filed:{detail},
		},
	}
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTPError code: %d message: %v", e.Code, e.ErrorsMessage)
}

// FromError try to convert an error to *HTTPError.
func FromError(err error) *HTTPError {
	if err == nil {
		return nil
	}
	if se := new(HTTPError); errors.As(err, &se) {
		return se
	}
	if se := new(errors.Error); errors.As(err, &se) {
		return NewHttpError(int(se.Code), se.Reason, se.Message)
	}
	return &HTTPError{Code: 500}
}

func ErrorEncoder(w stdhttp.ResponseWriter, r *stdhttp.Request, err error) {
	se := FromError(err)
	codec, _ := http.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(se)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/"+codec.Name())
	w.WriteHeader(se.Code)
	_, _ = w.Write(body)
}
