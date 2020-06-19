package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	auth "github.com/cleverswine/cliauthorizationflow"
)

const (
	SpotifyAuthURL  = "https://accounts.spotify.com/authorize"
	SpotifyTokenURL = "https://accounts.spotify.com/api/token"
	SpotifyAPIBase  = "https://api.spotify.com/v1"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := &auth.Config{
		ClientID:         os.Getenv("SPOTIFY_ID"),
		ClientSecret:     os.Getenv("SPOTIFY_SECRET"),
		AuthorizationURL: SpotifyAuthURL,
		TokenURL:         SpotifyTokenURL,
		Scopes:           []string{"user-top-read"}}

	client, err := auth.NewClient(ctx, config, auth.NewDefaultTokenStorage("spotify-cli"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Persist()

	// get my top tracks
	resp, err := client.Get(SpotifyAPIBase + "/me/top/tracks")
	if err != nil {
		log.Fatal(err)
	}
	// write them out to a json file
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("top-tracks.json", body, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
