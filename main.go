package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	tmdb "github.com/cyruzin/golang-tmdb"
	anime "github.com/dj-yacine-flutter/gojo-scraper/anime"
	"github.com/dj-yacine-flutter/gojo-scraper/movie"
)

type Server struct {
	*anime.AnimeScraper
	*movie.MovieTMDB
}

func main() {

	tmdbClient, err := tmdb.Init("cd74b33da8b164701b53cc22db416aea")
	if err != nil {
		log.Fatal(err)
		return
	}
	tmdbClient.SetClientAutoRetry()
	tmdbClient.SetAlternateBaseURL()

	proxy, err := url.Parse("http://127.0.0.1:8118")
	if err != nil {
		log.Fatal(err)
		return
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}

	Oimg := "https://www.themoviedb.org/t/p/original"
	Dimg := "https://www.themoviedb.org/t/p/w92"

	server := Server{
		AnimeScraper: anime.NewAnimeScraper(tmdbClient, httpClient, Oimg, Dimg),
		MovieTMDB:    movie.NewMovieTMDB(tmdbClient, httpClient, Oimg, Dimg),
	}

	http.HandleFunc("/anime", server.GetAnime)
	http.HandleFunc("/movie", server.GetMovie)

	fmt.Println("Server is running on port 3333")
	http.ListenAndServe(":3333", nil)
}
