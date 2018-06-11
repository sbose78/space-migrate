package migrate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type User struct {
	Data Data `json:"data"`
}

type Result struct {
	Tokens []TokenResult `json:"tokens"`
}

type TokenResult struct {
	UUID  string `json:"uuid"`
	Token string `json:"token"`
}

type Token struct {
	Data Data `json:"data"`
}

func LoadUser(id, env string) (User, error) {
	var user User
	url := fmt.Sprintf("%s/api/users/%s", getServerName(env, AUTHSERVICE), id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return user, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return user, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return user, err
	}
	if res.StatusCode != http.StatusOK {
		return user, errors.New("Status is not 200 OK: " + res.Status)
	}
	err = json.Unmarshal(body, &user)

	return user, err
}

func GetUserIDs(userIDLoc string) []string {

	userUUIDs, err := ioutil.ReadFile(userIDLoc)
	if err != nil {
		panic(err)
	}
	return strings.Split(string(userUUIDs), "\n")
}
