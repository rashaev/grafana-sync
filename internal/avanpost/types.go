package avanpost

import (
	"net/http"
	"net/url"
	"time"
)

type Group struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Info        struct {
		TotalUsers int `json:"total_users"`
		TotalApps  int `json:"total_apps"`
	}
}

type GroupsInfo struct {
	StartIndex   int     `json:"startIndex"`
	ItemsPerPage int     `json:"itemsPerPage"`
	TotalResults int     `json:"totalResults"`
	Resources    []Group `json:"resources"`
}

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
	Active    bool   `json:"active"`
}

type UsersInfo struct {
	StartIndex   int    `json:"startIndex"`
	ItemsPerPage int    `json:"itemsPerPage"`
	TotalResults int    `json:"totalResults"`
	Resources    []User `json:"resources"`
}

type Avanpost struct {
	URL          url.URL
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
	Client       *http.Client
}

type MonitorGroups map[string]string
