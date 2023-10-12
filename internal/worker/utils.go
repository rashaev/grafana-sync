package worker

import (
	"log"
	"strings"

	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/rashaev/grafana-sync/internal/keycloak"
)

// Check by email if Organization already contains the user. true = contains, false = don't contains
func orgContainsUser(orgUsers []gapi.OrgUser, email string) (int, bool) {
	for orgIdx, orgUser := range orgUsers {
		if orgUser.Email == email {
			return orgIdx, true
		}
	}
	return 0, false
}

// func parses a project code from a monitoring group name
func parseProjectCode(group string) string {
	s := strings.Split(group, "-")
	return s[len(s)-1]
}

// func check user Role in Grafana
func checkUserRole(orgUsers []gapi.OrgUser, user int, role string) bool {
	if orgUsers[user].Role == role {
		return true
	}
	return false
}

// merge users from RO and RW groups into map
func userFromGroups(kcloak *keycloak.Keycloak, kcloakGroups []keycloak.Group) map[string]string {

	var allProjectUsersSlice []keycloak.User
	var allProjectUsersMap = make(map[string]string)

	for _, group := range kcloakGroups {
		users, err := kcloak.UsersInGroup(group.ID)
		if err != nil {
			log.Println(err)
		}

		allProjectUsersSlice = append(allProjectUsersSlice, users...)
	}

	for _, value := range allProjectUsersSlice {
		allProjectUsersMap[value.Email] = value.ID
	}

	return allProjectUsersMap
}
