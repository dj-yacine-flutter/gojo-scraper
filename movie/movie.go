package movie

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dj-yacine-flutter/gojo-scraper/models"
	"github.com/dj-yacine-flutter/gojo-scraper/utils"
)

func (am *MovieTMDB) GetMovie(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Please provide an 'id' parameter", http.StatusBadRequest)
		return
	}

	movieID, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		http.Error(w, "provide a valid 'id'", http.StatusInternalServerError)
		return
	}

	anime, err := am.TMDB.GetMovieDetails(int(movieID), nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var releaseDate int
	if anime.ReleaseDate != "" {
		t, err := time.Parse(time.DateOnly, anime.ReleaseDate)
		if err != nil {
			fmt.Println(err)
		}
		releaseDate = t.Year()
	}

	var portriatBlurHash string
	if anime.PosterPath != "" {
		portriatBlurHash, err = utils.GetBlurHash(am.DecodeIMG, anime.PosterPath)
		if err != nil {
			http.Error(w, "cannot get portriatBlurHash blurhash", http.StatusInternalServerError)
			return
		}
	}

	var landscapeBlurHash string
	if anime.BackdropPath != "" {
		landscapeBlurHash, err = utils.GetBlurHash(am.DecodeIMG, anime.BackdropPath)
		if err != nil {
			http.Error(w, "cannot get landscapeBlurHash blurhash", http.StatusInternalServerError)
			return
		}
	}

	var duration string
	if anime.Runtime != 0 {
		duration = fmt.Sprintf("%dm", anime.Runtime)
	}

	var rating string
	md, err := am.TMDB.GetMovieReleaseDates(int(movieID))
	if err == nil {
		if len(md.MovieReleaseDatesResults.Results) > 0 {
			if len(md.MovieReleaseDatesResults.Results[0].ReleaseDates) > 0 {
				for _, r := range md.MovieReleaseDatesResults.Results {
					if strings.Contains(strings.ToLower(r.Iso3166_1), "us") {
						if len(r.ReleaseDates[0].Certification) >= 1 {
							rating = r.ReleaseDates[0].Certification
						}
						continue
					}
				}
			}
		}
	}

	animeData := models.Movie{
		OriginalTitle:     anime.OriginalTitle,
		Aired:             anime.ReleaseDate,
		ReleaseYear:       releaseDate,
		Rating:            rating,
		Duration:          duration,
		PortriatPoster:    am.OriginalIMG + anime.PosterPath,
		PortriatBlurHash:  portriatBlurHash,
		LandscapePoster:   am.OriginalIMG + anime.BackdropPath,
		LandscapeBlurHash: landscapeBlurHash,
	}

	if anime.Title != "" && anime.Overview != "" {
		animeData.AnimeMetas = make([]models.MetaData, len(models.Languages))

		translation, err := utils.Translate(am.HTTP, anime.Overview, "auto", "en")
		if err != nil {
			http.Error(w, fmt.Errorf("error when translate Overview to english: %w ", err).Error(), http.StatusInternalServerError)
			return
		}

		metaData := models.MetaData{
			Language: "en",
			Meta: models.Meta{
				Title:    anime.Title,
				Overview: translation.TranslatedText,
			},
		}

		var newTitle *models.LibreTranslate
		var newOverview *models.LibreTranslate
		for i, lang := range models.Languages {
			newTitle, err = utils.Translate(am.HTTP, metaData.Meta.Title, "en", lang)
			if err != nil {
				http.Error(w, fmt.Errorf("error when translate Title to %s: %w ", lang, err).Error(), http.StatusInternalServerError)
				return
			}

			newOverview, err = utils.Translate(am.HTTP, metaData.Meta.Overview, "en", lang)
			if err != nil {
				http.Error(w, fmt.Errorf("error when translate Overview to %s: %w ", lang, err).Error(), http.StatusInternalServerError)
				return
			}

			animeData.AnimeMetas[i] = models.MetaData{
				Language: lang,
				Meta: models.Meta{
					Title:    newTitle.TranslatedText,
					Overview: newOverview.TranslatedText,
				},
			}
		}
	}

	response, err := json.Marshal(animeData)
	if err != nil {
		http.Error(w, "Internal Server Error:", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
