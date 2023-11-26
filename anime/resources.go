package anime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dj-yacine-flutter/gojo-scraper/models"
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
				if strings.Contains(strings.ToLower(selection.Text()), strings.ToLower(title)) {
					found = true
					return
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
		query := fmt.Sprintf("https://www.livechart.me/search?q=%s", utils.CleanUnicode(title))
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
			block.Find(".grey").Each(func(index int, selection *goquery.Selection) {
				if strings.Contains(strings.ToLower(selection.Text()), strings.ToLower(originalTitle)) {
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
			anime := models.NotifyMoe{}
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

								anime := models.NotifyMoe{}
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