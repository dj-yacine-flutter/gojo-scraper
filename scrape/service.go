package scrape

import (
	"errors"
)

var (
	ErrNoDataFound = errors.New("no data found")
	ErrNotOK       = errors.New("status code not 200")
)

type Scraper struct{}
