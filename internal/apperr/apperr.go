package apperr

import (
	"net/http"
)

type StatusError struct {
	error
	code int
}

func WithErrors(fn func(http.ResponseWriter, *http.Request) *StatusError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serr := fn(w, r)

		if serr != nil {
			http.Error(w, serr.Error(), serr.code)
		}
	}
}

func ServeError(err error, code int) *StatusError {
	return &StatusError{
		error: err,
		code:  code,
	}
}
