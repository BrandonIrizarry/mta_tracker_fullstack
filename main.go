package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/BrandonIrizarry/mta_tracker_fullstack/internal/apperr"
	"github.com/BrandonIrizarry/mta_tracker_fullstack/internal/availableRoutes"
	"github.com/BrandonIrizarry/mta_tracker_fullstack/internal/geturl"
	"github.com/joho/godotenv"
)

var searchResultsHTML = template.Must(template.New("results").Parse(`
{{range .}}
<tr>
    <td>{{.}}</td>
</tr>
{{end}}
`))

// For now, the agency we use is fixed.
const routesForAgencyURL = "https://bustime.mta.info/api/where/routes-for-agency/MTA%20NYCT.json"

type config struct {
	apiKey string
	port   string
}

func (cfg *config) init() error {
	apiKey := os.Getenv("API_KEY")

	if apiKey == "" {
		return errors.New("Missing API key (API_KEY)")
	}

	cfg.apiKey = apiKey

	port := os.Getenv("PORT")

	if port == "" {
		return errors.New("Missing port configuration (PORT)")
	}

	cfg.port = port

	return nil
}

func (cfg config) searchHandler(w http.ResponseWriter, r *http.Request) *apperr.StatusError {
	routeQuery := r.FormValue("search")

	if routeQuery == "" {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	if err := r.ParseForm(); err != nil {
		return apperr.ServeError(err, http.StatusInternalServerError)
	}

	jsonBytes, callErr := geturl.Call(routesForAgencyURL, map[string]string{
		"key": cfg.apiKey,
	})

	if callErr != nil {
		return apperr.ServeError(callErr, http.StatusInternalServerError)
	}

	var routesInfo availableRoutes.AvailableRoutes

	if err := json.Unmarshal(jsonBytes, &routesInfo); err != nil {
		return apperr.ServeError(err, http.StatusInternalServerError)
	}

	if code := routesInfo.Code; code != 200 {
		codeErr := fmt.Errorf("Routes info reports non-200 code: %d", code)
		return apperr.ServeError(codeErr, code)
	}

	var results []string

	for _, route := range routesInfo.Data.List {
		// FIXME: The agency prefix is hardcoded here. If we
		// ever expand this to include subway, PATH, other
		// kinds of buses etc., we would to change this.
		baseID, found := strings.CutPrefix(route.ID, "MTA NYCT_")

		if !found {
			prefixErr := errors.New("Missing agency prefix: MTA NYCT_")
			return apperr.ServeError(prefixErr, http.StatusInternalServerError)
		}

		if strings.Contains(strings.ToLower(baseID), strings.ToLower(routeQuery)) {
			// I want to display the _original_ ID in the
			// HTML response; all the lowercasing is to
			// facilitate the search itself.
			results = append(results, baseID)
		}
	}

	if err := searchResultsHTML.Execute(w, results); err != nil {
		return apperr.ServeError(err, http.StatusInternalServerError)
	}

	return nil
}

func main() {
	godotenv.Load(".env")

	var cfg config

	if err := cfg.init(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	rootHandler := http.StripPrefix("/app/", http.FileServer(http.Dir("./templates")))
	mux.Handle("/app/", rootHandler)
	mux.HandleFunc("POST /search", apperr.WithErrors(cfg.searchHandler))

	srv := &http.Server{
		Addr:    ":" + cfg.port,
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())
}
