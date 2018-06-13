package main

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"github.com/sbose78/space-migrate/migrate"
)

func validateSpaceMigration(spaceID uuid.UUID, serviceAccountToken string, env *string) error {

	// Validate that all is good!
	fmt.Println("Starting validation")
	newAssignees, err := migrate.GetAssignees(spaceID.String(), serviceAccountToken, *env)
	if err != nil {
		fmt.Printf("Error fetching the new assignee list during validation - %s \n", err)
	}
	oldCollaborators, err := migrate.GetCollaborators(spaceID.String(), *env)
	if err != nil {
		fmt.Printf("Error fetching the old collab list during validation - %s \n", err)
	}
	if len(oldCollaborators) != len(newAssignees) {
		return fmt.Errorf("for space id %s, there are %d collaborators , but %d new assignees", spaceID.String(), len(oldCollaborators), len(newAssignees))
	}
	toBeAssigned := migrate.GetUsersToBeAddedToSpace(newAssignees, oldCollaborators)
	if len(toBeAssigned) != 0 {
		return fmt.Errorf("for space id %s, there are %d missing assignees", spaceID.String(), len(newAssignees))
	}
	fmt.Println("Validation successful")
	return nil
}

func migrateSpace(spaceID uuid.UUID, spaceOwnerID string, serviceAccountToken string, env *string) error {
	collaboratorList, err := migrate.GetCollaborators(spaceID.String(), *env)

	if err != nil {
		return err
	}
	fmt.Printf("\n The space %s has %d collaborators incl. creator \n", spaceID, len(collaboratorList))

	created, err := migrate.CreateSpace(spaceID.String(), spaceOwnerID, serviceAccountToken, *env)
	if err != nil {
		return err
	}

	if !created { // possible only if space resource already exists.

		originalListLen := len(collaboratorList)

		// This API Call finds out of if space resource exists and if yes, who are the members.
		existingAssignees, err := migrate.GetAssignees(spaceID.String(), serviceAccountToken, *env)
		if err != nil {
			// could be an API call error, or an error because 'Space resource not found'
			return err
		}
		if !migrate.IsUserAssigned(spaceOwnerID, existingAssignees) {
			// admin gets added at the time of space creation. S
			// So, if the space owner wasn't added as 'admin'in the space resource,
			// something is wrong if that has failed! HENCE -
			return err
		}

		// If all the users are already added to the space, we do nothing.
		collaboratorList = migrate.GetUsersToBeAddedToSpace(existingAssignees, collaboratorList)
		if len(collaboratorList) == 0 {
			fmt.Println("Assignees already present - skipping adding of contributors to this space")
			return nil
		}
		// else - just go ahead and add who already isn't present.
		fmt.Printf("\n unexpected : %d out of %d collaborators are missing in the new space resource, they will be added. \n", len(collaboratorList), originalListLen)
	}

	// Add the users to the space as a 'contributor'
	err = migrate.AddUsersToSpace(collaboratorList, spaceID.String(), spaceOwnerID, serviceAccountToken, *env)
	if err != nil {
		return err
	}

	err = validateSpaceMigration(spaceID, serviceAccountToken, env)
	if err != nil {
		return err
	}

	return err
}

func migrateSpaces(ids []string, env *string, sessionState *string, privateKey *rsa.PrivateKey, serviceAccountToken string) error {

	for userIDIndex, id := range ids {
		if len(id) > 0 {
			spaceOwnerID := strings.TrimSpace(id)

			user, err := migrate.LoadUser(spaceOwnerID, *env)
			if err != nil {
				return err
			}
			fmt.Printf("\n----- Starting space migration for identity ID %s ---- ( %d/%d ) \n", spaceOwnerID, userIDIndex+1, len(ids))
			fmt.Printf("%s, %s, %s\n", user.Data.Attributes.Username, user.Data.Attributes.Email, user.Data.Attributes.FullName)

			spacesList, err := migrate.GetSpacesOwnedByIdentity(user.Data.Attributes.Username, *env)
			if err != nil {
				return err
			}

			for c, space := range spacesList {
				fmt.Printf("\n--- Starting space migration for space ID %s for %s ---- ( %d/%d ) \n\n", (*space.ID).String(), user.Data.Attributes.Username, c+1, len(spacesList))
				time.Sleep(5 * time.Second)
				err = migrateSpace(*space.ID, spaceOwnerID, serviceAccountToken, env)
				if err != nil {
					return err
				}
			}
			fmt.Printf("\n\n----- Completed space migration for identity ID %s -----\n\n", spaceOwnerID)

		}
	}
	return nil

}

func main() {

	//keyLoc := flag.String("key", "foo", "private key location")
	serviceAccountTokenLoc := flag.String("satoken", "foo", "service account token location")
	//sessionState := flag.String("session", uuid.NewV4().String(), "session state")
	env := flag.String("env", "prod", "prod or prod-preview")
	userIDLoc := flag.String("users", "foo", "user UUIDs location")

	flag.Parse()
	fmt.Printf("service account token path: %s ; \n  env: %s ;\n userIDFile: %s \n", *serviceAccountTokenLoc, *env, *userIDLoc)

	serviceAccountToken, err := ioutil.ReadFile(*serviceAccountTokenLoc)
	if err != nil {
		panic(err)
	}
	serviceAccountTokenString := strings.TrimSpace(string(serviceAccountToken))

	ids := migrate.GetUserIDs(*userIDLoc)
	migrateSpaces(ids, env, nil, nil, serviceAccountTokenString)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
