package scrape

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dj-yacine-flutter/gojo-scraper/models"
)

type animeUnitySearch struct {
	Title string `json:"title"`
}

type animeUnityItem struct {
	Records []struct {
		ID       int    `json:"id,omitempty"`
		Date     string `json:"date,omitempty"`
		Type     string `json:"type,omitempty"`
		Slug     string `json:"slug,omitempty"`
		TitleEng string `json:"title_eng,omitempty"`
		MalID    int    `json:"mal_id,omitempty"`
	} `json:"records,omitempty"`
}

type animeUnityFrame struct {
	ID     int    `json:"id,omitempty"`
	Number string `json:"number,omitempty"`
	ScwsID int    `json:"scws_id,omitempty"`
}

func (s *Scraper) AnimeUnity(name string, isMovie bool, malID, year, ep int) ([]models.Iframe, error) {
	var err error

	re, err := http.Get("https://www.animeunity.to")
	if err != nil {
		return nil, err
	}
	defer re.Body.Close()

	if re.StatusCode != 200 {
		return nil, errors.New("StatusCode != 200")
	}

	var (
		xs string
		ss string
	)

	for _, v := range re.Cookies() {
		if strings.Contains(strings.ToUpper(v.Name), "XSRF-TOKEN") {
			xs = v.Value
		}

		if strings.Contains(strings.ToLower(v.Name), "animeunity_session") {
			ss = v.Value
		}
	}

	if xs == "" || ss == "" {
		return nil, errors.New("no cookies found")
	}

	doc, err := goquery.NewDocumentFromReader(re.Body)
	if err != nil {
		return nil, err
	}

	head := doc.Find("head")
	if head == nil {
		return nil, errors.New("no head found")
	}

	var cf string
	head.Find("meta").Each(func(i int, s *goquery.Selection) {
		if cf != "" {
			return
		}

		name, ok := s.Attr("name")
		if ok {
			if strings.Contains(strings.ToLower(name), "csrf-token") {
				cnt, ok1 := s.Attr("content")
				if ok1 {
					cf = cnt
				}
			}
		}
	})

	if cf == "" {
		return nil, errors.New("no token found")
	}

	searchData := animeUnitySearch{
		Title: name,
	}

	searchByte, err := json.Marshal(&searchData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://www.animeunity.to/livesearch", bytes.NewBuffer(searchByte))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	req.Header.Add("Referer", "https://www.animeunity.to/")
	req.Header.Add("Origin", "https://www.animeunity.to")
	req.Header.Add("Cookie", strconv.Quote(fmt.Sprintf("X-XSRF-TOKEN=%s; animeunity_session=%s", xs, ss)))
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Add("X-CSRF-TOKEN", cf)
	req.Header.Add("X-XSRF-TOKEN", xs)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("StatusCode != 200")
	}

	var result animeUnityItem

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, errors.New("failed to parse search data")
	}

	var data []string
	for _, v := range result.Records {
		if isMovie {
			if !strings.Contains(strings.ToLower(v.Type), "movie") {
				continue
			}
		} else {
			if !strings.Contains(strings.ToLower(v.Type), "tv") {
				continue
			}
		}

		if malID != 0 {
			if v.MalID != malID {
				continue
			}
		}

		if !strings.Contains(v.Date, fmt.Sprint(year)) {
			continue
		}

		data = append(data, fmt.Sprintf("https://www.animeunity.to/anime/%d-%s", v.ID, v.Slug))
	}

	var codes []struct {
		frame animeUnityFrame
		ref   string
		dub   bool
	}

	for _, v := range data {
		req2, err := http.NewRequest(http.MethodGet, v, nil)
		if err != nil {
			continue
		}

		req2.Header.Add("Accept", "*/*")
		req2.Header.Add("Accept-Language", "en-US,en;q=0.9")
		req2.Header.Add("Referer", "https://www.animeunity.to/")
		req2.Header.Add("Origin", "https://www.animeunity.to")
		req2.Header.Add("Cookie", strconv.Quote(fmt.Sprintf("X-XSRF-TOKEN=%s; animeunity_session=%s", xs, ss)))
		req2.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
		req2.Header.Add("X-CSRF-TOKEN", cf)
		req2.Header.Add("X-XSRF-TOKEN", xs)
		req2.Header.Add("X-Requested-With", "XMLHttpRequest")

		resp2, err := http.DefaultClient.Do(req2)
		if err != nil {
			continue
		}
		defer resp2.Body.Close()

		if resp2.StatusCode != 200 {
			continue
		}

		for _, v := range re.Cookies() {
			if strings.Contains(strings.ToUpper(v.Name), "XSRF-TOKEN") {
				if v.Value != "" {
					xs = v.Value
				}
			}

			if strings.Contains(strings.ToLower(v.Name), "animeunity_session") {
				if v.Value != "" {
					ss = v.Value
				}
			}
		}

		doc2, err := goquery.NewDocumentFromReader(resp2.Body)
		if err != nil {
			continue
		}

		player := doc2.Find("video-player")
		if player == nil {
			continue
		}

		eps, ok := player.Attr("episodes")
		if ok {
			//fmt.Println(eps)
			if eps != "" {
				eps = strings.ReplaceAll(eps, "&quot;", "\"")
				var code []animeUnityFrame

				err = json.Unmarshal(bytes.NewBufferString(eps).Bytes(), &code)
				if err != nil {
					continue
				}
				for _, c := range code {
					if !isMovie {
						if c.Number != fmt.Sprint(ep) {
							continue
						}
					}

					codes = append(codes, struct {
						frame animeUnityFrame
						ref   string
						dub   bool
					}{
						frame: c,
						ref:   fmt.Sprintf("%s/%d", v, c.ID),
						dub:   strings.Contains(v, "-ita"),
					})
				}
			}
		}
	}

	var iframes []models.Iframe
	for _, v := range codes {
		req3, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.animeunity.to/embed-url/%d", v.frame.ID), nil)
		if err != nil {
			continue
		}

		req3.Header.Add("Accept", "*/*")
		req3.Header.Add("Accept-Language", "en-US,en;q=0.9")
		req3.Header.Add("Referer", v.ref)
		req3.Header.Add("Origin", "https://www.animeunity.to")
		req3.Header.Add("Cookie", strconv.Quote(fmt.Sprintf("X-XSRF-TOKEN=%s; animeunity_session=%s", xs, ss)))
		req3.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
		req3.Header.Add("X-CSRF-TOKEN", cf)
		req3.Header.Add("X-XSRF-TOKEN", xs)
		req3.Header.Add("X-Requested-With", "XMLHttpRequest")

		resp3, err := http.DefaultClient.Do(req3)
		if err != nil {
			continue
		}
		defer resp3.Body.Close()

		if resp3.StatusCode != 200 {
			continue
		}

		b, err := io.ReadAll(resp3.Body)
		if err != nil {
			continue
		}

		if strings.Contains(string(b), fmt.Sprint(v.frame.ScwsID)) {
			var t string
			if v.dub {
				t = "dub"
			} else {
				t = "sub"
			}
			iframes = append(iframes, models.Iframe{
				Link:     strings.TrimSpace(string(b)),
				Referer:  "https://www.animeunity.to/",
				Type:     t,
				Quality:  "fhd",
				Language: "it",
			})
		}
	}

	return iframes, nil
}
