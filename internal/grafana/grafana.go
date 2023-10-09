package grafana

import (
	"net/url"

	gapi "github.com/grafana/grafana-api-golang-client"
)

type SyncUser struct {
}

func New(url url.URL, config gapi.Config) (*gapi.Client, error) {
	return gapi.New(url.String(), config)
}

func NewConfig(username, password string) gapi.Config {
	return gapi.Config{
		BasicAuth: url.UserPassword(username, password),
	}
}
