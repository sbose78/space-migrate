package main

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"github.com/sbose78/space-migrate/migrate"
)

func migrateSpaces(ids []string, env *string, sessionState *string, privateKey *rsa.PrivateKey, serviceAccountToken string) {

	for _, id := range ids {
		if len(id) > 0 {
			userID := strings.TrimSpace(id)
			user, err := migrate.LoadUser(userID, *env)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s, %s, %s\n", user.Data.Attributes.Username, user.Data.Attributes.Email, user.Data.Attributes.FullName)

			spacesList, err := migrate.GetSpacesOwnedByIdentity(user.Data.Attributes.Username, *env)
			if err != nil {
				panic(err)
			}

			for _, space := range spacesList {
				spaceID := *space.ID

				userList, err := migrate.GetCollaborators(spaceID.String(), *env)
				if err != nil {
					panic(err)
				}

				created, err := migrate.CreateSpace(spaceID.String(), userID, serviceAccountToken, *env)
				if err != nil {
					panic(err)
				}

				if !created {
					// TODO using the same API call :
					// 1. Check if the space really exists? and panic if not.
					// 2. If it does, Check if the assignees exist ? If not go to the next step.
				}

				// Add the users to the space as a 'contributor'
				err = migrate.AddUsersToSpace(userList, spaceID.String(), userID, serviceAccountToken, *env)
				if err != nil {
					panic(err)
				}

				// TODO:
				// Verify that the space assignees exist.

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

	ids := migrate.GetUserIDs(*userIDLoc)
	migrateSpaces(ids, env, sessionState, privateKey, string(serviceAccountToken))

}
