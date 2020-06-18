package cliauthorizationflow

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"
)

// TokenStorage is used to cache and retrieve tokens
// WIP - this doesn't yet support keying off of user or auth service
type TokenStorage interface {
	Get(string) (*oauth2.Token, error)
	Save(string, *oauth2.Token) error
}

// Config describes a typical 3-legged OAuth2 flow, with both the
// client application information and the server's endpoint URLs.
type Config struct {
	// ClientID is the application's ID.
	ClientID string
	// ClientSecret is the application's secret.
	ClientSecret string
	// AuthorizationURL contains the resource server's authorize endpoint URL
	AuthorizationURL string
	// TokenURL contains the resource server's token endpoint URL
	TokenURL string
	// Scope specifies optional requested permissions.
	Scopes []string
	// CallbackPort specifies which local port to use for the auth callback (default: 8080)
	CallbackPort int
}

func (c *Config) oauth2Config() *oauth2.Config {
	if c.CallbackPort == 0 {
		c.CallbackPort = 8080
	}
	return &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.AuthorizationURL,
			TokenURL: c.TokenURL,
		},
		RedirectURL: fmt.Sprintf("http://localhost:%d/callback", c.CallbackPort),
		Scopes:      c.Scopes,
	}
}

// Client contains an authorized http client.
// That is, the auth token is automatically set on all request headers.
type Client struct {
	*http.Client
	token           *oauth2.Token
	tokenStorage    TokenStorage
	tokenStorageKey string
}

// NewClient creates a new client with the specified authorization parameters
func NewClient(ctx context.Context, config *Config, tokenStorage TokenStorage) (*Client, error) {
	oauthConfig := config.oauth2Config()
	// try cache first
	var storageKey string
	if tokenStorage != nil {
		// for now, just use auth hostname as cache key
		authHost, err := url.Parse(oauthConfig.Endpoint.AuthURL)
		if err != nil {
			return nil, err
		}
		storageKey = authHost.Hostname()
		token, err := tokenStorage.Get(storageKey)
		if err != nil {
			return nil, err
		}
		// TODO: is the refresh token still valid?
		if token != nil {
			return &Client{
				Client:          oauthConfig.Client(ctx, token),
				token:           token,
				tokenStorage:    tokenStorage,
				tokenStorageKey: storageKey,
			}, nil
		}
	}
	// get via authorization flow
	state := randStringBytes(40)
	log.Printf("\nto continue, please log in and authorize this application at: \n%s\n\n", oauthConfig.AuthCodeURL(state))
	// start an http server and wait for callback
	queryValCh := make(chan url.Values)
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Thank you"))
		queryValCh <- r.URL.Query()
	})
	go http.ListenAndServe(":8080", nil)
	queryVals := <-queryValCh
	// verify response
	code := queryVals.Get("code")
	if code == "" {
		return nil, errors.New("didn't get access code")
	}
	if actualState := queryVals.Get("state"); actualState != state {
		return nil, errors.New("redirect state parameter doesn't match")
	}
	// exchange code for token
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return &Client{
		Client:          oauthConfig.Client(ctx, token),
		token:           token,
		tokenStorage:    tokenStorage,
		tokenStorageKey: storageKey,
	}, nil
}

// Persist will save the current token to storage
func (a *Client) Persist() {
	if a.token == nil || a.tokenStorage == nil || a.tokenStorageKey == "" {
		return
	}
	// token may have been updated via refresh token
	a.tokenStorage.Save(a.tokenStorageKey, a.token)
}

func randStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789_"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
