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

type MovieAnimeResources struct {
	LivechartID   int    `json:"livechartID"`
	AnimePlanetID string `json:"animePlanetID"`
	AnisearchID   int    `json:"anisearchID"`
	AnidbID       int    `json:"anidbID"`
	KitsuID       int    `json:"kitsuID"`
	MalID         int    `json:"malID"`
	NotifyMoeID   string `json:"notifyMoeID"`
	AnilistID     int    `json:"anilistID"`
	TVDbID        int    `json:"TVDBID"`
	IMDbID        string `json:"IMDBID"`
	TMDbID        int    `json:"TMDBID"`
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

type AnimeMovie struct {
	OriginalTitle       string              `json:"originalTitle"`
	Aired               string              `json:"aired"`
	ReleaseYear         int                 `json:"releaseYear"`
	Rating              string              `json:"rating"`
	Runtime             string              `json:"runtime"`
	PortraitPoster      string              `json:"portraitPoster"`
	PortraitBlurHash    string              `json:"portraitBlurHash"`
	LandscapePoster     string              `json:"landscapePoster"`
	LandscapeBlurHash   string              `json:"landscapeBlurHash"`
	AnimeResources      MovieAnimeResources `json:"animeResources"`
	Titles              Titles              `json:"titles"`
	Genres              []string            `json:"genres"`
	Studios             []string            `json:"studios"`
	ProductionCompanies []string            `json:"productionCompanies"`
	Tags                []string            `json:"tags"`
	Posters             []Image             `json:"posters"`
	Backdrops           []Image             `json:"backdrops"`
	Logos               []Image             `json:"logos"`
	Trailers            []Trailer           `json:"trailers"`
	AnimeMetas          []MetaData          `json:"animeMetas"`
}

type SerieAnimeResources struct {
	LivechartID   int    `json:"livechartID"`
	AnimePlanetID string `json:"animePlanetID"`
	AnisearchID   int    `json:"anisearchID"`
	AnidbID       int    `json:"anidbID"`
	KitsuID       int    `json:"kitsuID"`
	MalID         int    `json:"malID"`
	NotifyMoeID   string `json:"notifyMoeID"`
	AnilistID     int    `json:"anilistID"`
	SeasonTVDbID  int    `json:"seasonTVDBID"`
	SeasonTMDbID  int    `json:"seasonTMDBID"`
	Type          string `json:"type"`
}

type Season struct {
	OriginalTitle       string              `json:"seasonOriginalTitle"`
	Aired               string              `json:"aired"`
	ReleaseYear         int                 `json:"releaseYear"`
	Rating              string              `json:"rating"`
	PortraitPoster      string              `json:"portraitPoster"`
	PortraitBlurHash    string              `json:"portraitBlurHash"`
	AnimeResources      SerieAnimeResources `json:"animeResources"`
	Titles              Titles              `json:"titles"`
	Genres              []string            `json:"genres"`
	Studios             []string            `json:"studios"`
	ProductionCompanies []string            `json:"productionCompanies"`
	Tags                []string            `json:"tags"`
	Posters             []Image             `json:"posters"`
	Trailers            []Trailer           `json:"trailers"`
	AnimeMetas          []MetaData          `json:"animeMetas"`
}

type AnimeSerie struct {
	SerieMalID        int       `json:"serieMalID"`
	SerieName         string    `json:"serieName"`
	SerieTVDbID       int       `json:"serieTVDbID"`
	SerieTMDbID       int       `json:"serieTMDbID"`
	Aired             string    `json:"aired"`
	PortraitPoster    string    `json:"portraitPoster"`
	PortraitBlurHash  string    `json:"portraitBlurHash"`
	LandscapePoster   string    `json:"landscapePoster"`
	LandscapeBlurHash string    `json:"landscapeBlurHash"`
	Posters           []Image   `json:"posters"`
	Backdrops         []Image   `json:"backdrops"`
	Logos             []Image   `json:"logos"`
	Trailers          []Trailer `json:"trailers"`
	Season            Season    `json:"season"`
}

type AnimeEpisode struct {
	OriginalTitle      string     `json:"episodeOriginalTitle"`
	Aired              string     `json:"aired"`
	Rating             string     `json:"rating"`
	Runtime            string     `json:"runtime"`
	EpisodeNumber      uint       `json:"episodeNumber"`
	ThumbnailsPoster   string     `json:"thumbnailsPoster"`
	ThumbnailsBlurHash string     `json:"thumbnailsBlurHash"`
	EpisodeMetas       []MetaData `json:"episodeMetas"`
}

type Torrent struct {
	Name     string `json:"name"`
	InfoHash string `json:"info_hash"`
	Leechers string `json:"leechers"`
	Seeders  string `json:"seeders"`
	Size     string `json:"size"`
}

type PirateBayData struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	InfoHash string `json:"info_hash"`
	Leechers string `json:"leechers"`
	Seeders  string `json:"seeders"`
	NumFiles string `json:"num_files"`
	Size     string `json:"size"`
	Username string `json:"username"`
	Added    string `json:"added"`
	Status   string `json:"status"`
	Category string `json:"category"`
	Imdb     string `json:"imdb"`
}

type Iframe struct {
	Link     string `json:"link"`
	Referer  string `json:"referer"`
	Type     string `json:"type"`
	Quality  string `json:"quality"`
	Language string `json:"language"`
}
