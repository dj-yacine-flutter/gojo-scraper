package scrape

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dj-yacine-flutter/gojo-scraper/models"
	"github.com/dj-yacine-flutter/gojo-scraper/utils"
)

type jkAnimeSearch struct {
	Animes []struct {
		ID    string `json:"id,omitempty"`
		Slug  string `json:"slug,omitempty"`
		Title string `json:"title,omitempty"`
		Type  string `json:"type,omitempty"`
	} `json:"animes,omitempty"`
}

type jkAnimeItem []struct {
	Number string `json:"number,omitempty"`
	Title  string `json:"title,omitempty"`
}

type jkAnimeFrame []struct {
	Remote string `json:"remote,omitempty"`
}

func (s *Scraper) JKAnime(title string, isMovie bool, year, ep int) ([]models.Iframe, error) {

	var err error

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://jkanime.net/ajax/ajax_search/?q=%s", title), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Referer", "https://jkanime.net/")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", models.UserAgent)

	var (
		resp *http.Response
		rip  uint8
	)

	for rip < 5 {
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 {
			break
		}

		resp.Body.Close()

		rip++
		time.Sleep(750 * time.Millisecond)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, ErrNotOK
	}

	var search jkAnimeSearch
	err = json.NewDecoder(resp.Body).Decode(&search)
	if err != nil {
		return nil, errors.New("failed to parse search results")
	}

	if len(search.Animes) == 0 {
		return nil, ErrNoDataFound
	}

	var pages []struct {
		id   string
		slug string
	}
	for _, v := range search.Animes {
		if isMovie {
			if !strings.Contains(strings.ToLower(v.Type), "movie") {
				continue
			}

		} else {
			if !strings.Contains(strings.ToLower(v.Type), "tv") {
				continue
			}
		}

		if strings.Contains(utils.CleanTitle(v.Title), utils.CleanTitle(title)) {
			pages = append(pages, struct {
				id   string
				slug string
			}{
				slug: v.Slug,
				id:   v.ID,
			})
		}
	}

	if len(pages) == 0 {
		return nil, ErrNoDataFound
	}

	var links []string
	rip = 0
	for _, v := range pages {
		for rip < 5 {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://jkanime.net/%s", v.slug), nil)
			if err != nil {
				continue
			}

			req.Header.Set("Referer", "https://jkanime.net/")
			req.Header.Set("User-Agent", models.UserAgent)

			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				continue
			}

			if resp.StatusCode == 200 {
				break
			}

			resp.Body.Close()
			rip++
			time.Sleep(750 * time.Millisecond)
		}

		if resp.StatusCode != 200 {
			resp.Body.Close()
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			continue
		}
		resp.Body.Close()

		widget := doc.Find("div.anime__details__widget")
		if widget == nil {
			continue
		}

		var pass bool
		widget.Find("li").Each(func(i int, s *goquery.Selection) {
			if pass {
				return
			}
			p := s.Find("span")
			if p == nil {
				return
			}

			if strings.Contains(strings.ToLower(p.Text()), "emitido") {
				if strings.Contains(s.Text(), fmt.Sprint(year)) {
					pass = true
					return
				}
			}
		})

		if !pass {
			continue
		}

		var page int
		if isMovie {
			page = 1
		} else {
			page = (ep / 12) + 1
		}

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://jkanime.net/ajax/pagination_episodes/%s/%d", v.id, page), nil)
		if err != nil {
			continue
		}

		req.Header.Set("Referer", "https://jkanime.net/")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
		req.Header.Set("User-Agent", models.UserAgent)

		resp2, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}
		defer resp2.Body.Close()

		var items jkAnimeItem

		err = json.NewDecoder(resp2.Body).Decode(&items)
		if err != nil {
			continue
		}

		num := "1"
		if !isMovie {
			num = ""
			for _, i := range items {
				if strings.TrimSpace(i.Number) == strings.TrimSpace(fmt.Sprint(ep)) {
					num = i.Number
					break
				}
			}
		}

		links = append(links, fmt.Sprintf("https://jkanime.net/%s/%s", v.slug, num))
	}

	if len(links) == 0 {
		return nil, ErrNoDataFound
	}

	var ids []string
	for _, v := range links {
		req, err := http.NewRequest(http.MethodPost, v, nil)
		if err != nil {
			continue
		}

		req.Header.Set("Referer", "https://jkanime.net/")
		req.Header.Set("User-Agent", models.UserAgent)

		resp3, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}
		defer resp3.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp3.Body)
		if err != nil {
			continue
		}

		div := doc.Find("#guardar-capitulo")
		if div == nil {
			continue
		}

		id, ok := div.Attr("data-capitulo")
		if ok && id != "" {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		return nil, ErrNoDataFound
	}

	var iframes []models.Iframe
	for _, v := range ids {
		req, err := http.NewRequest(http.MethodGet, strings.TrimSpace(fmt.Sprintf("https://c4.jkdesu.com/servers/%s.js", v)), nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", models.UserAgent)

		resp4, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}
		defer resp4.Body.Close()

		b, err := io.ReadAll(resp4.Body)
		if err != nil {
			continue
		}

		var frms jkAnimeFrame
		err = json.NewDecoder(strings.NewReader(strings.ReplaceAll(string(b), "var servers = ", ""))).Decode(&frms)
		if err != nil {
			continue
		}

		for _, z := range frms {
			code, err := base64.StdEncoding.DecodeString(z.Remote)
			if err != nil || len(code) == 0 {
				continue
			}

			iframes = append(iframes, models.Iframe{
				Link:     string(code),
				Type:     "sub",
				Quality:  "fhd",
				Language: "es",
			})
		}
	}

	if len(iframes) == 0 {
		return nil, ErrNoDataFound
	}

	return iframes, nil
}
