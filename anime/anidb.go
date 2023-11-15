package anime

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type AniDB struct {
	XMLName      xml.Name `xml:"anime" json:"anime,omitempty"`
	Text         string   `xml:",chardata" json:"text,omitempty"`
	ID           string   `xml:"id,attr" json:"id,omitempty"`
	Restricted   string   `xml:"restricted,attr" json:"restricted,omitempty"`
	Type         string   `xml:"type"`
	Episodecount string   `xml:"episodecount"`
	Startdate    string   `xml:"startdate"`
	Enddate      string   `xml:"enddate"`
	Titles       struct {
		Text  string `xml:",chardata" json:"text,omitempty"`
		Title []struct {
			Text string `xml:",chardata" json:"text,omitempty"`
			Lang string `xml:"lang,attr" json:"lang,omitempty"`
			Type string `xml:"type,attr" json:"type,omitempty"`
		} `xml:"title" json:"title,omitempty"`
	} `xml:"titles" json:"titles,omitempty"`
	Relatedanime struct {
		Text  string `xml:",chardata" json:"text,omitempty"`
		Anime []struct {
			Text string `xml:",chardata" json:"text,omitempty"`
			ID   string `xml:"id,attr" json:"id,omitempty"`
			Type string `xml:"type,attr" json:"type,omitempty"`
		} `xml:"anime" json:"anime,omitempty"`
	} `xml:"relatedanime" json:"relatedanime,omitempty"`
	Similaranime struct {
		Text  string `xml:",chardata" json:"text,omitempty"`
		Anime []struct {
			Text     string `xml:",chardata" json:"text,omitempty"`
			ID       string `xml:"id,attr" json:"id,omitempty"`
			Approval string `xml:"approval,attr" json:"approval,omitempty"`
			Total    string `xml:"total,attr" json:"total,omitempty"`
		} `xml:"anime" json:"anime,omitempty"`
	} `xml:"similaranime" json:"similaranime,omitempty"`
	Recommendations struct {
		Text           string `xml:",chardata" json:"text,omitempty"`
		Total          string `xml:"total,attr" json:"total,omitempty"`
		Recommendation []struct {
			Text string `xml:",chardata" json:"text,omitempty"`
			Type string `xml:"type,attr" json:"type,omitempty"`
			Uid  string `xml:"uid,attr" json:"uid,omitempty"`
			Br   string `xml:"br"`
		} `xml:"recommendation" json:"recommendation,omitempty"`
	} `xml:"recommendations" json:"recommendations,omitempty"`
	URL      string `xml:"url"`
	Creators struct {
		Text string `xml:",chardata" json:"text,omitempty"`
		Name []struct {
			Text string `xml:",chardata" json:"text,omitempty"`
			ID   string `xml:"id,attr" json:"id,omitempty"`
			Type string `xml:"type,attr" json:"type,omitempty"`
		} `xml:"name" json:"name,omitempty"`
	} `xml:"creators" json:"creators,omitempty"`
	Description string `xml:"description"`
	Ratings     struct {
		Text      string `xml:",chardata" json:"text,omitempty"`
		Permanent struct {
			Text  string `xml:",chardata" json:"text,omitempty"`
			Count string `xml:"count,attr" json:"count,omitempty"`
		} `xml:"permanent" json:"permanent,omitempty"`
		Temporary struct {
			Text  string `xml:",chardata" json:"text,omitempty"`
			Count string `xml:"count,attr" json:"count,omitempty"`
		} `xml:"temporary" json:"temporary,omitempty"`
		Review struct {
			Text  string `xml:",chardata" json:"text,omitempty"`
			Count string `xml:"count,attr" json:"count,omitempty"`
		} `xml:"review" json:"review,omitempty"`
	} `xml:"ratings" json:"ratings,omitempty"`
	Picture   string `xml:"picture"`
	Resources struct {
		Text     string `xml:",chardata" json:"text,omitempty"`
		Resource []struct {
			Text           string `xml:",chardata" json:"text,omitempty"`
			Type           string `xml:"type,attr" json:"type,omitempty"`
			Externalentity struct {
				Text       string   `xml:",chardata" json:"text,omitempty"`
				Identifier []string `xml:"identifier"`
				URL        string   `xml:"url"`
			} `xml:"externalentity" json:"externalentity,omitempty"`
		} `xml:"resource" json:"resource,omitempty"`
	} `xml:"resources" json:"resources,omitempty"`
	Tags struct {
		Text string `xml:",chardata" json:"text,omitempty"`
		Tag  []struct {
			Text          string `xml:",chardata" json:"text,omitempty"`
			ID            string `xml:"id,attr" json:"id,omitempty"`
			Weight        string `xml:"weight,attr" json:"weight,omitempty"`
			Localspoiler  string `xml:"localspoiler,attr" json:"localspoiler,omitempty"`
			Globalspoiler string `xml:"globalspoiler,attr" json:"globalspoiler,omitempty"`
			Verified      string `xml:"verified,attr" json:"verified,omitempty"`
			Update        string `xml:"update,attr" json:"update,omitempty"`
			Parentid      string `xml:"parentid,attr" json:"parentid,omitempty"`
			Infobox       string `xml:"infobox,attr" json:"infobox,omitempty"`
			Name          string `xml:"name"`
			Description   string `xml:"description"`
			Picurl        string `xml:"picurl"`
		} `xml:"tag" json:"tag,omitempty"`
	} `xml:"tags" json:"tags,omitempty"`
	Characters struct {
		Text      string `xml:",chardata" json:"text,omitempty"`
		Character []struct {
			Text   string `xml:",chardata" json:"text,omitempty"`
			ID     string `xml:"id,attr" json:"id,omitempty"`
			Type   string `xml:"type,attr" json:"type,omitempty"`
			Update string `xml:"update,attr" json:"update,omitempty"`
			Rating struct {
				Text  string `xml:",chardata" json:"text,omitempty"`
				Votes string `xml:"votes,attr" json:"votes,omitempty"`
			} `xml:"rating" json:"rating,omitempty"`
			Name          string `xml:"name"`
			Gender        string `xml:"gender"`
			Charactertype struct {
				Text string `xml:",chardata" json:"text,omitempty"`
				ID   string `xml:"id,attr" json:"id,omitempty"`
			} `xml:"charactertype" json:"charactertype,omitempty"`
			Description string `xml:"description"`
			Picture     string `xml:"picture"`
			Seiyuu      struct {
				Text    string `xml:",chardata" json:"text,omitempty"`
				ID      string `xml:"id,attr" json:"id,omitempty"`
				Picture string `xml:"picture,attr" json:"picture,omitempty"`
			} `xml:"seiyuu" json:"seiyuu,omitempty"`
		} `xml:"character" json:"character,omitempty"`
	} `xml:"characters" json:"characters,omitempty"`
	Episodes struct {
		Text    string `xml:",chardata" json:"text,omitempty"`
		Episode []struct {
			Text   string `xml:",chardata" json:"text,omitempty"`
			ID     string `xml:"id,attr" json:"id,omitempty"`
			Update string `xml:"update,attr" json:"update,omitempty"`
			Epno   struct {
				Text string `xml:",chardata" json:"text,omitempty"`
				Type string `xml:"type,attr" json:"type,omitempty"`
			} `xml:"epno" json:"epno,omitempty"`
			Length string `xml:"length"`
			Rating struct {
				Text  string `xml:",chardata" json:"text,omitempty"`
				Votes string `xml:"votes,attr" json:"votes,omitempty"`
			} `xml:"rating" json:"rating,omitempty"`
			Title []struct {
				Text string `xml:",chardata" json:"text,omitempty"`
				Lang string `xml:"lang,attr" json:"lang,omitempty"`
			} `xml:"title" json:"title,omitempty"`
			Airdate string `xml:"airdate"`
			Summary string `xml:"summary"`
		} `xml:"episode" json:"episode,omitempty"`
	} `xml:"episodes" json:"episodes,omitempty"`
}

func (am *AnimeScraper) GetAniDBData(id int) (AniDB, error) {
	url := fmt.Sprintf("http://api.anidb.net:9001/httpapi?client=golangtest&clientver=1&protover=1&request=anime&aid=%d", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return AniDB{}, fmt.Errorf("aniDB creating request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

	resp, err := am.HTTP.Do(req)
	if err != nil {
		return AniDB{}, fmt.Errorf("aniDB fetching XML: %v", err)
	}
	defer resp.Body.Close()

	if !strings.Contains(resp.Header.Get("Content-Type"), "text/xml") {
		return AniDB{}, fmt.Errorf("aniDB content is not XML")
	}

	var animeXML AniDB
	err = xml.NewDecoder(resp.Body).Decode(&animeXML)
	if err != nil {
		return AniDB{}, fmt.Errorf("aniDB error decoding XML: %v", err)
	}

	return animeXML, nil
}

type Link struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (am *AnimeScraper) getAniDBID(links []Link) (int, error) {
	anidbPattern := regexp.MustCompile(`.*[&?]aid=(\d+).*`)

	for _, link := range links {
		if strings.Contains(strings.ToLower(link.Name), "anidb") {
			matches := anidbPattern.FindStringSubmatch(link.URL)
			if len(matches) >= 2 {
				id, err := strconv.ParseInt(matches[1], 0, 0)
				if err != nil {
					return 0, fmt.Errorf("AniDB URL with invalid id")
				}
				return int(id), nil
			}
			return 0, fmt.Errorf("AniDB URL found, but ID extraction failed")
		}
	}

	return 0, fmt.Errorf("AniDB URL not found in the provided links")
}
