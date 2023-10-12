package keycloak

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

func New(ctx context.Context, url url.URL, clientid, clientsecret, username, password string) (*Keycloak, error) {
	conf := &oauth2.Config{
		ClientID:     clientid,
		ClientSecret: clientsecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  url.String() + "/realms/master/protocol/openid-connect/auth",
			TokenURL: url.String() + "/realms/master/protocol/openid-connect/token",
		},
	}

	token, err := conf.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		return nil, err
	}

	client := conf.Client(ctx, token)

	return &Keycloak{
		URL:          url,
		Username:     username,
		Password:     password,
		ClientID:     clientid,
		ClientSecret: clientsecret,
		Client:       client,
	}, nil
}

func (a *Keycloak) Groups(groupRegexp string) ([]Group, error) {
	req, err := http.NewRequest("GET", a.URL.String()+"/admin/realms/master/groups?max=1000&search="+groupRegexp, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data []Group

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func (a *Keycloak) UsersInGroup(groupId string) ([]User, error) {
	req, err := http.NewRequest("GET", a.URL.String()+"/admin/realms/master/groups/"+groupId+"/members?briefRepresentation=true&max=1000", nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data []User

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
