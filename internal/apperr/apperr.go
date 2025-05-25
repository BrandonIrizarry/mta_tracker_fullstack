package apperr

import (
	"log/slog"
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
			slog.Error("Error serving request", "error", serr)
			return
		}

		slog.Info("Request served successfully")
	}
}

func ServeError(err error, code int) *StatusError {
	return &StatusError{
		error: err,
		code:  code,
	}
}
