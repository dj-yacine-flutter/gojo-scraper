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
	AnimeMetas          []MetaData     `json:"animeMetas"`
}

type Movie struct {
	OriginalTitle     string     `json:"originalTitle"`
	Aired             string     `json:"aired"`
	ReleaseYear       int        `json:"releaseYear"`
	Rating            string     `json:"rating"`
	Duration          string     `json:"duration"`
	PortriatPoster    string     `json:"portriatPoster"`
	PortriatBlurHash  string     `json:"portriatBlurHash"`
	LandscapePoster   string     `json:"landscapePoster"`
	LandscapeBlurHash string     `json:"landscapeBlurHash"`
	AnimeMetas        []MetaData `json:"animeMetas"`
}
