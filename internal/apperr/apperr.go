package apperr

import (
	"log"
	"net/http"
)

func WithErrors(fn func(http.ResponseWriter, *http.Request) (error, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err, code := fn(w, r)

		if err != nil {
			http.Error(w, "Oops, something went wrong", code)

			log.Printf("Error serving request: %v", err)
			return
		}

		log.Printf("%s %s: succeeded", r.Method, r.URL)
	}
}
