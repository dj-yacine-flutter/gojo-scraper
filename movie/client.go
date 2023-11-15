package movie

import (
	"net/http"

	tmdb "github.com/cyruzin/golang-tmdb"
)

type MovieTMDB struct {
	TMDB        *tmdb.Client
	HTTP        *http.Client
	OriginalIMG string
	DecodeIMG   string
}

func NewMovieTMDB(tmdb *tmdb.Client, http *http.Client, Oimg, Dimg string) *MovieTMDB {
	client := &MovieTMDB{
		TMDB:        tmdb,
		HTTP:        http,
		OriginalIMG: Oimg,
		DecodeIMG:   Dimg,
	}

	return client
}
