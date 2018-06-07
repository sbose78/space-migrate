package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

type SpaceList struct {
	Data  []*Space `form:"data" json:"data" xml:"data"`
	Links Links    `json:"links"`
}

// Space user type.
type Space struct {
	Attributes *SpaceAttributes `form:"attributes" json:"attributes" xml:"attributes"`
	// ID of the space
	ID   *uuid.UUID `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	Type string     `form:"type" json:"type" xml:"type"`
}

// SpaceAttributes user type.
type SpaceAttributes struct {
	// When the space was created
	CreatedAt *time.Time `form:"created-at,omitempty" json:"created-at,omitempty" xml:"created-at,omitempty"`
	// Description for the space
	Description *string `form:"description,omitempty" json:"description,omitempty" xml:"description,omitempty"`
	// Name for the space
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// When the space was updated
	UpdatedAt *time.Time `form:"updated-at,omitempty" json:"updated-at,omitempty" xml:"updated-at,omitempty"`
	// Version for optimistic concurrency control (optional during creating)
	Version *int `form:"version,omitempty" json:"version,omitempty" xml:"version,omitempty"`
}

func getSpacesOwnedByIdentity(username string, env string) ([]*Space, error) {
	client := &http.Client{}
	next := fmt.Sprintf("%s/api/namedspaces/%s?page[limit]=20", getServerName(env, WITSERVICE), username)
	url := ""

	var fullSpaceList []*Space

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
			return nil, fmt.Errorf("Namedspaces API call %s returned %d", url, resp.StatusCode)
		}

		returnedSpaceList := &SpaceList{}
		err = json.NewDecoder(resp.Body).Decode(returnedSpaceList)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		fullSpaceList = append(fullSpaceList, returnedSpaceList.Data...)
		next = returnedSpaceList.Links.Next
	}
	return fullSpaceList, nil
}
