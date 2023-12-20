package tvdb

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	EpisodeByIDExtandedPath = "/episodes/:id/extended"
	EpisodeByIDTrPath       = "/episodes/:id/translations/:language"
)

type EpisodeByIDExtanded struct {
	Data struct {
		Aired             string `json:"aired,omitempty"`
		AirsAfterSeason   int    `json:"airsAfterSeason,omitempty"`
		AirsBeforeEpisode int    `json:"airsBeforeEpisode,omitempty"`
		AirsBeforeSeason  int    `json:"airsBeforeSeason,omitempty"`
		Awards            []struct {
			ID   int    `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"awards,omitempty"`
		Characters []struct {
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
		Companies []struct {
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
		ContentRatings []struct {
			ID          int    `json:"id,omitempty"`
			Name        string `json:"name,omitempty"`
			Description string `json:"description,omitempty"`
			Country     string `json:"country,omitempty"`
			ContentType string `json:"contentType,omitempty"`
			Order       int    `json:"order,omitempty"`
			FullName    string `json:"fullName,omitempty"`
		} `json:"contentRatings,omitempty"`
		FinaleType       string   `json:"finaleType,omitempty"`
		ID               int      `json:"id,omitempty"`
		Image            string   `json:"image,omitempty"`
		ImageType        int      `json:"imageType,omitempty"`
		IsMovie          int      `json:"isMovie,omitempty"`
		LastUpdated      string   `json:"lastUpdated,omitempty"`
		LinkedMovie      int      `json:"linkedMovie,omitempty"`
		Name             string   `json:"name,omitempty"`
		NameTranslations []string `json:"nameTranslations,omitempty"`
		Networks         []struct {
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
		} `json:"networks,omitempty"`
		Nominations []struct {
			Character struct {
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
			} `json:"character,omitempty"`
			Details string `json:"details,omitempty"`
			Episode struct {
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
			} `json:"episode,omitempty"`
			ID       int  `json:"id,omitempty"`
			IsWinner bool `json:"isWinner,omitempty"`
			Movie    struct {
				Aliases []struct {
					Language string `json:"language,omitempty"`
					Name     string `json:"name,omitempty"`
				} `json:"aliases,omitempty"`
				ID                   int      `json:"id,omitempty"`
				Image                string   `json:"image,omitempty"`
				LastUpdated          string   `json:"lastUpdated,omitempty"`
				Name                 string   `json:"name,omitempty"`
				NameTranslations     []string `json:"nameTranslations,omitempty"`
				OverviewTranslations []string `json:"overviewTranslations,omitempty"`
				Score                int      `json:"score,omitempty"`
				Slug                 string   `json:"slug,omitempty"`
				Status               struct {
					ID          int    `json:"id,omitempty"`
					KeepUpdated bool   `json:"keepUpdated,omitempty"`
					Name        string `json:"name,omitempty"`
					RecordType  string `json:"recordType,omitempty"`
				} `json:"status,omitempty"`
				Runtime int    `json:"runtime,omitempty"`
				Year    string `json:"year,omitempty"`
			} `json:"movie,omitempty"`
			Series struct {
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
			} `json:"series,omitempty"`
			Year     string `json:"year,omitempty"`
			Category string `json:"category,omitempty"`
			Name     string `json:"name,omitempty"`
		} `json:"nominations,omitempty"`
		Number               int      `json:"number,omitempty"`
		Overview             string   `json:"overview,omitempty"`
		OverviewTranslations []string `json:"overviewTranslations,omitempty"`
		ProductionCode       string   `json:"productionCode,omitempty"`
		RemoteIds            []struct {
			ID         string `json:"id,omitempty"`
			Type       int    `json:"type,omitempty"`
			SourceName string `json:"sourceName,omitempty"`
		} `json:"remoteIds,omitempty"`
		Runtime      int `json:"runtime,omitempty"`
		SeasonNumber int `json:"seasonNumber,omitempty"`
		Seasons      []struct {
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
		SeriesID int `json:"seriesId,omitempty"`
		Studios  []struct {
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
		} `json:"studios,omitempty"`
		TagOptions []struct {
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

func (c *Client) GetEpisodeByIDExtanded(id int) (data *EpisodeByIDExtanded, err error) {
	resp, err := c.DoRequest(RequestArgs{
		Method: http.MethodGet,
		Path:   strings.Replace(EpisodeByIDExtandedPath, ":id", strconv.Itoa(id), 1),
	})
	if err != nil {
		return
	}

	data = new(EpisodeByIDExtanded)
	err = c.ParseResponse(resp.Body, data)
	return
}

type EpisodeByIDTr struct {
	Data struct {
		Aliases   []string `json:"aliases,omitempty"`
		IsAlias   bool     `json:"isAlias,omitempty"`
		IsPrimary bool     `json:"isPrimary,omitempty"`
		Language  string   `json:"language,omitempty"`
		Name      string   `json:"name,omitempty"`
		Overview  string   `json:"overview,omitempty"`
		Tagline   string   `json:"tagline,omitempty"`
	} `json:"data,omitempty"`
	Status string `json:"status,omitempty"`
}

func (c *Client) GetEpisodeByIDTr(id int, language string) (data *EpisodeByIDTr, err error) {
	resp, err := c.DoRequest(RequestArgs{
		Method: http.MethodGet,
		Path:   strings.Replace(strings.Replace(EpisodeByIDTrPath, ":id", strconv.Itoa(id), 1), ":language", language, 1),
	})
	if err != nil {
		return
	}

	data = new(EpisodeByIDTr)
	err = c.ParseResponse(resp.Body, data)
	return
}
