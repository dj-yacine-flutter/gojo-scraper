package anime

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	jikan "github.com/darenliang/jikan-go"
	"github.com/dj-yacine-flutter/gojo-scraper/models"
	"github.com/dj-yacine-flutter/gojo-scraper/utils"
)

type AnimePic struct {
	PortriatBlurHash  string
	PortriatPoster    string
	LandscapeBlurHash string
	LandscapePoster   string
}

func (server *AnimeScraper) postersFromTMDB(Poster, Backdrop string) (AnimePic AnimePic) {
	var err error
	if Poster != "" {
		AnimePic.PortriatBlurHash, err = utils.GetBlurHash(server.DecodeIMG, Poster)
		if err != nil {
			AnimePic.PortriatBlurHash = ""
			AnimePic.PortriatPoster = ""
		}
		AnimePic.PortriatPoster = server.OriginalIMG + Poster
	}

	if Backdrop != "" {
		AnimePic.LandscapeBlurHash, err = utils.GetBlurHash(server.DecodeIMG, Backdrop)
		if err != nil {
			AnimePic.LandscapeBlurHash = ""
			AnimePic.LandscapePoster = ""
		}
		AnimePic.LandscapePoster = server.OriginalIMG + Backdrop
	}

	return
}

func (server *AnimeScraper) postersFromMAL(Pic, JPG, WEBP string) (AnimePic AnimePic) {
	var err error
	img := ""
	if JPG != "" {
		img = JPG
	} else if WEBP != "" {
		img = WEBP
	} else {
		img = fmt.Sprint("https://cdn-eu.anidb.net/images/main/" + Pic)
	}
	AnimePic.PortriatBlurHash, err = utils.GetBlurHash(img, "")
	if err != nil {
		AnimePic.PortriatBlurHash = ""
	}

	AnimePic.PortriatPoster = img
	AnimePic.LandscapeBlurHash, err = utils.GetBlurHash(img, "")
	if err != nil {
		AnimePic.LandscapeBlurHash = ""
	}
	AnimePic.LandscapeBlurHash = ""
	return
}

func (server *AnimeScraper) findResourceByAniDBandmalID(anidbID, malID int) (AnimeResources, error) {
	for _, d := range GlobalAnimeResources {
		if anidbID == d.AnidbID {
			if d.Data.MalID != 0 && malID != 0 {
				if d.Data.MalID == malID {
					return d, nil
				}
			} else {
				return d, nil
			}
		}
	}
	return AnimeResources{}, fmt.Errorf("no resource found for this anime")
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

	anidbID, err := server.searchAniDBID(malData, links)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	aniDBData, err := server.GetAniDBData(anidbID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error: %s", err.Error()), http.StatusInternalServerError)
	}

	fmt.Printf("AniDB Episodes: %d\n", len(aniDBData.Episodes.Episode))

	var releaseDate int
	if malData.Data.Year != 0 {
		releaseDate = malData.Data.Year
	} else {
		releaseDate, err = utils.ExtractYear(aniDBData.Startdate)
		if err != nil {
			releaseDate = 0
		}
	}

	var ageRating string
	if malData.Data.Rating != "" {
		ageRating, err = utils.CleanRating(malData.Data.Rating)
		if err != nil {
			ageRating = ""
		}
	}

	var portriatBlurHash string
	var landscapeBlurHash string
	var portriatPoster string
	var landscapePoster string

	animeResources, err := server.findResourceByAniDBandmalID(anidbID, malData.Data.MalId)
	if err != nil {
		animePicMAL := server.postersFromMAL(aniDBData.Picture, malData.Data.Images.Jpg.SmallImageUrl, malData.Data.Images.Webp.SmallImageUrl)
		portriatBlurHash = animePicMAL.PortriatBlurHash
		portriatPoster = animePicMAL.PortriatPoster
		landscapeBlurHash = animePicMAL.LandscapeBlurHash
		landscapePoster = animePicMAL.LandscapePoster
	}

	var originalTitle string
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
					originalTitle = anime.OriginalName
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
					originalTitle = anime.OriginalTitle
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
						if aDate.Year() == qDate.Year() {
							if q.OriginalName != "" {
								originalTitle = q.OriginalName
							} else if q.OriginalTitle != "" {
								originalTitle = q.OriginalTitle
							}
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
	fmt.Printf("TMDB before 0 : %d\n", TMDbID)

	if TMDbID == 0 && portriatBlurHash == "" && landscapeBlurHash == "" {
		for _, r := range aniDBData.Resources.Resource {
			if strings.Contains(r.Type, "44") {
				fmt.Println("Externalentity in aniDB :", r.Externalentity.Text)
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
		fmt.Printf("TMDB after 0 : %d\n", TMDbID)
		if TMDbID != 0 {
			if strings.Contains(strings.ToLower(malData.Data.Type), "tv") && strings.Contains(strings.ToLower(aniDBData.Type), "tv") {
				anime, err := server.TMDB.GetTVDetails(TMDbID, nil)
				if err != nil {
					TMDbID = 0
				} else {
					originalTitle = anime.OriginalName
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
					originalTitle = anime.OriginalTitle
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
						if aDate.Year() == qDate.Year() {
							if q.OriginalName != "" {
								originalTitle = q.OriginalName
							} else if q.OriginalTitle != "" {
								originalTitle = q.OriginalTitle
							}
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
					fmt.Println("query id :", q.ID)
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

	fmt.Println("TVDBID", animeResources.Data.TheTVdbID)
	if animeResources.Data.TheTVdbID != 0 {
		if strings.Contains(strings.ToLower(malData.Data.Type), "tv") && strings.Contains(strings.ToLower(aniDBData.Type), "tv") {
			fmt.Println("TV TVDBID", animeResources.Data.TheTVdbID)

			ss, err := server.TVDB.GetSeriesByIDExtanded(animeResources.Data.TheTVdbID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Fatal(err)

				return
			}

			fmt.Println(ss.Data.Name)
			fmt.Println(ss.Data.Year)
			fmt.Printf("seasons : %d\n", len(ss.Data.Seasons))
			for _, s := range ss.Data.Seasons {
				if s.ID != 0 {
					rr, err := server.TVDB.GetSeasonsByIDExtended(s.ID)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						log.Fatal(err)

						return
					}
					if strings.Contains(rr.Data.Type.Type, "official") {
						fmt.Printf("season ID: %d\n", s.ID)
						fmt.Printf("season Number: %d\n", s.Number)
						fmt.Printf("season Type: %s\n", rr.Data.Type.Type)
						fmt.Printf("season Year: %s\n", rr.Data.Year)
						fmt.Printf("TVDB Episodes: %d\n", len(rr.Data.Episodes))
						aYear, err := utils.ExtractYear(aniDBData.Startdate)
						if err != nil {
							continue
						}

						rYear, err := utils.ExtractYear(rr.Data.Year)
						if err != nil {
							continue
						}
						if aYear == rYear {
							if landscapePoster == "" {
								if len(ss.Data.Artworks) > 0 {
									for i := len(ss.Data.Artworks) - 1; i >= 0; i-- {
										if ss.Data.Artworks[i].Image != "" && strings.Contains(ss.Data.Artworks[i].Image, "backgrounds") {
											landscapePoster = ss.Data.Artworks[i].Image
											landscapeBlurHash, err = utils.GetBlurHash(ss.Data.Artworks[i].Image, "")
											if err != nil {
												landscapeBlurHash = ""
											} else {
												break
											}
										}
									}
								}
							}
						}
					}
				}
			}
		} else if strings.Contains(strings.ToLower(malData.Data.Type), "movie") && strings.Contains(strings.ToLower(aniDBData.Type), "movie") {
			fmt.Println("Movie TVDBID", animeResources.Data.TheTVdbID)

			mm, err := server.TVDB.GetMovieByIDExtended(animeResources.Data.TheTVdbID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Fatal(err)
				return
			}

			fmt.Println("TVDB MOVIE: ", mm)
		}
	}

	fmt.Println("tmdbID", TMDbID)
	fmt.Println("releaseDate: ", releaseDate)
	fmt.Println("animeResources: ", animeResources)
	fmt.Println("portriatPoster: ", portriatPoster)
	fmt.Println("portriatBlurHash: ", portriatBlurHash)
	fmt.Println("landscapePoster: ", landscapePoster)
	fmt.Println("landscapeBlurHash: ", landscapeBlurHash)

	var animePlanetID string
	animePlanetByt, err := animeResources.Data.AnimePlanetID.MarshalJSON()
	if err != nil {
		animePlanetID = ""
	} else {
		animePlanetID = string(animePlanetByt)
		animePlanetID = strings.ReplaceAll(animePlanetID, "\"", "")
	}

	fmt.Println("OriginalTitle: ", originalTitle)
	if originalTitle == "" {
		originalTitle = malData.Data.TitleJapanese
	}

	animeData := models.Anime{
		OriginalTitle:     originalTitle,
		Aired:             malData.Data.Aired.From.String(),
		ReleaseYear:       releaseDate,
		Rating:            ageRating,
		PortriatPoster:    portriatPoster,
		PortriatBlurHash:  portriatBlurHash,
		LandscapePoster:   landscapePoster,
		LandscapeBlurHash: landscapeBlurHash,
		AnimeResources: models.AnimeResources{
			LivechartID:   animeResources.Data.AnisearchID,
			AnimePlanetID: animePlanetID,
			AnisearchID:   animeResources.Data.AnisearchID,
			AnidbID:       animeResources.Data.AnidbID,
			KitsuID:       animeResources.Data.KitsuID,
			MalID:         animeResources.Data.MalID,
			NotifyMoeID:   animeResources.Data.NotifyMoeID,
			AnilistID:     animeResources.Data.AnilistID,
			ThetvdbID:     animeResources.Data.TheTVdbID,
			ImdbID:        animeResources.Data.IMDbID,
			ThemoviedbID:  TMDbID,
			Type:          animeResources.Data.Type,
		},
	}

	for _, d := range GlobalAniDBTitles.Animes {
		if anidbID == d.Aid {
			for _, t := range d.Titles {
				if strings.Contains(t.Type, "offi") {
					animeData.Titles.Offical = append(animeData.Titles.Offical, t.Value)
				} else if strings.Contains(t.Type, "shor") {
					animeData.Titles.Short = append(animeData.Titles.Short, t.Value)
				} else {
					animeData.Titles.Others = append(animeData.Titles.Others, t.Value)
				}
			}
		}
	}

	/* 	if malData.Data.TitleEnglish != "" && malData.Data.Synopsis != "" {
		translation, err := gtranslate.TranslateWithParams(
			utils.CleanOverview(malData.Data.Synopsis),
			gtranslate.TranslationParams{
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
				Title:    malData.Data.TitleEnglish,
				Overview: translation,
			},
		}

		animeData.AnimeMetas = make([]models.MetaData, len(models.Languages))
		var newTitle string
		var newOverview string
		for i, lang := range models.Languages {
			newTitle, err = gtranslate.TranslateWithParams(
				metaData.Meta.Title,
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
				metaData.Meta.Overview,
				gtranslate.TranslationParams{
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
	} */

	response, err := json.Marshal(animeData)
	if err != nil {
		http.Error(w, "Internal Server Error:", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
