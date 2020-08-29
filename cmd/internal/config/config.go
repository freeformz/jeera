package config

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"path/filepath"

	"github.com/dghubble/oauth1"
)

func settingsData() (io.ReadCloser, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	return os.Open(filepath.Join(u.HomeDir, ".config", "jeera", "settings"))
}

// OAuth returns the oauth1 config and token from
// ~/.config/jeera/settings
func OAuth() (oauth1.Config, *oauth1.Token, error) {
	var cfg oauth1.Config

	f, err := settingsData()
	if err != nil {
		return cfg, nil, err
	}
	defer f.Close()
	d := json.NewDecoder(f)
	var s struct {
		Jira struct {
			URL string
		}
		OAuth struct {
			Token         string
			Secret        string
			PEM           string
			ConsumerKey   string
			ConsumerSeret string
		}
	}
	if err = d.Decode(&s); err != nil {
		return cfg, nil, fmt.Errorf("decoding json: %w", err)
	}
	cfg.ConsumerKey = s.OAuth.ConsumerKey
	cfg.ConsumerSecret = s.OAuth.ConsumerSeret
	cfg.Endpoint = oauth1.Endpoint{
		RequestTokenURL: path.Join(s.Jira.URL, "plugins/servlet/oauth/request-token"),
		AuthorizeURL:    path.Join(s.Jira.URL, "plugins/servlet/oauth/authorize"),
		AccessTokenURL:  path.Join(s.Jira.URL, "plugins/servlet/oauth/access-token"),
	}

	block, _ := pem.Decode([]byte(s.OAuth.PEM))
	if block == nil {
		return cfg, nil, fmt.Errorf("block is nil, OAuth.PEM is bad?")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return cfg, nil, fmt.Errorf("parsing private key: %w", err)
	}

	cfg.Signer = &oauth1.RSASigner{
		PrivateKey: key,
	}

	return cfg, oauth1.NewToken(s.OAuth.Token, s.OAuth.Secret), nil
}

// JiraURL to use
func JiraURL() (string, error) {
	f, err := settingsData()
	if err != nil {
		return "", err
	}
	defer f.Close()
	d := json.NewDecoder(f)
	var s struct {
		Jira struct {
			URL string
		}
	}
	if err = d.Decode(&s); err != nil {
		return "", fmt.Errorf("decoding json: %w", err)
	}
	return s.Jira.URL, nil
}

// JiraCustomEpicFieldID to use
func JiraCustomEpicFieldID() (string, error) {
	f, err := settingsData()
	if err != nil {
		return "", err
	}
	defer f.Close()
	d := json.NewDecoder(f)
	var s struct {
		Jira struct {
			EpicCustomFieldID string
		}
	}
	if err = d.Decode(&s); err != nil {
		return "", fmt.Errorf("decoding json: %w", err)
	}
	return s.Jira.EpicCustomFieldID, nil
}
