package anime

import (
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	go fetchAnimeTitles()
	go fetchAnimeResources()
}

var (
	GlobalAniDBTitles    AniDBTitles
	GlobalAnimeResources []AnimeResources
)

type AniTitle struct {
	Type  string `xml:"type,attr" json:"type"`
	Lang  string `xml:"lang,attr" json:"lang"`
	Value string `xml:",chardata" json:"title"`
}

type AnidbAnime struct {
	Aid    int        `xml:"aid,attr" json:"ID"`
	Titles []AniTitle `xml:"title" json:"titles"`
}

type AniDBTitles struct {
	XMLName xml.Name     `xml:"animetitles" json:"-"`
	Animes  []AnidbAnime `xml:"anime" json:"animeList,omitempty"`
}

func fetchAnimeTitles() {
	for {
		url := "http://anidb.net/api/anime-titles.xml.gz"

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error fetching XML:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusFound {
			redirectURL := resp.Header.Get("Location")
			resp.Body.Close()

			req, err := http.NewRequest("GET", redirectURL, nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}

			req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println("Error fetching redirected XML:", err)
				return
			}
			defer resp.Body.Close()
		}

		if resp.Header.Get("Content-Type") != "application/gzip" {
			fmt.Println("Content is not gzip")
			continue
		}

		gzipFile, err := os.Create("anidb-titles.xml.gz")
		if err != nil {
			fmt.Println("Error creating gzip file:", err)
			continue
		}
		defer gzipFile.Close()

		_, err = io.Copy(gzipFile, resp.Body)
		if err != nil {
			fmt.Println("Error writing gzip file:", err)
			continue
		}

		// Read the saved gzip file
		gz, err := os.Open("anidb-titles.xml.gz")
		if err != nil {
			fmt.Println("Error opening gzip file:", err)
			continue
		}
		defer gz.Close()

		reader, err := gzip.NewReader(gz)
		if err != nil {
			fmt.Println("Error reading gzip file:", err)
			continue
		}
		defer reader.Close()

		xmlData, err := io.ReadAll(reader)
		if err != nil {
			fmt.Println("Error reading XML data:", err)
			time.Sleep(5 * time.Minute)
			continue
		}

		// Parse the xml
		var animetitlesXML AniDBTitles
		err = xml.Unmarshal(xmlData, &animetitlesXML)
		if err != nil {
			fmt.Println("Error parsing XML:", err)
			time.Sleep(5 * time.Minute)
			continue
		}

		file, err := os.Create("anidb-titles.json")
		if err != nil {
			fmt.Println("Error creating JSON file:", err)
			return
		}
		defer file.Close()

		jsonData, err := json.Marshal(&animetitlesXML)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}

		err = json.Unmarshal(jsonData, &GlobalAniDBTitles)
		if err != nil {
			log.Fatal("Failed to unmarshal data")
			return
		}

		_, err = file.Write(jsonData)
		if err != nil {
			fmt.Println("Error writing JSON data:", err)
		}

		time.Sleep(5 * 24 * time.Hour)
	}
}

type AnimeResourcesResponse struct {
	LivechartID   int             `json:"livechart_id"`
	AnimePlanetID json.RawMessage `json:"anime-planet_id"`
	AnisearchID   int             `json:"anisearch_id"`
	AnidbID       int             `json:"anidb_id"`
	KitsuID       int             `json:"kitsu_id"`
	MalID         int             `json:"mal_id"`
	NotifyMoeID   string          `json:"notify.moe_id"`
	AnilistID     int             `json:"anilist_id"`
	TheTVdbID     int             `json:"thetvdb_id"`
	IMDbID        string          `json:"imdb_id"`
	TMDdID        json.RawMessage `json:"themoviedb_id"`
	Type          string          `json:"type"`
}

type AnimeResources struct {
	AnidbID int                    `json:"anidb_id"`
	Data    AnimeResourcesResponse `json:"anime_resources"`
}

func fetchAnimeResources() {
	for {
		url := "https://raw.githubusercontent.com/Fribb/anime-lists/master/anime-list-full.json"

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching anime resources:", err)
			return
		}
		defer resp.Body.Close()

		var animeResourcesList []AnimeResourcesResponse

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&animeResourcesList)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			continue
		}

		for _, res := range animeResourcesList {
			if res.AnidbID != 0 {
				GlobalAnimeResources = append(GlobalAnimeResources, AnimeResources{
					AnidbID: res.AnidbID,
					Data:    res,
				})
			}
		}

		jsonFile, err := os.Create("anime-resources.json")
		if err != nil {
			fmt.Println("Error creating anime-resources file:", err)
			continue
		}
		defer jsonFile.Close()

		encoder := json.NewEncoder(jsonFile)
		err = encoder.Encode(GlobalAnimeResources)
		if err != nil {
			fmt.Println("Error encoding anime-resources JSON:", err)
			continue
		}

		time.Sleep(5 * 24 * time.Hour)
	}
}
