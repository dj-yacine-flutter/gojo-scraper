package scrape

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/dj-yacine-flutter/gojo-scraper/models"
	"github.com/dj-yacine-flutter/gojo-scraper/utils"
)

func (s *Scraper) AnimeLek(title string, isMovie bool, malID, year, ep int) ([]models.Iframe, error) {
	var err error

	resp, err := http.Get(fmt.Sprintf("https://animelek.xyz/search/?s=%s", title))
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

	content := doc.Find(".anime-list-content")
	if content == nil {
		return nil, errors.New("no result found")
	}

	var paths []string

	content.Find(".anime-card-container").Each(func(i int, s *goquery.Selection) {
		details := s.Find(".anime-card-details")
		if details == nil {
			return
		}

		card := details.Find(".anime-card-title")
		if card == nil {
			return
		}

		at, ok := card.Attr("title")
		if ok {
			if strings.Contains(utils.CleanTitle(at), utils.CleanTitle(title)) {
				a := card.Find("a")
				if a == nil {
					return
				}

				href, ok := a.Attr("href")
				if ok {
					paths = append(paths, href)
				}
			}
		}
	})

	var (
		url   string
		movie bool
	)

	for _, path := range paths {
		if url != "" {
			break
		}

		resp2, err := http.Get(path)
		if err != nil {
			continue
		}
		defer resp2.Body.Close()

		doc2, err := goquery.NewDocumentFromReader(resp2.Body)
		if err != nil {
			continue
		}

		column := doc2.Find(".anime-container-infos")
		if column == nil {
			continue
		}

		column.Find(".full-list-info").Each(func(i int, s *goquery.Selection) {
			if url != "" {
				return
			}

			if malID != 0 {
				a := s.Find("a")
				if a == nil {
					return
				}

				href, ok := a.Attr("href")
				if ok {
					if strings.Contains(href, "myanimelist") {
						if strings.Contains(href, fmt.Sprint(malID)) {
							url = path
							return
						}
					}
				}

				return
			}

			if strings.Contains(s.Text(), "النوع") {
				if strings.Contains(strings.ToLower(s.Text()), "movie") {
					movie = true
					return
				}
			}

			if year != 0 {
				if strings.Contains(s.Text(), "بداية العرض") && strings.Contains(s.Text(), fmt.Sprint(year)) {
					if isMovie == movie {
						url = path
						return
					}
				}
			}

		})
	}

	if url == "" {
		return nil, errors.New("no data found")
	}

	resp3, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to : %s", url)
	}
	defer resp3.Body.Close()

	if resp3.StatusCode != 200 {
		return nil, fmt.Errorf("status code not 200: %s", url)
	}

	doc3, err := goquery.NewDocumentFromReader(resp3.Body)
	if err != nil {
		return nil, fmt.Errorf("error parse data: %s", url)
	}

	var page string

	doc3.Find(".episodes-card-container").Each(func(i int, s *goquery.Selection) {
		if page != "" {
			return
		}

		card := s.Find(".episodes-card-title")
		if card == nil {
			return
		}

		a := card.Find("a")
		if a == nil {
			return
		}

		href, ok := a.Attr("href")
		if ok {
			if isMovie {
				page = href
				return
			}
			if strings.Contains(strings.ToLower(a.Text()), fmt.Sprint(ep)) {
				page = href
				return
			}
		}
	})

	if page == "" {
		return nil, errors.New("no data found")
	}

	resp4, err := http.Get(page)
	if err != nil {
		return nil, fmt.Errorf("error making request to : %s", page)
	}
	defer resp4.Body.Close()

	if resp4.StatusCode != 200 {
		return nil, fmt.Errorf("status code not 200: %s", page)
	}

	doc4, err := goquery.NewDocumentFromReader(resp4.Body)
	if err != nil {
		return nil, fmt.Errorf("error parse data: %s", page)
	}

	content4 := doc4.Find(".tab-content")
	if content == nil {
		return nil, errors.New("no data found")
	}

	watch := content4.Find("#watch")
	if watch == nil {
		return nil, errors.New("no data found")
	}

	var iframes []models.Iframe
	watch.Find("li").Each(func(i int, s *goquery.Selection) {
		a := s.Find("a")
		if a == nil {
			return
		}

		url, ok := a.Attr("data-ep-url")
		if ok && url != "" {
			var quality string

			small := a.Find("small")
			if small != nil {
				quality = small.Text()
			}

			iframes = append(iframes, models.Iframe{
				Link:    url,
				Type:    "sub",
				Referer: "",
				Quality: quality,
				Language: "ara",
			})
		}
	})

	if len(iframes) == 0 {
		return nil, errors.New("no iframes found")
	}

	return iframes, nil
}
