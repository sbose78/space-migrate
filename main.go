package main

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
)

const (
	AUTHSERVICE = "auth"
	WITSERVICE  = "api"
)

func getServerName(env, service string) string {

	if env == "prod" {
		return fmt.Sprintf("https://%s.openshift.io", service)
	} else if env == "prod-preview" {
		return fmt.Sprintf("https://%s.%s.openshift.io", service, env)
	} else if service == AUTHSERVICE { // for localhost env
		return "http://localhost:8089"
	}
	return "http://localhost:8080"

}

func migrate(ids []string, env *string, sessionState *string, privateKey *rsa.PrivateKey, serviceAccountToken string) {

	for _, id := range ids {
		if len(id) > 0 {
			userID := strings.TrimSpace(id)
			user, err := loadUser(userID, *env)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s, %s, %s\n", user.Data.Attributes.Username, user.Data.Attributes.Email, user.Data.Attributes.FullName)

			spacesList, err := getSpacesOwnedByIdentity(user.Data.Attributes.Username, *env)
			if err != nil {
				panic(err)
			}

			token, err := generateToken(privateKey, user, userID, *sessionState, *env)
			if err != nil {
				panic(err)
			}

			for _, space := range spacesList {
				spaceID := *space.ID

				createSpace(spaceID.String(), userID, string(serviceAccountToken), *env)

				userList, err := getCollaborators(spaceID.String(), *env)
				if err != nil {
					panic(err)
				}

				addUsersToSpace(userList, spaceID.String(), userID, token, *env)

			}

		}
	}

}

func main() {

	keyLoc := flag.String("key", "foo", "private key location")
	serviceAccountTokenLoc := flag.String("satoken", "foo", "service account token location")
	sessionState := flag.String("session", uuid.NewV4().String(), "session state")
	env := flag.String("env", "prod", "prod or prod-preview")
	userIDLoc := flag.String("users", "foo", "user UUIDs location")

	flag.Parse()
	fmt.Printf("private key path: %s ; \n service account token path: %s ; \n sessionState: %s ; \n env: %s ;\n userIDFile: %s \n", *keyLoc, *serviceAccountTokenLoc, *sessionState, *env, *userIDLoc)

	key, err := ioutil.ReadFile(*keyLoc)
	if err != nil {
		panic(err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		panic(err)
	}

	serviceAccountToken, err := ioutil.ReadFile(*serviceAccountTokenLoc)
	if err != nil {
		panic(err)
	}

	ids := getUserIDs(*userIDLoc)
	migrate(ids, env, sessionState, privateKey, string(serviceAccountToken))

}
