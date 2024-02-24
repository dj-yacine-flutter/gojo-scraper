package scrape

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dj-yacine-flutter/gojo-scraper/models"
)

func (s *Scraper) AnimeRco(title string, isMovie bool, malID, year, ep int) ([]models.Iframe, error) {
	var err error

	uri, err := url.Parse("https://animerco.org/?s=" + title)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: "GET",
		URL:    uri,
		Header: make(http.Header),
	}

	req.Header.Add("authority", "animerco.org")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
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
	if doc == nil {
		return nil, errors.New("no doc found")
	}

	row := doc.Find(".row.gutter-small")
	if row == nil {
		return nil, errors.New("no row found")
	}

	var data []*goquery.Selection
	if isMovie {
		row.Find(".media-block.movies").Each(func(i int, s *goquery.Selection) {
			data = append(data, s)
		})
	} else {
		row.Find(".media-block.seasons").Each(func(i int, s *goquery.Selection) {
			data = append(data, s)
		})
	}

	var links []string
	for _, v := range data {
		a := v.Find("a")
		if a == nil {
			continue
		}

		href, ok := a.Attr("href")
		if ok {
			resp1, err := http.Get(href)
			if err != nil {
				continue
			}
			defer resp1.Body.Close()

			doc2, err := goquery.NewDocumentFromReader(resp1.Body)
			if err != nil {
				continue
			}

			div := doc2.Find("div.widget-sidebar")
			if div == nil {
				continue
			}

			var link string
			div.Find("a").Each(func(i int, s *goquery.Selection) {
				if link != "" {
					return
				}

				href2, ok := s.Attr("href")
				if ok {
					if strings.Contains(strings.ToLower(href2), "myanimelist") {
						if strings.Contains(href2, fmt.Sprint(malID)) {
							link = href
							return
						}
					}
				}
			})

			if link == "" {
				ul := div.Find("ul.media-info")
				if ul != nil {
					if strings.Contains(ul.Text(), fmt.Sprint(year)) {
						link = href
					}
				}
			}

			if link != "" {
				if isMovie {
					links = append(links, link)
					continue
				} else {
					ul := doc2.Find("ul.chapters-list")
					if ul != nil {
						var found bool
						ul.Find("li").Each(func(i int, s *goquery.Selection) {
							if found {
								return
							}
							if strings.Contains(s.Text(), fmt.Sprintf("الحلقة %d:", ep)) {
								a := s.Find("a")
								if a == nil {
									return
								}
								href3, ok := a.Attr("href")
								if ok {
									links = append(links, href3)
									found = true
									return
								}
							}
						})
					}
				}
			}
		}
	}

	if len(links) == 0 {
		return nil, errors.New("no links found")
	}

	var iframes []models.Iframe
	for _, v := range links {
		resp3, err := http.Get(v)
		if err != nil {
			continue
		}
		defer resp3.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp3.Body)
		if err != nil {
			continue
		}

		ul := doc.Find("ul.server-list")
		if ul == nil {
			continue
		}

		ul.Find("li").Each(func(i int, s *goquery.Selection) {
			typ, o1 := s.Attr("data-type")
			pst, o2 := s.Attr("data-post")
			num, o3 := s.Attr("data-nume")
			if o1 && o2 && o3 {
				bd := []byte(fmt.Sprintf("action=doo_player_ajax&post=%s&nume=%s&type=%s", pst, num, typ))
				req4, err := http.NewRequest("POST", "https://ww3.animerco.org/wp-admin/admin-ajax.php", bytes.NewBuffer(bd))
				if err != nil {
					return
				}

				req4.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
				req4.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
				req4.Header.Set("Accept", "*/*")
				req4.Header.Set("Accept-Language", "en")
				req4.Header.Set("X-Requested-With", "XMLHttpRequest")

				resp4, err := http.DefaultClient.Do(req4)
				if err != nil {
					return
				}
				defer resp4.Body.Close()

				if resp4.StatusCode != 200 {
					return
				}

				b, err := io.ReadAll(resp4.Body)
				if err != nil {
					return
				}

				var data map[string]interface{}

				err = json.Unmarshal(b, &data)
				if err != nil {
					return
				}

				embed, ok := data["embed_url"].(string)
				if !ok {
					return
				}

				if strings.Contains(strings.ToLower(embed), "<iframe") {
					ifr, err := goquery.NewDocumentFromReader(bytes.NewBuffer([]byte(embed)))
					if err != nil {
						return
					}

					embed = ""

					frm := ifr.Find("iframe")
					if frm != nil {
						src, ok := frm.Attr("src")
						if ok && src != "" {
							embed = strings.ReplaceAll(src, "/\\", "/")
							embed = strings.ReplaceAll(embed, "https:", "")
							embed = strings.ReplaceAll(embed, "http:", "")
							embed = strings.Replace(embed, "//", "https://", 1)
						}
					}
				}

				if embed != "" {
					iframes = append(iframes, models.Iframe{
						Link:    strings.ReplaceAll(embed, "/\\", "/"),
						Type:    "sub",
						Quality: "hd",
						Language: "ara",
					})
				}

			}
		})
	}

	if len(iframes) == 0 {
		return nil, errors.New("no iframes found")
	}

	return iframes, nil
}
