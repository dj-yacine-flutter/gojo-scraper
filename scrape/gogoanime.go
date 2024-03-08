package scrape

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dj-yacine-flutter/gojo-scraper/models"
	"github.com/dj-yacine-flutter/gojo-scraper/utils"
)

func (s *Scraper) GogoAnime(title string, isMovie bool, year, episode int) ([]models.Iframe, error) {
	var err error
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://ajax.gogocdn.net/site/loadAjaxSearch?keyword=%s&id=-1&link_web=https://gogoanime3.co/", utils.CleanQuery(title)), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, ErrNotOK
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var queries []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(utils.CleanTitle(s.Text()), utils.CleanTitle(title)) {
			if href, ok := s.Attr("href"); ok {
				href = strings.ReplaceAll(href, "\\/", "/")
				href = strings.ReplaceAll(href, "\\\"", "")
				queries = append(queries, href)
			}
		}
	})

	if len(queries) == 0 {
		return nil, ErrNoDataFound
	}

	var links []struct {
		path string
		id   string
	}
	for _, v := range queries {
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

		doc2, err := goquery.NewDocumentFromReader(resp2.Body)
		if err != nil {
			continue
		}

		sec := doc2.Find(".anime_info_body_bg")
		if sec == nil {
			continue
		}

		var (
			ch1  bool
			ch2  bool
			pass bool
		)

		sec.Find("p").Each(func(i int, p *goquery.Selection) {
			if pass {
				return
			}

			s := p.Find("span")
			if s == nil {
				return
			}

			txt := strings.ToLower(s.Text())
			if strings.Contains(txt, "type") {
				var t string

				if isMovie {
					t = "movie"
				} else {
					t = "anime"
				}

				if !strings.Contains(strings.ToLower(p.Text()), t) {
					return
				}
				ch1 = true
			}

			if strings.Contains(txt, "released") {
				if !strings.Contains(strings.ToLower(p.Text()), fmt.Sprint(year)) {
					return
				}
				ch2 = true
			}
			if ch1 && ch2 {
				pass = true
			}
		})

		if !pass {
			continue
		}

		var (
			id   string
			path string
		)

		eps := doc2.Find(".anime_info_episodes")
		if eps == nil {
			continue
		}

		eps.Find("input").Each(func(i int, s *goquery.Selection) {
			d, ok := s.Attr("id")
			if ok {
				if strings.Contains(d, "id") {
					if v, k := s.Attr("value"); k {
						id = v
					}
				} else if strings.Contains(d, "alias") {
					if v, k := s.Attr("value"); k {
						path = v
					}
				}
			}
		})

		links = append(links, struct {
			path string
			id   string
		}{
			path: path,
			id:   id,
		})
	}

	if len(links) == 0 {
		return nil, ErrNoDataFound
	}

	var pages []struct {
		path string
		dub  bool
	}
	for _, v := range links {
		req3, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://ajax.gogocdn.net/ajax/load-list-episode?ep_start=0&ep_end=9999&id=%s&default_ep=0&alias=%s", v.id, v.path), nil)
		if err != nil {
			continue
		}

		req3.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
		req3.Header.Set("Accept", "*/*")
		req3.Header.Set("Accept-Language", "en")

		resp3, err := http.DefaultClient.Do(req3)
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

		doc3.Find("li").Each(func(i int, s *goquery.Selection) {
			cate := s.Find(".cate")
			if cate == nil {
				return
			}

			a := s.Find("a")
			if a == nil {
				return
			}

			href, ok := a.Attr("href")
			if !ok {
				return
			}

			if isMovie {
				pages = append(pages, struct {
					path string
					dub  bool
				}{
					path: href,
					dub:  strings.Contains(strings.ToLower(cate.Text()), "dub"),
				})
			} else {
				nm := s.Find(".name")
				if nm == nil {
					return
				}

				ep := strings.ToLower(nm.Text())
				ep = strings.ReplaceAll(ep, "ep", "")
				ep = strings.ReplaceAll(ep, " ", "")

				if ep == fmt.Sprint(episode) {
					pages = append(pages, struct {
						path string
						dub  bool
					}{
						path: href,
						dub:  strings.Contains(strings.ToLower(cate.Text()), "dub"),
					})
				}
			}
		})
	}

	if len(pages) == 0 {
		return nil, ErrNoDataFound
	}

	var iframes []models.Iframe
	for _, v := range pages {
		resp4, err := http.Get(fmt.Sprintf("https://gogoanime3.co%s", strings.TrimSpace(v.path)))
		if err != nil {
			continue
		}
		defer resp4.Body.Close()

		if resp4.StatusCode != 200 {
			continue
		}

		doc4, err := goquery.NewDocumentFromReader(resp4.Body)
		if err != nil {
			continue
		}

		sec := doc4.Find(".anime_video_body")
		if sec == nil {
			continue
		}

		sec.Find("li").Each(func(i int, s *goquery.Selection) {
			a := s.Find("a")
			if a == nil {
				return
			}

			vid, ok := a.Attr("data-video")
			if ok {
				if strings.Contains(vid, "embtaku") {
					return
				}

				t := "sub"
				if v.dub {
					t = "dub"
				}

				iframes = append(iframes, models.Iframe{
					Link:     vid,
					Type:     t,
					Quality:  "fhd",
					Language: "en",
				})
			}
		})
	}

	if len(iframes) == 0 {
		return nil, ErrNoDataFound
	}

	return iframes, nil
}
