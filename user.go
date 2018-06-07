package main

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
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

func generateToken(key *rsa.PrivateKey, user User, userID, sessionState, env string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	if env == "prod" {
		token.Header["kid"] = "0lL0vXs9YRVqZMowyw8uNLR_yr0iFaozdQk9rzq2OVU"
		token.Claims.(jwt.MapClaims)["iss"] = "https://sso.openshift.io/auth/realms/fabric8"
	} else {
		token.Header["kid"] = "zD-57oBFIMUZsAWqUnIsVu_x71VIjd1irGkGUOiTsL8"
		token.Claims.(jwt.MapClaims)["iss"] = "https://sso.prod-preview.openshift.io/auth/realms/fabric8"
	}

	nowTime := time.Now().Unix()
	in3Days := nowTime + 3*24*60*60
	token.Claims.(jwt.MapClaims)["jti"] = uuid.NewV4().String()
	token.Claims.(jwt.MapClaims)["exp"] = in3Days
	token.Claims.(jwt.MapClaims)["nbf"] = 0
	token.Claims.(jwt.MapClaims)["iat"] = nowTime
	token.Claims.(jwt.MapClaims)["aud"] = "fabric8-online-platform"
	token.Claims.(jwt.MapClaims)["sub"] = userID
	token.Claims.(jwt.MapClaims)["typ"] = "Bearer"
	token.Claims.(jwt.MapClaims)["azp"] = "fabric8-online-platform"
	token.Claims.(jwt.MapClaims)["auth_time"] = nowTime
	token.Claims.(jwt.MapClaims)["session_state"] = sessionState
	token.Claims.(jwt.MapClaims)["acr"] = "1"
	token.Claims.(jwt.MapClaims)["approved"] = "true"
	token.Claims.(jwt.MapClaims)["name"] = user.Data.Attributes.FullName
	token.Claims.(jwt.MapClaims)["company"] = ""
	token.Claims.(jwt.MapClaims)["preferred_username"] = user.Data.Attributes.Username
	token.Claims.(jwt.MapClaims)["given_name"] = ""
	token.Claims.(jwt.MapClaims)["family_name"] = ""
	token.Claims.(jwt.MapClaims)["email"] = user.Data.Attributes.Email

	token.Claims.(jwt.MapClaims)["allowed-origins"] = []string{
		"https://auth.openshift.io",
		"https://auth.prod-preview.openshift.io",
		"https://api.openshift.io",
		"https://api.prod-preview.openshift.io",
		"https://openshift.io",
		"https://prod-preview.openshift.io",
		"http://localhost:3000"}

	realmAccess := make(map[string]interface{})
	realmAccess["roles"] = []string{"uma_authorization"}
	token.Claims.(jwt.MapClaims)["realm_access"] = realmAccess

	resourceAccess := make(map[string]interface{})
	broker := make(map[string]interface{})
	broker["roles"] = []string{"read-token"}
	resourceAccess["broker"] = broker

	account := make(map[string]interface{})
	account["roles"] = []string{"manage-account", "manage-account-links", "view-profile"}
	resourceAccess["account"] = account

	token.Claims.(jwt.MapClaims)["resource_access"] = resourceAccess

	tokenStr, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func loadUser(id, env string) (User, error) {
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

func getUserIDs(userIDLoc string) []string {

	userUUIDs, err := ioutil.ReadFile(userIDLoc)
	if err != nil {
		panic(err)
	}
	return strings.Split(string(userUUIDs), "\n")
}
