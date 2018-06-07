package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test the individual calls.

func TestLoadUser(t *testing.T) {

	// prod-preview
	identityID := "1d5a1cbd-00ca-4742-a3da-1fd0226cdc24"
	user, err := loadUser(identityID, "prod-preview")
	assert.NoError(t, err)
	assert.Equal(t, identityID, user.Data.Attributes.IdentityID)
	assert.Equal(t, "shbose-preview1", user.Data.Attributes.Username)

	// prod
	identityID = "3383826c-51e4-401b-9ccd-b898f7e2397d"
	user, err = loadUser(identityID, "prod")
	assert.NoError(t, err)
	assert.Equal(t, identityID, user.Data.Attributes.IdentityID)
	assert.Equal(t, "shbose", user.Data.Attributes.Username)

}

func TestNamedSpaces(t *testing.T) {
	spaceList, err := getSpacesOwnedByIdentity("shbose-preview", "prod-preview")
	assert.NoError(t, err)
	assert.Len(t, spaceList, 3)

	spaceList, err = getSpacesOwnedByIdentity("shbose", "prod")
	assert.NoError(t, err)
	assert.Len(t, spaceList, 2)
}

func TestGetCollaborators(t *testing.T) {
	spaceID := "a2deead7-fdc1-4953-ac86-ea9267868a8c"
	collabList, err := getCollaborators(spaceID, "prod-preview")
	assert.NoError(t, err)
	assert.NotNil(t, collabList)
	assert.Len(t, collabList, 2)
	assert.Contains(t, collabList[0].Attributes.Username, "bose")
}

// TODO: Test Create Space using a service account, using a local service account token.

// TODO: Add a test which does a full migration against local API servers
