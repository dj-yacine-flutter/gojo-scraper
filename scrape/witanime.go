package scrape

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dj-yacine-flutter/gojo-scraper/models"
	"github.com/dj-yacine-flutter/gojo-scraper/utils"
)

var (
	re = regexp.MustCompile(`go_to_player\('([^']+)'\)`).FindAllStringSubmatch
)

func (s *Scraper) WitAnime(title string, isMovie bool, malID, year, ep int) ([]models.Iframe, error) {
	var err error

	resp, err := http.Get(fmt.Sprintf("https://witanime.one/?search_param=animes&s=%s", title))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("status code not 200")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	row := doc.Find("div.row.display-flex")
	if row == nil {
		return nil, errors.New("no result found")
	}

	var links []string
	row.Find(".anime-card-details").Each(func(i int, s *goquery.Selection) {
		typ := s.Find(".anime-card-type")
		if typ == nil {
			return
		}

		if isMovie {
			if !strings.Contains(strings.ToLower(typ.Text()), "movie") {
				return
			}
		} else {
			if !strings.Contains(strings.ToLower(typ.Text()), "tv") {
				return
			}
		}

		at := s.Find(".anime-card-title")
		if at == nil {
			return
		}

		if strings.Contains(utils.CleanTitle(at.Text()), utils.CleanTitle(title)) {
			a := at.Find("a")

			href, ok := a.Attr("href")
			if ok {
				links = append(links, href)
			}
		}
	})

	var pages []string
	for _, v := range links {
		resp2, err := http.Get(v)
		if err != nil {
			continue
		}
		defer resp2.Body.Close()

		if resp2.StatusCode != 200 {
			continue
		}

		doc2, err := goquery.NewDocumentFromReader(resp2.Body)
		if err != nil {
			continue
		}

		ex := doc2.Find("div.anime-external-links")
		if ex == nil {
			continue
		}

		mal := ex.Find("a.anime-mal")
		if mal == nil {
			continue
		}

		if malID != 0 {
			href, ok := mal.Attr("href")
			if ok {
				if strings.Contains(href, fmt.Sprint(malID)) {
					if isMovie {
						card := doc2.Find("div.DivEpisodeContainer")
						if card == nil {
							continue
						}

						a := card.Find("a")
						if a == nil {
							continue
						}

						onclick, ok := a.Attr("onclick")
						if ok {
							data := strings.Replace(strings.Replace(onclick, "openEpisode('", "", 1), "')", "", 1)
							page, err := base64.StdEncoding.DecodeString(data)
							if err != nil {
								continue
							}
							pages = append(pages, string(page))

						}
					} else {
						list := doc2.Find("#DivEpisodesList")
						if list == nil {
							continue
						}

						list.Find("div.DivEpisodeContainer").Each(func(i int, s *goquery.Selection) {
							a := s.Find("a")
							if a == nil {
								return
							}

							if strings.TrimSpace(a.Text()) != fmt.Sprintf("الحلقة %d", ep) {
								return
							}

							onclick, ok := a.Attr("onclick")
							if ok {
								data := strings.Replace(strings.Replace(onclick, "openEpisode('", "", 1), "')", "", 1)
								page, err := base64.StdEncoding.DecodeString(data)
								if err != nil {
									return
								}
								pages = append(pages, string(page))

							}
						})
					}
				}
			}
		}
	}

	var iframes []models.Iframe
	for _, v := range pages {
		resp3, err := http.Get(v)
		if err != nil {
			continue
		}
		defer resp3.Body.Close()

		if resp3.StatusCode != 200 {
			continue
		}

		doc3, err := goquery.NewDocumentFromReader(resp3.Body)
		if err != nil {
			continue
		}

		ul := doc3.Find("ul#episode-servers")
		if ul == nil {
			continue
		}

		ul.Find("li").Each(func(i int, s *goquery.Selection) {
			a := s.Find("a")
			if a == nil {
				return
			}

			data, ok := a.Attr("data-url")
			if ok {
				link, err := base64.StdEncoding.DecodeString(data)
				if err != nil {
					return
				}
				url := string(link)
				if url != "" {
					if strings.Contains(url, "yonaplay") {
						req, err := http.NewRequest("GET", url, nil)
						if err != nil {
							return
						}

						req.Header.Add("Authority", "yonaplay.org")
						req.Header.Add("Referer", "https://witanime.one/")

						resp4, err := http.DefaultClient.Do(req)
						if err != nil {
							return
						}
						defer resp4.Body.Close()

						if resp4.StatusCode != 200 {
							return
						}

						doc4, err := goquery.NewDocumentFromReader(resp4.Body)
						if err != nil {
							return
						}

						doc4.Find("li").Each(func(i int, s *goquery.Selection) {
							onclick, ok := s.Attr("onclick")
							if ok {
								matches := re(onclick, -1)
								var frame string
								for _, match := range matches {
									frame = match[1]
									break
								}

								if frame != "" {
									m := models.Iframe{
										Link: frame,
										Type: "sub",
									}

									p := s.Find("p").Text()
									if p != "" {
										p = strings.TrimSpace(strings.ReplaceAll(p, "-", ""))
										m.Quality = p
									}

									iframes = append(iframes, m)
								}

							}
						})
					} else {
						iframes = append(iframes, models.Iframe{
							Link:    url,
							Type:    "sub",
							Quality: "hd",
							Language: "ara",
						})
					}
				}
			}
		})
	}

	return iframes, nil
}
