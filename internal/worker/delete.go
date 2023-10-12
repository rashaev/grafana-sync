package worker

import (
	"log"

	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/rashaev/grafana-sync/internal/keycloak"
)

func RemoveFromOrg(kcloak *keycloak.Keycloak, gclient *gapi.Client, groupRgxRO, groupRgxRW string) {
	orgs, _ := gclient.Orgs()

	for _, org := range orgs {
		// Skip Main Org
		if org.ID == 1 {
			continue
		}

		// Get list of users in Grafana Organization
		listOrgUsers, _ := gclient.OrgUsers(org.ID)

		if len(listOrgUsers) == 0 {
			log.Printf("Thera are no users in org %s, skipping", org.Name)
			return
		}

		var monitorProjectGroups []keycloak.Group

		kcloakROGroup, _ := kcloak.Groups(groupRgxRO + org.Name)
		kcloakRWGroup, _ := kcloak.Groups(groupRgxRW + org.Name)

		monitorProjectGroups = append(monitorProjectGroups, kcloakROGroup...)
		monitorProjectGroups = append(monitorProjectGroups, kcloakRWGroup...)

		allUsers := userFromGroups(kcloak, monitorProjectGroups)

		for _, orgUser := range listOrgUsers {
			// Skip admin user with ID = 1
			if orgUser.UserID == 1 {
				continue
			} else if _, ok := allUsers[orgUser.Email]; !ok {
				err := gclient.RemoveOrgUser(org.ID, orgUser.UserID)
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("User %s was deleted from Org %s", orgUser.Email, org.Name)
				}

			}
		}
	}
}
