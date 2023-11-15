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

type Anime struct {
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
