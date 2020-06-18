# Authorization Flow for Go CLI

This package makes it easier to perform the OAuth2 Authorization Flow from a command-line application.

*This is a beta/prototype level package and shouldn't be used in production without understanding what this code does.*

## Example

This example uses the Spotify API

```go
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
        Scopes:           []string{"user-top-read"},
    }

    client, err := auth.NewClient(ctx, config, nil)
    if err != nil {
        log.Fatal(err)
    }

    // get my top tracks
    resp, err := client.Get(SpotifyAPIBase + "/me/top/tracks")
    if err != nil {
        log.Fatal(err)
    }

    // do stuff with the resp...
}
```

## Token storage

You'll want to implement the `TokenStorage` interface and pass that to the `NewClient(...)` method in order to persist tokens. Otherwise the user will have to open a web browser and authenticate every time.

**Be very careful with storing tokens if you plan to authorize multiple users.**
