package worker

import (
	"log"

	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/rashaev/grafana-sync/internal/keycloak"
)

func SyncGroups(kcloak *keycloak.Keycloak, gclient *gapi.Client, groupRgx string, grole string) {

	groups, err := kcloak.Groups(groupRgx)
	if err != nil {
		log.Println(err)
	}

	for _, nameKcloakGroup := range groups {
		// Parse project code
		projectCode := parseProjectCode(nameKcloakGroup.Name)

		//Check if the Organization exists in Grafana
		org, err := gclient.OrgByName(projectCode)
		if err != nil {
			log.Printf("Organization %s not found, creating new\n", projectCode)
			orgID, _ := gclient.NewOrg(projectCode)
			org.ID = orgID
			org.Name = projectCode
		}

		log.Printf("synchronizing group %s", nameKcloakGroup.Name)

		// Get list of users in Keycloak group
		users, _ := kcloak.UsersInGroup(nameKcloakGroup.ID)

		if len(users) > 0 {
			for _, user := range users {
				createUser(gclient, user, org, grole)
			}
		} else {
			log.Printf("Group %s does not have users\n", nameKcloakGroup.Name)
		}
	}
}

// func create new user in Grafana
func createUser(gclient *gapi.Client, kcloakUser keycloak.User, gOrg gapi.Org, grole string) {

	// Get list of users in Grafana Organization
	listOrgUsers, _ := gclient.OrgUsers(gOrg.ID)

	// Check if the user already exists in Grafana
	_, err := gclient.UserByEmail(kcloakUser.Email)
	if err != nil {
		log.Printf("User %s not found in Grafana\n", kcloakUser.Email)
		newuser := gapi.User{
			Email:    kcloakUser.Email,
			Login:    kcloakUser.Username,
			Name:     kcloakUser.Firstname + " " + kcloakUser.Lastname,
			Password: kcloakUser.Firstname + kcloakUser.Lastname,
			IsAdmin:  false,
		}

		// Create new user
		_, err := gclient.CreateUser(newuser)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("User %s created\n", kcloakUser.Email)
		}

		//Add user to Organizaion
		errOrgUser := gclient.AddOrgUser(gOrg.ID, kcloakUser.Email, grole)
		if errOrgUser == nil {
			log.Printf("User %s was added to organization %s\n", kcloakUser.Email, gOrg.Name)
		} else {
			log.Println(errOrgUser)
		}

	} else {
		log.Printf("User %s already exists in Grafana", kcloakUser.Email)
		if useridx, ok := orgContainsUser(listOrgUsers, kcloakUser.Email); !ok {

			//Add user that already exits in Grafana to Organization
			errAddOrgUser := gclient.AddOrgUser(gOrg.ID, kcloakUser.Email, grole)
			if errAddOrgUser == nil {
				log.Printf("User %s was added to organization %s\n", kcloakUser.Email, gOrg.Name)
			} else {
				log.Println(errAddOrgUser)
			}
		} else if !checkUserRole(listOrgUsers, useridx, grole) {
			log.Printf("User %s is already member of %s organization but doesn't have correct Role. Updating", kcloakUser.Email, gOrg.Name)
			err := gclient.UpdateOrgUser(gOrg.ID, listOrgUsers[useridx].UserID, grole)
			if err != nil {
				log.Println(err)
			}
		} else {
			log.Printf("User %s is already member of %s organization and has correct Role", kcloakUser.Email, gOrg.Name)
		}
	}
}
