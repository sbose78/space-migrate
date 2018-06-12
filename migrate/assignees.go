package migrate

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Copied from fabric8-auth

// Identifier: application/vnd.identityroles+json; view=default
type Identityroles struct {
	Data []*IdentityRolesData `form:"data" json:"data" xml:"data"`
}

// IdentityRolesData user type.
type IdentityRolesData struct {
	// The ID of the assignee
	AssigneeID string `form:"assignee_id" json:"assignee_id" xml:"assignee_id"`
	// The type of assignee, example: user,group,team
	AssigneeType string `form:"assignee_type" json:"assignee_type" xml:"assignee_type"`
	Inherited    bool   `form:"inherited" json:"inherited" xml:"inherited"`
	// The ID of the resource from this role was inherited
	InheritedFrom *string `form:"inherited_from,omitempty" json:"inherited_from,omitempty" xml:"inherited_from,omitempty"`
	// The name of the role
	RoleName string `form:"role_name" json:"role_name" xml:"role_name"`
}

func GetAssignees(resourceID string, token string, env string) ([]*IdentityRolesData, error) {
	assignedIdentityRoles, code, err := getAssignees(resourceID, token, env)
	if err != nil {
		return nil, err
	}
	if code != http.StatusOK {
		// if space resource is absent, 404 would be returned.
		return nil, fmt.Errorf("resourceRoles API call returned %d", code)
	}

	return assignedIdentityRoles.Data, nil
}

func getAssignees(resourceID string, token string, env string) (*Identityroles, int, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/api/resources/%s/roles", getServerName(env, AUTHSERVICE), resourceID)

	fmt.Println("Calling ", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, -1, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, -1, err
	}
	defer resp.Body.Close()

	returnedUserList := &Identityroles{}
	err = json.NewDecoder(resp.Body).Decode(returnedUserList)
	if err != nil {
		fmt.Println(err)
		return nil, resp.StatusCode, err
	}
	return returnedUserList, resp.StatusCode, nil
}
