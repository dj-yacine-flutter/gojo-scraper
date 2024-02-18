package scrape

import "net/http"

type Scraper struct {
	HTTP *http.Client
}
