package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func createSpace(spaceID string, creatorID string, serviceAccountToken string, env string) error {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/spaces/%s?creator=%s", getServerName(env, AUTHSERVICE), spaceID, creatorID)
	fmt.Printf("Calling %s ", url)
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", serviceAccountToken))
	if err != nil {
		fmt.Println(err)
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Space API call %s returned %d", url, resp.StatusCode)
	}

	return nil
}

func addUsersToSpace(userList []*Data, spaceID string, spaceManagerIdentityID string, usertoken string, env string) error {
	var identityIDs []string
	for _, user := range userList {
		identityIDs = append(identityIDs, user.Attributes.IdentityID)
	}

	payload := AssignRoleResourceRolesPayload{
		Data: []*AssignRoleData{
			&AssignRoleData{
				Ids:  identityIDs,
				Role: "contributor",
			},
		},
	}

	client := &http.Client{}
	url := fmt.Sprintf("%s/api/resources/%s/roles", getServerName(env, AUTHSERVICE), spaceID)
	fmt.Printf("Calling %s ", url)

	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	err := enc.Encode(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, buf)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", usertoken))
	if err != nil {
		fmt.Println(err)
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Space API call %s returned %d", url, resp.StatusCode)
	}
	return nil
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
