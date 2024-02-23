package scrape

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dj-yacine-flutter/gojo-scraper/models"
	"github.com/dj-yacine-flutter/gojo-scraper/utils"
)

type arabAnimeSearch struct {
	SearchResaults []string `json:"SearchResaults,omitempty"`
}

type arabAnimeSearchItem struct {
	AnimeName        string `json:"anime_name,omitempty"`
	AnimeReleaseDate string `json:"anime_release_date,omitempty"`
	AnimeType        string `json:"anime_type,omitempty"`
	InfoURL          string `json:"info_url,omitempty"`
}

type arabAnimeEp struct {
	Eps []struct {
		EpisodeName   string `json:"episode_name,omitempty"`
		EpisodeNumber int    `json:"episode_number,omitempty"`
		InfoSrc       string `json:"info-src,omitempty"`
	} `json:"EPS,omitempty"`
}

type arabAnimeStream struct {
	EpInfo []struct {
		StreamServers []string `json:"stream_servers,omitempty"`
	} `json:"ep_info,omitempty"`
}

func (s *Scraper) ArabAnime(title string, isMovie bool, year, ep int) ([]models.Iframe, error) {
	var err error

	resp, err := http.Get(fmt.Sprintf("https://www.arabanime.net/api/search?q=%s", title))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("status code not 200")
	}

	var search arabAnimeSearch

	err = json.NewDecoder(resp.Body).Decode(&search)
	if err != nil {
		return nil, errors.New("cannot parse search results")
	}

	var results []string
	for _, v := range search.SearchResaults {
		x, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			continue
		}

		var item arabAnimeSearchItem

		err = json.Unmarshal(x, &item)
		if err != nil {
			continue
		}

		if isMovie {
			if !strings.Contains(strings.ToLower(item.AnimeType), "movie") {
				continue
			}
		} else {
			if !strings.Contains(strings.ToLower(item.AnimeType), "serie") {
				continue
			}
		}

		if strings.Contains(item.AnimeReleaseDate, fmt.Sprint(year)) {
			if strings.Contains(utils.CleanTitle(item.AnimeName), utils.CleanTitle(title)) {
				results = append(results, item.InfoURL)
			}
		}
	}

	var links []string
	for _, v := range results {
		resp2, err := http.Get(v)
		if err != nil {
			continue
		}
		defer resp2.Body.Close()

		if resp2.StatusCode != 200 {
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp2.Body)
		if err != nil {
			continue
		}

		code := doc.Find("div#data")
		if code == nil {
			continue
		}

		txt := code.Text()
		if txt == "" {
			continue
		}

		dec, err := base64.StdEncoding.DecodeString(txt)
		if err != nil {
			continue
		}

		var data arabAnimeEp
		err = json.Unmarshal(dec, &data)
		if err != nil {
			continue
		}

		if len(data.Eps) == 0 {
			continue
		}

		if isMovie {
			links = append(links, data.Eps[0].InfoSrc)
		} else {
			for _, z := range data.Eps {
				if z.EpisodeNumber == ep {
					links = append(links, z.InfoSrc)
				}
			}
		}
	}

	if len(links) == 0 {
		return nil, errors.New("no data found")
	}

	var embeds []string
	for _, v := range links {
		resp3, err := http.Get(v)
		if err != nil {
			continue
		}
		defer resp3.Body.Close()

		if resp3.StatusCode != 200 {
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp3.Body)
		if err != nil {
			continue
		}

		code := doc.Find("div#datawatch")
		if code == nil {
			continue
		}

		txt := code.Text()
		if txt == "" {
			continue
		}

		dec, err := base64.StdEncoding.DecodeString(txt)
		if err != nil {
			continue
		}

		var data arabAnimeStream
		err = json.Unmarshal(dec, &data)
		if err != nil {
			continue
		}

		for _, i := range data.EpInfo {
			for _, j := range i.StreamServers {
				embed, err := base64.StdEncoding.DecodeString(j)
				if err != nil {
					continue
				}
				embeds = append(embeds, string(embed))
			}
		}
	}

	if len(embeds) == 0 {
		return nil, errors.New("no data found")
	}

	var iframes []models.Iframe
	for _, v := range embeds {
		resp4, err := http.Get(v)
		if err != nil {
			continue
		}
		defer resp4.Body.Close()

		if resp4.StatusCode != 200 {
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp4.Body)
		if err != nil {
			continue
		}

		doc.Find("option").Each(func(i int, s *goquery.Selection) {
			src, ok := s.Attr("data-src")
			if ok {
				frm, err := base64.StdEncoding.DecodeString(src)
				if err != nil {
					return
				}

				if len(frm) != 0 {
					iframes = append(iframes, models.Iframe{
						Link:    string(frm),
						Type:    "sub",
						Quality: "hd",
					})
				}
			}
		})
	}

	if len(iframes) == 0 {
		return nil, errors.New("no data found")
	}

	return iframes, nil
}
