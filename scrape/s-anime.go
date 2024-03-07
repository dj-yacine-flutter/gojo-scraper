package scrape

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dj-yacine-flutter/gojo-scraper/models"
	"github.com/dj-yacine-flutter/gojo-scraper/utils"
)

type sAnimeSearch []struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Year string `json:"year,omitempty"`
}

type sAnimeItem struct {
	Name string       `json:"name,omitempty"`
	Type string       `json:"type,omitempty"`
	Ep   [][]sAnimeEp `json:"ep,omitempty"`
}

type sAnimeEp struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	EpName any    `json:"epName,omitempty"`
	Date   string `json:"date,omitempty"`
}

type sAnimeVideo struct {
	Sd string `json:"sd,omitempty"`
	Hd string `json:"hd,omitempty"`
}

func (s *Scraper) SAnime(title string, isMovie bool, year, ep int) ([]models.Iframe, error) {
	var err error

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://app.sanime.net/function/h10.php?page=search&name=%s", utils.CleanQuery(title)), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Referer", "https://ios.sanime.net/")
	req.Header.Set("X-Requested-With", "com.sanimenew.apk")
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

	var search sAnimeSearch
	err = json.NewDecoder(resp.Body).Decode(&search)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, v := range search {
		if v.Year == fmt.Sprint(year) {
			if strings.Contains(utils.CleanTitle(v.Name), utils.CleanTitle(title)) {
				ids = append(ids, v.ID)
			}
		}
	}

	if len(ids) == 0 {
		return nil, ErrNoDataFound
	}

	var (
		resp2 *http.Response
		eps   []sAnimeEp
	)
	for _, v := range ids {
		rip = 0
		req1, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://app.sanime.net/function/h10.php?page=info&id=%s", v), nil)
		if err != nil {
			continue
		}

		req1.Header.Set("Referer", fmt.Sprintf("https://ios.sanime.net/info?id=%s", v))
		req1.Header.Set("User-Agent", models.UserAgent)
		req1.Header.Set("X-Requested-With", "com.sanimenew.apk")

		for rip < 10 {
			resp2, err = http.DefaultClient.Do(req1)
			if err != nil {
				continue
			}

			if resp2.StatusCode == 200 {
				break
			}

			resp2.Body.Close()

			rip++
			time.Sleep(750 * time.Millisecond)
		}
		defer resp2.Body.Close()

		var item sAnimeItem
		err = json.NewDecoder(resp2.Body).Decode(&item)
		if err != nil {
			continue
		}

		if isMovie {
			if !strings.Contains(strings.ToLower(item.Type), "movie") {
				continue
			}

			for _, x := range item.Ep {
				eps = append(eps, x...)
			}
		} else {
			if !strings.Contains(item.Type, "مسلسل") {
				continue
			}
			for _, x := range item.Ep {
				for _, z := range x {
					if !strings.Contains(z.Name, "الحلقة الخاصة") {
						switch z.EpName.(type) {
						case float64:
							if z.EpName == float64(ep) {
								eps = append(eps, z)
							}
						}
					}
				}
			}
		}
	}

	if len(eps) == 0 {
		return nil, ErrNoDataFound
	}

	var (
		resp3   *http.Response
		iframes []models.Iframe
	)
	for _, v := range eps {
		djs, err := json.Marshal(&v)
		if err != nil {
			continue
		}

		dbs := base64.StdEncoding.EncodeToString(djs)

		rip = 0
		req2, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://app.sanime.net/function/h10.php?page=openAnd&id=%s", dbs), nil)
		if err != nil {
			continue
		}

		req2.Header.Set("Referer", fmt.Sprintf("https://ios.sanime.net/info?id=%s", v))
		req2.Header.Set("User-Agent", models.UserAgent)
		req2.Header.Set("X-Requested-With", "com.sanimenew.apk")

		for rip < 10 {
			resp3, err = http.DefaultClient.Do(req2)
			if err != nil {
				continue
			}

			if resp3.StatusCode == 200 {
				break
			}

			resp3.Body.Close()

			rip++
			time.Sleep(750 * time.Millisecond)
		}
		defer resp3.Body.Close()

		var video sAnimeVideo

		err = json.NewDecoder(resp3.Body).Decode(&video)
		if err != nil {
			continue
		}

		if strings.Contains(video.Sd, "mp4") {
			iframes = append(iframes, models.Iframe{
				Link:     video.Sd,
				Type:     "sub",
				Quality:  "sd",
				Language: "ara",
			})
		}

		if strings.Contains(video.Hd, "mp4") {
			iframes = append(iframes, models.Iframe{
				Link:     video.Sd,
				Type:     "sub",
				Quality:  "hd",
				Language: "ara",
			})
		}
	}

	if len(iframes) == 0 {
		return nil, ErrNoDataFound
	}

	return iframes, nil
}
