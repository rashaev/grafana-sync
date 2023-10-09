package worker

import (
	"strings"

	gapi "github.com/grafana/grafana-api-golang-client"
)

// Check by email if Organization already contains the user. true = contains, false = don't contains
func orgContainsUser(orgUsers []gapi.OrgUser, email string) bool {
	for _, orgUser := range orgUsers {
		if orgUser.Email == email {
			return true
		}
	}
	return false
}

// func parses a project code from a monitoring group name
func parseProjectCode(group string) string {
	s := strings.Split(group, "-")
	return s[len(s)-1]
}
