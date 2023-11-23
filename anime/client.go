package anime

import (
	"net/http"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/dj-yacine-flutter/gojo-scraper/tvdb"
)

type AnimeScraper struct {
	TMDB        *tmdb.Client
	HTTP        *http.Client
	TVDB        *tvdb.Client
	Proxy       string
	OriginalIMG string
	DecodeIMG   string
}

func NewAnimeScraper(tmdb *tmdb.Client, http *http.Client, tvdb *tvdb.Client, Oimg, Dimg, Proxy string) *AnimeScraper {
	client := &AnimeScraper{
		TMDB:        tmdb,
		HTTP:        http,
		TVDB:        tvdb,
		OriginalIMG: Oimg,
		DecodeIMG:   Dimg,
		Proxy:       Proxy,
	}

	return client
}
