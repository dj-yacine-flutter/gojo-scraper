package anime

import (
	"net/http"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/dj-yacine-flutter/gojo-scraper/scrape"
	"github.com/dj-yacine-flutter/gojo-scraper/tvdb"
	"github.com/rs/zerolog"
)

type AnimeScraper struct {
	TMDB        *tmdb.Client
	HTTP        *http.Client
	TVDB        *tvdb.Client
	LOG         *zerolog.Logger
	SCRP        *scrape.Scraper
	OriginalIMG string
	DecodeIMG   string
}

func NewAnimeScraper(tmdb *tmdb.Client, http *http.Client, tvdb *tvdb.Client, logger *zerolog.Logger, Oimg, Dimg string) *AnimeScraper {
	return &AnimeScraper{
		TMDB:        tmdb,
		HTTP:        http,
		TVDB:        tvdb,
		LOG:         logger,
		SCRP:        &scrape.Scraper{},
		OriginalIMG: Oimg,
		DecodeIMG:   Dimg,
	}
}
