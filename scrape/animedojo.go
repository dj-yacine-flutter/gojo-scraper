package scrape

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dj-yacine-flutter/gojo-scraper/models"
	"github.com/dj-yacine-flutter/gojo-scraper/utils"
)

func (s *Scraper) AnimeDojo(title string, isMovie bool, year, ep int) ([]models.Iframe, error) {
	root := "https://animedojo.net"
	var err error
	resp, err := http.Get(fmt.Sprintf("%s/search?keyword=%s", root, utils.CleanQuery(title)))
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

	content := doc.Find(".tab-content")
	if content == nil {
		return nil, errors.New("no data found")
	}

	wrap := content.Find(".film_list-wrap")
	if wrap == nil {
		return nil, errors.New("no data found")
	}

	var pages []string
	wrap.Find(".film-detail").Each(func(i int, s *goquery.Selection) {
		a := s.Find("a")
		if a == nil {
			return
		}

		if strings.Contains(utils.CleanTitle(a.Text()), utils.CleanTitle(title)) {
			if isMovie {
				yt := s.Find(".fd-infor")
				if !strings.Contains(yt.Text(), fmt.Sprint(year)) {
					return
				}
			}
			href, ok := a.Attr("href")
			if ok {
				pages = append(pages, strings.TrimSpace(root+href))
			}
		}
	})

	links := make(map[string]string, len(pages))
	for _, page := range pages {
		resp2, err := http.Get(page)
		if err != nil {
			continue
		}
		defer resp2.Body.Close()

		doc2, err := goquery.NewDocumentFromReader(resp2.Body)
		if err != nil {
			continue
		}

		if !isMovie {
			prm := doc2.Find("div.anisc-info")
			if prm == nil {
				continue
			}

			if !strings.Contains(prm.Text(), fmt.Sprint(year)) {
				continue
			}
		}

		stats := doc2.Find(".film-stats")
		if stats == nil {
			continue
		}

		typ := stats.Find(".item")
		if typ == nil {
			continue
		}

		if isMovie {
			if !strings.Contains(strings.ToLower(typ.Text()), "movie") {
				continue
			}
		} else {
			if strings.Contains(strings.ToLower(typ.Text()), "special") {
				continue
			}
		}

		card := doc2.Find(".film-buttons")
		if card == nil {
			continue
		}

		a := card.Find("a")
		if a == nil {
			continue
		}

		href, ok := a.Attr("href")
		if ok {
			q := stats.Find(".tac.tick-item.tick-dub")
			links[strings.TrimSpace(root+href)] = strings.TrimSpace(q.Text())
		}
	}

	var iframes []models.Iframe
	for k, v := range links {
		resp3, err := http.Get(k)
		if err != nil {
			continue
		}
		defer resp3.Body.Close()

		doc3, err := goquery.NewDocumentFromReader(resp3.Body)
		if err != nil {
			continue
		}

		if isMovie {
			servers := doc3.Find("#servers-content")
			if servers == nil {
				continue
			}

			servers.Find(".item").Each(func(i int, s *goquery.Selection) {
				a := s.Find("a")
				if a == nil {
					return
				}

				href, ok := a.Attr("href")
				if ok {
					iframes = append(iframes, models.Iframe{
						Link:    href,
						Quality: "hd",
						Referer: strings.TrimSpace(k),
						Type:    strings.ToLower(v),
					})
				}
			})
		} else {
			episodes := doc3.Find("#episodes-content")
			if episodes == nil {
				continue
			}

			var url string
			episodes.Find("a").Each(func(i int, s *goquery.Selection) {
				if url != "" {
					return
				}

				number, ok := s.Attr("data-number")
				if ok {
					num, _ := strconv.Atoi(number)
					if num == ep {
						href, ok := s.Attr("href")
						if ok {
							url = strings.TrimSpace(root + href)
							return
						}
					}
				}
			})

			if url == "" {
				continue
			}

			resp4, err := http.Get(url)
			if err != nil {
				continue
			}
			defer resp4.Body.Close()

			doc4, err := goquery.NewDocumentFromReader(resp4.Body)
			if err != nil {
				continue
			}

			servers := doc4.Find("#servers-content")
			if servers == nil {
				continue
			}

			servers.Find(".item").Each(func(i int, s *goquery.Selection) {
				a := s.Find("a")
				if a == nil {
					return
				}

				href, ok := a.Attr("href")
				if ok {
					iframes = append(iframes, models.Iframe{
						Link:    href,
						Quality: "hd",
						Referer: strings.TrimSpace(url),
						Type:    strings.ToLower(v),
					})
				}
			})
		}
	}

	if len(iframes) == 0 {
		return nil, errors.New("no iframe found")
	}

	return iframes, nil
}
