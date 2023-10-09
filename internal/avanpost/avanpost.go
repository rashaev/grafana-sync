package avanpost

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"golang.org/x/oauth2"
)

func New(ctx context.Context, url url.URL, clientid, clientsecret, username, password string) (*Avanpost, error) {
	conf := &oauth2.Config{
		ClientID:     clientid,
		ClientSecret: clientsecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  url.String() + "/oauth2/auth",
			TokenURL: url.String() + "/oauth2/token",
		},
	}

	token, err := conf.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		return nil, err
	}

	client := conf.Client(ctx, token)

	return &Avanpost{
		URL:          url,
		Username:     username,
		Password:     password,
		ClientID:     clientid,
		ClientSecret: clientsecret,
		Client:       client,
	}, nil
}

func (a *Avanpost) Groups(groupRegexp string) (MonitorGroups, error) {
	req, err := http.NewRequest("GET", a.URL.String()+"/api/v1/groups?pageSize=1000", nil)
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

	var data GroupsInfo

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	monitoringGroups := make(MonitorGroups)

	re := regexp.MustCompile(groupRegexp)
	for _, v := range data.Resources {
		if re.MatchString(v.Name) {
			monitoringGroups[v.ID] = v.Name
		}

	}
	return monitoringGroups, nil

}

func (a *Avanpost) UsersInGroup(groupId string) ([]User, error) {
	req, err := http.NewRequest("GET", a.URL.String()+"/api/v1/groups/"+groupId+"/users?pageSize=1000", nil)
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

	var data UsersInfo

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	var users []User

	users = append(users, data.Resources...)

	return users, nil
}
