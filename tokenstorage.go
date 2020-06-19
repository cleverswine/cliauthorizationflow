package cliauthorizationflow

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"

	"golang.org/x/oauth2"
)

// TokenStorage is used to cache and retrieve tokens
type TokenStorage interface {
	Get(key string) (*oauth2.Token, error)
	Save(key string, token *oauth2.Token) error
}

// DefaultTokenStorage stores tokens in a user's Home directory
type DefaultTokenStorage struct {
	fsPath string
}

// NewDefaultTokenStorage returns a configured DefaultTokenStorage with a storage path of ~/.config/{appName}
func NewDefaultTokenStorage(appName string) *DefaultTokenStorage {
	homeDir := ""
	osUser, err := user.Current()
	if err == nil {
		homeDir = osUser.HomeDir
	}
	p := path.Join(homeDir, ".config", appName)
	fmt.Printf("TOK_DEBUG:: setting token storage path to %s\n", p)
	return &DefaultTokenStorage{fsPath: p}
}

// Get gets a token from the filesystem
func (t *DefaultTokenStorage) Get(key string) (*oauth2.Token, error) {
	fn := path.Join(t.fsPath, key)
	fmt.Printf("TOK_DEBUG:: looking for token in %s\n", fn)
	if _, err := os.Stat(fn); err != nil {
		fmt.Printf("TOK_DEBUG:: no token found in %s\n", fn)
		return nil, nil
	}
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	if err != nil {
		return nil, err
	}
	fmt.Printf("TOK_DEBUG:: got token from %s\n", fn)
	return token, nil
}

// Save saves a token to the filesystem
func (t *DefaultTokenStorage) Save(key string, token *oauth2.Token) error {
	if _, err := os.Stat(t.fsPath); err != nil {
		fmt.Printf("TOK_DEBUG:: making dir %s\n", t.fsPath)
		err = os.MkdirAll(t.fsPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	fn := path.Join(t.fsPath, key)
	fmt.Printf("TOK_DEBUG:: saving token in %s\n", fn)
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}
