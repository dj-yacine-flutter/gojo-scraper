package anime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bregydoc/gtranslate"
	"github.com/darenliang/jikan-go"
	"github.com/dj-yacine-flutter/gojo-scraper/models"
	"github.com/dj-yacine-flutter/gojo-scraper/utils"
)

func (server *AnimeScraper) GetAnimeEpisode(w http.ResponseWriter, r *http.Request) {
	seasonMal := r.URL.Query().Get("mal")
	if seasonMal == "" {
		http.Error(w, "Please provide an 'mal' parservereter", http.StatusBadRequest)
		return
	}

	MaleID, err := strconv.Atoi(seasonMal)
	if err != nil {
		http.Error(w, "provide a valid 'mal'", http.StatusInternalServerError)
		return
	}

	seasonTVDB := r.URL.Query().Get("tvdb")
	if seasonTVDB == "" {
		http.Error(w, "Please provide an 'tvdb' parservereter", http.StatusBadRequest)
		return
	}

	TVDBID, err := strconv.Atoi(seasonTVDB)
	if err != nil {
		http.Error(w, "provide a valid 'tvdb'", http.StatusInternalServerError)
		return
	}

	if TVDBID == 0 && MaleID == 0 {
		http.Error(w, "pu at least one ID", http.StatusBadRequest)
		return
	}

	server.LOG.Info().Msgf("Mal ID : %d", MaleID)
	server.LOG.Info().Msgf("TVDB Id : %d", TVDBID)

	start := time.Now()

	var episodes []models.AnimeEpisode

	if TVDBID != 0 {
		season, err := server.TVDB.GetSeasonsByIDExtended(TVDBID)
		if err != nil {
			http.Error(w, fmt.Sprintf("there is no tvdb season with this id : %s", err.Error()), http.StatusNotFound)
			return
		}

		for _, e := range season.Data.Episodes {
			data, err := server.TVDB.GetEpisodeByIDExtanded(e.ID)
			if err != nil {
				http.Error(w, fmt.Sprintf("there is no tvdb episode with this id : %s", err.Error()), http.StatusNotFound)
				return
			}

			var AirDate time.Time
			aired, err := time.Parse(time.DateOnly, data.Data.Aired)
			if err == nil {
				AirDate = aired
			}

			if data.Data.Name == "TBA" || AirDate.After(time.Now()) {
				continue
			}

			var tt string
			tt, err = utils.GetBlurHash(data.Data.Image, "")
			if err != nil {
				tt = ""
			}

			var rating string
			for _, item := range data.Data.ContentRatings {
				if strings.Contains(item.Name, "us") {
					rating = item.Name
					break
				}
			}

			if rating == "" {
				rating = data.Data.ContentRatings[len(data.Data.ContentRatings)-1].Name
			}

			ep := models.AnimeEpisode{
				OriginalTitle:      data.Data.Name,
				Rating:             rating,
				Runtime:            fmt.Sprintf("%dm", data.Data.Runtime),
				Aired:              AirDate.String(),
				EpisodeNumber:      uint(data.Data.Number),
				ThumbnailsPoster:   data.Data.Image,
				ThumbnailsBlurHash: tt,
			}

			var trn string
			Nt := make(map[string]bool)
			for _, item := range data.Data.NameTranslations {
				Nt[item] = true
			}

			if !Nt["eng"] {
				trn, _ = gtranslate.TranslateWithParams(
					utils.CleanUnicode(data.Data.Name),
					gtranslate.TranslationParams{
						From: "auto",
						To:   "en",
					},
				)
			} else {
				dn, err := server.TVDB.GetEpisodeByIDTr(e.ID, "eng")
				if err != nil {
					trn = ""
				} else {
					trn = dn.Data.Name
				}
			}

			var tro string
			Ot := make(map[string]bool)
			for _, item := range data.Data.OverviewTranslations {
				Ot[item] = true
			}

			if !Ot["eng"] {
				tro, _ = gtranslate.TranslateWithParams(
					utils.CleanUnicode(data.Data.Name),
					gtranslate.TranslationParams{
						From: "auto",
						To:   "en",
					},
				)
			} else {
				do, err := server.TVDB.GetEpisodeByIDTr(e.ID, "eng")
				if err != nil {
					tro = ""
				} else {
					tro = do.Data.Overview
				}
			}

			meta := models.MetaData{
				Language: "en",
				Meta: models.Meta{
					Title:    trn,
					Overview: tro,
				},
			}

			ep.EpisodeMetas = make([]models.MetaData, len(models.Languages))

			var newTitle string
			var newOverview string
			for i, lang := range models.Languages {
				newTitle, err = gtranslate.TranslateWithParams(
					meta.Meta.Title,
					gtranslate.TranslationParams{
						From: "en",
						To:   lang,
					},
				)
				if err != nil {
					http.Error(w, fmt.Errorf("error when translate Title to %s: %w ", lang, err).Error(), http.StatusInternalServerError)
					return
				}

				newOverview, err = gtranslate.TranslateWithParams(
					meta.Meta.Overview,
					gtranslate.TranslationParams{
						From: "en",
						To:   lang,
					},
				)
				if err != nil {
					http.Error(w, fmt.Errorf("error when translate Overview to %s: %w ", lang, err).Error(), http.StatusInternalServerError)
					return
				}

				ep.EpisodeMetas[i] = models.MetaData{
					Language: lang,
					Meta: models.Meta{
						Title:    utils.CleanUnicode(newTitle),
						Overview: utils.CleanUnicode(newOverview),
					},
				}
			}

			episodes = append(episodes, ep)
		}
	} else {
		season, err := jikan.GetAnimeEpisodes(MaleID, 1)
		if err != nil {
			http.Error(w, fmt.Sprintf("there is no episodes data with this id : %s", err.Error()), http.StatusNotFound)
			return
		}

		var allEp jikan.AnimeEpisodes
		allEp.Data = append(allEp.Data, season.Data...)

		for i := 2; i <= season.Pagination.LastVisiblePage; i++ {
			time.Sleep(350 * time.Millisecond)

			n, err := jikan.GetAnimeEpisodes(MaleID, i)
			if err != nil {
				http.Error(w, fmt.Sprintf("there is no episodes data : %s", err.Error()), http.StatusNotFound)
				return
			}

			allEp.Data = append(allEp.Data, n.Data...)
		}

		for _, e := range allEp.Data {

			var tt string
			if e.TitleJapanese != "" {
				tt = e.TitleJapanese
			} else {
				tt = e.Title
			}

			ep := models.AnimeEpisode{
				OriginalTitle: tt,
				Aired:         e.Aired.String(),
				EpisodeNumber: uint(e.MalId),
			}

			var tro string
			tro, _ = gtranslate.TranslateWithParams(
				utils.CleanUnicode(tt),
				gtranslate.TranslationParams{
					From: "auto",
					To:   "en",
				},
			)

			ep.EpisodeMetas = make([]models.MetaData, len(models.Languages))

			var newTitle string
			for i, lang := range models.Languages {
				newTitle, err = gtranslate.TranslateWithParams(
					tro,
					gtranslate.TranslationParams{
						From: "en",
						To:   lang,
					},
				)
				if err != nil {
					http.Error(w, fmt.Errorf("error when translate Title to %s: %w ", lang, err).Error(), http.StatusInternalServerError)
					return
				}

				ep.EpisodeMetas[i] = models.MetaData{
					Language: lang,
					Meta: models.Meta{
						Title: utils.CleanUnicode(newTitle),
					},
				}
			}

			episodes = append(episodes, ep)
		}
	}

	end := time.Since(start)

	response, err := json.Marshal(episodes)
	if err != nil {
		http.Error(w, "Internal Server Error:", http.StatusInternalServerError)
		return
	}

	server.LOG.Info().Msgf("Full Time : %v", end)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
