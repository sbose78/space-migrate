package migrate

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// https://api.openshift.io/api/spaces/fb28c2cb-4fc1-465b-b0df-fcd930e6afc2/collaborators

type UserList struct {
	Data  []*Data `form:"data" json:"data" xml:"data"`
	Links Links   `json:"links"`
}

type Data struct {
	Attributes Attributes `json:"attributes"`
}
type Links struct {
	Next string `json:"next"`
}

type Attributes struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	FullName   string `json:"fullName"`
	IdentityID string `json:"identityID"`
}

func GetCollaborators(spaceID string, env string) ([]*Data, error) {
	client := &http.Client{}

	next := fmt.Sprintf("%s/api/spaces/%s/collaborators?page[limit]=20", getServerName(env, WITSERVICE), spaceID)
	url := ""

	var fullReturnedUserList []*Data

	for next != "" {
		url = next
		fmt.Println("Calling ", url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Collaborators API call %s returned %d", url, resp.StatusCode)
		}

		returnedUserList := &UserList{}
		err = json.NewDecoder(resp.Body).Decode(returnedUserList)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		fullReturnedUserList = append(fullReturnedUserList, returnedUserList.Data...)
		next = returnedUserList.Links.Next
	}
	return fullReturnedUserList, nil
}
