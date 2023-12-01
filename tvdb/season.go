package tvdb

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	SeasonsByIDPath         = "/seasons/:id"
	SeasonsByIDExtendedPath = "/seasons/:id/extended"
)

type SeasonsByIDExtended struct {
	Status string `json:"status,omitempty"`
	Data   struct {
		ID       int `json:"id,omitempty"`
		SeriesID int `json:"seriesId,omitempty"`
		Type     struct {
			ID            int    `json:"id,omitempty"`
			Name          string `json:"name,omitempty"`
			Type          string `json:"type,omitempty"`
			AlternateName any    `json:"alternateName,omitempty"`
		} `json:"type,omitempty"`
		Number               int      `json:"number,omitempty"`
		NameTranslations     []any    `json:"nameTranslations,omitempty"`
		OverviewTranslations []string `json:"overviewTranslations,omitempty"`
		Image                string   `json:"image,omitempty"`
		ImageType            int      `json:"imageType,omitempty"`
		Companies            struct {
			Studio         []any `json:"studio,omitempty"`
			Network        []any `json:"network,omitempty"`
			Production     []any `json:"production,omitempty"`
			Distributor    []any `json:"distributor,omitempty"`
			SpecialEffects []any `json:"special_effects,omitempty"`
		} `json:"companies,omitempty"`
		LastUpdated string `json:"lastUpdated,omitempty"`
		Year        string `json:"year,omitempty"`
		Episodes    []struct {
			ID                   int      `json:"id,omitempty"`
			SeriesID             int      `json:"seriesId,omitempty"`
			Name                 string   `json:"name,omitempty"`
			Aired                string   `json:"aired,omitempty"`
			Runtime              int      `json:"runtime,omitempty"`
			NameTranslations     []string `json:"nameTranslations,omitempty"`
			Overview             string   `json:"overview,omitempty"`
			OverviewTranslations []string `json:"overviewTranslations,omitempty"`
			Image                string   `json:"image,omitempty"`
			ImageType            int      `json:"imageType,omitempty"`
			IsMovie              int      `json:"isMovie,omitempty"`
			Seasons              any      `json:"seasons,omitempty"`
			Number               int      `json:"number,omitempty"`
			SeasonNumber         int      `json:"seasonNumber,omitempty"`
			LastUpdated          string   `json:"lastUpdated,omitempty"`
			FinaleType           any      `json:"finaleType,omitempty"`
			Year                 string   `json:"year,omitempty"`
		} `json:"episodes,omitempty"`
		Trailers []struct {
			ID       int    `json:"id,omitempty"`
			Language string `json:"language,omitempty"`
			Name     string `json:"name,omitempty"`
			Runtime  int    `json:"runtime,omitempty"`
			URL      string `json:"url,omitempty"`
		} `json:"trailers,omitempty"`
		Artwork []struct {
			ID           int    `json:"id,omitempty"`
			Image        string `json:"image,omitempty"`
			Thumbnail    string `json:"thumbnail,omitempty"`
			Language     string `json:"language,omitempty"`
			Type         int    `json:"type,omitempty"`
			Score        int    `json:"score,omitempty"`
			Width        int    `json:"width,omitempty"`
			Height       int    `json:"height,omitempty"`
			IncludesText bool   `json:"includesText,omitempty"`
		} `json:"artwork,omitempty"`
		TagOptions any `json:"tagOptions,omitempty"`
	} `json:"data,omitempty"`
}

func (c *Client) GetSeasonsByIDExtended(id int) (data *SeasonsByIDExtended, err error) {
	resp, err := c.DoRequest(RequestArgs{
		Method: http.MethodGet,
		Path:   strings.Replace(SeasonsByIDExtendedPath, ":id", strconv.Itoa(id), 1),
	})
	if err != nil {
		return
	}
	data = new(SeasonsByIDExtended)
	err = c.ParseResponse(resp.Body, data)
	return
}
