package keycloak

import (
	"net/http"
	"net/url"
)

type Group struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Path      string  `json:"path"`
	SubGroups []Group `json:"subgroups"`
}

type User struct {
	ID               string `json:"id"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	Lastname         string `json:"lastName"`
	Firstname        string `json:"firstName"`
	Enabled          bool   `json:"enabled"`
	EmailVerified    bool   `json:"emailVerified"`
	CreatedTimestamp int64  `json:"createdTimestamp"`
}

type Keycloak struct {
	URL          url.URL
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
	Client       *http.Client
}
