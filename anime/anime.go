package anime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jikan "github.com/darenliang/jikan-go"
	"github.com/dj-yacine-flutter/gojo-scraper/utils"
)

type TmdbPic struct {
	PortriatBlurHash  string
	PortriatPoster    string
	LandscapeBlurHash string
	LandscapePoster   string
}

func (server *AnimeScraper) postersFromTMDB(Poster, Backdrop string) (TmdbPic TmdbPic) {
	var err error
	if Poster != "" {
		TmdbPic.PortriatBlurHash, err = utils.GetBlurHash(server.DecodeIMG, Poster)
		if err != nil {
			TmdbPic.PortriatBlurHash = ""
			TmdbPic.PortriatPoster = ""
		}
		TmdbPic.PortriatPoster = server.OriginalIMG + Poster
	}

	if Backdrop != "" {
		TmdbPic.LandscapeBlurHash, err = utils.GetBlurHash(server.DecodeIMG, Backdrop)
		if err != nil {
			TmdbPic.LandscapeBlurHash = ""
			TmdbPic.LandscapePoster = ""
		}
		TmdbPic.LandscapePoster = server.OriginalIMG + Backdrop
	}

	return
}

func (server *AnimeScraper) postersFromMAL(Pic, JPG, WEBP string) (TmdbPic TmdbPic) {
	var err error
	img := ""
	if JPG != "" {
		img = JPG
	} else if WEBP != "" {
		img = WEBP
	} else {
		img = fmt.Sprint("https://cdn-eu.anidb.net/images/main/" + Pic)
	}
	TmdbPic.PortriatBlurHash, err = utils.GetBlurHash(img, "")
	if err != nil {
		TmdbPic.PortriatBlurHash = ""
	}

	TmdbPic.PortriatPoster = img
	TmdbPic.LandscapeBlurHash, err = utils.GetBlurHash(img, "")
	if err != nil {
		TmdbPic.LandscapeBlurHash = ""
	}
	TmdbPic.LandscapeBlurHash = ""
	return
}

func (server *AnimeScraper) findResourceByAniDB(anidbID int) (AnimeResources, error) {
	for _, d := range GlobalAnimeResources {
		if anidbID == d.AnidbID {
			return d, nil
		}
	}
	return AnimeResources{}, fmt.Errorf("no resource found with this aniDB ID")
}

func (server *AnimeScraper) findAniDBIDByTitle(malData *jikan.AnimeById) (int, error) {
	for _, v := range GlobalAniDBTitles.Animes {
		for _, title := range v.Titles {
			for _, mt := range append(malData.Data.TitleSynonyms, malData.Data.Title, malData.Data.TitleEnglish) {
				titleMatches := strings.Contains(strings.ToLower(title.Value), strings.ToLower(mt))
				if titleMatches {
					//					fmt.Printf("AniDB title : %s || Mal title : %s\n", title.Value, malData.Data.TitleEnglish)
					aniDBData, err := server.GetAniDBData(v.Aid)
					if err != nil {
						return 0, err
					}
					typeM := strings.Contains(strings.ToLower(aniDBData.Type), strings.ToLower(malData.Data.Type))
					//					fmt.Printf("AniDB type : %s || Mal type : %s\n", aniDBData.Type, malData.Data.Type)
					aniY, err := utils.ExtractYear(aniDBData.Startdate)
					if err != nil {
						return 0, err
					}
					//					fmt.Printf("AniDB year : %d || Mal year : %d\n", aniY, malData.Data.Aired.From.Year())
					yearM := malData.Data.Aired.From.Year() == aniY
					if typeM && yearM {
						return v.Aid, nil
					}
				}
			}
		}
	}
	return 0, nil
}

func (server *AnimeScraper) resolveAniDBID(w http.ResponseWriter, r *http.Request, malData *jikan.AnimeById, links []Link) (int, error) {

	anidbID, err := server.getAniDBID(links)
	if err != nil {
		anidbID, err = server.findAniDBIDByTitle(malData)
		if err != nil {
			return 0, err
		}
	}

	if anidbID == 0 {
		http.Error(w, "there is no AniDB ID for this anime", http.StatusNotFound)
		return 0, nil
	}

	return anidbID, nil
}

func (server *AnimeScraper) GetAnime(w http.ResponseWriter, r *http.Request) {
	mal := r.URL.Query().Get("mal")
	if mal == "" {
		http.Error(w, "Please provide an 'mal' parservereter", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(mal, 0, 0)
	if err != nil {
		http.Error(w, "provide a valid 'mal'", http.StatusInternalServerError)
		return
	}

	malData, err := jikan.GetAnimeById(int(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("there no data with this id: %s", err.Error()), http.StatusNotFound)
		return
	}

	malExt, err := jikan.GetAnimeExternal(int(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("there no external data with this id : %s", err.Error()), http.StatusNotFound)
		return
	}

	var links []Link
	for _, d := range malExt.Data {
		links = append(links, Link{
			URL:  d.Url,
			Name: d.Name,
		})
	}

	anidbID, err := server.resolveAniDBID(w, r, malData, links)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	aniDBData, err := server.GetAniDBData(anidbID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error: %s", err.Error()), http.StatusInternalServerError)
	}

	var releaseDate int
	if malData.Data.Year != 0 {
		releaseDate = malData.Data.Year
	} else {
		releaseDate, err = utils.ExtractYear(aniDBData.Startdate)
		if err != nil {
			releaseDate = 0
		}
	}

	var portriatBlurHash string
	var landscapeBlurHash string
	var portriatPoster string
	var landscapePoster string

	animeResources, err := server.findResourceByAniDB(anidbID)
	if err != nil {
		animePicMAL := server.postersFromMAL(aniDBData.Picture, malData.Data.Images.Jpg.SmallImageUrl, malData.Data.Images.Webp.SmallImageUrl)
		portriatBlurHash = animePicMAL.PortriatBlurHash
		portriatPoster = animePicMAL.PortriatPoster
		landscapeBlurHash = animePicMAL.LandscapeBlurHash
		landscapePoster = animePicMAL.LandscapePoster
	}

	var TMDbID int
	if animeResources.Data.TMDdID != nil {
		tt, err := animeResources.Data.TMDdID.MarshalJSON()
		if err != nil {
			TMDbID = 0
		} else {
			ti, err := strconv.ParseInt(string(tt), 0, 0)
			if err != nil {
				TMDbID = 0
			} else {
				TMDbID = int(ti)
			}
		}
		if TMDbID != 0 {
			if strings.Contains(strings.ToLower(malData.Data.Type), "tv") && strings.Contains(strings.ToLower(aniDBData.Type), "tv") {
				anime, err := server.TMDB.GetTVDetails(TMDbID, nil)
				if err != nil {
					portriatBlurHash = ""
					landscapeBlurHash = ""
				} else {
					animePicTmdb := server.postersFromTMDB(anime.PosterPath, anime.BackdropPath)
					portriatBlurHash = animePicTmdb.PortriatBlurHash
					portriatPoster = animePicTmdb.PortriatPoster
					landscapeBlurHash = animePicTmdb.LandscapeBlurHash
					landscapePoster = animePicTmdb.LandscapePoster
				}
			} else if strings.Contains(strings.ToLower(malData.Data.Type), "movie") && strings.Contains(strings.ToLower(aniDBData.Type), "movie") {
				anime, err := server.TMDB.GetMovieDetails(TMDbID, nil)
				if err != nil {
					portriatBlurHash = ""
					landscapeBlurHash = ""
				} else {
					animePicTmdb := server.postersFromTMDB(anime.PosterPath, anime.BackdropPath)
					portriatBlurHash = animePicTmdb.PortriatBlurHash
					portriatPoster = animePicTmdb.PortriatPoster
					landscapeBlurHash = animePicTmdb.LandscapeBlurHash
					landscapePoster = animePicTmdb.LandscapePoster
				}
			} else {
				querys, err := server.TMDB.GetSearchMulti(malData.Data.TitleEnglish, nil)
				if err != nil {
					animePicMAL := server.postersFromMAL(aniDBData.Picture, malData.Data.Images.Jpg.LargeImageUrl, malData.Data.Images.Webp.LargeImageUrl)
					portriatBlurHash = animePicMAL.PortriatBlurHash
					portriatPoster = animePicMAL.PortriatPoster
					landscapeBlurHash = animePicMAL.LandscapeBlurHash
					landscapePoster = animePicMAL.LandscapePoster
				}

				if querys != nil {
					for _, q := range querys.Results {
						aDate, err := time.Parse(time.DateOnly, aniDBData.Startdate)
						if err != nil {
							continue
						}
						qDate, err := time.Parse(time.DateOnly, q.ReleaseDate)
						if err != nil {
							continue
						}
						if aDate.Equal(qDate) {
							animePicTmdb := server.postersFromTMDB(q.PosterPath, q.BackdropPath)
							portriatBlurHash = animePicTmdb.PortriatBlurHash
							portriatPoster = animePicTmdb.PortriatPoster
							landscapeBlurHash = animePicTmdb.LandscapeBlurHash
							landscapePoster = animePicTmdb.LandscapePoster
						}
					}
				}
			}
		}
	}

	if TMDbID == 0 && portriatBlurHash == "" && landscapeBlurHash == "" {
		for _, r := range aniDBData.Resources.Resource {
			if strings.Contains(r.Type, "44") {
				fmt.Print(r.Externalentity)
				id, err := strconv.ParseInt(r.Externalentity.Identifier[0], 0, 0)
				if err != nil {
					TMDbID = 0
					break
				}
				TMDbID = int(id)
			}
			if TMDbID != 0 {
				break
			}
		}
		if TMDbID != 0 {
			if strings.Contains(strings.ToLower(malData.Data.Type), "tv") && strings.Contains(strings.ToLower(aniDBData.Type), "tv") {
				anime, err := server.TMDB.GetTVDetails(TMDbID, nil)
				if err != nil {
					TMDbID = 0
				} else {
					animePicTmdb := server.postersFromTMDB(anime.PosterPath, anime.BackdropPath)
					portriatBlurHash = animePicTmdb.PortriatBlurHash
					portriatPoster = animePicTmdb.PortriatPoster
					landscapeBlurHash = animePicTmdb.LandscapeBlurHash
					landscapePoster = animePicTmdb.LandscapePoster
				}
			} else if strings.Contains(strings.ToLower(malData.Data.Type), "movie") && strings.Contains(strings.ToLower(aniDBData.Type), "movie") {
				anime, err := server.TMDB.GetMovieDetails(TMDbID, nil)
				if err != nil {
					portriatBlurHash = ""
					landscapeBlurHash = ""
				} else {
					animePicTmdb := server.postersFromTMDB(anime.PosterPath, anime.BackdropPath)
					portriatBlurHash = animePicTmdb.PortriatBlurHash
					portriatPoster = animePicTmdb.PortriatPoster
					landscapeBlurHash = animePicTmdb.LandscapeBlurHash
					landscapePoster = animePicTmdb.LandscapePoster
				}
			} else {
				querys, err := server.TMDB.GetSearchMulti(malData.Data.TitleEnglish, nil)
				if err != nil {
					animePicMAL := server.postersFromMAL(aniDBData.Picture, malData.Data.Images.Jpg.LargeImageUrl, malData.Data.Images.Webp.LargeImageUrl)
					portriatBlurHash = animePicMAL.PortriatBlurHash
					portriatPoster = animePicMAL.PortriatPoster
					landscapeBlurHash = animePicMAL.LandscapeBlurHash
					landscapePoster = animePicMAL.LandscapePoster
				}
				if querys != nil {
					for _, q := range querys.Results {
						aDate, err := time.Parse(time.DateOnly, aniDBData.Startdate)
						if err != nil {
							continue
						}
						qDate, err := time.Parse(time.DateOnly, q.ReleaseDate)
						if err != nil {
							continue
						}
						if aDate.Equal(qDate) {
							animePicTmdb := server.postersFromTMDB(q.PosterPath, q.BackdropPath)
							portriatBlurHash = animePicTmdb.PortriatBlurHash
							portriatPoster = animePicTmdb.PortriatPoster
							landscapeBlurHash = animePicTmdb.LandscapeBlurHash
							landscapePoster = animePicTmdb.LandscapePoster
						}
					}
				}
			}
		} else {
			querys, err := server.TMDB.GetSearchMulti(malData.Data.TitleEnglish, nil)
			if err != nil {
				animePicMAL := server.postersFromMAL(aniDBData.Picture, malData.Data.Images.Jpg.LargeImageUrl, malData.Data.Images.Webp.LargeImageUrl)
				portriatBlurHash = animePicMAL.PortriatBlurHash
				portriatPoster = animePicMAL.PortriatPoster
				landscapeBlurHash = animePicMAL.LandscapeBlurHash
				landscapePoster = animePicMAL.LandscapePoster
			}
			if querys != nil {
				for _, q := range querys.Results {
					fmt.Println(q.ID)
					aDate, err := time.Parse(time.DateOnly, aniDBData.Startdate)
					if err != nil {
						continue
					}
					var qDate time.Time
					if q.ReleaseDate != "" {
						qDate, err = time.Parse(time.DateOnly, q.ReleaseDate)
						if err != nil {
							continue
						}
					} else {
						qDate, err = time.Parse(time.DateOnly, q.FirstAirDate)
						if err != nil {
							continue
						}
					}
					fmt.Println("query Year: ", qDate.String())
					fmt.Println("aniDB Year: ", aDate.String())
					if aDate.Year() == qDate.Year() {
						TMDbID = int(q.ID)
						animePicTmdb := server.postersFromTMDB(q.PosterPath, q.BackdropPath)
						portriatBlurHash = animePicTmdb.PortriatBlurHash
						portriatPoster = animePicTmdb.PortriatPoster
						landscapeBlurHash = animePicTmdb.LandscapeBlurHash
						landscapePoster = animePicTmdb.LandscapePoster
						break
					}
				}
			}
			if TMDbID == 0 {
				animePicMAL := server.postersFromMAL(aniDBData.Picture, malData.Data.Images.Jpg.LargeImageUrl, malData.Data.Images.Webp.LargeImageUrl)
				portriatBlurHash = animePicMAL.PortriatBlurHash
				portriatPoster = animePicMAL.PortriatPoster
				landscapeBlurHash = animePicMAL.LandscapeBlurHash
				landscapePoster = animePicMAL.LandscapePoster
			}

		}
	}

	fmt.Println("tmdbID", TMDbID)
	fmt.Println("releaseDate: ", releaseDate)
	fmt.Println("animeResources: ", animeResources)
	fmt.Println("portriatPoster: ", portriatPoster)
	fmt.Println("portriatBlurHash: ", portriatBlurHash)
	fmt.Println("landscapePoster: ", landscapePoster)
	fmt.Println("landscapeBlurHash: ", landscapeBlurHash)

	response, err := json.Marshal(aniDBData)
	if err != nil {
		http.Error(w, "Internal Server Error:", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	/*


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
			portriatBlurHash, err = utils.GetBlurHash(server.DecodeIMG, anime.PosterPath)
			if err != nil {
				http.Error(w, "cannot get portriatBlurHash blurhash", http.StatusInternalServerError)
				return
			}
		}

		var landscapeBlurHash string
		if anime.BackdropPath != "" {
			landscapeBlurHash, err = utils.GetBlurHash(server.DecodeIMG, anime.BackdropPath)
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
		md, err := server.TMDB.GetMovieReleaseDates(int(id))
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

		animeData := models.Anime{
			OriginalTitle:     anime.OriginalTitle,
			Aired:             anime.ReleaseDate,
			ReleaseYear:       releaseDate,
			Rating:            rating,
			Duration:          duration,
			PortriatPoster:    server.OriginalIMG + anime.PosterPath,
			PortriatBlurHash:  portriatBlurHash,
			LandscapePoster:   server.OriginalIMG + anime.BackdropPath,
			LandscapeBlurHash: landscapeBlurHash,
		}

		if anime.Title != "" && anime.Overview != "" {
			translation, err := gtranslate.TranslateWithParservers(
				anime.Overview,
				gtranslate.TranslationParservers{
					From: "auto",
					To:   "en",
				},
			)
			if err != nil {
				http.Error(w, fmt.Errorf("error when translate Overview to default english: %w ", err).Error(), http.StatusInternalServerError)
				return
			}

			metaData := models.MetaData{
				Language: "en",
				Meta: models.Meta{
					Title:    anime.Title,
					Overview: translation,
				},
			}

			animeData.AnimeMetas = make([]models.MetaData, len(models.Languages))
			var newTitle string
			var newOverview string
			for i, lang := range models.Languages {
				newTitle, err = gtranslate.TranslateWithParservers(
					metaData.Meta.Title,
					gtranslate.TranslationParservers{
						From: "en",
						To:   lang,
					},
				)
				if err != nil {
					http.Error(w, fmt.Errorf("error when translate Title to %s: %w ", lang, err).Error(), http.StatusInternalServerError)
					return
				}

				newOverview, err = gtranslate.TranslateWithParservers(
					metaData.Meta.Overview,
					gtranslate.TranslationParservers{
						From: "en",
						To:   lang,
					},
				)
				if err != nil {
					http.Error(w, fmt.Errorf("error when translate Overview to %s: %w ", lang, err).Error(), http.StatusInternalServerError)
					return
				}

				animeData.AnimeMetas[i] = models.MetaData{
					Language: lang,
					Meta: models.Meta{
						Title:    newTitle,
						Overview: newOverview,
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
	*/
}

/*  translation, err := utils.Translate(server.HTTP, anime.Overview, "auto", "en")
if err != nil {
	http.Error(w, fmt.Errorf("error when translate Overview to default english: %w ", err).Error(), http.StatusInternalServerError)
	return
}

metaData := models.MetaData{
	Language: "en",
	Meta: models.Meta{
		Title:    anime.Title,
		Overview: translation.TranslatedText,
	},
}

animeData.AnimeMetas = make([]models.MetaData, len(models.Languages))
		var newTitle *models.LibreTranslate
   		var newOverview *models.LibreTranslate
   		for i, lang := range models.Languages {
   			newTitle, err = utils.Translate(server.HTTP, metaData.Meta.Title, "en", lang)
   			if err != nil {
   				http.Error(w, fmt.Errorf("error when translate Title to %s: %w ", lang, err).Error(), http.StatusInternalServerError)
   				return
   			}

   			newOverview, err = utils.Translate(server.HTTP, metaData.Meta.Overview, "en", lang)
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
   		} */
