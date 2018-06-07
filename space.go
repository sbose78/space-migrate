package main

import (
	"fmt"
	"net/http"
)

func createSpace(spaceID string, creatorID string, serviceAccountToken string, env string) error {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/spaces/%s?creator=%s", getServerName(env, AUTHSERVICE), spaceID, creatorID)
	fmt.Printf("Calling %s ", url)
	req, err := http.NewRequest("GET", url, nil)
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
	for _, user := range userList {
		err := addUserToSpace(user.Attributes.IdentityID, spaceID, spaceManagerIdentityID, usertoken, env)
		if err != nil {
			return err
		}
	}
	return nil
}

func addUserToSpace(userID string, spaceID string, spaceManagerIdentityID string, usertoken string, env string) error {
	return nil
}
