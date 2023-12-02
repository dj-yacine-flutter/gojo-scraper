package anime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jikan "github.com/darenliang/jikan-go"
	"github.com/dj-yacine-flutter/gojo-scraper/models"
	"github.com/dj-yacine-flutter/gojo-scraper/tvdb"
	"github.com/dj-yacine-flutter/gojo-scraper/utils"
)

func (server *AnimeScraper) GetAnimeMovie(w http.ResponseWriter, r *http.Request) {
	mal := r.URL.Query().Get("mal")
	if mal == "" {
		http.Error(w, "Please provide an 'mal' parservereter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(mal)
	if err != nil {
		http.Error(w, "provide a valid 'mal'", http.StatusInternalServerError)
		return
	}

	var (
		ReleaseYear       int
		AgeRating         string
		PortriatPoster    string
		PortriatBlurHash  string
		LandscapePoster   string
		LandscapeBlurHash string
		AnimePlanetID     string
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

	MalID = malData.Data.MalId
	animeResources, err := server.getResourceByIDs(AniDBID, MalID)
	if err != nil {
		err = nil
		animeResources = AnimeResources{}
	}

	if animeResources.Data.IMDbID != "" && strings.Contains(animeResources.Data.IMDbID, "tt") {
		IMDbID = animeResources.Data.IMDbID
	}

	server.getMalPic(AniDBData.Picture, malData.Data.Images.Jpg.LargeImageUrl, malData.Data.Images.Webp.LargeImageUrl, &PortriatBlurHash, &PortriatPoster)

	if malData.Data.TitleJapanese != "" {
		OriginalTitle = malData.Data.TitleJapanese
	} else if malData.Data.TitleEnglish != "" {
		OriginalTitle = malData.Data.TitleEnglish
	} else {
		OriginalTitle = malData.Data.Title
	}

	if len(AniDBData.Creators.Name) > 0 {
		for _, c := range AniDBData.Creators.Name {
			if (strings.Contains(strings.ToLower(c.Type), "work") || (strings.Contains(strings.ToLower(c.Type), "animation") && strings.Contains(strings.ToLower(c.Type), "work"))) && !strings.Contains(strings.ToLower(c.Type), "original") {
				Studios = append(Studios, c.Text)
			}
		}
	}

	if len(malData.Data.Studios) > 0 {
		for _, s := range malData.Data.Studios {
			if s.Name != "" {
				Studios = append(Studios, s.Name)
			}
		}
	}

	var licensors []string
	if len(malData.Data.Licensors) > 0 {
		for _, p := range malData.Data.Licensors {
			licensors = append(licensors, p.Name)
		}
	}
	if len(malData.Data.Producers) > 0 {
		for _, p := range malData.Data.Producers {
			licensors = append(licensors, p.Name)
		}
	}

	var queries []string
	var totalSearch tvdb.Search
	queries = append(queries, malData.Data.TitleEnglish, malData.Data.Title)
	queries = append(queries, malData.Data.TitleSynonyms...)

	if len(queries) > 0 {
		for _, t := range queries {
			movies, err := server.TVDB.GetSearch(t, ReleaseYear)
			if err != nil {
				continue
			}
			totalSearch.Data = append(totalSearch.Data, movies.Data...)
		}
	}

	for _, a := range totalSearch.Data {
		server.LOG.Info().Msgf("search ID: %s", a.ID)
		server.LOG.Info().Msgf("search TVDB: %s", a.TvdbID)
		server.LOG.Info().Msgf("search Name: %s", a.Name)
		server.LOG.Info().Msgf("search Year: %s", a.Year)
		server.LOG.Info().Msgf("search ExtendedTitle: %s", a.ExtendedTitle)
		server.LOG.Info().Msgf("search FirstAirTime: %s", a.FirstAirTime)

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

					if len(movie.Data.Genres) > 0 {
						for _, g := range movie.Data.Genres {
							Genres = append(Genres, g.Name)
						}
					}
					if len(Studios) == 0 {
						if len(movie.Data.Companies.Production) > 0 {
							for _, p := range movie.Data.Companies.Production {
								if p.Name != "" {
									Studios = append(Studios, p.Name)
								}
							}
						}
					} else {
						if len(movie.Data.Companies.Production) > 0 {
							for _, p := range movie.Data.Companies.Production {
								if p.Name != "" {
									licensors = append(licensors, p.Name)
								}
							}
						}
					}

					if len(movie.Data.Companies.Distributor) > 0 {
						for _, d := range movie.Data.Companies.Distributor {
							if d.Name != "" {
								licensors = append(licensors, d.Name)
							}
						}
					}
					if len(movie.Data.Artworks) > 0 {
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
					}
					if len(movie.Data.Trailers) > 0 {
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
					if len(serie.Data.Genres) > 0 {
						for _, g := range serie.Data.Genres {
							Genres = append(Genres, g.Name)
						}
					}
					if len(serie.Data.Companies) > 0 {
						for _, p := range serie.Data.Companies {
							if p.Name != "" {
								licensors = append(licensors, p.Name)
							}
						}
					}
					if serie.Data.OriginalNetwork.Name != "" {
						licensors = append(licensors, serie.Data.OriginalNetwork.Name)
					}
					if serie.Data.LatestNetwork.Name != "" {
						licensors = append(licensors, serie.Data.LatestNetwork.Name)
					}
					if len(serie.Data.Artworks) > 0 {
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
					}
					if len(serie.Data.Trailers) > 0 {
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
				}
				break
			}
		}
	}

	if TVDbID == 0 {
		if animeResources.Data.TheTVdbID != 0 {
			movie, err := server.TVDB.GetMovieByIDExtended(animeResources.Data.TheTVdbID)
			if err != nil {
				TVDbID = animeResources.Data.TheTVdbID
			}

			if movie != nil {
				for _, r := range movie.Data.RemoteIds {
					if strings.Contains(strings.ToLower(r.SourceName), "imdb") && r.SourceName != "" {
						IMDbID = r.ID
					}
				}

				if len(movie.Data.Genres) > 0 {
					for _, g := range movie.Data.Genres {
						Genres = append(Genres, g.Name)
					}
				}
				if len(Studios) == 0 {
					if len(movie.Data.Companies.Production) > 0 {
						for _, p := range movie.Data.Companies.Production {
							if p.Name != "" {
								Studios = append(Studios, p.Name)
							}
						}
					}
				} else {
					if len(movie.Data.Companies.Production) > 0 {
						for _, p := range movie.Data.Companies.Production {
							if p.Name != "" {
								licensors = append(licensors, p.Name)
							}
						}
					}
				}
				if len(movie.Data.Companies.Distributor) > 0 {
					for _, d := range movie.Data.Companies.Distributor {
						if d.Name != "" {
							licensors = append(licensors, d.Name)
						}
					}
				}
				if len(movie.Data.Artworks) > 0 {
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
				}
				if len(movie.Data.Trailers) > 0 {
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
		}
	}

	var TMDBIDs []int
	for _, r := range AniDBData.Resources.Resource {
		if strings.Contains(r.Type, "44") {
			if len(r.Externalentity) > 0 {
				for _, f := range r.Externalentity {
					for _, v := range f.Identifier {
						id, err := strconv.Atoi(v)
						if err != nil {
							continue
						}
						TMDBIDs = append(TMDBIDs, id)
					}
				}
			}
		}
	}

	if animeResources.Data.TMDdID != nil {
		tt, err := animeResources.Data.TMDdID.MarshalJSON()
		if err != nil {
			TMDbID = 0
		} else {
			for _, d := range strings.Split(string(tt), ",") {
				ti, err := strconv.Atoi(d)
				if err != nil {
					TMDbID = 0
				} else {
					TMDBIDs = append(TMDBIDs, int(ti))
				}
			}
		}
	}

	var TMDBRuntime string
	var TMDBTitle string
	if len(TMDBIDs) > 0 {
		for _, l := range TMDBIDs {
			TMDbID = l
			anime, err := server.TMDB.GetMovieDetails(TMDbID, nil)
			if err != nil {
				PortriatBlurHash = ""
				LandscapeBlurHash = ""
				TMDbID = 0
			} else {
				var rd bool
				if anime.ReleaseDate != "" {
					eDate, err := time.Parse(time.DateOnly, anime.ReleaseDate)
					if err != nil {
						PortriatBlurHash = ""
						LandscapeBlurHash = ""
						TMDbID = 0
					}
					qDate := malData.Data.Aired.From
					if eDate.Year() == qDate.Year() && eDate.Month() == qDate.Month() {
						rd = true
					}

				} else {
					newAnime, err := server.TMDB.GetMovieReleaseDates(TMDbID)
					if err != nil {
						PortriatBlurHash = ""
						LandscapeBlurHash = ""
						TMDbID = 0
					}
					if len(newAnime.Results) > 0 {
						var rs []string
						for _, e := range newAnime.Results {
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
									rd = true
									break
								}
							}
						} else {
							PortriatBlurHash = ""
							LandscapeBlurHash = ""
							TMDbID = 0
						}

					}
				}

				if rd {
					if OriginalTitle == "" {
						OriginalTitle = anime.OriginalTitle
					}
					TMDBRuntime = fmt.Sprintf("%dm", anime.Runtime)
					TMDBTitle = anime.Title
					server.getTMDBPic(anime.PosterPath, anime.BackdropPath, &PortriatBlurHash, &PortriatPoster, &LandscapeBlurHash, &LandscapePoster)
					server.getTMDBRating(TMDbID, &AgeRating)
					if len(anime.Genres) > 0 {
						for _, g := range anime.Genres {
							if g.Name != "" {
								Genres = append(Genres, g.Name)
							}
						}
					}
					if len(anime.ProductionCompanies) > 0 {
						for _, p := range anime.ProductionCompanies {
							if p.Name != "" {
								licensors = append(licensors, p.Name)
							}
						}
					}

					amimg, _ := server.TMDB.GetMovieImages(TMDbID, nil)
					if err == nil {
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
					}
					tttt, _ := server.TMDB.GetMovieVideos(TMDbID, nil)
					if len(tttt.Results) > 0 {
						for _, t := range tttt.Results {
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
					break
				}
			}
		}
	}

	if TMDbID != 0 && PortriatBlurHash == "" && LandscapeBlurHash == "" {
		anime, _ := server.TMDB.GetMovieDetails(TMDbID, nil)
		if anime != nil {
			OriginalTitle = anime.OriginalTitle
			TMDBTitle = anime.Title
			TMDBRuntime = fmt.Sprintf("%dm", anime.Runtime)
			server.getTMDBPic(anime.PosterPath, anime.BackdropPath, &PortriatBlurHash, &PortriatPoster, &LandscapeBlurHash, &LandscapePoster)
			server.getTMDBRating(TMDbID, &AgeRating)
			if len(anime.Genres) > 0 {
				for _, g := range anime.Genres {
					if g.Name != "" {
						Genres = append(Genres, g.Name)
					}
				}
			}
			if len(anime.ProductionCompanies) > 0 {
				for _, p := range anime.ProductionCompanies {
					if p.Name != "" {
						licensors = append(licensors, p.Name)
					}
				}
			}
			amimg, _ := server.TMDB.GetMovieImages(TMDbID, nil)
			if err == nil {
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
			}
			tttt, _ := server.TMDB.GetMovieVideos(TMDbID, nil)
			if len(tttt.Results) > 0 {
				for _, t := range tttt.Results {
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
	} else if TMDbID == 0 && PortriatBlurHash == "" && LandscapeBlurHash == "" {
		querys, _ := server.TMDB.GetSearchMulti(malData.Data.TitleEnglish, nil)
		if querys != nil {
			for _, q := range querys.Results {
				server.LOG.Info().Msgf("query id: %d\n", q.ID)
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
						if OriginalTitle == "" {
							OriginalTitle = q.OriginalTitle
						}
						server.getTMDBPic(q.PosterPath, q.BackdropPath, &PortriatBlurHash, &PortriatPoster, &LandscapeBlurHash, &LandscapePoster)
						server.getTMDBRating(TMDbID, &AgeRating)
						results, _ := server.TMDB.GetGenreMovieList(nil)
						if results != nil {
							for _, f := range results.Genres {
								if len(q.GenreIDs) > 0 {
									for _, h := range q.GenreIDs {
										if int64(f.ID) == h {
											Genres = append(Genres, f.Name)
										}
									}
								}
							}
						}

						anime, _ := server.TMDB.GetMovieDetails(int(q.ID), nil)
						if len(anime.ProductionCompanies) > 0 {
							for _, p := range anime.ProductionCompanies {
								if p.Name != "" {
									licensors = append(licensors, p.Name)
								}
							}
						}
						amimg, _ := server.TMDB.GetMovieImages(TMDbID, nil)
						if err == nil {
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
						}
						tttt, _ := server.TMDB.GetMovieVideos(TMDbID, nil)
						if len(tttt.Results) > 0 {
							for _, t := range tttt.Results {
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
						break
					}
				}
			}
		}
	}

	if TMDbID == 0 && (LandscapePoster == "" || PortriatPoster == "" || PortriatBlurHash == "" || LandscapeBlurHash == "") {
		server.getMalPic(AniDBData.Picture, malData.Data.Images.Jpg.LargeImageUrl, malData.Data.Images.Webp.LargeImageUrl, &PortriatBlurHash, &PortriatPoster)
	}

	if len(malData.Data.Genres) > 0 {
		for _, g := range malData.Data.Genres {
			Genres = append(Genres, g.Name)
		}
	}

	if len(malData.Data.ExplicitGenres) > 0 {
		for _, g := range malData.Data.ExplicitGenres {
			Genres = append(Genres, g.Name)
		}
	}

	if len(malData.Data.Demographics) > 0 {
		for _, g := range malData.Data.Demographics {
			Genres = append(Genres, g.Name)
		}
	}

	if len(AniDBData.Tags.Tag) > 0 {
		for _, t := range AniDBData.Tags.Tag {
			if utils.CleanTag(t.Name) != "" {
				Tags = append(Tags, strings.ToLower(t.Name))
			}
		}
	}

	if malData.Data.Rating == "" {
		AgeRating, err = utils.CleanRating(malData.Data.Rating)
		if err != nil {
			AgeRating = ""
		}
	}

	animePlanetByte, err := animeResources.Data.AnimePlanetID.MarshalJSON()
	if err != nil {
		AnimePlanetID = ""
	} else {
		AnimePlanetID = string(animePlanetByte)
		AnimePlanetID = strings.ReplaceAll(AnimePlanetID, "\"", "")
	}

	Aired = utils.CleanDates([]string{malData.Data.Aired.From.Format(time.DateOnly), AniDBData.Startdate, AniDBData.Enddate})
	if Aired.IsZero() {
		Aired = malData.Data.Aired.From
	}

	if malData.Data.Year != 0 {
		ReleaseYear = malData.Data.Year
	} else {
		ReleaseYear = Aired.Year()
	}

	if Runtime == "" {
		var titles []string
		titles = append(titles, TMDBTitle, malData.Data.Title, malData.Data.TitleEnglish, malData.Data.TitleJapanese)
		if len(AniDBData.Episodes.Episode) > 0 {
			var h int
			for _, e := range AniDBData.Episodes.Episode {
				if e.Epno.Text != "" {
					if strings.Contains(e.Epno.Type, "1") {
						if len(e.Title) > 0 {
							for _, u := range e.Title {
								for _, t := range titles {
									if strings.Contains(utils.CleanTitle(t), utils.CleanTitle(u.Text)) {
										if e.Length != "" {
											b, err := strconv.Atoi(e.Length)
											if err != nil {
												continue
											}
											if h < int(b) && b != 0 {
												h = int(b)
											}
											break
										}
									} else {
										airdate, err := time.Parse(time.DateOnly, e.Airdate)
										if err == nil {
											if malData.Data.Aired.From.Year() == airdate.Year() && malData.Data.Aired.From.Month() == airdate.Month() {
												if e.Length != "" {
													b, err := strconv.Atoi(e.Length)
													if err != nil {
														continue
													}
													if h < int(b) && b != 0 {
														h = int(b)
													}
													break
												}
											}
										}
									}
								}
								if h != 0 {
									break
								}
							}
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
					if err != nil {
						Runtime = ""
					} else {
						Runtime = fmt.Sprintf("%fm", run.Minutes())
					}
				} else {
					Runtime = fmt.Sprintf("%fm", run.Minutes())
				}
			}
		} else {
			run, err := time.ParseDuration(utils.CleanRuntime(malData.Data.Duration))
			if err != nil {
				run, err = time.ParseDuration(utils.CleanRuntime(TMDBRuntime))
				if err != nil {
					Runtime = ""
				} else {
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

	LivechartID := server.Livechart(animeResources.Data.LivechartID, OriginalTitle, Aired)
	AnysearchID := server.Anysearch(animeResources.Data.AnisearchID, malData.Data.TitleEnglish, OriginalTitle, Aired)
	KitsuID := server.Kitsu(animeResources.Data.KitsuID, OriginalTitle, Aired)
	NotifyMoeID := server.NotifyMoe(utils.CleanResText(animeResources.Data.NotifyMoeID), malData.Data.Title, Aired)
	AnilistID := server.Anylist(malData.Data.MalId)

	animeData := models.AnimeMovie{
		OriginalTitle:       OriginalTitle,
		Aired:               Aired.Format(time.DateOnly),
		Runtime:             Runtime,
		ReleaseYear:         ReleaseYear,
		Rating:              AgeRating,
		PortriatPoster:      PortriatPoster,
		PortriatBlurHash:    PortriatBlurHash,
		LandscapePoster:     LandscapePoster,
		LandscapeBlurHash:   LandscapeBlurHash,
		Genres:              utils.CleanStringArray(Genres),
		Studios:             utils.CleanDuplicates(utils.CleanStringArray(Studios)),
		Tags:                utils.CleanStringArray(Tags),
		ProductionCompanies: utils.CleanDuplicates(PsCs),
		Titles:              Titles,
		Backdrops:           Backdrops,
		Posters:             Posters,
		Logos:               Logos,
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
	server.LOG.Info().Msgf("ReleaseYear: %d", ReleaseYear)
	server.LOG.Info().Msgf("PortriatPoster: %s", PortriatPoster)
	server.LOG.Info().Msgf("PortriatBlurHash: %s", PortriatBlurHash)
	server.LOG.Info().Msgf("LandscapePoster: %s", LandscapePoster)
	server.LOG.Info().Msgf("LandscapeBlurHash: %s", LandscapeBlurHash)

	/* if malData.Data.TitleEnglish != "" && malData.Data.Synopsis != "" {
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
					Title:    utils.CleanUnicode(newTitle),
					Overview: utils.CleanUnicode(newOverview),
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
