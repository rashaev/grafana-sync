package worker

import (
	"log"

	gapi "github.com/grafana/grafana-api-golang-client"
)

func DeleteGrafanaUser(gclient *gapi.Client) {
	orgs, err := gclient.Orgs()
	if err != nil {
		log.Println(err)
	}

	for _, org := range orgs {
		users, _ := gclient.OrgUsers(org.ID)
		log.Printf("Organization %s contains user %v", org.Name, users)
	}
}
