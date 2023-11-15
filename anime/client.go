package anime

import (
	"net/http"

	tmdb "github.com/cyruzin/golang-tmdb"
)

type AnimeScraper struct {
	TMDB        *tmdb.Client
	HTTP        *http.Client
	OriginalIMG string
	DecodeIMG   string
}

func NewAnimeScraper(tmdb *tmdb.Client, http *http.Client, Oimg, Dimg string) *AnimeScraper {
	client := &AnimeScraper{
		TMDB:        tmdb,
		HTTP:        http,
		OriginalIMG: Oimg,
		DecodeIMG:   Dimg,
	}

	return client
}
