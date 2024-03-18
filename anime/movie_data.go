package anime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bregydoc/gtranslate"
	jikan "github.com/darenliang/jikan-go"
	"github.com/dj-yacine-flutter/gojo-scraper/models"
	"github.com/dj-yacine-flutter/gojo-scraper/tvdb"
	"github.com/dj-yacine-flutter/gojo-scraper/utils"
)

func (server *AnimeScraper) GetAnimeMovie(w http.ResponseWriter, r *http.Request) {
	mal := r.URL.Query().Get("mal")
	if mal == "" {
		http.Error(w, "Please provide an 'mal' query", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(mal)
	if err != nil {
		http.Error(w, "provide a valid 'mal' id", http.StatusInternalServerError)
		return
	}

	var (
		//AgeRating         string
		PortraitPoster    string
		PortraitBlurHash  string
		LandscapePoster   string
		LandscapeBlurHash string
		TVDbID            int
		OriginalTitle     string
		TMDbID            int
		MalID             int
		IMDbID            string
		Aired             time.Time
		Runtime           string
		Genres            []string
		Studios           []string
		Tags              []string
		PsCs              []string
		Titles            models.Titles
		Posters           []models.Image
		Backdrops         []models.Image
		Logos             []models.Image
		Trailers          []models.Trailer
	)

	malData, err := jikan.GetAnimeById(int(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("there no data with this id: %s", err.Error()), http.StatusNotFound)
		return
	}

	if !strings.Contains(strings.ToLower(malData.Data.Type), "movie") {
		http.Error(w, "this not a anime movie", http.StatusBadRequest)
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

	AniDBID, err := server.searchAniDBID(malData, links)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	AniDBData, err := server.GetAniDBData(AniDBID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	animeResources, _ := server.getResourceByIDs(AniDBID, MalID)

	if animeResources.Data.IMDbID != "" && strings.Contains(animeResources.Data.IMDbID, "tt") {
		IMDbID = animeResources.Data.IMDbID
	}

	for _, c := range AniDBData.Creators.Name {
		if (strings.Contains(strings.ToLower(c.Type), "work") || (strings.Contains(strings.ToLower(c.Type), "animation") && strings.Contains(strings.ToLower(c.Type), "work"))) && !strings.Contains(strings.ToLower(c.Type), "original") {
			Studios = append(Studios, c.Text)
		}
	}

	MalID = malData.Data.MalId
	if malData.Data.TitleJapanese != "" {
		OriginalTitle = malData.Data.TitleJapanese
	} else if malData.Data.TitleEnglish != "" {
		OriginalTitle = malData.Data.TitleEnglish
	} else {
		OriginalTitle = malData.Data.Title
	}

	for _, s := range malData.Data.Studios {
		if s.Name != "" {
			Studios = append(Studios, s.Name)
		}
	}

	var licensors []string
	for _, p := range malData.Data.Licensors {
		licensors = append(licensors, p.Name)
	}

	for _, p := range malData.Data.Producers {
		licensors = append(licensors, p.Name)
	}

	for _, g := range malData.Data.Genres {
		Genres = append(Genres, g.Name)
	}

	for _, g := range malData.Data.ExplicitGenres {
		Genres = append(Genres, g.Name)
	}

	for _, g := range malData.Data.Demographics {
		Genres = append(Genres, g.Name)
	}

	for _, t := range AniDBData.Tags.Tag {
		if utils.CleanTag(t.Name) != "" {
			Tags = append(Tags, strings.ToLower(t.Name))
		}
	}

	Aired = utils.CleanDates([]string{malData.Data.Aired.From.Format(time.DateOnly), AniDBData.Startdate, AniDBData.Enddate})
	if Aired.IsZero() {
		Aired = malData.Data.Aired.From
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if len(malData.Data.TitleSynonyms) > 0 {
			Titles.Others = append(Titles.Others, malData.Data.TitleSynonyms...)
		}
		if malData.Data.TitleJapanese != "" {
			Titles.Offical = append(Titles.Offical, malData.Data.TitleJapanese)
		} else if malData.Data.TitleEnglish != "" {
			Titles.Offical = append(Titles.Offical, malData.Data.TitleEnglish)
		} else if malData.Data.Title != "" {
			Titles.Offical = append(Titles.Offical, malData.Data.Title)
		}
		for _, d := range GlobalAniDBTitles.Animes {
			if AniDBID == d.Aid {
				for _, t := range d.Titles {
					if strings.Contains(t.Type, "main") {
						Titles.Offical = append(Titles.Offical, t.Value)
					} else if strings.Contains(t.Type, "sho") {
						Titles.Short = append(Titles.Short, t.Value)
					} else {
						Titles.Others = append(Titles.Others, t.Value)
					}
				}
			}
		}
	}()

	queries := malData.Data.TitleSynonyms
	queries = append(queries, malData.Data.TitleEnglish, malData.Data.Title)

	if len(queries) > 0 {
		var totalSearch tvdb.Search

		for _, t := range queries {
			movies, err := server.TVDB.GetSearch(t, Aired.Year())
			if err != nil {
				continue
			}
			totalSearch.Data = append(totalSearch.Data, movies.Data...)
		}

		for _, a := range totalSearch.Data {
			server.LOG.Info().Msgf("[TVDB] search ID: %s", a.ID)
			server.LOG.Info().Msgf("[TVDB] search Name: %s", a.Name)
			server.LOG.Info().Msgf("[TVDB] search FirstAirTime: %s", a.FirstAirTime)

			qDate, err := time.Parse(time.DateOnly, a.FirstAirTime)
			if err != nil {
				continue
			}

			var aDate time.Time
			if malData.Data.Aired.From.String() != "" {
				aDate = malData.Data.Aired.From
			} else {
				aDate, err = time.Parse(time.DateOnly, AniDBData.Startdate)
				if err != nil {
					continue
				}
			}

			if aDate.Year() == qDate.Year() && aDate.Month() == qDate.Month() {
				if strings.Contains(a.Type, "movie") {
					newTVDBid, err := strconv.Atoi(a.TvdbID)
					if err != nil {
						continue
					}

					TVDbID = int(newTVDBid)
					movie, err := server.TVDB.GetMovieByIDExtended(TVDbID)
					if err != nil {
						continue
					}
					if movie != nil {
						for _, r := range movie.Data.RemoteIds {
							if strings.Contains(strings.ToLower(r.SourceName), "imdb") && r.SourceName != "" {
								IMDbID = r.ID
							}
						}

						for _, g := range movie.Data.Genres {
							Genres = append(Genres, g.Name)
						}

						if len(Studios) > 0 {
							for _, p := range movie.Data.Companies.Production {
								if p.Name != "" {
									licensors = append(licensors, p.Name)
								}
							}

						} else {
							for _, p := range movie.Data.Companies.Production {
								if p.Name != "" {
									Studios = append(Studios, p.Name)
								}
							}

						}

						for _, d := range movie.Data.Companies.Distributor {
							if d.Name != "" {
								licensors = append(licensors, d.Name)
							}
						}

						for _, d := range movie.Data.Artworks {
							if d.Image != "" {
								if d.Type == 15 {
									bb, _ := utils.GetBlurHash(d.Thumbnail, "")
									Backdrops = append(Backdrops, models.Image{
										Height:    d.Height,
										Width:     d.Width,
										Image:     d.Image,
										Thumbnail: d.Thumbnail,
										BlurHash:  bb,
									})
								} else if d.Type == 14 {
									pp, _ := utils.GetBlurHash(d.Thumbnail, "")
									Posters = append(Posters, models.Image{
										Height:    d.Height,
										Width:     d.Width,
										Image:     d.Image,
										Thumbnail: d.Thumbnail,
										BlurHash:  pp,
									})
								} else if d.Type == 25 {
									ll, _ := utils.GetBlurHash(d.Thumbnail, "")
									Logos = append(Logos, models.Image{
										Height:    d.Height,
										Width:     d.Width,
										Image:     d.Image,
										Thumbnail: d.Thumbnail,
										BlurHash:  ll,
									})
								}
							}
						}

						for _, t := range movie.Data.Trailers {
							if utils.ExtractYTKey(t.URL) != "" {
								Trailers = append(Trailers, models.Trailer{
									Official: true,
									Host:     "YouTube",
									Key:      utils.ExtractYTKey(t.URL),
								})
							}
						}

					}
					break
				} else if strings.Contains(a.Type, "tv") {
					newTVDBid, err := strconv.Atoi(a.TvdbID)
					if err != nil {
						continue
					}

					TVDbID = int(newTVDBid)
					serie, err := server.TVDB.GetSeriesByIDExtanded(TVDbID)
					if err != nil {
						continue
					}

					if serie != nil {
						for _, r := range serie.Data.RemoteIds {
							if strings.Contains(strings.ToLower(r.SourceName), "imdb") && r.SourceName != "" {
								IMDbID = r.ID
							}
						}
						for _, g := range serie.Data.Genres {
							Genres = append(Genres, g.Name)
						}

						for _, p := range serie.Data.Companies {
							if p.Name != "" {
								licensors = append(licensors, p.Name)
							}
						}

						if serie.Data.OriginalNetwork.Name != "" {
							licensors = append(licensors, serie.Data.OriginalNetwork.Name)
						}
						if serie.Data.LatestNetwork.Name != "" {
							licensors = append(licensors, serie.Data.LatestNetwork.Name)
						}
						for _, d := range serie.Data.Artworks {
							if d.Image != "" {
								if d.Type == 15 {
									bb, _ := utils.GetBlurHash(d.Thumbnail, "")
									Backdrops = append(Backdrops, models.Image{
										Height:    d.Height,
										Width:     d.Width,
										Image:     d.Image,
										Thumbnail: d.Thumbnail,
										BlurHash:  bb,
									})
								} else if d.Type == 14 {
									pp, _ := utils.GetBlurHash(d.Thumbnail, "")
									Posters = append(Posters, models.Image{
										Height:    d.Height,
										Width:     d.Width,
										Image:     d.Image,
										Thumbnail: d.Thumbnail,
										BlurHash:  pp,
									})
								} else if d.Type == 25 {
									ll, _ := utils.GetBlurHash(d.Thumbnail, "")
									Logos = append(Logos, models.Image{
										Height:    d.Height,
										Width:     d.Width,
										Image:     d.Image,
										Thumbnail: d.Thumbnail,
										BlurHash:  ll,
									})
								}
							}
						}

						for _, t := range serie.Data.Trailers {
							if utils.ExtractYTKey(t.URL) != "" {
								Trailers = append(Trailers, models.Trailer{
									Official: true,
									Host:     "YouTube",
									Key:      utils.ExtractYTKey(t.URL),
								})
							}
						}
					}
					break
				}
			}
		}
	}

	if TVDbID == 0 && animeResources.Data.TheTVdbID != 0 {
		TVDbID = animeResources.Data.TheTVdbID
		movie, err := server.TVDB.GetMovieByIDExtended(TVDbID)
		if movie != nil && err == nil {
			for _, r := range movie.Data.RemoteIds {
				if strings.Contains(strings.ToLower(r.SourceName), "imdb") && r.SourceName != "" {
					IMDbID = r.ID
				}
			}

			for _, g := range movie.Data.Genres {
				Genres = append(Genres, g.Name)
			}

			if len(Studios) > 0 {
				if len(movie.Data.Companies.Production) > 0 {
					for _, p := range movie.Data.Companies.Production {
						if p.Name != "" {
							licensors = append(licensors, p.Name)
						}
					}
				}
			} else {
				if len(movie.Data.Companies.Production) > 0 {
					for _, p := range movie.Data.Companies.Production {
						if p.Name != "" {
							Studios = append(Studios, p.Name)
						}
					}
				}
			}
			for _, d := range movie.Data.Companies.Distributor {
				if d.Name != "" {
					licensors = append(licensors, d.Name)
				}
			}

			for _, d := range movie.Data.Artworks {
				if d.Image != "" {
					if d.Type == 15 {
						bb, _ := utils.GetBlurHash(d.Thumbnail, "")
						Backdrops = append(Backdrops, models.Image{
							Height:    d.Height,
							Width:     d.Width,
							Image:     d.Image,
							Thumbnail: d.Thumbnail,
							BlurHash:  bb,
						})
					} else if d.Type == 14 {
						pp, _ := utils.GetBlurHash(d.Thumbnail, "")
						Posters = append(Posters, models.Image{
							Height:    d.Height,
							Width:     d.Width,
							Image:     d.Image,
							Thumbnail: d.Thumbnail,
							BlurHash:  pp,
						})
					} else if d.Type == 25 {
						ll, _ := utils.GetBlurHash(d.Thumbnail, "")
						Logos = append(Logos, models.Image{
							Height:    d.Height,
							Width:     d.Width,
							Image:     d.Image,
							Thumbnail: d.Thumbnail,
							BlurHash:  ll,
						})
					}
				}
			}

			for _, t := range movie.Data.Trailers {
				if utils.ExtractYTKey(t.URL) != "" {
					Trailers = append(Trailers, models.Trailer{
						Official: true,
						Host:     "YouTube",
						Key:      utils.ExtractYTKey(t.URL),
					})
				}
			}
		}
	}

	var tmdbIds []int
	for _, r := range AniDBData.Resources.Resource {
		if strings.Contains(r.Type, "44") {
			for _, f := range r.Externalentity {
				for _, v := range f.Identifier {
					id, err := strconv.Atoi(v)
					if err != nil {
						continue
					}
					tmdbIds = append(tmdbIds, id)
				}
			}
		}
	}

	if animeResources.Data.TMDdID != nil {
		tt, err := animeResources.Data.TMDdID.MarshalJSON()
		if err == nil {
			for _, d := range strings.Split(string(tt), ",") {
				ti, err := strconv.Atoi(d)
				if err == nil {
					tmdbIds = append(tmdbIds, int(ti))
				}
			}
		}
	}

	for _, l := range tmdbIds {
		var found bool
		rls, err := server.TMDB.GetMovieReleaseDates(l)
		if err != nil {
			continue
		}
		if len(rls.Results) > 0 {
			var rs []string
			for _, e := range rls.Results {
				if len(e.ReleaseDates) > 0 {
					for _, k := range e.ReleaseDates {
						rs = append(rs, k.ReleaseDate)
					}
				} else {
					continue
				}
			}
			if len(rs) > 0 {
				for _, f := range rs {
					if strings.Contains(f, malData.Data.Aired.From.Format(time.DateOnly)) {
						found = true
						break
					}
				}
			} else {
				continue
			}
		} else {
			continue
		}

		if found {
			TMDbID = l
			break
		}
	}

	if TMDbID == 0 {
		querys, _ := server.TMDB.GetSearchMulti(malData.Data.TitleEnglish, nil)
		if querys != nil {
			for _, q := range querys.Results {
				server.LOG.Info().Msgf("[TMDB] query id: %d\n", q.ID)
				aDate, err := time.Parse(time.DateOnly, AniDBData.Startdate)
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

				if aDate.Year() == qDate.Year() {
					if strings.Contains(strings.ToLower(q.MediaType), "movie") {
						TMDbID = int(q.ID)
						break
					}
				}
			}
		}
	}

	var TMDBRuntime string
	var TMDBTitle string
	if TMDbID != 0 {
		anime, _ := server.TMDB.GetMovieDetails(TMDbID, nil)
		if anime != nil {
			OriginalTitle = anime.OriginalTitle
			TMDBTitle = anime.Title
			TMDBRuntime = fmt.Sprintf("%dm", anime.Runtime)
			PortraitPoster, PortraitBlurHash = server.getTMDBPic(anime.PosterPath)
			LandscapePoster, LandscapeBlurHash = server.getTMDBPic(anime.BackdropPath)

			//AgeRating = server.getTMDBRating(TMDbID)
			for _, g := range anime.Genres {
				if g.Name != "" {
					Genres = append(Genres, g.Name)
				}
			}

			for _, p := range anime.ProductionCompanies {
				if p.Name != "" {
					licensors = append(licensors, p.Name)
				}
			}

			amimg, _ := server.TMDB.GetMovieImages(TMDbID, nil)
			if amimg != nil {
				for _, l := range amimg.Logos {
					if l.FilePath != "" {
						ll, _ := utils.GetBlurHash("https://image.tmdb.org/t/p/w45"+l.FilePath, "")
						Logos = append(Logos, models.Image{
							Height:    l.Height,
							Width:     l.Width,
							Image:     strings.TrimSpace(fmt.Sprintf(server.OriginalIMG + l.FilePath)),
							Thumbnail: strings.TrimSpace("https://image.tmdb.org/t/p/w300" + l.FilePath),
							BlurHash:  ll,
						})
					}
				}
				for _, b := range amimg.Backdrops {
					if b.FilePath != "" {
						bb, _ := utils.GetBlurHash(server.DecodeIMG+b.FilePath, "")
						Backdrops = append(Backdrops, models.Image{
							Height:    b.Height,
							Width:     b.Width,
							Image:     strings.TrimSpace(fmt.Sprintf(server.OriginalIMG + b.FilePath)),
							Thumbnail: strings.TrimSpace("https://image.tmdb.org/t/p/w500" + b.FilePath),
							BlurHash:  bb,
						})
					}
				}
				for _, p := range amimg.Posters {
					if p.FilePath != "" {
						pp, _ := utils.GetBlurHash(server.DecodeIMG+p.FilePath, "")
						Posters = append(Posters, models.Image{
							Height:    p.Height,
							Width:     p.Width,
							Image:     strings.TrimSpace(fmt.Sprintf(server.OriginalIMG + p.FilePath)),
							Thumbnail: strings.TrimSpace("https://image.tmdb.org/t/p/w342" + p.FilePath),
							BlurHash:  pp,
						})

					}
				}
			}

			trl, _ := server.TMDB.GetMovieVideos(TMDbID, nil)
			for _, t := range trl.Results {
				if strings.Contains(strings.ToLower(t.Site), "youtube") {
					if t.Key != "" {
						Trailers = append(Trailers, models.Trailer{
							Official: true,
							Host:     "YouTube",
							Key:      t.Key,
						})
					}

				}
			}

		}
	}

	if PortraitBlurHash == "" {
		PortraitPoster, PortraitBlurHash = server.getMainPic(AniDBData.Picture, malData.Data.Images)
	}

	if Runtime == "" {
		var titles []string
		titles = append(titles, TMDBTitle, malData.Data.Title, malData.Data.TitleEnglish, malData.Data.TitleJapanese)
		var h int
		for _, e := range AniDBData.Episodes.Episode {
			if strings.Contains(e.Epno.Type, "1") {
				for _, u := range e.Title {
					for _, t := range titles {
						var el string
						if strings.Contains(utils.CleanTitle(t), utils.CleanTitle(u.Text)) {
							el = e.Length
						} else {
							airdate, err := time.Parse(time.DateOnly, e.Airdate)
							if err == nil {
								if malData.Data.Aired.From.Year() == airdate.Year() && malData.Data.Aired.From.Month() == airdate.Month() {
									el = e.Length
								}
							}
						}
						if el != "" {
							b, err := strconv.Atoi(el)
							if err != nil {
								continue
							}
							if h < int(b) && b != 0 {
								h = int(b)
							}
							break
						}
					}
					if h != 0 {
						break
					}
				}
			}
		}
		if h != 0 {
			Runtime = utils.CleanRuntime(fmt.Sprintf("%dm", h))
		} else {
			run, err := time.ParseDuration(utils.CleanRuntime(TMDBRuntime))
			if err != nil {
				run, err = time.ParseDuration(utils.CleanRuntime(malData.Data.Duration))
				if err == nil {
					Runtime = fmt.Sprintf("%fm", run.Minutes())
				}
			} else {
				Runtime = fmt.Sprintf("%fm", run.Minutes())
			}
		}
	}

	for _, s := range utils.CleanDuplicates(utils.CleanStringArray(Studios)) {
		for _, r := range utils.CleanDuplicates(utils.CleanStringArray(licensors)) {
			if !strings.Contains(utils.CleanTitle(r), utils.CleanTitle(s)) {
				if strings.TrimSpace(r) != "ltd." {
					PsCs = append(PsCs, r)
				}
			}
		}
	}

	var AnimePlanetID string
	animePlanetByte, err := animeResources.Data.AnimePlanetID.MarshalJSON()
	if err == nil {
		AnimePlanetID = string(animePlanetByte)
		AnimePlanetID = strings.ReplaceAll(AnimePlanetID, "\"", "")
	}

	var (
		LivechartID int
		AnysearchID int
		KitsuID     int
		NotifyMoeID string
		AnilistID   int
	)

	wg.Add(5)
	go func() {
		defer wg.Done()
		LivechartID = server.Livechart(animeResources.Data.LivechartID, OriginalTitle, Aired)
	}()
	go func() {
		defer wg.Done()
		AnysearchID = server.Anysearch(animeResources.Data.AnisearchID, malData.Data.TitleEnglish, OriginalTitle, Aired)
	}()
	go func() {
		defer wg.Done()
		KitsuID = server.Kitsu(animeResources.Data.KitsuID, OriginalTitle, Aired)
	}()
	go func() {
		defer wg.Done()
		NotifyMoeID = server.NotifyMoe(utils.CleanResText(animeResources.Data.NotifyMoeID), malData.Data.Title, Aired)
	}()
	go func() {
		defer wg.Done()
		AnilistID = server.Anylist(malData.Data.MalId)
	}()
	wg.Wait()

	animeData := models.AnimeMovie{
		OriginalTitle:       OriginalTitle,
		Aired:               Aired.Format(time.DateOnly),
		Runtime:             Runtime,
		ReleaseYear:         Aired.Year(),
		Rating:              utils.CleanUnicode(malData.Data.Rating),
		PortraitPoster:      PortraitPoster,
		PortraitBlurHash:    PortraitBlurHash,
		LandscapePoster:     LandscapePoster,
		LandscapeBlurHash:   LandscapeBlurHash,
		Genres:              utils.CleanStringArray(Genres),
		Studios:             utils.CleanDuplicates(utils.CleanStringArray(Studios)),
		Tags:                utils.CleanStringArray(Tags),
		ProductionCompanies: utils.CleanDuplicates(PsCs),
		Titles:              Titles,
		Backdrops:           utils.CleanImages(Backdrops),
		Posters:             utils.CleanImages(Posters),
		Logos:               utils.CleanImages(Logos),
		Trailers:            utils.CleanTrailers(Trailers),
		AnimeResources: models.MovieAnimeResources{
			LivechartID:   LivechartID,
			AnimePlanetID: utils.CleanResText(AnimePlanetID),
			AnisearchID:   AnysearchID,
			AnidbID:       AniDBID,
			KitsuID:       KitsuID,
			MalID:         MalID,
			NotifyMoeID:   NotifyMoeID,
			AnilistID:     AnilistID,
			TVDbID:        TVDbID,
			IMDbID:        utils.CleanResText(IMDbID),
			TMDbID:        TMDbID,
			Type:          utils.CleanResText(animeResources.Data.Type),
		},
	}

	server.LOG.Info().Msgf("Licensors: %v", licensors)
	server.LOG.Info().Msgf("TMDBID: %d", TMDbID)
	server.LOG.Info().Msgf("TVDBID: %d", TVDbID)
	server.LOG.Info().Msgf("Aired: %v", Aired)
	server.LOG.Info().Msgf("Runtime: %s", Runtime)
	server.LOG.Info().Msgf("AniDB Episodes: %d", len(AniDBData.Episodes.Episode))
	server.LOG.Info().Msgf("OriginalTitle: %s", OriginalTitle)
	server.LOG.Info().Msgf("PortraitPoster: %s", PortraitPoster)
	server.LOG.Info().Msgf("PortraitBlurHash: %s", PortraitBlurHash)
	server.LOG.Info().Msgf("LandscapePoster: %s", LandscapePoster)
	server.LOG.Info().Msgf("LandscapeBlurHash: %s", LandscapeBlurHash)

	var TTitle string
	var enT bool
	if malData.Data.TitleEnglish != "" {
		enT = true
		TTitle = malData.Data.TitleEnglish
	} else {
		TTitle = malData.Data.Title
	}

	if TTitle != "" && malData.Data.Synopsis != "" {
		var translationTitle string
		if !enT {
			translationTitle, err = gtranslate.TranslateWithParams(
				utils.CleanUnicode(TTitle),
				gtranslate.TranslationParams{
					From: "auto",
					To:   "en",
				},
			)
			if err != nil {
				http.Error(w, fmt.Errorf("error when translate TTitle to default english: %w ", err).Error(), http.StatusInternalServerError)
				return
			}
		} else {
			translationTitle = utils.CleanUnicode(TTitle)
		}

		translationOverview, err := gtranslate.TranslateWithParams(
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
				Title:    translationTitle,
				Overview: translationOverview,
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
					Title:    utils.CleanUnicode(newTitle),
					Overview: utils.CleanUnicode(newOverview),
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
