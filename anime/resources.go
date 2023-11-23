package anime

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dj-yacine-flutter/gojo-scraper/utils"
)

var (
	UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"
)

func (server *AnimeScraper) Livechart(id int, title string, date time.Time) int {
	url := fmt.Sprintf("https://www.livechart.me/anime/%d", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0
	}
	req.Header.Set("Authority", "www.livechart.me")
	req.Header.Set("Referer", url)
	req.Header.Set("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			id = 0
		}

		block := doc.Find(".lc-poster-col")
		if block != nil {
			var found bool
			block.Find(".text-sm").Each(func(index int, selection *goquery.Selection) {
				if strings.Contains(strings.ToLower(selection.Text()), strings.ToLower(title)) {
					found = true
					return
				}
			})
			if found {
				return id
			}
		}
		id = 0
	} else {
		id = 0
	}
	if id == 0 {
		query := fmt.Sprintf("https://www.livechart.me/search?q=%s", utils.CleanQuery(title))
		req, err := http.NewRequest("GET", query, nil)
		if err != nil {
			return 0
		}

		req.Header.Set("User-Agent", UserAgent)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return 0
		}

		defer resp.Body.Close()
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		block := doc.Find(".anime-list")
		if block != nil {
			var found bool
			var ID int
			block.Find(".info").Each(func(index int, selection *goquery.Selection) {
				if selection != nil {
					span := selection.Find("span").First()
					if span != nil {
						layout := "January 2, 2006"
						parsedTime, err := time.Parse(layout, span.Text())
						fmt.Println(parsedTime.String())
						if err == nil {
							if parsedTime.Year() == date.Year() && parsedTime.Month() == date.Month() {
								title := block.Find("strong").First()
								fmt.Println(parsedTime.String())
								fmt.Println(title.Html())
								if title != nil {
									a := title.Find("a").First()
									if a != nil {
										link, ok := a.Attr("href")
										if ok {
											ID = utils.ExtractID(link)
											if ID != 0 {
												found = true
												return
											}
										}
									}
								}
							}
						}
					}
				}
			})
			if found {
				return ID
			}
		}
	}
	return 0
}

