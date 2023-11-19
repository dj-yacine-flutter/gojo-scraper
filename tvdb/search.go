package tvdb

import (
	"fmt"
	"net/http"
)

const (
	SearchPath = "/search"
)

type Search struct {
	Data []struct {
		ObjectID        string   `json:"objectID,omitempty"`
		Aliases         []string `json:"aliases,omitempty"`
		Country         string   `json:"country,omitempty"`
		ExtendedTitle   string   `json:"extended_title,omitempty"`
		Genres          []string `json:"genres,omitempty"`
		Studios         []string `json:"studios,omitempty"`
		ID              string   `json:"id,omitempty"`
		ImageURL        string   `json:"image_url,omitempty"`
		Name            string   `json:"name,omitempty"`
		FirstAirTime    string   `json:"first_air_time,omitempty"`
		Overview        string   `json:"overview,omitempty"`
		PrimaryLanguage string   `json:"primary_language,omitempty"`
		PrimaryType     string   `json:"primary_type,omitempty"`
		Status          string   `json:"status,omitempty"`
		Type            string   `json:"type,omitempty"`
		TvdbID          string   `json:"tvdb_id,omitempty"`
		Year            string   `json:"year,omitempty"`
		Slug            string   `json:"slug,omitempty"`
		Overviews       struct {
			Ara string `json:"ara,omitempty"`
			Deu string `json:"deu,omitempty"`
			Eng string `json:"eng,omitempty"`
			Fra string `json:"fra,omitempty"`
			Ita string `json:"ita,omitempty"`
			Jpn string `json:"jpn,omitempty"`
			Pol string `json:"pol,omitempty"`
			Por string `json:"por,omitempty"`
			Pt  string `json:"pt,omitempty"`
			Rus string `json:"rus,omitempty"`
			Spa string `json:"spa,omitempty"`
		} `json:"overviews,omitempty"`
		Translations struct {
			Ara string `json:"ara,omitempty"`
			Deu string `json:"deu,omitempty"`
			Eng string `json:"eng,omitempty"`
			Fra string `json:"fra,omitempty"`
			Ita string `json:"ita,omitempty"`
			Jpn string `json:"jpn,omitempty"`
			Pol string `json:"pol,omitempty"`
			Por string `json:"por,omitempty"`
			Pt  string `json:"pt,omitempty"`
			Rus string `json:"rus,omitempty"`
			Spa string `json:"spa,omitempty"`
		} `json:"translations,omitempty"`
		RemoteIds []struct {
			ID         string `json:"id,omitempty"`
			Type       int    `json:"type,omitempty"`
			SourceName string `json:"sourceName,omitempty"`
		} `json:"remote_ids,omitempty"`
		Thumbnail string `json:"thumbnail,omitempty"`
	} `json:"data,omitempty"`
	Status string `json:"status,omitempty"`
	Links  struct {
		Prev       string `json:"prev,omitempty"`
		Self       string `json:"self,omitempty"`
		Next       string `json:"next,omitempty"`
		TotalItems int    `json:"total_items,omitempty"`
		PageSize   int    `json:"page_size,omitempty"`
	} `json:"links,omitempty"`
}

func (c *Client) GetSearch(query string, year int) (data *Search, err error) {
	req, err := http.NewRequest("GET", c.BuildUrlPath(SearchPath), nil)
	if err != nil {
		return
	}

	q := req.URL.Query()
	q.Add("query", query)
	if year != 0 {
		q.Add("year", fmt.Sprint(year))
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", ApplicationJson)
	req.Header.Set("Accept", ApplicationJson)
	req.Header.Set("Accept-Language", c.Language)
	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	data = new(Search)
	err = c.ParseResponse(resp.Body, data)
	return
}
