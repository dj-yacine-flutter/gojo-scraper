package tvdb

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	SeriesByIDPath         = "/series/:id"
	SeriesByIDArtworksPath = "/series/:id/artworks"
	SeriesByIDExtendedPath = "/series/:id/extended"
)

type SeriesByIDExtended struct {
	Data struct {
		Abbreviation string `json:"abbreviation,omitempty"`
		AirsDays     struct {
			Friday    bool `json:"friday,omitempty"`
			Monday    bool `json:"monday,omitempty"`
			Saturday  bool `json:"saturday,omitempty"`
			Sunday    bool `json:"sunday,omitempty"`
			Thursday  bool `json:"thursday,omitempty"`
			Tuesday   bool `json:"tuesday,omitempty"`
			Wednesday bool `json:"wednesday,omitempty"`
		} `json:"airsDays,omitempty"`
		AirsTime string `json:"airsTime,omitempty"`
		Aliases  []struct {
			Language string `json:"language,omitempty"`
			Name     string `json:"name,omitempty"`
		} `json:"aliases,omitempty"`
		Artworks []struct {
			EpisodeID      int    `json:"episodeId,omitempty"`
			Height         int    `json:"height,omitempty"`
			ID             int    `json:"id,omitempty"`
			Image          string `json:"image,omitempty"`
			IncludesText   bool   `json:"includesText,omitempty"`
			Language       string `json:"language,omitempty"`
			MovieID        int    `json:"movieId,omitempty"`
			NetworkID      int    `json:"networkId,omitempty"`
			PeopleID       int    `json:"peopleId,omitempty"`
			Score          int    `json:"score,omitempty"`
			SeasonID       int    `json:"seasonId,omitempty"`
			SeriesID       int    `json:"seriesId,omitempty"`
			SeriesPeopleID int    `json:"seriesPeopleId,omitempty"`
			Status         struct {
				ID   int    `json:"id,omitempty"`
				Name string `json:"name,omitempty"`
			} `json:"status,omitempty"`
			TagOptions []struct {
				HelpText string `json:"helpText,omitempty"`
				ID       int    `json:"id,omitempty"`
				Name     string `json:"name,omitempty"`
				Tag      int    `json:"tag,omitempty"`
				TagName  string `json:"tagName,omitempty"`
			} `json:"tagOptions,omitempty"`
			Thumbnail       string `json:"thumbnail,omitempty"`
			ThumbnailHeight int    `json:"thumbnailHeight,omitempty"`
			ThumbnailWidth  int    `json:"thumbnailWidth,omitempty"`
			Type            int    `json:"type,omitempty"`
			UpdatedAt       int    `json:"updatedAt,omitempty"`
			Width           int    `json:"width,omitempty"`
		} `json:"artworks,omitempty"`
		AverageRuntime int `json:"averageRuntime,omitempty"`
		Characters     []struct {
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
		ContentRatings []struct {
			ID          int    `json:"id,omitempty"`
			Name        string `json:"name,omitempty"`
			Description string `json:"description,omitempty"`
			Country     string `json:"country,omitempty"`
			ContentType string `json:"contentType,omitempty"`
			Order       int    `json:"order,omitempty"`
			FullName    string `json:"fullName,omitempty"`
		} `json:"contentRatings,omitempty"`
		Country           string `json:"country,omitempty"`
		DefaultSeasonType int    `json:"defaultSeasonType,omitempty"`
		Episodes          []struct {
			Aired                string   `json:"aired,omitempty"`
			AirsAfterSeason      int      `json:"airsAfterSeason,omitempty"`
			AirsBeforeEpisode    int      `json:"airsBeforeEpisode,omitempty"`
			AirsBeforeSeason     int      `json:"airsBeforeSeason,omitempty"`
			FinaleType           string   `json:"finaleType,omitempty"`
			ID                   int      `json:"id,omitempty"`
			Image                string   `json:"image,omitempty"`
			ImageType            int      `json:"imageType,omitempty"`
			IsMovie              int      `json:"isMovie,omitempty"`
			LastUpdated          string   `json:"lastUpdated,omitempty"`
			LinkedMovie          int      `json:"linkedMovie,omitempty"`
			Name                 string   `json:"name,omitempty"`
			NameTranslations     []string `json:"nameTranslations,omitempty"`
			Number               int      `json:"number,omitempty"`
			Overview             string   `json:"overview,omitempty"`
			OverviewTranslations []string `json:"overviewTranslations,omitempty"`
			Runtime              int      `json:"runtime,omitempty"`
			SeasonNumber         int      `json:"seasonNumber,omitempty"`
			Seasons              []struct {
				ID                   int      `json:"id,omitempty"`
				Image                string   `json:"image,omitempty"`
				ImageType            int      `json:"imageType,omitempty"`
				LastUpdated          string   `json:"lastUpdated,omitempty"`
				Name                 string   `json:"name,omitempty"`
				NameTranslations     []string `json:"nameTranslations,omitempty"`
				Number               int      `json:"number,omitempty"`
				OverviewTranslations []string `json:"overviewTranslations,omitempty"`
				Companies            struct {
					Studio struct {
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
					Network struct {
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
					Production struct {
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
					Distributor struct {
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
					SpecialEffects struct {
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
				SeriesID int `json:"seriesId,omitempty"`
				Type     struct {
					AlternateName string `json:"alternateName,omitempty"`
					ID            int    `json:"id,omitempty"`
					Name          string `json:"name,omitempty"`
					Type          string `json:"type,omitempty"`
				} `json:"type,omitempty"`
				Year string `json:"year,omitempty"`
			} `json:"seasons,omitempty"`
			SeriesID   int    `json:"seriesId,omitempty"`
			SeasonName string `json:"seasonName,omitempty"`
			Year       string `json:"year,omitempty"`
		} `json:"episodes,omitempty"`
		FirstAired string `json:"firstAired,omitempty"`
		Lists      []struct {
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
		Genres []struct {
			ID   int    `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
			Slug string `json:"slug,omitempty"`
		} `json:"genres,omitempty"`
		ID                int      `json:"id,omitempty"`
		Image             string   `json:"image,omitempty"`
		IsOrderRandomized bool     `json:"isOrderRandomized,omitempty"`
		LastAired         string   `json:"lastAired,omitempty"`
		LastUpdated       string   `json:"lastUpdated,omitempty"`
		Name              string   `json:"name,omitempty"`
		NameTranslations  []string `json:"nameTranslations,omitempty"`
		Companies         []struct {
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
		} `json:"companies,omitempty"`
		NextAired        string `json:"nextAired,omitempty"`
		OriginalCountry  string `json:"originalCountry,omitempty"`
		OriginalLanguage string `json:"originalLanguage,omitempty"`
		OriginalNetwork  struct {
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
		} `json:"originalNetwork,omitempty"`
		Overview      string `json:"overview,omitempty"`
		LatestNetwork struct {
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
		} `json:"latestNetwork,omitempty"`
		OverviewTranslations []string `json:"overviewTranslations,omitempty"`
		RemoteIds            []struct {
			ID         string `json:"id,omitempty"`
			Type       int    `json:"type,omitempty"`
			SourceName string `json:"sourceName,omitempty"`
		} `json:"remoteIds,omitempty"`
		Score   int `json:"score,omitempty"`
		Seasons []struct {
			ID                   int      `json:"id,omitempty"`
			Image                string   `json:"image,omitempty"`
			ImageType            int      `json:"imageType,omitempty"`
			LastUpdated          string   `json:"lastUpdated,omitempty"`
			Name                 string   `json:"name,omitempty"`
			NameTranslations     []string `json:"nameTranslations,omitempty"`
			Number               int      `json:"number,omitempty"`
			OverviewTranslations []string `json:"overviewTranslations,omitempty"`
			Companies            struct {
				Studio struct {
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
				Network struct {
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
				Production struct {
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
				Distributor struct {
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
				SpecialEffects struct {
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
			SeriesID int `json:"seriesId,omitempty"`
			Type     struct {
				AlternateName string `json:"alternateName,omitempty"`
				ID            int    `json:"id,omitempty"`
				Name          string `json:"name,omitempty"`
				Type          string `json:"type,omitempty"`
			} `json:"type,omitempty"`
			Year string `json:"year,omitempty"`
		} `json:"seasons,omitempty"`
		SeasonTypes []struct {
			AlternateName string `json:"alternateName,omitempty"`
			ID            int    `json:"id,omitempty"`
			Name          string `json:"name,omitempty"`
			Type          string `json:"type,omitempty"`
		} `json:"seasonTypes,omitempty"`
		Slug   string `json:"slug,omitempty"`
		Status struct {
			ID          int    `json:"id,omitempty"`
			KeepUpdated bool   `json:"keepUpdated,omitempty"`
			Name        string `json:"name,omitempty"`
			RecordType  string `json:"recordType,omitempty"`
		} `json:"status,omitempty"`
		Tags []struct {
			HelpText string `json:"helpText,omitempty"`
			ID       int    `json:"id,omitempty"`
			Name     string `json:"name,omitempty"`
			Tag      int    `json:"tag,omitempty"`
			TagName  string `json:"tagName,omitempty"`
		} `json:"tags,omitempty"`
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

func (c *Client) GetSeriesByIDExtanded(id int) (data *SeriesByIDExtended, err error) {
	resp, err := c.DoRequest(RequestArgs{
		Method: http.MethodGet,
		Path:   strings.Replace(SeriesByIDExtendedPath, ":id", strconv.Itoa(id), 1),
	})
	if err != nil {
		return
	}
	data = new(SeriesByIDExtended)
	err = c.ParseResponse(resp.Body, data)
	return
}

type SeriesByID struct {
	Data struct {
		Aliases []struct {
			Language string `json:"language,omitempty"`
			Name     string `json:"name,omitempty"`
		} `json:"aliases,omitempty"`
		AverageRuntime    int    `json:"averageRuntime,omitempty"`
		Country           string `json:"country,omitempty"`
		DefaultSeasonType int    `json:"defaultSeasonType,omitempty"`
		Episodes          []struct {
			Aired                string   `json:"aired,omitempty"`
			AirsAfterSeason      int      `json:"airsAfterSeason,omitempty"`
			AirsBeforeEpisode    int      `json:"airsBeforeEpisode,omitempty"`
			AirsBeforeSeason     int      `json:"airsBeforeSeason,omitempty"`
			FinaleType           string   `json:"finaleType,omitempty"`
			ID                   int      `json:"id,omitempty"`
			Image                string   `json:"image,omitempty"`
			ImageType            int      `json:"imageType,omitempty"`
			IsMovie              int      `json:"isMovie,omitempty"`
			LastUpdated          string   `json:"lastUpdated,omitempty"`
			LinkedMovie          int      `json:"linkedMovie,omitempty"`
			Name                 string   `json:"name,omitempty"`
			NameTranslations     []string `json:"nameTranslations,omitempty"`
			Number               int      `json:"number,omitempty"`
			Overview             string   `json:"overview,omitempty"`
			OverviewTranslations []string `json:"overviewTranslations,omitempty"`
			Runtime              int      `json:"runtime,omitempty"`
			SeasonNumber         int      `json:"seasonNumber,omitempty"`
			Seasons              []struct {
				ID                   int      `json:"id,omitempty"`
				Image                string   `json:"image,omitempty"`
				ImageType            int      `json:"imageType,omitempty"`
				LastUpdated          string   `json:"lastUpdated,omitempty"`
				Name                 string   `json:"name,omitempty"`
				NameTranslations     []string `json:"nameTranslations,omitempty"`
				Number               int      `json:"number,omitempty"`
				OverviewTranslations []string `json:"overviewTranslations,omitempty"`
				Companies            struct {
					Studio struct {
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
					Network struct {
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
					Production struct {
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
					Distributor struct {
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
					SpecialEffects struct {
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
				SeriesID int `json:"seriesId,omitempty"`
				Type     struct {
					AlternateName string `json:"alternateName,omitempty"`
					ID            int    `json:"id,omitempty"`
					Name          string `json:"name,omitempty"`
					Type          string `json:"type,omitempty"`
				} `json:"type,omitempty"`
				Year string `json:"year,omitempty"`
			} `json:"seasons,omitempty"`
			SeriesID   int    `json:"seriesId,omitempty"`
			SeasonName string `json:"seasonName,omitempty"`
			Year       string `json:"year,omitempty"`
		} `json:"episodes,omitempty"`
		FirstAired           string   `json:"firstAired,omitempty"`
		ID                   int      `json:"id,omitempty"`
		Image                string   `json:"image,omitempty"`
		IsOrderRandomized    bool     `json:"isOrderRandomized,omitempty"`
		LastAired            string   `json:"lastAired,omitempty"`
		LastUpdated          string   `json:"lastUpdated,omitempty"`
		Name                 string   `json:"name,omitempty"`
		NameTranslations     []string `json:"nameTranslations,omitempty"`
		NextAired            string   `json:"nextAired,omitempty"`
		OriginalCountry      string   `json:"originalCountry,omitempty"`
		OriginalLanguage     string   `json:"originalLanguage,omitempty"`
		OverviewTranslations []string `json:"overviewTranslations,omitempty"`
		Score                int      `json:"score,omitempty"`
		Slug                 string   `json:"slug,omitempty"`
		Status               struct {
			ID          int    `json:"id,omitempty"`
			KeepUpdated bool   `json:"keepUpdated,omitempty"`
			Name        string `json:"name,omitempty"`
			RecordType  string `json:"recordType,omitempty"`
		} `json:"status,omitempty"`
		Year string `json:"year,omitempty"`
	} `json:"data,omitempty"`
	Status string `json:"status,omitempty"`
}

func (c *Client) GetSeriesByID(id int) (data *SeriesByID, err error) {
	resp, err := c.DoRequest(RequestArgs{
		Method: http.MethodGet,
		Path:   strings.Replace(SeriesByIDPath, ":id", strconv.Itoa(id), 1),
	})
	if err != nil {
		return
	}
	data = new(SeriesByID)
	err = c.ParseResponse(resp.Body, data)
	return
}
