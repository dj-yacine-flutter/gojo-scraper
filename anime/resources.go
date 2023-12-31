package anime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dj-yacine-flutter/gojo-scraper/utils"
)

var (
	UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"
)

func (server *AnimeScraper) Livechart(id int, title string, date time.Time) int {
	url := fmt.Sprintf("https://www.livechart.me/anime/%d", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0
	}
	req.Header.Set("Authority", "www.livechart.me")
	req.Header.Set("Referer", url)
	req.Header.Set("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			id = 0
		}

		block := doc.Find(".lc-poster-col")
		if block != nil {
			var found bool
			block.Find(".text-sm").Each(func(index int, selection *goquery.Selection) {
				if selection != nil {
					hover := selection.Find(".link-hover")
					if hover != nil {
						href, ok := hover.Attr("href")
						if ok {
							if strings.Contains(href, "timetable") {
								if hover.Text() != "" {
									lay := "January 02, 2006"
									parsedTime, err := time.Parse(lay, hover.Text())
									if err == nil {
										if parsedTime.Year() == date.Year() && parsedTime.Month() == date.Month() {
											found = true
											return
										}
									}
								}
							}
						}
					}
				}
			})
			if found {
				return id
			}
		}
		id = 0
	} else {
		id = 0
	}
	if id == 0 {
		titled := utils.CleanUnicode(title)
		titled = strings.ReplaceAll(strings.ToLower(titled), "movie", "")
		titled = strings.ReplaceAll(strings.ToLower(titled), "gekijouban", "")

		query := fmt.Sprintf("https://www.livechart.me/search?q=%s", titled)
		req, err := http.NewRequest("GET", query, nil)
		if err != nil {
			return 0
		}

		req.Header.Set("User-Agent", UserAgent)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return 0
		}

		defer resp.Body.Close()
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return 0
		}
		block := doc.Find(".anime-list")
		if block != nil {
			var found bool
			var ID int
			block.Find(".info").Each(func(index int, selection *goquery.Selection) {
				if selection != nil {
					span := selection.Find("span").First()
					if span != nil {
						layout := "January 2, 2006"
						parsedTime, err := time.Parse(layout, span.Text())
						if err == nil {
							if parsedTime.Year() == date.Year() && parsedTime.Month() == date.Month() {
								title := block.Find("strong").First()
								if title != nil {
									a := title.Find("a").First()
									if a != nil {
										link, ok := a.Attr("href")
										if ok {
											ID = utils.ExtractID(link)
											if ID != 0 {
												found = true
												return
											}
										}
									}
								}
							}
						}
					}
				}
			})
			if found {
				return ID
			}
		}
	}
	return 0
}

func (server *AnimeScraper) Anysearch(id int, title, originalTitle string, date time.Time) int {
	url := fmt.Sprintf("https://www.anisearch.com/anime/%d", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0
	}
	req.Header.Set("Authority", "www.anisearch.com")
	req.Header.Set("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			id = 0
		}
		block := doc.Find(".infoblock")
		if block != nil {
			var found bool
			block.Find(".title").Each(func(index int, selection *goquery.Selection) {
				grey := selection.Find(".grey")
				if (grey != nil && strings.Contains(strings.ToLower(grey.Text()), strings.ToLower(originalTitle))) || strings.Contains(strings.ToLower(selection.Text()), strings.ToLower(title)) {
					found = true
					return
				} else {
					release := block.Find(".released").First()
					if release != nil {
						if !strings.Contains(release.Text(), "‑") {
							layout := "02.01.2006"
							parsedTime, err := time.Parse(layout, release.Text())
							if err == nil {
								if parsedTime.Year() == date.Year() && parsedTime.Month() == date.Month() {
									found = true
									return
								}
							}
						} else {
							nr := strings.ReplaceAll(strings.ToLower(release.Text()), "release date: ", "")
							dateRange := strings.Split(nr, "‑")
							layout := "02.01.2006"
							for _, d := range dateRange {
								ntime := strings.ReplaceAll(d, " ", "")
								ptime, err := time.Parse(layout, ntime)
								if err == nil {
									if ptime.Year() == date.Year() && ptime.Month() == date.Month() {
										found = true
										return
									}
								}
							}
						}
					}
				}
			})
			if found {
				return id
			}
		}
		id = 0
	}
	return 0
}

type KitsuData struct {
	Data struct {
		ID    string `json:"id,omitempty"`
		Type  string `json:"type,omitempty"`
		Links struct {
			Self string `json:"self,omitempty"`
		} `json:"links,omitempty"`
		Attributes struct {
			Slug   string `json:"slug,omitempty"`
			Titles struct {
				En   string `json:"en,omitempty"`
				EnJp string `json:"en_jp,omitempty"`
				JaJp string `json:"ja_jp,omitempty"`
			} `json:"titles,omitempty"`
			CanonicalTitle    string   `json:"canonicalTitle,omitempty"`
			AbbreviatedTitles []string `json:"abbreviatedTitles,omitempty"`
			StartDate         string   `json:"startDate,omitempty"`
			AgeRating         string   `json:"ageRating,omitempty"`
			Status            string   `json:"status,omitempty"`
			EpisodeCount      int      `json:"episodeCount,omitempty"`
			EpisodeLength     int      `json:"episodeLength,omitempty"`
			ShowType          string   `json:"showType,omitempty"`
		} `json:"attributes,omitempty"`
	} `json:"data,omitempty"`
}

type KitsuSearch struct {
	Data []struct {
		ID    string `json:"id,omitempty"`
		Type  string `json:"type,omitempty"`
		Links struct {
			Self string `json:"self,omitempty"`
		} `json:"links,omitempty"`
		Attributes struct {
			Slug   string `json:"slug,omitempty"`
			Titles struct {
				En   string `json:"en,omitempty"`
				EnJp string `json:"en_jp,omitempty"`
				JaJp string `json:"ja_jp,omitempty"`
			} `json:"titles,omitempty"`
			CanonicalTitle    string   `json:"canonicalTitle,omitempty"`
			AbbreviatedTitles []string `json:"abbreviatedTitles,omitempty"`
			StartDate         string   `json:"startDate,omitempty"`
			AgeRating         string   `json:"ageRating,omitempty"`
			Status            string   `json:"status,omitempty"`
			EpisodeCount      int      `json:"episodeCount,omitempty"`
			EpisodeLength     int      `json:"episodeLength,omitempty"`
			ShowType          string   `json:"showType,omitempty"`
		} `json:"attributes,omitempty"`
	} `json:"data,omitempty"`
}

func (server *AnimeScraper) Kitsu(id int, title string, date time.Time) int {
	if id != 0 {
		url := fmt.Sprintf("https://kitsu.io/api/edge/anime/%d", id)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return 0
		}
		req.Header.Set("authority", "kitsu.io")
		req.Header.Set("origin", "https://kitsu.io")
		req.Header.Set("referer", "https://kitsu.io")
		req.Header.Set("user-agent", UserAgent)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return 0
		}
		defer resp.Body.Close()
		anime := KitsuData{}
		err = json.NewDecoder(resp.Body).Decode(&anime)
		if err != nil {
			return 0
		}

		ptime, err := time.Parse(time.DateOnly, anime.Data.Attributes.StartDate)
		if err == nil {
			if ptime.Year() == date.Year() && ptime.Month() == date.Month() {
				ID, err := strconv.Atoi(anime.Data.ID)
				if err == nil {
					return ID
				}
			}
		}
		var titles []string
		titles = append(titles, anime.Data.Attributes.Titles.En)
		titles = append(titles, anime.Data.Attributes.Titles.EnJp)
		titles = append(titles, anime.Data.Attributes.Titles.JaJp)
		titles = append(titles, anime.Data.Attributes.CanonicalTitle)
		titles = append(titles, anime.Data.Attributes.AbbreviatedTitles...)

		for _, t := range titles {
			if strings.Contains(utils.CleanTitle(t), utils.CleanTitle(title)) {
				return id
			}
		}
	} else {
		query := fmt.Sprintf("https://kitsu.io/api/edge/anime?filter[text]=%s", title)
		req, err := http.NewRequest("GET", query, nil)
		if err != nil {
			return 0
		}
		req.Header.Set("authority", "kitsu.io")
		req.Header.Set("origin", "https://kitsu.io")
		req.Header.Set("referer", "https://kitsu.io")
		req.Header.Set("user-agent", UserAgent)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return 0
		}
		defer resp.Body.Close()

		animes := KitsuSearch{}
		err = json.NewDecoder(resp.Body).Decode(&animes)
		if err != nil {
			return 0
		}

		for _, d := range animes.Data {
			ptime, err := time.Parse(time.DateOnly, d.Attributes.StartDate)
			if err == nil {
				if ptime.Year() == date.Year() && ptime.Month() == date.Month() {
					ID, err := strconv.Atoi(d.ID)
					if err == nil {
						return ID
					}
				}
			}

			var titles []string
			titles = append(titles, d.Attributes.Titles.En)
			titles = append(titles, d.Attributes.Titles.EnJp)
			titles = append(titles, d.Attributes.Titles.JaJp)
			titles = append(titles, d.Attributes.CanonicalTitle)
			titles = append(titles, d.Attributes.AbbreviatedTitles...)

			for _, t := range titles {
				if strings.Contains(utils.CleanTitle(t), utils.CleanTitle(title)) {
					ID, err := strconv.Atoi(d.ID)
					if err == nil {
						return ID
					}
				}
			}
		}
	}
	return 0
}

type NotifyMoe struct {
	ID    string `json:"id,omitempty"`
	Type  string `json:"type,omitempty"`
	Title struct {
		Canonical string   `json:"canonical,omitempty"`
		Romaji    string   `json:"romaji,omitempty"`
		English   string   `json:"english,omitempty"`
		Japanese  string   `json:"japanese,omitempty"`
		Hiragana  string   `json:"hiragana,omitempty"`
		Synonyms  []string `json:"synonyms,omitempty"`
	} `json:"title,omitempty"`
	Summary       string   `json:"summary,omitempty"`
	Status        string   `json:"status,omitempty"`
	Genres        []string `json:"genres,omitempty"`
	StartDate     string   `json:"startDate,omitempty"`
	EndDate       string   `json:"endDate,omitempty"`
	EpisodeCount  int      `json:"episodeCount,omitempty"`
	EpisodeLength int      `json:"episodeLength,omitempty"`
	Source        string   `json:"source,omitempty"`
	Image         struct {
		Extension    string `json:"extension,omitempty"`
		Width        int    `json:"width,omitempty"`
		Height       int    `json:"height,omitempty"`
		AverageColor struct {
			Hue        float64 `json:"hue,omitempty"`
			Saturation float64 `json:"saturation,omitempty"`
			Lightness  float64 `json:"lightness,omitempty"`
		} `json:"averageColor,omitempty"`
		LastModified int `json:"lastModified,omitempty"`
	} `json:"image,omitempty"`
	FirstChannel string `json:"firstChannel,omitempty"`
	Rating       struct {
		Overall    float64 `json:"overall,omitempty"`
		Story      float64 `json:"story,omitempty"`
		Visuals    float64 `json:"visuals,omitempty"`
		Soundtrack float64 `json:"soundtrack,omitempty"`
		Count      struct {
			Overall    int `json:"overall,omitempty"`
			Story      int `json:"story,omitempty"`
			Visuals    int `json:"visuals,omitempty"`
			Soundtrack int `json:"soundtrack,omitempty"`
		} `json:"count,omitempty"`
	} `json:"rating,omitempty"`
	Popularity struct {
		Watching  int `json:"watching,omitempty"`
		Completed int `json:"completed,omitempty"`
		Planned   int `json:"planned,omitempty"`
		Hold      int `json:"hold,omitempty"`
		Dropped   int `json:"dropped,omitempty"`
	} `json:"popularity,omitempty"`
	Trailers []struct {
		Service   string `json:"service,omitempty"`
		ServiceID string `json:"serviceId,omitempty"`
	} `json:"trailers,omitempty"`
	Episodes []string `json:"episodes,omitempty"`
	Mappings []struct {
		Service   string `json:"service,omitempty"`
		ServiceID string `json:"serviceId,omitempty"`
	} `json:"mappings,omitempty"`
	Posts     []string `json:"posts,omitempty"`
	Likes     any      `json:"likes,omitempty"`
	Created   string   `json:"created,omitempty"`
	CreatedBy string   `json:"createdBy,omitempty"`
	Edited    string   `json:"edited,omitempty"`
	EditedBy  string   `json:"editedBy,omitempty"`
	IsDraft   bool     `json:"isDraft,omitempty"`
	Studios   []string `json:"studios,omitempty"`
	Producers []string `json:"producers,omitempty"`
	Licensors any      `json:"licensors,omitempty"`
	Links     []struct {
		Title string `json:"title,omitempty"`
		URL   string `json:"url,omitempty"`
	} `json:"links,omitempty"`
}

func (server *AnimeScraper) NotifyMoe(id, title string, date time.Time) string {
	if id != "" {
		url := fmt.Sprintf("https://notify.moe/api/anime/%s", id)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			id = ""
		}
		req.Header.Set("authority", "notify.moe")
		req.Header.Set("referer", url)
		req.Header.Set("user-agent", UserAgent)

		resp, err := server.HTTP.Do(req)
		if err != nil {
			id = ""
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			anime := NotifyMoe{}
			err = json.NewDecoder(resp.Body).Decode(&anime)
			if err != nil {
				id = ""
			}

			var created time.Time
			if anime.StartDate != "" {
				created, err = time.Parse(time.DateOnly, anime.StartDate)
				if err != nil {
					id = ""
				}
			} else {
				created, err = time.Parse(time.DateOnly, anime.EndDate)
				if err != nil {
					id = ""
				}
			}

			if created.String() != "" {
				if created.Year() == date.Year() && created.Month() == date.Month() {
					id = anime.ID
				} else {
					id = ""
				}
			} else {
				var titles []string
				titles = append(titles, anime.Title.Canonical)
				titles = append(titles, anime.Title.English)
				titles = append(titles, anime.Title.Japanese)
				titles = append(titles, anime.Title.Hiragana)

				for _, t := range titles {
					if strings.Contains(utils.CleanTitle(t), utils.CleanTitle(title)) {
						id = anime.ID
						break
					}
				}
			}
		}
	}

	if id == "" {
		query := fmt.Sprintf("https://notify.moe/_/anime-search/%s", utils.CleanQuery(title))
		req2, err := http.NewRequest("GET", query, nil)
		if err != nil {
			return ""
		}

		req2.Header.Set("authority", "notify.moe")
		req2.Header.Set("user-agent", UserAgent)

		resp2, err := server.HTTP.Do(req2)
		if err != nil {
			return ""
		}
		defer resp2.Body.Close()

		if resp2.StatusCode == 200 {
			doc, err := goquery.NewDocumentFromReader(resp2.Body)
			if err != nil {
				return ""
			}

			block := doc.Find(".anime-search")
			if block != nil {
				var found bool
				block.Find(".profile-watching-list-item").Each(func(index int, selection *goquery.Selection) {
					if found {
						return
					}
					data, ok := selection.Attr("aria-label")
					if ok {
						href, ok := selection.Attr("href")
						if ok {
							s := strings.Split(href, "/")
							id = s[len(s)-1]
							if id != "" {
								url := fmt.Sprintf("https://notify.moe/api/anime/%s", id)
								req3, err := http.NewRequest("GET", url, nil)
								if err != nil {
									return
								}

								req3.Header.Set("authority", "notify.moe")
								req3.Header.Set("referer", url)
								req3.Header.Set("user-agent", UserAgent)

								resp3, err := server.HTTP.Do(req3)
								if err != nil {
									return

								}
								defer resp3.Body.Close()

								anime := NotifyMoe{}
								err = json.NewDecoder(resp3.Body).Decode(&anime)
								if err != nil {
									return
								}

								var created time.Time
								if anime.StartDate != "" {
									created, err = time.Parse(time.DateOnly, anime.StartDate)
									if err != nil {
										return
									}
								} else {
									created, err = time.Parse(time.DateOnly, anime.EndDate)
									if err != nil {
										return
									}
								}

								if created.Year() == date.Year() && created.Month() == date.Month() {
									id = anime.ID
									found = true
									return
								}

								if id == "" {
									var titles []string
									titles = append(titles, anime.Title.Canonical)
									titles = append(titles, anime.Title.English)
									titles = append(titles, anime.Title.Japanese)
									titles = append(titles, anime.Title.Hiragana)

									for _, t := range titles {
										if strings.Contains(utils.CleanTitle(t), utils.CleanTitle(title)) {
											id = anime.ID
											found = true
											break
										}
									}
								}
							}
							if id == "" {
								if strings.Contains(utils.CleanTitle(data), utils.CleanTitle(title)) {
									found = true
									return
								}
							}
						}
					}
				})
				if found {
					return id
				}
			} else {
				id = ""
			}
		}
	}

	return id
}

func (server *AnimeScraper) Anylist(mal int) int {
	type Variables struct {
		MALID int    `json:"malId"`
		Type  string `json:"type"`
	}

	body := struct {
		Query     string    `json:"query"`
		Variables Variables `json:"variables"`
	}{
		Query: `
			query ($malId: Int, $type: MediaType) {
				Media(idMal: $malId, type: $type) {
					id
				}
			}
		`,
		Variables: Variables{
			MALID: mal,
			Type:  "ANIME",
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return 0
	}

	req, err := http.NewRequest("POST", "https://graphql.anilist.co", bytes.NewBuffer(data))
	if err != nil {
		return 0
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("user-agent", UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	response := new(struct {
		Data struct {
			Media struct {
				ID int `json:"id"`
			} `json:"Media"`
		} `json:"data"`
	})

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return 0
	}

	return response.Data.Media.ID
}
