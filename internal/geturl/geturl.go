// Geturl manages GET requests to MTA API endpoints, along with
// queries.
package geturl

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

func Call(_url string, queries map[string]string) ([]byte, error) {
	url, err := url.Parse(_url)

	if err != nil {
		return nil, err
	}

	// Prepare the required queries.
	queryValues := url.Query()

	for queryParam, value := range queries {
		queryValues.Add(queryParam, value)
	}

	url.RawQuery = queryValues.Encode()

	// Make the GET request.
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", url.String(), nil)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
