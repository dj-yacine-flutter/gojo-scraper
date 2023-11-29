package models

type LibreTranslate struct {
	DetectedLanguage DetectedLanguage `json:"detectedLanguage"`
	TranslatedText   string           `json:"translatedText"`
}

type DetectedLanguage struct {
	Confidence float32 `json:"confidence"`
	Language   string  `json:"language"`
}

type Meta struct {
	Title    string `json:"title"`
	Overview string `json:"overview"`
}

type MetaData struct {
	Language string `json:"language"`
	Meta     Meta   `json:"meta"`
}

type AnimeResources struct {
	LivechartID   int    `json:"livechartID"`
	AnimePlanetID string `json:"animePlanetID"`
	AnisearchID   int    `json:"anisearchID"`
	AnidbID       int    `json:"anidbID"`
	KitsuID       int    `json:"kitsuID"`
	MalID         int    `json:"malID"`
	NotifyMoeID   string `json:"notifyMoeID"`
	AnilistID     int    `json:"anilistID"`
	ThetvdbID     int    `json:"TVDBID"`
	ImdbID        string `json:"IMDBID"`
	ThemoviedbID  int    `json:"TMDBID"`
	Type          string `json:"type"`
}

type Titles struct {
	Offical []string `json:"official"`
	Short   []string `json:"short"`
	Others  []string `json:"others"`
}

type Image struct {
	Height    int    `json:"height"`
	Width     int    `json:"width"`
	Image     string `json:"image"`
	Thumbnail string `json:"thumbnail"`
	BlurHash  string `json:"blurHash"`
}

type Trailer struct {
	Official bool   `json:"official"`
	Host     string `json:"host"`
	Key      string `json:"key"`
}

type Anime struct {
	OriginalTitle       string         `json:"originalTitle"`
	Aired               string         `json:"aired"`
	ReleaseYear         int            `json:"releaseYear"`
	Rating              string         `json:"rating"`
	Runtime             string         `json:"runtime"`
	PortriatPoster      string         `json:"portriatPoster"`
	PortriatBlurHash    string         `json:"portriatBlurHash"`
	LandscapePoster     string         `json:"landscapePoster"`
	LandscapeBlurHash   string         `json:"landscapeBlurHash"`
	AnimeResources      AnimeResources `json:"animeResources"`
	Titles              Titles         `json:"titles"`
	Genres              []string       `json:"genres"`
	Studios             []string       `json:"studios"`
	ProductionCompanies []string       `json:"productionCompanies"`
	Tags                []string       `json:"tags"`
	Posters             []Image        `json:"posters"`
	Backdrops           []Image        `json:"backdrops"`
	Logos               []Image        `json:"logos"`
	Trailers            []Trailer      `json:"trailers"`
	AnimeMetas          []MetaData     `json:"animeMetas"`
}

type NotifyMoe struct {
	ID    string `json:"id,omitempty"`
	Type  string `json:"type,omitempty"`
	Title struct {
		Canonical string   `json:"canonical,omitempty"`
		Romaji    string   `json:"romaji,omitempty"`
		English   string   `json:"english,omitempty"`
		Japanese  string   `json:"japanese,omitempty"`
		Hiragana  string   `json:"hiragana,omitempty"`
		Synonyms  []string `json:"synonyms,omitempty"`
	} `json:"title,omitempty"`
	Summary       string   `json:"summary,omitempty"`
	Status        string   `json:"status,omitempty"`
	Genres        []string `json:"genres,omitempty"`
	StartDate     string   `json:"startDate,omitempty"`
	EndDate       string   `json:"endDate,omitempty"`
	EpisodeCount  int      `json:"episodeCount,omitempty"`
	EpisodeLength int      `json:"episodeLength,omitempty"`
	Source        string   `json:"source,omitempty"`
	Image         struct {
		Extension    string `json:"extension,omitempty"`
		Width        int    `json:"width,omitempty"`
		Height       int    `json:"height,omitempty"`
		AverageColor struct {
			Hue        float64 `json:"hue,omitempty"`
			Saturation float64 `json:"saturation,omitempty"`
			Lightness  float64 `json:"lightness,omitempty"`
		} `json:"averageColor,omitempty"`
		LastModified int `json:"lastModified,omitempty"`
	} `json:"image,omitempty"`
	FirstChannel string `json:"firstChannel,omitempty"`
	Rating       struct {
		Overall    float64 `json:"overall,omitempty"`
		Story      float64 `json:"story,omitempty"`
		Visuals    float64 `json:"visuals,omitempty"`
		Soundtrack float64 `json:"soundtrack,omitempty"`
		Count      struct {
			Overall    int `json:"overall,omitempty"`
			Story      int `json:"story,omitempty"`
			Visuals    int `json:"visuals,omitempty"`
			Soundtrack int `json:"soundtrack,omitempty"`
		} `json:"count,omitempty"`
	} `json:"rating,omitempty"`
	Popularity struct {
		Watching  int `json:"watching,omitempty"`
		Completed int `json:"completed,omitempty"`
		Planned   int `json:"planned,omitempty"`
		Hold      int `json:"hold,omitempty"`
		Dropped   int `json:"dropped,omitempty"`
	} `json:"popularity,omitempty"`
	Trailers []struct {
		Service   string `json:"service,omitempty"`
		ServiceID string `json:"serviceId,omitempty"`
	} `json:"trailers,omitempty"`
	Episodes []string `json:"episodes,omitempty"`
	Mappings []struct {
		Service   string `json:"service,omitempty"`
		ServiceID string `json:"serviceId,omitempty"`
	} `json:"mappings,omitempty"`
	Posts     []string `json:"posts,omitempty"`
	Likes     any      `json:"likes,omitempty"`
	Created   string   `json:"created,omitempty"`
	CreatedBy string   `json:"createdBy,omitempty"`
	Edited    string   `json:"edited,omitempty"`
	EditedBy  string   `json:"editedBy,omitempty"`
	IsDraft   bool     `json:"isDraft,omitempty"`
	Studios   []string `json:"studios,omitempty"`
	Producers []string `json:"producers,omitempty"`
	Licensors any      `json:"licensors,omitempty"`
	Links     []struct {
		Title string `json:"title,omitempty"`
		URL   string `json:"url,omitempty"`
	} `json:"links,omitempty"`
}
