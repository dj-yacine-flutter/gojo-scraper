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
		ReleaseYear       int
		AgeRating         string
		PortriatPoster    string
		PortriatBlurHash  string
		LandscapePoster   string
		LandscapeBlurHash string
		AnimePlanetID     string
		OriginalTitle     string
		Aired             time.Time
		Genres            []string
		Studios           []string
		Tags              []string
		PsCs              []string
		Titles            models.Titles
		Posters           []models.Image
		Backdrops         []models.Image
		Logos             []models.Image
		Trailers          []models.Trailer
		Licensors         []string
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

	server.getMalPic(AniDBData.Picture, MyAnimeListData.Data.Images.Jpg.LargeImageUrl, MyAnimeListData.Data.Images.Webp.LargeImageUrl, &PortriatBlurHash, &PortriatPoster)

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

											if Query.papaSerieTVDbID == 0 {
												Query.papaSerieTVDbID = season.Data.SeriesID
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
										Query.papaSerieAired = fair
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

	/*	var queries []string
		var totalSearch tvdb.Search
		queries = append(queries, MyAnimeListData.Data.TitleEnglish, MyAnimeListData.Data.Title)
		queries = append(queries, MyAnimeListData.Data.TitleSynonyms...)

		if len(queries) > 0 {
			for _, t := range queries {
				Series, err := server.TVDB.GetSearch(t, ReleaseYear)
				if err != nil {
					continue
				}
				totalSearch.Data = append(totalSearch.Data, Series.Data...)
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
			if MyAnimeListData.Data.Aired.From.String() != "" {
				aDate = MyAnimeListData.Data.Aired.From

			} else {
				aDate, err = time.Parse(time.DateOnly, AniDBData.Startdate)
				if err != nil {
					continue
				}
			}

			if aDate.Year() == qDate.Year() && aDate.Month() == qDate.Month() {
				if strings.Contains(a.Type, "tv") {
					newTVDBid, err := strconv.Atoi(a.TvdbID)
					if err != nil {
						continue
					}
					TVDbID = int(newTVDBid)
					Serie, err := server.TVDB.GetSeriesByIDExtanded(TVDbID)
					if err != nil {
						continue
					}
					if Serie != nil {
						for _, r := range Serie.Data.RemoteIds {
							if strings.Contains(strings.ToLower(r.SourceName), "imdb") && r.SourceName != "" {
								IMDbID = r.ID
							}
						}

						if len(Serie.Data.Genres) > 0 {
							for _, g := range Serie.Data.Genres {
								Genres = append(Genres, g.Name)
							}
						}

						if len(Serie.Data.Companies) > 0 {
							for _, d := range Serie.Data.Companies {
								if d.Name != "" {
									Studios = append(Studios, d.Name)
								}
							}
						}
						if len(Serie.Data.Artworks) > 0 {
							for _, d := range Serie.Data.Artworks {
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
						if len(Serie.Data.Trailers) > 0 {
							for _, t := range Serie.Data.Trailers {
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
				} else if strings.Contains(a.Type, "movie") {
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
										Licensors = append(Licensors, p.Name)
									}
								}
							}
						}

						if len(movie.Data.Companies.Distributor) > 0 {
							for _, d := range movie.Data.Companies.Distributor {
								if d.Name != "" {
									Licensors = append(Licensors, d.Name)
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
				}
			}
		}

		if TVDbID == 0 {
			if animeResources.Data.TheTVdbID != 0 {
				Serie, err := server.TVDB.GetSeriesByIDExtanded(animeResources.Data.TheTVdbID)
				if err != nil {
					TVDbID = animeResources.Data.TheTVdbID
				}

				if Serie != nil {
					for _, r := range Serie.Data.RemoteIds {
						if strings.Contains(strings.ToLower(r.SourceName), "imdb") && r.SourceName != "" {
							IMDbID = r.ID
						}
					}

					if len(Serie.Data.Genres) > 0 {
						for _, g := range Serie.Data.Genres {
							Genres = append(Genres, g.Name)
						}
					}

					if len(Serie.Data.Companies) > 0 {
						for _, d := range Serie.Data.Companies {
							if d.Name != "" {
								Licensors = append(Licensors, d.Name)
							}
						}
					}
					if len(Serie.Data.Artworks) > 0 {
						for _, d := range Serie.Data.Artworks {
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
					if len(Serie.Data.Trailers) > 0 {
						for _, t := range Serie.Data.Trailers {
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
				anime, err := server.TMDB.GetTVDetails(TMDbID, nil)
				if err != nil {
					PortriatBlurHash = ""
					LandscapeBlurHash = ""
					TMDbID = 0
				} else {
					var rd bool
					if anime.FirstAirDate != "" {
						eDate, err := time.Parse(time.DateOnly, anime.FirstAirDate)
						if err != nil {
							PortriatBlurHash = ""
							LandscapeBlurHash = ""
							TMDbID = 0
						}
						qDate := MyAnimeListData.Data.Aired.From
						if eDate.Year() == qDate.Year() && eDate.Month() == qDate.Month() {
							rd = true
						}
					}

					if rd {
						if OriginalTitle == "" {
							OriginalTitle = anime.OriginalName
						}
						TMDBRuntime = fmt.Sprintf("%dm", anime.EpisodeRunTime[0])
						TMDBTitle = anime.Name
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
									Licensors = append(Licensors, p.Name)
								}
							}
						}

						amimg, _ := server.TMDB.GetTVImages(TMDbID, nil)
						if err == nil {
							if amimg != nil {
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
						tttt, _ := server.TMDB.GetTVVideos(TMDbID, nil)
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
			anime, _ := server.TMDB.GetTVDetails(TMDbID, nil)
			if anime != nil {
				OriginalTitle = anime.OriginalName
				TMDBTitle = anime.Name
				TMDBRuntime = fmt.Sprintf("%dm", anime.EpisodeRunTime[0])
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
							Licensors = append(Licensors, p.Name)
						}
					}
				}
				amimg, _ := server.TMDB.GetTVImages(TMDbID, nil)
				if err == nil {
					if amimg != nil {
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
				tttt, _ := server.TMDB.GetTVVideos(TMDbID, nil)
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
			querys, _ := server.TMDB.GetSearchMulti(MyAnimeListData.Data.TitleEnglish, nil)
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
						if strings.Contains(strings.ToLower(q.MediaType), "Serie") {
							TMDbID = int(q.ID)
							if OriginalTitle == "" {
								OriginalTitle = q.OriginalTitle
							}
							server.getTMDBPic(q.PosterPath, q.BackdropPath, &PortriatBlurHash, &PortriatPoster, &LandscapeBlurHash, &LandscapePoster)
							server.getTMDBRating(TMDbID, &AgeRating)
							results, _ := server.TMDB.GetGenreTVList(nil)
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

							anime, _ := server.TMDB.GetTVDetails(int(q.ID), nil)
							if len(anime.ProductionCompanies) > 0 {
								for _, p := range anime.ProductionCompanies {
									if p.Name != "" {
										Licensors = append(Licensors, p.Name)
									}
								}
							}
							amimg, _ := server.TMDB.GetTVImages(TMDbID, nil)
							if err == nil {
								if amimg != nil {
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
							tttt, _ := server.TMDB.GetTVVideos(TMDbID, nil)
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
			server.getMalPic(AniDBData.Picture, MyAnimeListData.Data.Images.Jpg.LargeImageUrl, MyAnimeListData.Data.Images.Webp.LargeImageUrl, &PortriatBlurHash, &PortriatPoster, &LandscapeBlurHash, &LandscapePoster)
		}
	*/

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

	if MyAnimeListData.Data.Rating == "" {
		AgeRating, err = utils.CleanRating(MyAnimeListData.Data.Rating)
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
		SerieMalID:  Query.papaSerieID,
		SerieName:   Query.papaSerieName,
		SerieTVDbID: Query.papaSerieTVDbID,
		SerieTMDbID: Query.papaSerieTMDbID,
		Aired:       utils.CleanResText(Query.papaSerieAired.Format(time.DateOnly)),
		Season: models.Season{
			OriginalTitle:       OriginalTitle,
			Aired:               Aired.Format(time.DateOnly),
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

	/* if MyAnimeListData.Data.TitleEnglish != "" && MyAnimeListData.Data.Synopsis != "" {
		translation, err := gtranslate.TranslateWithParams(
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
				Title:    MyAnimeListData.Data.TitleEnglish,
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
