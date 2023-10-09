package worker

import (
	"log"

	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/rashaev/grafana-sync/internal/avanpost"
)

func SyncGroups(avp *avanpost.Avanpost, gclient *gapi.Client, groupRgx string) {

	groups, err := avp.Groups(groupRgx)
	if err != nil {
		log.Println(err)
	}

	for idxAvpGroup, nameAvpGroup := range groups {
		// Parse project code
		projectCode := parseProjectCode(nameAvpGroup)

		//Check if the Organization exists in Grafana
		org, err := gclient.OrgByName(projectCode)
		if err != nil {
			log.Printf("Organization %s not found, creating new\n", projectCode)
			orgID, _ := gclient.NewOrg(projectCode)
			org.ID = orgID
			org.Name = projectCode
		}

		// Get list of users in Avanpost group
		users, _ := avp.UsersInGroup(idxAvpGroup)

		if len(users) > 0 {
			for _, user := range users {
				createUser(gclient, user, org)
			}
		} else {
			log.Printf("Group %s does not have users\n", nameAvpGroup)
		}
	}
}

// func create new user in Grafana
func createUser(gclient *gapi.Client, avpUser avanpost.User, gOrg gapi.Org) {

	// Get list of users in Grafana Organization
	listOrgUsers, _ := gclient.OrgUsers(gOrg.ID)

	// Check if the user already exists in Grafana
	_, err := gclient.UserByEmail(avpUser.Email)
	if err != nil {
		log.Printf("User %s not found in Grafana\n", avpUser.Email)
		newuser := gapi.User{
			Email:    avpUser.Email,
			Login:    avpUser.Username,
			Name:     avpUser.Firstname + " " + avpUser.Lastname,
			Password: avpUser.Firstname + avpUser.Lastname,
			IsAdmin:  false,
		}

		// Create new user
		_, err := gclient.CreateUser(newuser)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("User %s created\n", avpUser.Email)
		}

		//Add user to Organizaion
		errOrgUser := gclient.AddOrgUser(gOrg.ID, avpUser.Email, "Editor")
		if errOrgUser == nil {
			log.Printf("User %s was added to organization %s\n", avpUser.Email, gOrg.Name)
		} else {
			log.Println(errOrgUser)
		}

	} else {
		log.Printf("User %s already exists in Grafana", avpUser.Email)
		if !orgContainsUser(listOrgUsers, avpUser.Email) {

			//Add user that already exits in Grafana to Organization
			errAddOrgUser := gclient.AddOrgUser(gOrg.ID, avpUser.Email, "Editor")
			if errAddOrgUser == nil {
				log.Printf("User %s was added to organization %s\n", avpUser.Email, gOrg.Name)
			} else {
				log.Println(errAddOrgUser)
			}
		} else {
			log.Printf("User %s already is member of %s organization", avpUser.Email, gOrg.Name)
		}
	}
}
