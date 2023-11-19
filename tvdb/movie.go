package tvdb

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	MovieByIDPath         = "/movies/:id"
	MovieByIDArtworksPath = "/movies/:id/artworks"
	MovieByIDExtendedPath = "/movies/:id/extended"
)

type MovieByIDExtended struct {
	Data struct {
		Aliases []struct {
			Language string `json:"language,omitempty"`
			Name     string `json:"name,omitempty"`
		} `json:"aliases,omitempty"`
		Artworks []struct {
			Height       int    `json:"height,omitempty"`
			ID           int    `json:"id,omitempty"`
			Image        string `json:"image,omitempty"`
			IncludesText bool   `json:"includesText,omitempty"`
			Language     string `json:"language,omitempty"`
			Score        int    `json:"score,omitempty"`
			Thumbnail    string `json:"thumbnail,omitempty"`
			Type         int    `json:"type,omitempty"`
			Width        int    `json:"width,omitempty"`
		} `json:"artworks,omitempty"`
		AudioLanguages []string `json:"audioLanguages,omitempty"`
		Awards         []struct {
			ID   int    `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"awards,omitempty"`
		BoxOffice   string `json:"boxOffice,omitempty"`
		BoxOfficeUS string `json:"boxOfficeUS,omitempty"`
		Budget      string `json:"budget,omitempty"`
		Characters  []struct {
			Aliases []struct {
				Language string `json:"language,omitempty"`
				Name     string `json:"name,omitempty"`
			} `json:"aliases,omitempty"`
			Episode struct {
				Image string `json:"image,omitempty"`
				Name  string `json:"name,omitempty"`
				Year  string `json:"year,omitempty"`
			} `json:"episode,omitempty"`
			EpisodeID  int    `json:"episodeId,omitempty"`
			ID         int    `json:"id,omitempty"`
			Image      string `json:"image,omitempty"`
			IsFeatured bool   `json:"isFeatured,omitempty"`
			MovieID    int    `json:"movieId,omitempty"`
			Movie      struct {
				Image string `json:"image,omitempty"`
				Name  string `json:"name,omitempty"`
				Year  string `json:"year,omitempty"`
			} `json:"movie,omitempty"`
			Name                 string   `json:"name,omitempty"`
			NameTranslations     []string `json:"nameTranslations,omitempty"`
			OverviewTranslations []string `json:"overviewTranslations,omitempty"`
			PeopleID             int      `json:"peopleId,omitempty"`
			PersonImgURL         string   `json:"personImgURL,omitempty"`
			PeopleType           string   `json:"peopleType,omitempty"`
			SeriesID             int      `json:"seriesId,omitempty"`
			Series               struct {
				Image string `json:"image,omitempty"`
				Name  string `json:"name,omitempty"`
				Year  string `json:"year,omitempty"`
			} `json:"series,omitempty"`
			Sort       int `json:"sort,omitempty"`
			TagOptions []struct {
				HelpText string `json:"helpText,omitempty"`
				ID       int    `json:"id,omitempty"`
				Name     string `json:"name,omitempty"`
				Tag      int    `json:"tag,omitempty"`
				TagName  string `json:"tagName,omitempty"`
			} `json:"tagOptions,omitempty"`
			Type       int    `json:"type,omitempty"`
			URL        string `json:"url,omitempty"`
			PersonName string `json:"personName,omitempty"`
		} `json:"characters,omitempty"`
		Companies struct {
			Studio []struct {
				ActiveDate string `json:"activeDate,omitempty"`
				Aliases    []struct {
					Language string `json:"language,omitempty"`
					Name     string `json:"name,omitempty"`
				} `json:"aliases,omitempty"`
				Country              string   `json:"country,omitempty"`
				ID                   int      `json:"id,omitempty"`
				InactiveDate         string   `json:"inactiveDate,omitempty"`
				Name                 string   `json:"name,omitempty"`
				NameTranslations     []string `json:"nameTranslations,omitempty"`
				OverviewTranslations []string `json:"overviewTranslations,omitempty"`
				PrimaryCompanyType   int      `json:"primaryCompanyType,omitempty"`
				Slug                 string   `json:"slug,omitempty"`
				ParentCompany        struct {
					ID       int    `json:"id,omitempty"`
					Name     string `json:"name,omitempty"`
					Relation struct {
						ID       int    `json:"id,omitempty"`
						TypeName string `json:"typeName,omitempty"`
					} `json:"relation,omitempty"`
				} `json:"parentCompany,omitempty"`
				TagOptions []struct {
					HelpText string `json:"helpText,omitempty"`
					ID       int    `json:"id,omitempty"`
					Name     string `json:"name,omitempty"`
					Tag      int    `json:"tag,omitempty"`
					TagName  string `json:"tagName,omitempty"`
				} `json:"tagOptions,omitempty"`
			} `json:"studio,omitempty"`
			Network []struct {
				ActiveDate string `json:"activeDate,omitempty"`
				Aliases    []struct {
					Language string `json:"language,omitempty"`
					Name     string `json:"name,omitempty"`
				} `json:"aliases,omitempty"`
				Country              string   `json:"country,omitempty"`
				ID                   int      `json:"id,omitempty"`
				InactiveDate         string   `json:"inactiveDate,omitempty"`
				Name                 string   `json:"name,omitempty"`
				NameTranslations     []string `json:"nameTranslations,omitempty"`
				OverviewTranslations []string `json:"overviewTranslations,omitempty"`
				PrimaryCompanyType   int      `json:"primaryCompanyType,omitempty"`
				Slug                 string   `json:"slug,omitempty"`
				ParentCompany        struct {
					ID       int    `json:"id,omitempty"`
					Name     string `json:"name,omitempty"`
					Relation struct {
						ID       int    `json:"id,omitempty"`
						TypeName string `json:"typeName,omitempty"`
					} `json:"relation,omitempty"`
				} `json:"parentCompany,omitempty"`
				TagOptions []struct {
					HelpText string `json:"helpText,omitempty"`
					ID       int    `json:"id,omitempty"`
					Name     string `json:"name,omitempty"`
					Tag      int    `json:"tag,omitempty"`
					TagName  string `json:"tagName,omitempty"`
				} `json:"tagOptions,omitempty"`
			} `json:"network,omitempty"`
			Production []struct {
				ActiveDate string `json:"activeDate,omitempty"`
				Aliases    []struct {
					Language string `json:"language,omitempty"`
					Name     string `json:"name,omitempty"`
				} `json:"aliases,omitempty"`
				Country              string   `json:"country,omitempty"`
				ID                   int      `json:"id,omitempty"`
				InactiveDate         string   `json:"inactiveDate,omitempty"`
				Name                 string   `json:"name,omitempty"`
				NameTranslations     []string `json:"nameTranslations,omitempty"`
				OverviewTranslations []string `json:"overviewTranslations,omitempty"`
				PrimaryCompanyType   int      `json:"primaryCompanyType,omitempty"`
				Slug                 string   `json:"slug,omitempty"`
				ParentCompany        struct {
					ID       int    `json:"id,omitempty"`
					Name     string `json:"name,omitempty"`
					Relation struct {
						ID       int    `json:"id,omitempty"`
						TypeName string `json:"typeName,omitempty"`
					} `json:"relation,omitempty"`
				} `json:"parentCompany,omitempty"`
				TagOptions []struct {
					HelpText string `json:"helpText,omitempty"`
					ID       int    `json:"id,omitempty"`
					Name     string `json:"name,omitempty"`
					Tag      int    `json:"tag,omitempty"`
					TagName  string `json:"tagName,omitempty"`
				} `json:"tagOptions,omitempty"`
			} `json:"production,omitempty"`
			Distributor []struct {
				ActiveDate string `json:"activeDate,omitempty"`
				Aliases    []struct {
					Language string `json:"language,omitempty"`
					Name     string `json:"name,omitempty"`
				} `json:"aliases,omitempty"`
				Country              string   `json:"country,omitempty"`
				ID                   int      `json:"id,omitempty"`
				InactiveDate         string   `json:"inactiveDate,omitempty"`
				Name                 string   `json:"name,omitempty"`
				NameTranslations     []string `json:"nameTranslations,omitempty"`
				OverviewTranslations []string `json:"overviewTranslations,omitempty"`
				PrimaryCompanyType   int      `json:"primaryCompanyType,omitempty"`
				Slug                 string   `json:"slug,omitempty"`
				ParentCompany        struct {
					ID       int    `json:"id,omitempty"`
					Name     string `json:"name,omitempty"`
					Relation struct {
						ID       int    `json:"id,omitempty"`
						TypeName string `json:"typeName,omitempty"`
					} `json:"relation,omitempty"`
				} `json:"parentCompany,omitempty"`
				TagOptions []struct {
					HelpText string `json:"helpText,omitempty"`
					ID       int    `json:"id,omitempty"`
					Name     string `json:"name,omitempty"`
					Tag      int    `json:"tag,omitempty"`
					TagName  string `json:"tagName,omitempty"`
				} `json:"tagOptions,omitempty"`
			} `json:"distributor,omitempty"`
			SpecialEffects []struct {
				ActiveDate string `json:"activeDate,omitempty"`
				Aliases    []struct {
					Language string `json:"language,omitempty"`
					Name     string `json:"name,omitempty"`
				} `json:"aliases,omitempty"`
				Country              string   `json:"country,omitempty"`
				ID                   int      `json:"id,omitempty"`
				InactiveDate         string   `json:"inactiveDate,omitempty"`
				Name                 string   `json:"name,omitempty"`
				NameTranslations     []string `json:"nameTranslations,omitempty"`
				OverviewTranslations []string `json:"overviewTranslations,omitempty"`
				PrimaryCompanyType   int      `json:"primaryCompanyType,omitempty"`
				Slug                 string   `json:"slug,omitempty"`
				ParentCompany        struct {
					ID       int    `json:"id,omitempty"`
					Name     string `json:"name,omitempty"`
					Relation struct {
						ID       int    `json:"id,omitempty"`
						TypeName string `json:"typeName,omitempty"`
					} `json:"relation,omitempty"`
				} `json:"parentCompany,omitempty"`
				TagOptions []struct {
					HelpText string `json:"helpText,omitempty"`
					ID       int    `json:"id,omitempty"`
					Name     string `json:"name,omitempty"`
					Tag      int    `json:"tag,omitempty"`
					TagName  string `json:"tagName,omitempty"`
				} `json:"tagOptions,omitempty"`
			} `json:"special_effects,omitempty"`
		} `json:"companies,omitempty"`
		ContentRatings []struct {
			ID          int    `json:"id,omitempty"`
			Name        string `json:"name,omitempty"`
			Description string `json:"description,omitempty"`
			Country     string `json:"country,omitempty"`
			ContentType string `json:"contentType,omitempty"`
			Order       int    `json:"order,omitempty"`
			FullName    string `json:"fullName,omitempty"`
		} `json:"contentRatings,omitempty"`
		FirstRelease struct {
			Country string `json:"country,omitempty"`
			Date    string `json:"date,omitempty"`
			Detail  string `json:"detail,omitempty"`
		} `json:"first_release,omitempty"`
		Genres []struct {
			ID   int    `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
			Slug string `json:"slug,omitempty"`
		} `json:"genres,omitempty"`
		ID           int    `json:"id,omitempty"`
		Image        string `json:"image,omitempty"`
		Inspirations []struct {
			ID       int    `json:"id,omitempty"`
			Type     string `json:"type,omitempty"`
			TypeName string `json:"type_name,omitempty"`
			URL      string `json:"url,omitempty"`
		} `json:"inspirations,omitempty"`
		LastUpdated string `json:"lastUpdated,omitempty"`
		Lists       []struct {
			Aliases []struct {
				Language string `json:"language,omitempty"`
				Name     string `json:"name,omitempty"`
			} `json:"aliases,omitempty"`
			ID                   int      `json:"id,omitempty"`
			Image                string   `json:"image,omitempty"`
			ImageIsFallback      bool     `json:"imageIsFallback,omitempty"`
			IsOfficial           bool     `json:"isOfficial,omitempty"`
			Name                 string   `json:"name,omitempty"`
			NameTranslations     []string `json:"nameTranslations,omitempty"`
			Overview             string   `json:"overview,omitempty"`
			OverviewTranslations []string `json:"overviewTranslations,omitempty"`
			RemoteIds            []struct {
				ID         string `json:"id,omitempty"`
				Type       int    `json:"type,omitempty"`
				SourceName string `json:"sourceName,omitempty"`
			} `json:"remoteIds,omitempty"`
			Tags []struct {
				HelpText string `json:"helpText,omitempty"`
				ID       int    `json:"id,omitempty"`
				Name     string `json:"name,omitempty"`
				Tag      int    `json:"tag,omitempty"`
				TagName  string `json:"tagName,omitempty"`
			} `json:"tags,omitempty"`
			Score int    `json:"score,omitempty"`
			URL   string `json:"url,omitempty"`
		} `json:"lists,omitempty"`
		Name                 string   `json:"name,omitempty"`
		NameTranslations     []string `json:"nameTranslations,omitempty"`
		OriginalCountry      string   `json:"originalCountry,omitempty"`
		OriginalLanguage     string   `json:"originalLanguage,omitempty"`
		OverviewTranslations []string `json:"overviewTranslations,omitempty"`
		ProductionCountries  []struct {
			ID      int    `json:"id,omitempty"`
			Country string `json:"country,omitempty"`
			Name    string `json:"name,omitempty"`
		} `json:"production_countries,omitempty"`
		Releases []struct {
			Country string `json:"country,omitempty"`
			Date    string `json:"date,omitempty"`
			Detail  string `json:"detail,omitempty"`
		} `json:"releases,omitempty"`
		RemoteIds []struct {
			ID         string `json:"id,omitempty"`
			Type       int    `json:"type,omitempty"`
			SourceName string `json:"sourceName,omitempty"`
		} `json:"remoteIds,omitempty"`
		Runtime         int      `json:"runtime,omitempty"`
		Score           int      `json:"score,omitempty"`
		Slug            string   `json:"slug,omitempty"`
		SpokenLanguages []string `json:"spoken_languages,omitempty"`
		Status          struct {
			ID          int    `json:"id,omitempty"`
			KeepUpdated bool   `json:"keepUpdated,omitempty"`
			Name        string `json:"name,omitempty"`
			RecordType  string `json:"recordType,omitempty"`
		} `json:"status,omitempty"`
		Studios []struct {
			ID           int    `json:"id,omitempty"`
			Name         string `json:"name,omitempty"`
			ParentStudio int    `json:"parentStudio,omitempty"`
		} `json:"studios,omitempty"`
		SubtitleLanguages []string `json:"subtitleLanguages,omitempty"`
		TagOptions        []struct {
			HelpText string `json:"helpText,omitempty"`
			ID       int    `json:"id,omitempty"`
			Name     string `json:"name,omitempty"`
			Tag      int    `json:"tag,omitempty"`
			TagName  string `json:"tagName,omitempty"`
		} `json:"tagOptions,omitempty"`
		Trailers []struct {
			ID       int    `json:"id,omitempty"`
			Language string `json:"language,omitempty"`
			Name     string `json:"name,omitempty"`
			URL      string `json:"url,omitempty"`
			Runtime  int    `json:"runtime,omitempty"`
		} `json:"trailers,omitempty"`
		Translations struct {
			NameTranslations []struct {
				Aliases   []string `json:"aliases,omitempty"`
				IsAlias   bool     `json:"isAlias,omitempty"`
				IsPrimary bool     `json:"isPrimary,omitempty"`
				Language  string   `json:"language,omitempty"`
				Name      string   `json:"name,omitempty"`
				Overview  string   `json:"overview,omitempty"`
				Tagline   string   `json:"tagline,omitempty"`
			} `json:"nameTranslations,omitempty"`
			OverviewTranslations []struct {
				Aliases   []string `json:"aliases,omitempty"`
				IsAlias   bool     `json:"isAlias,omitempty"`
				IsPrimary bool     `json:"isPrimary,omitempty"`
				Language  string   `json:"language,omitempty"`
				Name      string   `json:"name,omitempty"`
				Overview  string   `json:"overview,omitempty"`
				Tagline   string   `json:"tagline,omitempty"`
			} `json:"overviewTranslations,omitempty"`
			Alias []string `json:"alias,omitempty"`
		} `json:"translations,omitempty"`
		Year string `json:"year,omitempty"`
	} `json:"data,omitempty"`
	Status string `json:"status,omitempty"`
}

func (c *Client) GetMovieByIDExtended(id int) (data *MovieByIDExtended, err error) {
	resp, err := c.DoRequest(RequestArgs{
		Method: http.MethodGet,
		Path:   strings.Replace(MovieByIDExtendedPath, ":id", strconv.Itoa(id), 1),
	})
	if err != nil {
		return
	}

	data = new(MovieByIDExtended)
	err = c.ParseResponse(resp.Body, data)
	return
}

/* func (c *Client) GetMovieByIDExtended(id int) (data *MovieByIDExtended, err error) {
	url := c.BuildUrlPath(strings.Replace(MovieByIDExtendedPath, ":id", strconv.Itoa(id), 1))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", ApplicationJson)
	req.Header.Set("Accept", ApplicationJson)
	req.Header.Set("Accept-Language", c.Language)
	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println(string(responseBody))

	data = new(MovieByIDExtended)
	err = c.ParseResponse(resp.Body, data)
	return
}
*/
