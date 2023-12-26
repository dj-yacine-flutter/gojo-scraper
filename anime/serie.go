package anime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bregydoc/gtranslate"
	jikan "github.com/darenliang/jikan-go"
	"github.com/dj-yacine-flutter/gojo-scraper/models"
	"github.com/dj-yacine-flutter/gojo-scraper/utils"
)

type SQuery struct {
	Title         string
	OriginalTitle string
	EnglishTitle  string
	MalID         int
	AnidbID       int
	TVDbID        int
	TMDbID        int
	Aired         time.Time

	seasonNumber    int
	papaSerieID     int
	papaSerieTVDbID int
	papaSerieTMDbID int
	papaSerieName   string
	papaSerieAired  time.Time
}

func (server *AnimeScraper) GetAnimeSerie(w http.ResponseWriter, r *http.Request) {
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
		ReleaseYear            int
		PortraitPoster         string
		PortraitBlurHash       string
		SeriePortraitPoster    string
		SeriePortraitBlurHash  string
		SerieLandscapePoster   string
		SerieLandscapeBlurHash string
		AnimePlanetID          string
		OriginalTitle          string
		Aired                  time.Time
		Genres                 []string
		Studios                []string
		Tags                   []string
		PsCs                   []string
		Titles                 models.Titles
		Posters                []models.Image
		Trailers               []models.Trailer
		Licensors              []string
		SeriePosters           []models.Image
		SerieBackdrops         []models.Image
		SerieLogos             []models.Image
		SerieTrailers          []models.Trailer
	)

	var (
		SerieQueries []SQuery
		Query        SQuery
	)

	qmalData, err := jikan.GetAnimeById(int(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("there no data with this id: %s", err.Error()), http.StatusNotFound)
		return
	}

	if !strings.Contains(strings.ToLower(qmalData.Data.Type), "tv") {
		http.Error(w, "this not a anime Serie", http.StatusBadRequest)
		return
	}

	time.Sleep(700 * time.Millisecond)

	malRelation, err := jikan.GetAnimeRelations(int(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get relation with this id: %s", err.Error()), http.StatusNotFound)
		return
	}

	if len(malRelation.Data) > 0 {

		server.LOG.Info().Msgf("MalID: %v", qmalData.Data.MalId)

		SerieQueries = append(SerieQueries, SQuery{
			MalID: qmalData.Data.MalId,
			Title: qmalData.Data.Title,
		})
		for _, e := range malRelation.Data {
			if strings.Contains(strings.ToLower(e.Relation), "sequel") {
				for _, q := range e.Entry {

					server.LOG.Info().Msgf("MalID: %v", q.MalId)

					if strings.Contains(strings.ToLower(q.Type), "anime") {
						SerieQueries = append(SerieQueries, SQuery{
							MalID: q.MalId,
							Title: q.Name,
						})
					}
				}
			} else if strings.Contains(strings.ToLower(e.Relation), "prequel") {
				for _, q := range e.Entry {
					if strings.Contains(strings.ToLower(q.Type), "anime") {

						server.LOG.Info().Msgf("MalID: %v", q.MalId)

						SerieQueries = append(SerieQueries, SQuery{
							MalID: q.MalId,
							Title: q.Name,
						})
					}
				}
			}
		}
	} else {
		server.LOG.Info().Msgf("MalID: %v", qmalData.Data.MalId)

		SerieQueries = append(SerieQueries, SQuery{
			MalID: qmalData.Data.MalId,
			Title: qmalData.Data.Title,
		})
	}

	var getIt bool

	var animeResources = AnimeResources{}
	var AniDBData = AniDB{}
	var MyAnimeListData = &jikan.AnimeById{}

	for _, sq := range SerieQueries {
		time.Sleep(500 * time.Millisecond)

		myanimelistData, err := jikan.GetAnimeById(int(sq.MalID))
		if err != nil {
			continue
		}

		if !strings.Contains(strings.ToLower(myanimelistData.Data.Type), "tv") {
			continue
		}

		server.LOG.Info().Msgf("MalID: %v, Title: %s", sq.MalID, myanimelistData.Data.Title)

		time.Sleep(500 * time.Millisecond)

		malExt, err := jikan.GetAnimeExternal(int(sq.MalID))
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

		sq.AnidbID, err = server.searchAniDBID(myanimelistData, links)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		server.LOG.Info().Msgf("MalID: %v, Title: %s, AniDB: %v", sq.MalID, myanimelistData.Data.Title, sq.AnidbID)

		AniDBData, err = server.GetAniDBData(sq.AnidbID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		server.LOG.Info().Msgf("MalID: %v, Title: %s, AniDB: %v, AniDate: %s", sq.MalID, myanimelistData.Data.Title, sq.AnidbID, AniDBData.Startdate)

		var aired bool
		if AniDBData.Startdate != "" {
			stratDate, err := time.Parse(time.DateOnly, AniDBData.Startdate)
			if err != nil {
				continue
			}
			if qmalData.Data.Aired.From.Year() == stratDate.Year() {
				aired = true
			}
		}

		var rs bool
		animeRes, err := server.getResourceByIDs(sq.AnidbID, sq.MalID)
		if err == nil {
			rs = true
		}

		if aired {
			getIt = true
			MyAnimeListData = myanimelistData
			if rs {
				animeResources = animeRes
			}
			Query = SQuery{
				MalID:         sq.MalID,
				AnidbID:       sq.AnidbID,
				Title:         sq.Title,
				OriginalTitle: qmalData.Data.TitleJapanese,
				EnglishTitle:  qmalData.Data.TitleEnglish,
				Aired:         qmalData.Data.Aired.From,
			}
			break
		}
		AniDBData = AniDB{}
	}

	if !getIt {
		http.Error(w, "No Anime Data", http.StatusNotFound)
		return
	}

	server.getMalPic(AniDBData.Picture, MyAnimeListData.Data.Images.Jpg.LargeImageUrl, MyAnimeListData.Data.Images.Webp.LargeImageUrl, &PortraitBlurHash, &PortraitPoster)

	if MyAnimeListData.Data.TitleEnglish != "" {
		OriginalTitle = MyAnimeListData.Data.TitleEnglish
	} else {
		OriginalTitle = MyAnimeListData.Data.Title
	}

	if len(AniDBData.Creators.Name) > 0 {
		for _, c := range AniDBData.Creators.Name {
			if (strings.Contains(strings.ToLower(c.Type), "work") || (strings.Contains(strings.ToLower(c.Type), "animation") && strings.Contains(strings.ToLower(c.Type), "work"))) && !strings.Contains(strings.ToLower(c.Type), "original") {
				Studios = append(Studios, c.Text)
			}
		}
	}
	if len(MyAnimeListData.Data.Studios) > 0 {
		for _, s := range MyAnimeListData.Data.Studios {
			if s.Name != "" {
				Studios = append(Studios, s.Name)
			}
		}
	}
	if len(MyAnimeListData.Data.Licensors) > 0 {
		for _, p := range MyAnimeListData.Data.Licensors {
			Licensors = append(Licensors, p.Name)
		}
	}
	if len(MyAnimeListData.Data.Producers) > 0 {
		for _, p := range MyAnimeListData.Data.Producers {
			Licensors = append(Licensors, p.Name)
		}
	}
	if MyAnimeListData.Data.Trailer.YoutubeID != "" {
		Trailers = append(Trailers, models.Trailer{
			Official: true,
			Host:     "YouTube",
			Key:      MyAnimeListData.Data.Trailer.YoutubeID,
		})
	}

	time.Sleep(500 * time.Millisecond)
	Query.papaSerieID = server.getMALOriginalID(qmalData.Data.MalId)

	if Query.TVDbID == 0 {
		if Query.papaSerieID != 0 {
			time.Sleep(700 * time.Millisecond)

			mld, err := jikan.GetAnimeById(Query.papaSerieID)
			if err != nil {
				server.LOG.Error().Msgf("TVDB In MAL Error: %v", err.Error())
			}

			if mld != nil {
				Query.papaSerieName = mld.Data.Title
				Query.papaSerieAired = mld.Data.Aired.From

				time.Sleep(500 * time.Millisecond)

				query, err := server.TVDB.GetSearch(mld.Data.Title, 0)
				if err == nil {
					if query != nil {
						if len(query.Data) > 0 {
							for _, d := range query.Data {
								if strings.Contains(strings.ToLower(d.Type), "serie") {

									id, err := strconv.Atoi(d.TvdbID)
									if err != nil {
										server.LOG.Error().Msgf("Papa Serie TVDB ID Loop Error: %v", err.Error())
										continue
									}

									serie, err := server.TVDB.GetSeriesByIDExtanded(id)
									if err != nil {
										server.LOG.Error().Msgf("Papa Serie TVDB Data Loop Error: %v", err.Error())
										continue
									}

									if serie != nil {
										server.LOG.Info().Msgf("Papa Serie TVDB Name: %v", serie.Data.Name)

										if len(serie.Data.Seasons) > 0 {
											for _, s := range serie.Data.Seasons {
												if strings.Contains(strings.ToLower(s.Type.Type), "official") && s.Number != 0 {

													season, err := server.TVDB.GetSeasonsByIDExtended(s.ID)
													if err != nil {
														server.LOG.Error().Msgf("Papa Serie TVDB Season Loop Error: %v", err.Error())
														continue
													}

													year, err := utils.ExtractYear(season.Data.Year)
													if err != nil {
														server.LOG.Error().Msgf("Papa Serie TVDB Year Loop Error: %v", err.Error())

														continue
													}

													if mld.Data.Aired.From.Year() == year {
														Query.papaSerieTVDbID = serie.Data.ID
														if len(serie.Data.Artworks) > 0 {
															for _, d := range serie.Data.Artworks {
																if d.Image != "" {
																	if d.Type == 3 || d.Type == 22 {
																		bb, _ := utils.GetBlurHash(d.Thumbnail, "")
																		SerieBackdrops = append(SerieBackdrops, models.Image{
																			Height:    d.Height,
																			Width:     d.Width,
																			Image:     d.Image,
																			Thumbnail: d.Thumbnail,
																			BlurHash:  bb,
																		})
																	} else if d.Type == 2 || d.Type == 7 {
																		pp, _ := utils.GetBlurHash(d.Thumbnail, "")
																		SeriePosters = append(SeriePosters, models.Image{
																			Height:    d.Height,
																			Width:     d.Width,
																			Image:     d.Image,
																			Thumbnail: d.Thumbnail,
																			BlurHash:  pp,
																		})
																	} else if d.Type == 23 || d.Type == 5 {
																		ll, _ := utils.GetBlurHash(d.Thumbnail, "")
																		SerieLogos = append(SerieLogos, models.Image{
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
																	SerieTrailers = append(SerieTrailers, models.Trailer{
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
										if Query.papaSerieTVDbID != 0 {
											break
										}
									}
								}

							}
						}
					}
				} else {
					server.LOG.Error().Msgf("TVDB Query Loop Error: %v", err.Error())
				}
			}
		}

		if animeResources.Data.TheTVdbID != 0 && Query.papaSerieTVDbID == 0 {
			serie, err := server.TVDB.GetSeriesByIDExtanded(animeResources.Data.TheTVdbID)
			if err == nil && serie != nil {
				if serie.Data.FirstAired != "" {

					server.LOG.Info().Msgf("(1) TVDB AirDate: %v", serie.Data.FirstAired)

					aired, err := time.Parse(time.DateOnly, serie.Data.FirstAired)
					if err == nil {
						if qmalData.Data.Aired.From.Year() == aired.Year() && qmalData.Data.Aired.From.Month() == aired.Month() {
							if len(serie.Data.Seasons) > 0 {
								for _, s := range serie.Data.Seasons {
									if strings.Contains(strings.ToLower(s.Type.Type), "official") && s.Number != 0 {

										server.LOG.Info().Msgf("(1) TVDB Season Name: %v", s.Name)
										server.LOG.Info().Msgf("(1) TVDB Season Year: %v", s.Year)

										season, err := server.TVDB.GetSeasonsByIDExtended(s.ID)
										if err != nil {
											continue
										}

										year, err := utils.ExtractYear(season.Data.Year)
										if err != nil {
											continue
										}

										if MyAnimeListData.Data.Aired.From.Year() == year {

											Query.papaSerieTVDbID = season.Data.SeriesID
											if len(serie.Data.Artworks) > 0 {
												for _, d := range serie.Data.Artworks {
													if d.Image != "" {
														if d.Type == 3 || d.Type == 22 {
															bb, _ := utils.GetBlurHash(d.Thumbnail, "")
															SerieBackdrops = append(SerieBackdrops, models.Image{
																Height:    d.Height,
																Width:     d.Width,
																Image:     d.Image,
																Thumbnail: d.Thumbnail,
																BlurHash:  bb,
															})
														} else if d.Type == 2 || d.Type == 7 {
															pp, _ := utils.GetBlurHash(d.Thumbnail, "")
															SeriePosters = append(SeriePosters, models.Image{
																Height:    d.Height,
																Width:     d.Width,
																Image:     d.Image,
																Thumbnail: d.Thumbnail,
																BlurHash:  pp,
															})
														} else if d.Type == 23 || d.Type == 5 {
															ll, _ := utils.GetBlurHash(d.Thumbnail, "")
															SerieLogos = append(SerieLogos, models.Image{
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
														SerieTrailers = append(SerieTrailers, models.Trailer{
															Official: true,
															Host:     "YouTube",
															Key:      utils.ExtractYTKey(t.URL),
														})
													}
												}
											}
											Query.TVDbID = season.Data.ID
											Query.seasonNumber = season.Data.Number

											if len(season.Data.Artwork) > 0 {
												for _, d := range season.Data.Artwork {
													if d.Image != "" {
														if d.Type == 7 {
															pp, _ := utils.GetBlurHash(d.Thumbnail, "")
															Posters = append(Posters, models.Image{
																Height:    d.Height,
																Width:     d.Width,
																Image:     d.Image,
																Thumbnail: d.Thumbnail,
																BlurHash:  pp,
															})
														}
													}
												}
											}
											if len(season.Data.Trailers) > 0 {
												for _, t := range season.Data.Trailers {
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
										//gg, err := json.Marshal(&season)
										//if err != nil {
										//	continue
										//}
										//
										//server.LOG.Info().Msgf("TVDB Json: %s ", string(gg))
									}
								}
							}

						}
					}
				}
			}
		}

		if Query.papaSerieTVDbID != 0 {
			server.LOG.Info().Msgf("Papa TVDB Query Tv: ID: --%d--", Query.papaSerieTVDbID)

			serie, err := server.TVDB.GetSeriesByIDExtanded(Query.papaSerieTVDbID)
			if err == nil {
				if serie != nil {
					if len(serie.Data.Seasons) > 0 {
						for _, s := range serie.Data.Seasons {
							if strings.Contains(strings.ToLower(s.Type.Type), "official") && s.Number != 0 {

								season, err := server.TVDB.GetSeasonsByIDExtended(s.ID)
								if err != nil {
									server.LOG.Error().Msgf("Papa TVDB Season Loop Error: %v", err.Error())
									continue
								}

								if season != nil {
									if len(season.Data.Episodes) > 0 {

										year, err := time.Parse(time.DateOnly, season.Data.Episodes[0].Aired)
										if err != nil {
											server.LOG.Error().Msgf("Papa TVDB EP Year Loop Error: %v", err.Error())
											continue
										}

										if MyAnimeListData.Data.Aired.From.Year() == year.Year() && MyAnimeListData.Data.Aired.From.Month() == year.Month() {

											if Query.papaSerieID == 0 {
												aired, _ := time.Parse(time.DateOnly, serie.Data.FirstAired)
												Query.papaSerieAired = aired
												Query.papaSerieName = serie.Data.Name
											}

											Query.TVDbID = season.Data.ID
											Query.seasonNumber = season.Data.Number

											if len(season.Data.Artwork) > 0 {
												for _, d := range season.Data.Artwork {
													if d.Image != "" {
														if d.Type == 7 {
															pp, _ := utils.GetBlurHash(d.Thumbnail, "")
															Posters = append(Posters, models.Image{
																Height:    d.Height,
																Width:     d.Width,
																Image:     d.Image,
																Thumbnail: d.Thumbnail,
																BlurHash:  pp,
															})
														}
													}
												}
											}
											if len(season.Data.Trailers) > 0 {
												for _, t := range season.Data.Trailers {
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
										} else {
											dd, err := time.Parse(time.DateOnly, season.Data.Episodes[len(season.Data.Episodes)-1].Aired)
											if err != nil {
												server.LOG.Error().Msgf("Papa TVDB EP -1 Year Loop Error: %v", err.Error())
												continue
											}

											if MyAnimeListData.Data.Aired.From.Year() == dd.Year() {

												for _, d := range season.Data.Episodes {
													dd2, err := time.Parse(time.DateOnly, d.Aired)
													if err != nil {
														server.LOG.Error().Msgf("Papa TVDB EP -1 (2) Year Loop Error: %v", err.Error())
														continue
													}

													if MyAnimeListData.Data.Aired.From.Year() == dd2.Year() && MyAnimeListData.Data.Aired.From.Month() == dd2.Month() {
														if Query.papaSerieID == 0 {
															aired, _ := time.Parse(time.DateOnly, serie.Data.FirstAired)
															Query.papaSerieAired = aired
															Query.papaSerieName = serie.Data.Name
														}

														Query.TVDbID = season.Data.ID
														Query.seasonNumber = season.Data.Number

														if len(season.Data.Artwork) > 0 {
															for _, d := range season.Data.Artwork {
																if d.Image != "" {
																	if d.Type == 7 {
																		pp, _ := utils.GetBlurHash(d.Thumbnail, "")
																		Posters = append(Posters, models.Image{
																			Height:    d.Height,
																			Width:     d.Width,
																			Image:     d.Image,
																			Thumbnail: d.Thumbnail,
																			BlurHash:  pp,
																		})
																	}
																}
															}
														}
														if len(season.Data.Trailers) > 0 {
															for _, t := range season.Data.Trailers {
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
													if Query.TVDbID != 0 {
														break
													}
												}
											}
										}

									} else {
										year, err := utils.ExtractYear(season.Data.Year)
										if err != nil {
											server.LOG.Error().Msgf("Papa TVDB Year Loop Error: %v", err.Error())

											continue
										}

										if MyAnimeListData.Data.Aired.From.Year() == year {

											if Query.papaSerieID == 0 {
												aired, _ := time.Parse(time.DateOnly, serie.Data.FirstAired)
												Query.papaSerieAired = aired
												Query.papaSerieName = serie.Data.Name
											}

											Query.TVDbID = season.Data.ID
											Query.seasonNumber = season.Data.Number

											if len(season.Data.Artwork) > 0 {
												for _, d := range season.Data.Artwork {
													if d.Image != "" {
														if d.Type == 7 {
															pp, _ := utils.GetBlurHash(d.Thumbnail, "")
															Posters = append(Posters, models.Image{
																Height:    d.Height,
																Width:     d.Width,
																Image:     d.Image,
																Thumbnail: d.Thumbnail,
																BlurHash:  pp,
															})
														}
													}
												}
											}
											if len(season.Data.Trailers) > 0 {
												for _, t := range season.Data.Trailers {
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
						}
					}
				}
			}
		}

		if Query.TVDbID == 0 {
			var AlternativeQuries []SQuery
			for _, sq := range SerieQueries {

				server.LOG.Info().Msgf("In TVDB (2) Title: %s", sq.Title)

				query, err := server.TVDB.GetSearch(sq.Title, 0)
				if err != nil {
					server.LOG.Error().Msgf("TVDB Query Loop Error: %v", err.Error())
					continue
				}

				if query != nil {
					if len(query.Data) > 0 {
						for _, d := range query.Data {
							if strings.Contains(strings.ToLower(d.Type), "serie") {

								server.LOG.Info().Msgf("(2) TVDB Query Tv: --%s-- with ID: --%s--", d.Name, d.TvdbID)

								id, err := strconv.Atoi(d.TvdbID)
								if err != nil {
									server.LOG.Error().Msgf("(2) TVDB ID Loop Error: %v", err.Error())
									continue
								}

								serie, err := server.TVDB.GetSeriesByIDExtanded(id)
								if err != nil {
									server.LOG.Error().Msgf("(2) TVDB Serie Loop Error: %v", err.Error())
									continue
								}

								if serie != nil {
									if len(serie.Data.Seasons) > 0 {
										for _, s := range serie.Data.Seasons {
											if strings.Contains(strings.ToLower(s.Type.Type), "official") && s.Number != 0 {

												season, err := server.TVDB.GetSeasonsByIDExtended(s.ID)
												if err != nil {
													server.LOG.Error().Msgf("(2) TVDB Season Loop Error: %v", err.Error())
													continue
												}

												if season != nil {
													if len(season.Data.Episodes) > 0 {

														year, err := time.Parse(time.DateOnly, season.Data.Episodes[0].Aired)
														if err != nil {
															server.LOG.Error().Msgf("(2) Papa TVDB EP Year Loop Error: %v", err.Error())
															continue
														}

														if MyAnimeListData.Data.Aired.From.Year() == year.Year() && MyAnimeListData.Data.Aired.From.Month() == year.Month() {

															if Query.papaSerieID == 0 {
																aired, _ := time.Parse(time.DateOnly, serie.Data.FirstAired)
																Query.papaSerieAired = aired
																Query.papaSerieID = server.getMALOriginalID(sq.MalID)
																Query.papaSerieName = serie.Data.Name
															}

															if Query.papaSerieTVDbID == 0 {
																Query.papaSerieTVDbID = season.Data.SeriesID
															}
															if len(serie.Data.Artworks) > 0 {
																for _, d := range serie.Data.Artworks {
																	if d.Image != "" {
																		if d.Type == 3 || d.Type == 22 {
																			bb, _ := utils.GetBlurHash(d.Thumbnail, "")
																			SerieBackdrops = append(SerieBackdrops, models.Image{
																				Height:    d.Height,
																				Width:     d.Width,
																				Image:     d.Image,
																				Thumbnail: d.Thumbnail,
																				BlurHash:  bb,
																			})
																		} else if d.Type == 2 || d.Type == 7 {
																			pp, _ := utils.GetBlurHash(d.Thumbnail, "")
																			SeriePosters = append(SeriePosters, models.Image{
																				Height:    d.Height,
																				Width:     d.Width,
																				Image:     d.Image,
																				Thumbnail: d.Thumbnail,
																				BlurHash:  pp,
																			})
																		} else if d.Type == 23 || d.Type == 5 {
																			ll, _ := utils.GetBlurHash(d.Thumbnail, "")
																			SerieLogos = append(SerieLogos, models.Image{
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
																		SerieTrailers = append(SerieTrailers, models.Trailer{
																			Official: true,
																			Host:     "YouTube",
																			Key:      utils.ExtractYTKey(t.URL),
																		})
																	}
																}
															}
															Query.TVDbID = season.Data.ID
															Query.seasonNumber = season.Data.Number

															if len(season.Data.Artwork) > 0 {
																for _, d := range season.Data.Artwork {
																	if d.Image != "" {
																		if d.Type == 7 {
																			pp, _ := utils.GetBlurHash(d.Thumbnail, "")
																			Posters = append(Posters, models.Image{
																				Height:    d.Height,
																				Width:     d.Width,
																				Image:     d.Image,
																				Thumbnail: d.Thumbnail,
																				BlurHash:  pp,
																			})
																		}
																	}
																}
															}
															if len(season.Data.Trailers) > 0 {
																for _, t := range season.Data.Trailers {
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
														} else {
															dd, err := time.Parse(time.DateOnly, season.Data.Episodes[len(season.Data.Episodes)-1].Aired)
															if err != nil {
																server.LOG.Error().Msgf("(2) Papa TVDB EP -1 Year Loop Error: %v", err.Error())
																continue
															}

															if MyAnimeListData.Data.Aired.From.Year() == dd.Year() {

																for _, d := range season.Data.Episodes {
																	dd2, err := time.Parse(time.DateOnly, d.Aired)
																	if err != nil {
																		server.LOG.Error().Msgf("(2) Papa TVDB EP -1 (2) Year Loop Error: %v", err.Error())
																		continue
																	}

																	if MyAnimeListData.Data.Aired.From.Year() == dd2.Year() && MyAnimeListData.Data.Aired.From.Month() == dd2.Month() {
																		if Query.papaSerieID == 0 {
																			aired, _ := time.Parse(time.DateOnly, serie.Data.FirstAired)
																			Query.papaSerieAired = aired
																			Query.papaSerieName = serie.Data.Name
																		}

																		Query.TVDbID = season.Data.ID
																		Query.seasonNumber = season.Data.Number

																		if len(serie.Data.Artworks) > 0 {
																			for _, d := range serie.Data.Artworks {
																				if d.Image != "" {
																					if d.Type == 3 || d.Type == 22 {
																						bb, _ := utils.GetBlurHash(d.Thumbnail, "")
																						SerieBackdrops = append(SerieBackdrops, models.Image{
																							Height:    d.Height,
																							Width:     d.Width,
																							Image:     d.Image,
																							Thumbnail: d.Thumbnail,
																							BlurHash:  bb,
																						})
																					} else if d.Type == 2 || d.Type == 7 {
																						pp, _ := utils.GetBlurHash(d.Thumbnail, "")
																						SeriePosters = append(SeriePosters, models.Image{
																							Height:    d.Height,
																							Width:     d.Width,
																							Image:     d.Image,
																							Thumbnail: d.Thumbnail,
																							BlurHash:  pp,
																						})
																					} else if d.Type == 23 || d.Type == 5 {
																						ll, _ := utils.GetBlurHash(d.Thumbnail, "")
																						SerieLogos = append(SerieLogos, models.Image{
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
																					SerieTrailers = append(SerieTrailers, models.Trailer{
																						Official: true,
																						Host:     "YouTube",
																						Key:      utils.ExtractYTKey(t.URL),
																					})
																				}
																			}
																		}
																		if len(season.Data.Artwork) > 0 {
																			for _, d := range season.Data.Artwork {
																				if d.Image != "" {
																					if d.Type == 7 {
																						pp, _ := utils.GetBlurHash(d.Thumbnail, "")
																						Posters = append(Posters, models.Image{
																							Height:    d.Height,
																							Width:     d.Width,
																							Image:     d.Image,
																							Thumbnail: d.Thumbnail,
																							BlurHash:  pp,
																						})
																					}
																				}
																			}
																		}
																		if len(season.Data.Trailers) > 0 {
																			for _, t := range season.Data.Trailers {
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
																	if Query.TVDbID != 0 {
																		break
																	}
																}
															}
														}

													} else {
														year, err := utils.ExtractYear(season.Data.Year)
														if err != nil {
															server.LOG.Error().Msgf("Papa TVDB Year Loop Error: %v", err.Error())

															continue
														}

														if MyAnimeListData.Data.Aired.From.Year() == year {

															if Query.papaSerieID == 0 {
																aired, _ := time.Parse(time.DateOnly, serie.Data.FirstAired)
																Query.papaSerieAired = aired
																Query.papaSerieID = server.getMALOriginalID(sq.MalID)
																Query.papaSerieName = serie.Data.Name
															}

															if Query.papaSerieTVDbID == 0 {
																Query.papaSerieTVDbID = season.Data.SeriesID
															}
															Query.TVDbID = season.Data.ID
															Query.seasonNumber = season.Data.Number
															if len(serie.Data.Artworks) > 0 {
																for _, d := range serie.Data.Artworks {
																	if d.Image != "" {
																		if d.Type == 3 || d.Type == 22 {
																			bb, _ := utils.GetBlurHash(d.Thumbnail, "")
																			SerieBackdrops = append(SerieBackdrops, models.Image{
																				Height:    d.Height,
																				Width:     d.Width,
																				Image:     d.Image,
																				Thumbnail: d.Thumbnail,
																				BlurHash:  bb,
																			})
																		} else if d.Type == 2 || d.Type == 7 {
																			pp, _ := utils.GetBlurHash(d.Thumbnail, "")
																			SeriePosters = append(SeriePosters, models.Image{
																				Height:    d.Height,
																				Width:     d.Width,
																				Image:     d.Image,
																				Thumbnail: d.Thumbnail,
																				BlurHash:  pp,
																			})
																		} else if d.Type == 23 || d.Type == 5 {
																			ll, _ := utils.GetBlurHash(d.Thumbnail, "")
																			SerieLogos = append(SerieLogos, models.Image{
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
																		SerieTrailers = append(SerieTrailers, models.Trailer{
																			Official: true,
																			Host:     "YouTube",
																			Key:      utils.ExtractYTKey(t.URL),
																		})
																	}
																}
															}
															if len(season.Data.Artwork) > 0 {
																for _, d := range season.Data.Artwork {
																	if d.Image != "" {
																		if d.Type == 7 {
																			pp, _ := utils.GetBlurHash(d.Thumbnail, "")
																			Posters = append(Posters, models.Image{
																				Height:    d.Height,
																				Width:     d.Width,
																				Image:     d.Image,
																				Thumbnail: d.Thumbnail,
																				BlurHash:  pp,
																			})
																		}
																	}
																}
															}
															if len(season.Data.Trailers) > 0 {
																for _, t := range season.Data.Trailers {
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
										}
									}
								}

							}
						}
					}
				}
				if Query.TVDbID != 0 {
					break
				}

				for _, t := range strings.Split(sq.Title, ":") {
					AlternativeQuries = append(AlternativeQuries, SQuery{
						Title: t,
						MalID: sq.MalID,
					})
				}
			}
			if Query.TVDbID == 0 {
				for _, sq := range AlternativeQuries {

					server.LOG.Info().Msgf("(3)  TVDB Title: %s", sq.Title)

					query, err := server.TVDB.GetSearch(sq.Title, 0)
					if err != nil {
						server.LOG.Error().Msgf("(3) TVDB Query Loop Error: %v", err.Error())
						continue
					}

					if query != nil {
						if len(query.Data) > 0 {
							for _, d := range query.Data {
								if strings.Contains(strings.ToLower(d.Type), "serie") {

									server.LOG.Info().Msgf("(3) TVDB Query Tv: --%s-- with ID: --%s--", d.Name, d.TvdbID)

									id, err := strconv.Atoi(d.TvdbID)
									if err != nil {
										server.LOG.Error().Msgf("(3) TVDB ID Loop Error: %v", err.Error())
										continue
									}

									serie, err := server.TVDB.GetSeriesByIDExtanded(id)
									if err != nil {
										server.LOG.Error().Msgf("(3) TVDB Serie Loop Error: %v", err.Error())
										continue
									}
									if serie != nil {
										if len(serie.Data.Seasons) > 0 {
											for _, s := range serie.Data.Seasons {
												if strings.Contains(strings.ToLower(s.Type.Type), "official") && s.Number != 0 {

													season, err := server.TVDB.GetSeasonsByIDExtended(s.ID)
													if err != nil {
														server.LOG.Error().Msgf("(3) TVDB Season Loop Error: %v", err.Error())
														continue
													}

													if season != nil {
														if len(season.Data.Episodes) > 0 {

															year, err := time.Parse(time.DateOnly, season.Data.Episodes[0].Aired)
															if err != nil {
																server.LOG.Error().Msgf("(3) Papa TVDB EP Year Loop Error: %v", err.Error())

																continue
															}

															if MyAnimeListData.Data.Aired.From.Year() == year.Year() && MyAnimeListData.Data.Aired.From.Month() == year.Month() {

																if Query.papaSerieID == 0 {
																	aired, _ := time.Parse(time.DateOnly, serie.Data.FirstAired)
																	Query.papaSerieAired = aired
																	Query.papaSerieID = server.getMALOriginalID(sq.MalID)
																	Query.papaSerieName = serie.Data.Name
																}

																if Query.papaSerieTVDbID == 0 {
																	Query.papaSerieTVDbID = season.Data.SeriesID
																}
																Query.TVDbID = season.Data.ID
																Query.seasonNumber = season.Data.Number
																if len(serie.Data.Artworks) > 0 {
																	for _, d := range serie.Data.Artworks {
																		if d.Image != "" {
																			if d.Type == 3 || d.Type == 22 {
																				bb, _ := utils.GetBlurHash(d.Thumbnail, "")
																				SerieBackdrops = append(SerieBackdrops, models.Image{
																					Height:    d.Height,
																					Width:     d.Width,
																					Image:     d.Image,
																					Thumbnail: d.Thumbnail,
																					BlurHash:  bb,
																				})
																			} else if d.Type == 2 || d.Type == 7 {
																				pp, _ := utils.GetBlurHash(d.Thumbnail, "")
																				SeriePosters = append(SeriePosters, models.Image{
																					Height:    d.Height,
																					Width:     d.Width,
																					Image:     d.Image,
																					Thumbnail: d.Thumbnail,
																					BlurHash:  pp,
																				})
																			} else if d.Type == 23 || d.Type == 5 {
																				ll, _ := utils.GetBlurHash(d.Thumbnail, "")
																				SerieLogos = append(SerieLogos, models.Image{
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
																			SerieTrailers = append(SerieTrailers, models.Trailer{
																				Official: true,
																				Host:     "YouTube",
																				Key:      utils.ExtractYTKey(t.URL),
																			})
																		}
																	}
																}
																if len(season.Data.Artwork) > 0 {
																	for _, d := range season.Data.Artwork {
																		if d.Image != "" {
																			if d.Type == 7 {
																				pp, _ := utils.GetBlurHash(d.Thumbnail, "")
																				Posters = append(Posters, models.Image{
																					Height:    d.Height,
																					Width:     d.Width,
																					Image:     d.Image,
																					Thumbnail: d.Thumbnail,
																					BlurHash:  pp,
																				})
																			}
																		}
																	}
																}
																if len(season.Data.Trailers) > 0 {
																	for _, t := range season.Data.Trailers {
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
															} else {
																dd, err := time.Parse(time.DateOnly, season.Data.Episodes[len(season.Data.Episodes)-1].Aired)
																if err != nil {
																	server.LOG.Error().Msgf("(3) Papa TVDB EP -1 Year Loop Error: %v", err.Error())
																	continue
																}

																if MyAnimeListData.Data.Aired.From.Year() == dd.Year() {

																	for _, d := range season.Data.Episodes {
																		dd2, err := time.Parse(time.DateOnly, d.Aired)
																		if err != nil {
																			server.LOG.Error().Msgf("(3) Papa TVDB EP -1 (2) Year Loop Error: %v", err.Error())
																			continue
																		}

																		if MyAnimeListData.Data.Aired.From.Year() == dd2.Year() && MyAnimeListData.Data.Aired.From.Month() == dd2.Month() {
																			if Query.papaSerieID == 0 {
																				aired, _ := time.Parse(time.DateOnly, serie.Data.FirstAired)
																				Query.papaSerieAired = aired
																				Query.papaSerieName = serie.Data.Name
																			}

																			Query.TVDbID = season.Data.ID
																			Query.seasonNumber = season.Data.Number
																			if len(serie.Data.Artworks) > 0 {
																				for _, d := range serie.Data.Artworks {
																					if d.Image != "" {
																						if d.Type == 3 || d.Type == 22 {
																							bb, _ := utils.GetBlurHash(d.Thumbnail, "")
																							SerieBackdrops = append(SerieBackdrops, models.Image{
																								Height:    d.Height,
																								Width:     d.Width,
																								Image:     d.Image,
																								Thumbnail: d.Thumbnail,
																								BlurHash:  bb,
																							})
																						} else if d.Type == 2 || d.Type == 7 {
																							pp, _ := utils.GetBlurHash(d.Thumbnail, "")
																							SeriePosters = append(SeriePosters, models.Image{
																								Height:    d.Height,
																								Width:     d.Width,
																								Image:     d.Image,
																								Thumbnail: d.Thumbnail,
																								BlurHash:  pp,
																							})
																						} else if d.Type == 23 || d.Type == 5 {
																							ll, _ := utils.GetBlurHash(d.Thumbnail, "")
																							SerieLogos = append(SerieLogos, models.Image{
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
																						SerieTrailers = append(SerieTrailers, models.Trailer{
																							Official: true,
																							Host:     "YouTube",
																							Key:      utils.ExtractYTKey(t.URL),
																						})
																					}
																				}
																			}
																			if len(season.Data.Artwork) > 0 {
																				for _, d := range season.Data.Artwork {
																					if d.Image != "" {
																						if d.Type == 7 {
																							pp, _ := utils.GetBlurHash(d.Thumbnail, "")
																							Posters = append(Posters, models.Image{
																								Height:    d.Height,
																								Width:     d.Width,
																								Image:     d.Image,
																								Thumbnail: d.Thumbnail,
																								BlurHash:  pp,
																							})
																						}
																					}
																				}
																			}
																			if len(season.Data.Trailers) > 0 {
																				for _, t := range season.Data.Trailers {
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
																		if Query.seasonNumber != 0 {
																			break
																		}
																	}
																}
															}
														} else {
															year, err := utils.ExtractYear(season.Data.Year)
															if err != nil {
																server.LOG.Error().Msgf("(3) Papa TVDB Year Loop Error: %v", err.Error())

																continue
															}

															if MyAnimeListData.Data.Aired.From.Year() == year {

																if Query.papaSerieID == 0 {
																	aired, _ := time.Parse(time.DateOnly, serie.Data.FirstAired)
																	Query.papaSerieAired = aired
																	Query.papaSerieID = server.getMALOriginalID(sq.MalID)
																	Query.papaSerieName = serie.Data.Name
																}

																if Query.papaSerieTVDbID == 0 {
																	Query.papaSerieTVDbID = season.Data.SeriesID
																}
																Query.TVDbID = season.Data.ID
																Query.seasonNumber = season.Data.Number
																if len(serie.Data.Artworks) > 0 {
																	for _, d := range serie.Data.Artworks {
																		if d.Image != "" {
																			if d.Type == 3 || d.Type == 22 {
																				bb, _ := utils.GetBlurHash(d.Thumbnail, "")
																				SerieBackdrops = append(SerieBackdrops, models.Image{
																					Height:    d.Height,
																					Width:     d.Width,
																					Image:     d.Image,
																					Thumbnail: d.Thumbnail,
																					BlurHash:  bb,
																				})
																			} else if d.Type == 2 || d.Type == 7 {
																				pp, _ := utils.GetBlurHash(d.Thumbnail, "")
																				SeriePosters = append(SeriePosters, models.Image{
																					Height:    d.Height,
																					Width:     d.Width,
																					Image:     d.Image,
																					Thumbnail: d.Thumbnail,
																					BlurHash:  pp,
																				})
																			} else if d.Type == 23 || d.Type == 5 {
																				ll, _ := utils.GetBlurHash(d.Thumbnail, "")
																				SerieLogos = append(SerieLogos, models.Image{
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
																			SerieTrailers = append(SerieTrailers, models.Trailer{
																				Official: true,
																				Host:     "YouTube",
																				Key:      utils.ExtractYTKey(t.URL),
																			})
																		}
																	}
																}
																if len(season.Data.Artwork) > 0 {
																	for _, d := range season.Data.Artwork {
																		if d.Image != "" {
																			if d.Type == 7 {
																				pp, _ := utils.GetBlurHash(d.Thumbnail, "")
																				Posters = append(Posters, models.Image{
																					Height:    d.Height,
																					Width:     d.Width,
																					Image:     d.Image,
																					Thumbnail: d.Thumbnail,
																					BlurHash:  pp,
																				})
																			}
																		}
																	}
																}
																if len(season.Data.Trailers) > 0 {
																	for _, t := range season.Data.Trailers {
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
											}
										}
									}
								}
							}

						}
					}

					if Query.TVDbID != 0 {
						break
					}
				}
			}
		}
	}

	var tmdbTime time.Time
	if Query.TMDbID == 0 {
		if Query.papaSerieTVDbID != 0 {
			serie, err := server.TMDB.GetFindByID(fmt.Sprint(Query.papaSerieTVDbID), map[string]string{"external_source": "tvdb_id"})
			if err != nil {
				server.LOG.Error().Msgf("error when find serie in tmdb: %v", err.Error())
				http.Error(w, "error when find serie in tmdb", http.StatusInternalServerError)
			}

			if serie != nil {
				if len(serie.TvResults) > 0 {
					for _, s := range serie.TvResults {
						server.LOG.Info().Msgf("Papa Serie TMDB Name: %v", s.Name)
						server.LOG.Info().Msgf("Papa Serie TMDB First AirDate: %v", s.FirstAirDate)

						fair, err := time.Parse(time.DateOnly, s.FirstAirDate)
						if err != nil {
							server.LOG.Error().Msgf("error when get airDate serie in tmdb: %v", err.Error())
							http.Error(w, "error when get airDate serie in tmdb", http.StatusInternalServerError)
						}

						if Query.papaSerieAired.Year() == fair.Year() {
							Query.papaSerieTMDbID = int(s.ID)
							server.getTMDBPic(s.PosterPath, s.BackdropPath, &SeriePortraitBlurHash, &SeriePortraitPoster, &SerieLandscapeBlurHash, &SerieLandscapePoster)

							samimg, _ := server.TMDB.GetTVImages(int(s.ID), nil)
							if err == nil {
								if samimg != nil {
									for _, b := range samimg.Backdrops {
										if b.FilePath != "" {
											bb, _ := utils.GetBlurHash(server.DecodeIMG+b.FilePath, "")
											SerieBackdrops = append(SerieBackdrops, models.Image{
												Height:    b.Height,
												Width:     b.Width,
												Image:     strings.TrimSpace(fmt.Sprintf(server.OriginalIMG + b.FilePath)),
												Thumbnail: strings.TrimSpace("https://image.tmdb.org/t/p/w500" + b.FilePath),
												BlurHash:  bb,
											})
										}
									}
									for _, p := range samimg.Posters {
										if p.FilePath != "" {
											pp, _ := utils.GetBlurHash(server.DecodeIMG+p.FilePath, "")
											SeriePosters = append(SeriePosters, models.Image{
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

							tttt, _ := server.TMDB.GetTVVideos(int(s.ID), nil)
							if len(tttt.Results) > 0 {
								for _, t := range tttt.Results {
									if strings.Contains(strings.ToLower(t.Site), "youtube") {
										if t.Key != "" {
											SerieTrailers = append(SerieTrailers, models.Trailer{
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

			if Query.papaSerieTMDbID != 0 {
				data, err := server.TMDB.GetTVDetails(Query.papaSerieTMDbID, nil)
				if err != nil {
					server.LOG.Error().Msgf("error when get serie seasons in tmdb: %v", err.Error())
					http.Error(w, "error when get serie seasons in tmdb", http.StatusInternalServerError)
				}

				if data != nil {
					for _, s := range data.Seasons {
						if s.SeasonNumber != 0 {
							server.LOG.Info().Msgf("Papa Serie Season TMDB Name: %v", s.Name)
							server.LOG.Info().Msgf("Papa Serie Season TMDB First AirDate: %v", s.AirDate)

							air, err := time.Parse(time.DateOnly, s.AirDate)
							if err != nil {
								server.LOG.Error().Msgf("error when get airDate serie season in tmdb: %v", err.Error())
								http.Error(w, "error when get airDate serie season in tmdb", http.StatusInternalServerError)
							}

							if Query.Aired.Year() == air.Year() {
								tmdbTime = air
								Query.TMDbID = int(s.ID)

								amimg, _ := server.TMDB.GetTVSeasonImages(Query.papaSerieTMDbID, s.SeasonNumber, nil)
								if amimg != nil {
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

								server.getTMDBPic(s.PosterPath, "", &PortraitBlurHash, &PortraitPoster, nil, nil)

								tttt, _ := server.TMDB.GetTVSeasonVideos(Query.papaSerieTMDbID, s.SeasonNumber, nil)
								if tttt != nil {
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

								break
							}

						}
					}
				}
			}
		}

		if Query.TVDbID != 0 && Query.papaSerieTMDbID == 0 {
			season, err := server.TMDB.GetFindByID(fmt.Sprint(Query.TVDbID), map[string]string{"external_source": "tvdb_id"})
			if err != nil {
				server.LOG.Error().Msgf("error when find serie in tmdb: %v", err.Error())
				http.Error(w, "error when find serie in tmdb", http.StatusInternalServerError)
			}

			if season != nil {
				if len(season.TvSeasonResults) > 0 {
					for _, s := range season.TvSeasonResults {
						if s.SeasonNumber != 0 {
							server.LOG.Info().Msgf("Papa Season TMDB Name: %v", s.Name)
							server.LOG.Info().Msgf("Papa Season TMDB First AirDate: %v", s.AirDate)

							air, err := time.Parse(time.DateOnly, s.AirDate)
							if err != nil {
								server.LOG.Error().Msgf("error when get airDate season in tmdb: %v", err.Error())
								http.Error(w, "error when get airDate season in tmdb", http.StatusInternalServerError)
							}

							if Query.Aired.Year() == air.Year() {
								tmdbTime = air
								Query.TMDbID = int(s.ID)
								Query.papaSerieTMDbID = int(s.ShowID)

								samimg, _ := server.TMDB.GetTVImages(int(s.ShowID), nil)
								if err == nil {
									if samimg != nil {
										for _, b := range samimg.Backdrops {
											if b.FilePath != "" {
												bb, _ := utils.GetBlurHash(server.DecodeIMG+b.FilePath, "")
												SerieBackdrops = append(SerieBackdrops, models.Image{
													Height:    b.Height,
													Width:     b.Width,
													Image:     strings.TrimSpace(fmt.Sprintf(server.OriginalIMG + b.FilePath)),
													Thumbnail: strings.TrimSpace("https://image.tmdb.org/t/p/w500" + b.FilePath),
													BlurHash:  bb,
												})
											}
										}
										for _, p := range samimg.Posters {
											if p.FilePath != "" {
												pp, _ := utils.GetBlurHash(server.DecodeIMG+p.FilePath, "")
												SeriePosters = append(SeriePosters, models.Image{
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
								stttt, _ := server.TMDB.GetTVVideos(int(s.ShowID), nil)
								if len(stttt.Results) > 0 {
									for _, t := range stttt.Results {
										if strings.Contains(strings.ToLower(t.Site), "youtube") {
											if t.Key != "" {
												SerieTrailers = append(SerieTrailers, models.Trailer{
													Official: true,
													Host:     "YouTube",
													Key:      t.Key,
												})
											}

										}
									}
								}
								//server.getTMDBPic(s.PosterPath, s.BackdropPath, &SeriePortraitBlurHash, &SeriePortraitPoster, &SerieLandscapeBlurHash, &SerieLandscapePoster)

								amimg, _ := server.TMDB.GetTVSeasonImages(int(s.ShowID), s.SeasonNumber, nil)
								if amimg != nil {
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

								anime, _ := server.TMDB.GetTVSeasonDetails(int(s.ShowID), s.SeasonNumber, nil)
								if anime != nil {
									server.getTMDBPic(anime.PosterPath, "", &PortraitBlurHash, &PortraitPoster, nil, nil)
								}

								tttt, _ := server.TMDB.GetTVSeasonVideos(int(s.ShowID), s.SeasonNumber, nil)
								if tttt != nil {
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

								break
							}

						}
					}
				}
			}
		}

		if Query.TMDbID == 0 {
			querys, _ := server.TMDB.GetSearchTVShow(MyAnimeListData.Data.TitleEnglish, map[string]string{"first_air_date_year": fmt.Sprint(MyAnimeListData.Data.Aired.From.Year())})
			if querys != nil {
				for _, q := range querys.Results {

					fair, err := time.Parse(time.DateOnly, q.FirstAirDate)
					if err != nil {
						server.LOG.Error().Msgf("error in airdate when search serie in tmdb: %v", err.Error())
						continue
					}

					if fair.Year() == MyAnimeListData.Data.Aired.From.Year() && fair.Month() == MyAnimeListData.Data.Aired.From.Month() {
						data, err := server.TMDB.GetTVDetails(int(q.ID), nil)
						if err != nil {
							server.LOG.Error().Msgf("error in get when search serie in tmdb: %v", err.Error())
							http.Error(w, "error in get when search serie in tmdb", http.StatusInternalServerError)
						}

						if data != nil {
							for _, s := range data.Seasons {
								if s.SeasonNumber != 0 {
									server.LOG.Info().Msgf("Papa search serie TMDB Name: %v", s.Name)
									server.LOG.Info().Msgf("Papa search serie TMDB First AirDate: %v", s.AirDate)

									air, err := time.Parse(time.DateOnly, s.AirDate)
									if err != nil {
										server.LOG.Error().Msgf("error when get airDate search serie in tmdb: %v", err.Error())
										http.Error(w, "error when get airDate search serie in tmdb", http.StatusInternalServerError)
									}

									if Query.Aired.Year() == air.Year() {
										tmdbTime = air
										Query.TMDbID = int(s.ID)

										Query.papaSerieTMDbID = int(q.ID)
										samimg, _ := server.TMDB.GetTVImages(int(q.ID), nil)
										if err == nil {
											if samimg != nil {
												for _, b := range samimg.Backdrops {
													if b.FilePath != "" {
														bb, _ := utils.GetBlurHash(server.DecodeIMG+b.FilePath, "")
														SerieBackdrops = append(SerieBackdrops, models.Image{
															Height:    b.Height,
															Width:     b.Width,
															Image:     strings.TrimSpace(fmt.Sprintf(server.OriginalIMG + b.FilePath)),
															Thumbnail: strings.TrimSpace("https://image.tmdb.org/t/p/w500" + b.FilePath),
															BlurHash:  bb,
														})
													}
												}
												for _, p := range samimg.Posters {
													if p.FilePath != "" {
														pp, _ := utils.GetBlurHash(server.DecodeIMG+p.FilePath, "")
														SeriePosters = append(SeriePosters, models.Image{
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

										stttt, _ := server.TMDB.GetTVVideos(int(q.ID), nil)
										if len(stttt.Results) > 0 {
											for _, t := range stttt.Results {
												if strings.Contains(strings.ToLower(t.Site), "youtube") {
													if t.Key != "" {
														SerieTrailers = append(SerieTrailers, models.Trailer{
															Official: true,
															Host:     "YouTube",
															Key:      t.Key,
														})
													}

												}
											}
										}
										server.getTMDBPic(q.PosterPath, q.BackdropPath, &SeriePortraitBlurHash, &SeriePortraitPoster, &SerieLandscapeBlurHash, &SerieLandscapePoster)
										Query.papaSerieAired = fair

										amimg, _ := server.TMDB.GetTVSeasonImages(int(q.ID), s.SeasonNumber, nil)
										if amimg != nil {
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

										server.getTMDBPic(s.PosterPath, "", &PortraitBlurHash, &PortraitPoster, nil, nil)

										tttt, _ := server.TMDB.GetTVSeasonVideos(int(q.ID), s.SeasonNumber, nil)
										if tttt != nil {
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

	if len(MyAnimeListData.Data.Genres) > 0 {
		for _, g := range MyAnimeListData.Data.Genres {
			Genres = append(Genres, g.Name)
		}
	}

	if len(MyAnimeListData.Data.ExplicitGenres) > 0 {
		for _, g := range MyAnimeListData.Data.ExplicitGenres {
			Genres = append(Genres, g.Name)
		}
	}

	if len(MyAnimeListData.Data.Demographics) > 0 {
		for _, g := range MyAnimeListData.Data.Demographics {
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

	animePlanetByte, err := animeResources.Data.AnimePlanetID.MarshalJSON()
	if err != nil {
		AnimePlanetID = ""
	} else {
		AnimePlanetID = string(animePlanetByte)
		AnimePlanetID = strings.ReplaceAll(AnimePlanetID, "\"", "")
	}

	for _, s := range utils.CleanDuplicates(utils.CleanStringArray(Studios)) {
		for _, r := range utils.CleanDuplicates(utils.CleanStringArray(Licensors)) {
			if !strings.Contains(utils.CleanTitle(r), utils.CleanTitle(s)) {
				if strings.TrimSpace(r) != "ltd." {
					PsCs = append(PsCs, r)
				}
			}
		}
	}

	if len(MyAnimeListData.Data.TitleSynonyms) > 0 {
		Titles.Others = append(Titles.Others, MyAnimeListData.Data.TitleSynonyms...)
	}
	if MyAnimeListData.Data.TitleJapanese != "" {
		Titles.Offical = append(Titles.Offical, MyAnimeListData.Data.TitleJapanese)
	} else if MyAnimeListData.Data.TitleEnglish != "" {
		Titles.Offical = append(Titles.Offical, MyAnimeListData.Data.TitleEnglish)
	} else if MyAnimeListData.Data.Title != "" {
		Titles.Offical = append(Titles.Offical, MyAnimeListData.Data.Title)
	}

	for _, d := range GlobalAniDBTitles.Animes {
		if Query.AnidbID == d.Aid {
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

	Aired = utils.CleanDates([]string{MyAnimeListData.Data.Aired.From.Format(time.DateOnly), tmdbTime.Format(time.DateTime), AniDBData.Startdate})

	if Aired.IsZero() {
		Aired = MyAnimeListData.Data.Aired.From
	}
	if MyAnimeListData.Data.Year != 0 {
		ReleaseYear = MyAnimeListData.Data.Year
	} else {
		ReleaseYear = Aired.Year()
	}

	LivechartID := server.Livechart(animeResources.Data.LivechartID, OriginalTitle, Aired)
	AnysearchID := server.Anysearch(animeResources.Data.AnisearchID, MyAnimeListData.Data.TitleEnglish, OriginalTitle, Aired)
	KitsuID := server.Kitsu(animeResources.Data.KitsuID, OriginalTitle, Aired)
	NotifyMoeID := server.NotifyMoe(utils.CleanResText(animeResources.Data.NotifyMoeID), MyAnimeListData.Data.Title, Aired)
	AnilistID := server.Anylist(MyAnimeListData.Data.MalId)

	animeData := models.AnimeSerie{
		SerieMalID:        Query.papaSerieID,
		SerieName:         Query.papaSerieName,
		SerieTVDbID:       Query.papaSerieTVDbID,
		SerieTMDbID:       Query.papaSerieTMDbID,
		Aired:             utils.CleanResText(Query.papaSerieAired.Format(time.DateOnly)),
		PortraitPoster:    SeriePortraitPoster,
		PortraitBlurHash:  SeriePortraitBlurHash,
		LandscapePoster:   SerieLandscapePoster,
		LandscapeBlurHash: SerieLandscapeBlurHash,
		Backdrops:         utils.CleanImages(SerieBackdrops),
		Posters:           utils.CleanImages(SeriePosters),
		Logos:             utils.CleanImages(SerieLogos),
		Trailers:          utils.CleanTrailers(SerieTrailers),
		Season: models.Season{
			OriginalTitle:       OriginalTitle,
			Aired:               Aired.Format(time.DateOnly),
			ReleaseYear:         ReleaseYear,
			Rating:              utils.CleanUnicode(MyAnimeListData.Data.Rating),
			PortraitPoster:      PortraitPoster,
			PortraitBlurHash:    PortraitBlurHash,
			Genres:              utils.CleanStringArray(Genres),
			Studios:             utils.CleanDuplicates(utils.CleanStringArray(Studios)),
			Tags:                utils.CleanStringArray(Tags),
			ProductionCompanies: utils.CleanDuplicates(PsCs),
			Titles:              Titles,
			Posters:             utils.CleanImages(Posters),
			Trailers:            utils.CleanTrailers(Trailers),
			AnimeResources: models.SerieAnimeResources{
				LivechartID:   LivechartID,
				AnimePlanetID: utils.CleanResText(AnimePlanetID),
				AnisearchID:   AnysearchID,
				AnidbID:       Query.AnidbID,
				KitsuID:       KitsuID,
				MalID:         Query.MalID,
				NotifyMoeID:   NotifyMoeID,
				AnilistID:     AnilistID,
				SeasonTVDbID:  Query.TVDbID,
				SeasonTMDbID:  Query.TMDbID,
				Type:          utils.CleanResText(animeResources.Data.Type),
			},
		},
	}

	var TTitle string
	var enT bool
	if MyAnimeListData.Data.TitleEnglish != "" {
		enT = true
		TTitle = MyAnimeListData.Data.TitleEnglish
	} else {
		TTitle = MyAnimeListData.Data.Title
	}

	if TTitle != "" && MyAnimeListData.Data.Synopsis != "" {
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
			utils.CleanOverview(MyAnimeListData.Data.Synopsis),
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

		animeData.Season.AnimeMetas = make([]models.MetaData, len(models.Languages))
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

			animeData.Season.AnimeMetas[i] = models.MetaData{
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
