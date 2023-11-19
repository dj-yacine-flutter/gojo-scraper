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
	Status string `json:"status"`
	Data   struct {
		ID       int `json:"id"`
		SeriesID int `json:"seriesId"`
		Type     struct {
			ID            int    `json:"id"`
			Name          string `json:"name"`
			Type          string `json:"type"`
			AlternateName any    `json:"alternateName"`
		} `json:"type"`
		Number               int    `json:"number"`
		NameTranslations     []any  `json:"nameTranslations"`
		OverviewTranslations []any  `json:"overviewTranslations"`
		Image                string `json:"image"`
		ImageType            int    `json:"imageType"`
		Companies            struct {
			Studio         []any `json:"studio"`
			Network        []any `json:"network"`
			Production     []any `json:"production"`
			Distributor    []any `json:"distributor"`
			SpecialEffects []any `json:"special_effects"`
		} `json:"companies"`
		LastUpdated string `json:"lastUpdated"`
		Year        string `json:"year"`
		Episodes    []struct {
			ID                   int      `json:"id"`
			SeriesID             int      `json:"seriesId"`
			Name                 string   `json:"name"`
			Aired                string   `json:"aired"`
			Runtime              int      `json:"runtime"`
			NameTranslations     []string `json:"nameTranslations"`
			Overview             any      `json:"overview"`
			OverviewTranslations []string `json:"overviewTranslations"`
			Image                string   `json:"image"`
			ImageType            int      `json:"imageType"`
			IsMovie              int      `json:"isMovie"`
			Seasons              any      `json:"seasons"`
			Number               int      `json:"number"`
			SeasonNumber         int      `json:"seasonNumber"`
			LastUpdated          string   `json:"lastUpdated"`
			FinaleType           any      `json:"finaleType"`
			AirsBeforeSeason     int      `json:"airsBeforeSeason,omitempty"`
			AirsBeforeEpisode    int      `json:"airsBeforeEpisode,omitempty"`
			Year                 string   `json:"year"`
		} `json:"episodes"`
		Trailers []any `json:"trailers"`
		Artwork  []struct {
			ID           int    `json:"id"`
			Image        string `json:"image"`
			Thumbnail    string `json:"thumbnail"`
			Language     string `json:"language"`
			Type         int    `json:"type"`
			Score        int    `json:"score"`
			Width        int    `json:"width"`
			Height       int    `json:"height"`
			IncludesText bool   `json:"includesText"`
		} `json:"artwork"`
		TagOptions any `json:"tagOptions"`
	} `json:"data"`
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
