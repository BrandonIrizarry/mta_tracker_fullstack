package main

import (
	"database/sql"
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
	_ "github.com/mattn/go-sqlite3"
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
}

func (cfg *config) init() error {
	apiKey := os.Getenv("API_KEY")

	if apiKey == "" {
		return errors.New("Missing API key (API_KEY)")
	}

	cfg.apiKey = apiKey

	return nil
}

var routesInfo availableRoutes.AvailableRoutes

func (cfg config) getRoutes(w http.ResponseWriter, r *http.Request) (error, int) {
	routesResponse, callErr := geturl.Call(routesForAgencyURL, map[string]string{
		"key": cfg.apiKey,
	})

	if callErr != nil {
		return fmt.Errorf("Failed to fetch routes: %w", callErr), http.StatusInternalServerError
	}

	if err := json.Unmarshal(routesResponse, &routesInfo); err != nil {
		return fmt.Errorf("Failed to unmarshal routes response: %w", err), http.StatusInternalServerError
	}

	if code := routesInfo.Code; code != 200 {
		return fmt.Errorf("Routes info reports non-200 code: %d", code), http.StatusInternalServerError
	}

	return nil, 0
}

func (cfg config) searchHandler(w http.ResponseWriter, r *http.Request) (error, int) {
	// Test whether routesInfo is empty (not initialized) by
	// checking for a zero-valued Version
	if routesInfo.Version == 0 {
		return errors.New("Routes not initialized"), http.StatusInternalServerError
	}

	routeQuery := r.FormValue("search")

	if routeQuery == "" {
		w.WriteHeader(http.StatusOK)
		return nil, 0
	}

	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("Failed to parse form: %w", err), http.StatusInternalServerError
	}

	var results []string

	for _, route := range routesInfo.Data.List {
		// FIXME: The agency prefix is hardcoded here. If we
		// ever expand this to include subway, PATH, other
		// kinds of buses etc., we would to change this.
		baseID, found := strings.CutPrefix(route.ID, "MTA NYCT_")

		if !found {
			return errors.New("Missing agency prefix: MTA NYCT_"), http.StatusInternalServerError
		}

		if strings.Contains(strings.ToLower(baseID), strings.ToLower(routeQuery)) {
			// I want to display the _original_ ID in the
			// HTML response; all the lowercasing is to
			// facilitate the search itself.
			results = append(results, baseID)
		}
	}

	w.Header().Set("Content-Type", "text/html")

	if err := searchResultsHTML.Execute(w, results); err != nil {
		return fmt.Errorf("Failed to write results to HTML: %w", err), http.StatusInternalServerError
	}

	return nil, 0
}

func main() {
	godotenv.Load(".env")

	var cfg config

	if err := cfg.init(); err != nil {
		log.Fatal(err)
	}

	databaseFilename := os.Getenv("GOOSE_DBSTRING")

	if databaseFilename == "" {
		log.Fatal("Failed to open database connection: missing GOOSE_DBSTRING")
	}

	db, err := sql.Open("sqlite3", databaseFilename)

	if err != nil {
		log.Fatalf("Failed to open connection to database file '%s'", databaseFilename)
	}

	// Enable SQLite foreign keys.
	if _, err := db.Exec("PRAGMA foreign_keys=ON;"); err != nil {
		log.Fatal("Failed to enable SQLite foreign keys")
	}

	mux := http.NewServeMux()

	rootHandler := http.StripPrefix("/app/", http.FileServer(http.Dir("./templates")))
	mux.Handle("/app/", rootHandler)
	mux.HandleFunc("POST /search", apperr.WithErrors(cfg.searchHandler))
	mux.HandleFunc("GET /routes", apperr.WithErrors(cfg.getRoutes))

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Missing port configuration (PORT)")
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())
}
