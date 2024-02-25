package scrape

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dj-yacine-flutter/gojo-scraper/models"
)

type animeSaturnSearch []struct {
	Name    string `json:"name,omitempty"`
	Link    string `json:"link,omitempty"`
	Release string `json:"release,omitempty"`
}

func (s *Scraper) AnimeSaturn(title string, isMovie bool, malID, year, ep int) ([]models.Iframe, error) {
	var err error

	/* 	req, err := http.NewRequest(http.MethodGet, "https://www.animesaturn.tv", nil)
	   	if err != nil {
	   		return nil, errors.New("cannot make request")
	   	}

	   	req.Header.Add("Connection", "keep-alive")
	   	req.Header.Add("Referer", "https://www.animesaturn.tv/")
	   	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")

	   	resp, err := http.DefaultClient.Do(req)
	   	if err != nil {
	   		return nil, errors.New("cannot GET response")
	   	}
	   	defer resp.Body.Close()


	   	b1, _ := io.ReadAll(resp.Body)

	   	fmt.Println(string(b1))
	   	fmt.Println(resp.StatusCode)
	   	fmt.Println(resp.Status)
	   	//if resp.StatusCode != 200 {
	   	//	return nil, errors.New("1- StatusCode != 200 ")
	   	//}

	   	var (
	   		as1 string
	   		as2 string
	   		php string
	   	)

	   	for _, v := range resp.Cookies() {
	   		fmt.Println(v.Raw)
	   		if strings.Contains(strings.ToLower(v.Name), "astest-es") {
	   			as1 = v.Value
	   			continue
	   		}
	   		if strings.Contains(strings.ToLower(v.Name), "astest-2v") {
	   			as2 = v.Value
	   			continue
	   		}
	   		if strings.Contains(strings.ToUpper(v.Name), "PHPSESSID") {
	   			php = v.Value
	   			continue
	   		}
	   	} */

	req2, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://www.animesaturn.tv/index.php?search=1&key=%s&d=1", title), nil)
	if err != nil {
		return nil, errors.New("cannot make request")
	}

	req2.Header.Add("Accept", "*/*")
	req2.Header.Add("Connection", "keep-alive")
	req2.Header.Add("Referer", "https://www.animesaturn.tv/")
	req2.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req2.Header.Add("X-Requested-With", "XMLHttpRequest")

	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		return nil, errors.New("cannot GET response")
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != 200 {
		return nil, errors.New("2- StatusCode != 200 ")
	}

	var result animeSaturnSearch
	err = json.NewDecoder(resp2.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
		return nil, errors.New("cannot parse search result")
	}

	var pages []string
	for _, v := range result {
		if !strings.Contains(v.Release, fmt.Sprint(year)) {
			continue
		}

		pages = append(pages, fmt.Sprintf("https://www.animesaturn.tv/anime/%s", v.Link))
	}

	if len(pages) == 0 {
		return nil, errors.New("no results found")
	}

	re := regexp.MustCompile(`\d+`)
	var links []string
	for _, v := range pages {
		req2, err := http.NewRequest(http.MethodGet, v, nil)
		if err != nil {
			continue
		}

		resp2, err := http.DefaultClient.Do(req2)
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

		var found bool
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			if found {
				return
			}
			href, ok := s.Attr("href")
			if ok {
				if strings.Contains(strings.ToLower(href), "myanimelist.net/anime/") {
					if strings.Contains(href, fmt.Sprint(malID)) {
						found = true
					}
				}
			}
		})

		if !found {
			continue
		}

		rng := doc.Find("div#range-anime-0")
		if rng == nil {
			continue
		}

		rng.Find(".episodes-button").Each(func(i int, s *goquery.Selection) {
			if !isMovie {
				match := re.FindString(s.Text())

				num, err := strconv.Atoi(match)
				if err != nil || num != ep {
					return
				}
			}

			a := s.Find("a")
			if a == nil {
				return
			}

			href, ok := a.Attr("href")
			if ok {
				if href != "" {
					links = append(links, href)
				}
			}
		})

	}

	if len(links) == 0 {
		return nil, errors.New("no links found")
	}

	var watchers []struct {
		url string
		dub bool
	}
	for _, v := range links {
		req5, err := http.NewRequest(http.MethodGet, v, nil)
		if err != nil {
			continue
		}

		resp5, err := http.DefaultClient.Do(req5)
		if err != nil {
			continue
		}
		defer resp5.Body.Close()

		if resp5.StatusCode != 200 {
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp5.Body)
		if err != nil {
			continue
		}

		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if ok {
				if strings.Contains(href, "watch") {
					watchers = append(watchers, struct {
						url string
						dub bool
					}{
						url: href,
						dub: strings.Contains(strings.ToLower(v), "-ita"),
					})
				}
			}
		})

	}

	var embeds []struct {
		url string
		dub bool
	}
	for _, v := range watchers {
		req3, err := http.NewRequest(http.MethodGet, v.url, nil)
		if err != nil {
			continue
		}

		resp3, err := http.DefaultClient.Do(req3)
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

		div := doc.Find("div.dropdown-menu")
		if div == nil {
			continue
		}

		div.Find("a").Each(func(i int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if ok {
				if href != "" {
					embeds = append(embeds, struct {
						url string
						dub bool
					}{
						url: href,
						dub: v.dub,
					})
				}
			}
		})
	}

	if len(embeds) == 0 {
		return nil, errors.New("no embeds found")
	}

	var iframes []models.Iframe
	for _, v := range embeds {
		req4, err := http.NewRequest(http.MethodGet, v.url, nil)
		if err != nil {
			continue
		}

		resp4, err := http.DefaultClient.Do(req4)
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

		embed := doc.Find(".embed-container")
		if embed == nil {
			continue
		}

		frame := embed.Find("iframe")
		if frame == nil {
			continue
		}

		src, ok := frame.Attr("src")
		if ok {
			if src != "" {
				t := "sub"
				if v.dub {
					t = "dub"
				}
				iframes = append(iframes, models.Iframe{
					Link:     src,
					Referer:  v.url,
					Type:     t,
					Quality:  "fhd",
					Language: "it",
				})
			}

		}
	}

	return iframes, nil
}
