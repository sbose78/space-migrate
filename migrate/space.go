package migrate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func CreateSpace(spaceID string, creatorID string, serviceAccountToken string, env string) (bool, error) {
	fmt.Printf("Creating Space resource %s with 'admin' %s in %s ... \n", spaceID, creatorID, env)
	code, err := createSpace(spaceID, creatorID, serviceAccountToken, env)
	if err != nil {
		return false, err
	}
	if code == http.StatusConflict {
		fmt.Printf("Space Resource %s in native auth service already exists \n", spaceID)
		return false, nil
	}
	if code != http.StatusOK {
		return false, fmt.Errorf("space resource creation for spaceID %s failed with http response code : %d", spaceID, code)
	}
	// TODO: let's verify that the space resource actually exists ?
	return true, nil
}

func createSpace(spaceID string, creatorID string, serviceAccountToken string, env string) (int, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/spaces/%s?creator=%s", getServerName(env, AUTHSERVICE), spaceID, creatorID)
	fmt.Printf("Calling POST %s \n", url)
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", serviceAccountToken))
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Printf("Space API call %s returned %d \n", url, resp.StatusCode)

	return resp.StatusCode, nil
}

func AddUsersToSpace(userList []*Data, spaceID string, spaceManagerIdentityID string, usertoken string, env string) error {
	var identityIDs []string
	for _, user := range userList {
		if user.Attributes.IdentityID != spaceManagerIdentityID {
			// no need to add the space owner to the space,
			// she would already have been added as part of the space creation.
			identityIDs = append(identityIDs, user.Attributes.IdentityID)
		}

	}

	payload := AssignRoleResourceRolesPayload{
		Data: []*AssignRoleData{
			&AssignRoleData{
				Ids:  identityIDs,
				Role: "contributor",
			},
		},
	}

	fmt.Printf("Adding %d users to the space %s ... \n", len(payload.Data[0].Ids), spaceID)
	if len(payload.Data[0].Ids) == 0 {
		return nil
	}
	code, err := addUsersToSpace(payload, spaceID, spaceManagerIdentityID, usertoken, env)
	if err != nil {
		return err
	}

	if code != http.StatusNoContent {
		return fmt.Errorf("AddUsersToSpace API call for spaceID %s returned %d", spaceID, code)
	}

	return nil

}

func addUsersToSpace(payload AssignRoleResourceRolesPayload, spaceID string, spaceManagerIdentityID string, usertoken string, env string) (int, error) {

	client := &http.Client{}
	url := fmt.Sprintf("%s/api/resources/%s/roles", getServerName(env, AUTHSERVICE), spaceID)
	fmt.Printf("Calling PUT %s \n", url)

	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	err := enc.Encode(payload)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("PUT", url, buf)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", usertoken))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Printf("Authorization API call %s returned %d \n", url, resp.StatusCode)
	return resp.StatusCode, nil
}

// AssignRoleResourceRolesPayload is the resource_roles assignRole action payload.
type AssignRoleResourceRolesPayload struct {
	Data []*AssignRoleData `json:"data"`
}

// AssignRoleData user type.
type AssignRoleData struct {
	// identity ids to assign role to
	Ids []string `json:"ids"`
	// name of the role to assign
	Role string `json:"role"`
}

// GetUsersToBeAddedToSpace finds out the diff between the 'who has been already assigned' and 'who needs to be assigned'
func GetUsersToBeAddedToSpace(alreadyAssigned []*IdentityRolesData, collaborators []*Data) []*Data {
	var toBeAdded []*Data
	for _, collab := range collaborators {
		if !IsUserAssigned(collab.Attributes.IdentityID, alreadyAssigned) {
			toBeAdded = append(toBeAdded, collab)
		}
	}
	return toBeAdded
}

func IsUserAssigned(identityId string, alreadyAssigned []*IdentityRolesData) bool {
	for _, collab := range alreadyAssigned {
		if collab.AssigneeID == identityId {
			return true
		}
	}
	return false
}
