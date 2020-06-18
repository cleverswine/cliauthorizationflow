package cliauthorizationflow

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"golang.org/x/oauth2"
)

func TestConfig_oauth2Config(t *testing.T) {
	type fields struct {
		ClientID         string
		ClientSecret     string
		AuthorizationURL string
		TokenURL         string
		Scopes           []string
		CallbackPort     int
	}
	tests := []struct {
		name   string
		fields fields
		want   *oauth2.Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				ClientID:         tt.fields.ClientID,
				ClientSecret:     tt.fields.ClientSecret,
				AuthorizationURL: tt.fields.AuthorizationURL,
				TokenURL:         tt.fields.TokenURL,
				Scopes:           tt.fields.Scopes,
				CallbackPort:     tt.fields.CallbackPort,
			}
			if got := c.oauth2Config(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.oauth2Config() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		ctx    context.Context
		config *Config
		cache  TokenStorage
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.ctx, tt.args.config, tt.args.cache)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Persist(t *testing.T) {
	type fields struct {
		Client   *http.Client
		token    *oauth2.Token
		cache    TokenStorage
		cacheKey string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Client{
				Client:   tt.fields.Client,
				token:    tt.fields.token,
				tokenStorage:    tt.fields.cache,
				tokenStorageKey: tt.fields.cacheKey,
			}
			a.Persist()
		})
	}
}

func Test_randStringBytes(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := randStringBytes(tt.args.n); got != tt.want {
				t.Errorf("randStringBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
