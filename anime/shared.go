package anime

import (
	"fmt"
	"strings"
	"time"

	jikan "github.com/darenliang/jikan-go"
	"github.com/dj-yacine-flutter/gojo-scraper/utils"
)

func (server *AnimeScraper) getTMDBRating(TMDbID int, AgeRating *string) {
	results, err := server.TMDB.GetMovieReleaseDates(TMDbID)
	if err != nil {
		return
	}
	if results != nil {
		for _, r := range results.Results {
			if strings.Contains(strings.ToLower(r.Iso3166_1), "us") {
				for _, t := range r.ReleaseDates {
					if t.Certification != "" {
						*AgeRating, err = utils.CleanRating(t.Certification)
						if err != nil {
							*AgeRating = ""
							continue
						}
						break
					}
				}
			}
		}
	}
}

func (server *AnimeScraper) getTMDBPic(posterPath, backdropPath string, PortriatBlurHash, PortriatPoster, LandscapeBlurHash, LandscapePoster *string) {
	var err error
	if posterPath != "" {
		*PortriatBlurHash, err = utils.GetBlurHash(server.DecodeIMG, posterPath)
		if err != nil {
			*PortriatBlurHash = ""
			*PortriatPoster = ""
		}
		*PortriatPoster = server.OriginalIMG + posterPath
	}

	if backdropPath != "" {
		*LandscapeBlurHash, err = utils.GetBlurHash(server.DecodeIMG, backdropPath)
		if err != nil {
			*LandscapeBlurHash = ""
			*LandscapePoster = ""
		}
		*LandscapePoster = server.OriginalIMG + backdropPath
	}
}

func (server *AnimeScraper) getMalPic(Pic, JPG, WEBP string, PortriatBlurHash, PortriatPoster *string) {
	var err error
	img := ""
	if JPG != "" {
		img = JPG
	} else if WEBP != "" {
		img = WEBP
	} else {
		img = fmt.Sprint("https://cdn-eu.anidb.net/images/main/" + Pic)
	}
	*PortriatBlurHash, err = utils.GetBlurHash(img, "")
	if err != nil {
		*PortriatBlurHash = ""
	}

	*PortriatPoster = img
}

func (server *AnimeScraper) getAniDBIDFromTitles(malData *jikan.AnimeById) (int, error) {
	for _, v := range GlobalAniDBTitles.Animes {
		for _, title := range v.Titles {
			titles := append(malData.Data.TitleSynonyms, malData.Data.Title, malData.Data.TitleEnglish)
			for _, mt := range titles {
				titleMatches := strings.Contains(strings.ToLower(title.Value), strings.ToLower(mt))
				if titleMatches {

					aniDBData, err := server.GetAniDBData(v.Aid)
					if err != nil {
						return 0, err
					}
					typeM := strings.Contains(strings.ToLower(aniDBData.Type), strings.ToLower(malData.Data.Type))

					aniY, err := utils.ExtractYear(aniDBData.Startdate)
					if err != nil {
						return 0, err
					}

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

func (server *AnimeScraper) searchAniDBID(malData *jikan.AnimeById, links []Link) (int, error) {
	anidbID, err := server.getAniDBID(links)
	if err != nil {
		anidbID, err = server.getAniDBIDFromTitles(malData)
		if err != nil {
			return 0, err
		}
	}
	if anidbID == 0 {
		return 0, fmt.Errorf("there is no AniDB ID for this anime")
	}
	return anidbID, nil
}

func (server *AnimeScraper) getResourceByIDs(anidbID, malID int) (AnimeResources, error) {
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

func (server *AnimeScraper) getMALOriginalID(malID int) int {
	id := -1
	for {
		time.Sleep(500 * time.Millisecond)
		if id != malID {
			if id == -1 {
				id = malID
			}
			server.LOG.Info().Msgf("Mal try: %d", id)

			data, err := jikan.GetAnimeRelations(id)
			if err != nil {
				server.LOG.Error().Msgf("Mal try Error: %v", err.Error())
				continue
			}
			//gg, _ := json.Marshal(&data)
			//server.LOG.Info().Msg(string(gg))
			if len(data.Data) > 0 {
				var p bool
				for _, e := range data.Data {
					if strings.Contains(strings.ToLower(e.Relation), "prequel") {
						for _, q := range e.Entry {
							if strings.Contains(strings.ToLower(q.Type), "anime") {
								id = q.MalId
								p = true
								continue
							}
						}
						if p {
							continue
						}
					}
				}

				if p {
					continue
				}

				for _, e := range data.Data {
					if strings.Contains(strings.ToLower(e.Relation), "sequel") {
						for _, q := range e.Entry {
							if strings.Contains(strings.ToLower(q.Type), "anime") {
								return id
							}
						}
					}
				}

				return malID
			} else {
				break
			}
		} else {
			break
		}
	}

	return 0
}
