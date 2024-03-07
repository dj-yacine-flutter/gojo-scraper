package scrape

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/dj-yacine-flutter/gojo-scraper/models"
	"github.com/dj-yacine-flutter/gojo-scraper/utils"
)

type animeSlayerSearch struct {
	Response struct {
		Data []struct {
			AnimeID          string `json:"anime_id,omitempty"`
			AnimeName        string `json:"anime_name,omitempty"`
			AnimeType        string `json:"anime_type,omitempty"`
			AnimeReleaseYear string `json:"anime_release_year,omitempty"`
		} `json:"data,omitempty"`
	} `json:"response,omitempty"`
}

type animeSlayerEp struct {
	Response struct {
		Data []struct {
			EpisodeID     string `json:"episode_id,omitempty"`
			EpisodeNumber string `json:"episode_number,omitempty"`
		} `json:"data,omitempty"`
	} `json:"response,omitempty"`
}

type animeSlayerData struct {
	Response struct {
		Data []struct {
			EpisodeUrls []struct {
				EpisodeURL string `json:"episode_url,omitempty"`
			} `json:"episode_urls,omitempty"`
		} `json:"data,omitempty"`
	} `json:"response,omitempty"`
}

func (s *Scraper) AnimeSlayer(title string, isMovie bool, year, ep int) ([]models.Iframe, error) {
	clID := "android-app2"
	clSecret := "7befba6263cc14c90d2f1d6da2c5cf9b251bfbbd"

	req, err := http.NewRequest(http.MethodGet, strings.ReplaceAll("https://anslayer.com/anime/public/animes/get-published-animes?json="+url.PathEscape(fmt.Sprintf(`{"_offset":0,"_limit":100,"_order_by":"latest_first","list_type":"filter","anime_name":"%s\n","just_info":"Yes"}`, utils.CleanQuery(title))), " ", ""), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", models.UserAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en")
	req.Header.Set("Client-Id", clID)
	req.Header.Set("Client-Secret", clSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, ErrNotOK
	}

	var search animeSlayerSearch
	err = json.NewDecoder(resp.Body).Decode(&search)
	if err != nil {
		return nil, err
	}

	if len(search.Response.Data) == 0 {
		return nil, ErrNoDataFound
	}

	var ids []string
	for _, v := range search.Response.Data {
		if isMovie {
			if !strings.Contains(strings.ToLower(v.AnimeType), "movie") {
				continue
			}
		} else {
			if !strings.Contains(strings.ToLower(v.AnimeType), "tv") {
				continue
			}
		}

		if !strings.Contains(v.AnimeReleaseYear, fmt.Sprint(year)) {
			continue
		}

		if strings.Contains(utils.CleanTitle(v.AnimeName), utils.CleanTitle(title)) {
			ids = append(ids, v.AnimeID)
		}
	}

	if len(ids) == 0 {
		return nil, ErrNoDataFound
	}

	var eps []struct {
		aid string
		eid string
	}
	for _, v := range ids {
		data := url.Values{
			"inf":  {""},
			"json": {fmt.Sprintf(`{"more_info":"No","anime_id":%s}`, v)},
		}

		req2, err := http.NewRequest(http.MethodPost, "https://anslayer.com/anime/public/episodes/get-episodes-new", strings.NewReader(data.Encode()))
		if err != nil {
			continue
		}

		req2.Header.Set("Accept", "*/*")
		req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req2.Header.Set("Client-Id", clID)
		req2.Header.Set("Client-Secret", clSecret)
		req2.Header.Set("User-Agent", models.UserAgent)

		resp2, err := http.DefaultClient.Do(req2)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		defer resp2.Body.Close()

		if resp2.StatusCode != 200 {
			continue
		}

		var episodes animeSlayerEp
		err = json.NewDecoder(resp2.Body).Decode(&episodes)
		if err != nil {
			continue
		}

		for _, x := range episodes.Response.Data {
			if isMovie {
				eps = append(eps, struct {
					aid string
					eid string
				}{
					aid: v,
					eid: x.EpisodeID,
				})
			} else {
				if x.EpisodeNumber == fmt.Sprint(ep) {
					eps = append(eps, struct {
						aid string
						eid string
					}{
						aid: v,
						eid: x.EpisodeID,
					})
					break
				}
			}
		}
	}

	if len(eps) == 0 {
		return nil, ErrNoDataFound
	}

	var servers []struct {
		url  string
		data url.Values
	}
	for _, v := range eps {
		data := url.Values{
			"inf":  {""},
			"json": {fmt.Sprintf(`{"anime_id":%s,"episode_id":"%s"}`, v.aid, v.eid)},
		}

		req3, err := http.NewRequest(http.MethodPost, "https://anslayer.com/anime/public/episodes/get-episodes-new", strings.NewReader(data.Encode()))
		if err != nil {
			continue
		}

		req3.Header.Set("Accept", "*/*")
		req3.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req3.Header.Set("Client-Id", clID)
		req3.Header.Set("Client-Secret", clSecret)
		req3.Header.Set("User-Agent", models.UserAgent)

		resp3, err := http.DefaultClient.Do(req3)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		defer resp3.Body.Close()

		if resp3.StatusCode != 200 {
			continue
		}

		var item animeSlayerData
		err = json.NewDecoder(resp3.Body).Decode(&item)
		if err != nil {
			continue
		}

		for _, x := range item.Response.Data {
			for _, y := range x.EpisodeUrls {
				if strings.Contains(y.EpisodeURL, "v-qs.php") {
					vls, err := url.ParseQuery(y.EpisodeURL)
					if err != nil {
						continue
					}

					vls.Add("inf", "")

					servers = append(servers, struct {
						url  string
						data url.Values
					}{
						url:  "https://anslayer.com/anime/public/v-qs.php",
						data: vls,
					})

				}
				if strings.Contains(y.EpisodeURL, "api/f2") {
					vls, err := url.ParseQuery(y.EpisodeURL)
					if err != nil {
						continue
					}

					vls.Add("inf", "")

					servers = append(servers, struct {
						url  string
						data url.Values
					}{
						url:  "https://anslayer.com/la/public/api/fw",
						data: vls,
					})
				}
			}
		}
	}

	if len(servers) == 0 {
		return nil, ErrNoDataFound
	}

	var iframes []models.Iframe
	re := regexp.MustCompile(`\?(.*)`)
	for _, v := range servers {
		dvls := url.Values{}
		for k, s := range v.data {
			submatches := re.FindStringSubmatch(k)
			x := ""
			if len(s) > 0 {
				x = s[len(s)-1]
			}

			if len(submatches) > 1 {
				dvls.Add(submatches[1], x)
				continue
			}
			dvls.Add(k, x)
		}

		req4, err := http.NewRequest(http.MethodPost, v.url, strings.NewReader(dvls.Encode()))
		if err != nil {
			continue
		}

		req4.Header.Set("Accept", "*/*")
		req4.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req4.Header.Set("User-Agent", models.UserAgent)

		resp4, err := http.DefaultClient.Do(req4)
		if err != nil {
			continue
		}
		defer resp4.Body.Close()

		if resp4.StatusCode != 200 {
			continue
		}

		var result []string
		err = json.NewDecoder(resp4.Body).Decode(&result)
		if err != nil {
			continue
		}

		for _, z := range result {
			iframes = append(iframes, models.Iframe{
				Link:     z,
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
