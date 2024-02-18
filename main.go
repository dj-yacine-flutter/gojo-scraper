package main

import (
	"net/http"
	"net/url"
	"os"

	tmdb "github.com/cyruzin/golang-tmdb"
	anime "github.com/dj-yacine-flutter/gojo-scraper/anime"
	"github.com/dj-yacine-flutter/gojo-scraper/tvdb"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	*anime.AnimeScraper
}

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	tmdbClient, err := tmdb.Init("cd74b33da8b164701b53cc22db416aea")
	if err != nil {
		log.Fatal().Err(err)
		return
	}
	tmdbClient.SetClientAutoRetry()
	tmdbClient.SetAlternateBaseURL()

	proxy, err := url.Parse("http://127.0.0.1:8118")
	if err != nil {
		log.Fatal().Err(err)
		return
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}

	Oimg := "https://www.themoviedb.org/t/p/original"
	Dimg := "https://www.themoviedb.org/t/p/w92"

	tvdbClient := tvdb.NewClient(httpClient)
	err = tvdbClient.Login(&tvdb.AuthenticationRequest{
		ApiKey: "84f7322d-6bfa-4a67-b4e7-855b56db2239",
		Pin:    "",
	})
	if err != nil {
		log.Fatal().Err(err)
		return
	}

	server := Server{
		AnimeScraper: anime.NewAnimeScraper(tmdbClient, httpClient, tvdbClient, &log.Logger, Oimg, Dimg),
	}

	http.HandleFunc("/anime/movie", server.GetAnimeMovie)
	http.HandleFunc("/anime/iframe/movie", server.GetAnimeMovieIframes)
	http.HandleFunc("/anime/serie", server.GetAnimeSerie)
	http.HandleFunc("/anime/iframe/serie", server.GetAnimeSerieIframes)
	http.HandleFunc("/anime/episode", server.GetAnimeEpisode)

	log.Info().Msg("Server is running on port 3333")
	err = http.ListenAndServe(":3333", nil)
	if err != nil {
		log.Err(err).Msg("cannot start the server")
	}
}
