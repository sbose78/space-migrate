package migrate_test

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/sbose78/space-migrate/migrate"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test the individual calls.

func TestLoadUser(t *testing.T) {

	// prod-preview
	identityID := "1d5a1cbd-00ca-4742-a3da-1fd0226cdc24"
	user, err := migrate.LoadUser(identityID, "prod-preview")
	assert.NoError(t, err)
	assert.Equal(t, identityID, user.Data.Attributes.IdentityID)
	assert.Equal(t, "shbose-preview1", user.Data.Attributes.Username)

	// prod
	identityID = "3383826c-51e4-401b-9ccd-b898f7e2397d"
	user, err = migrate.LoadUser(identityID, "prod")
	assert.NoError(t, err)
	assert.Equal(t, identityID, user.Data.Attributes.IdentityID)
	assert.Equal(t, "shbose", user.Data.Attributes.Username)

}

func TestNamedSpaces(t *testing.T) {
	spaceList, err := migrate.GetSpacesOwnedByIdentity("shbose-preview", "prod-preview")
	assert.NoError(t, err)
	assert.Len(t, spaceList, 3)

	spaceList, err = migrate.GetSpacesOwnedByIdentity("shbose", "prod")
	assert.NoError(t, err)
	assert.Len(t, spaceList, 2)
}

func TestGetCollaborators(t *testing.T) {
	spaceID := "a2deead7-fdc1-4953-ac86-ea9267868a8c"
	collabList, err := migrate.GetCollaborators(spaceID, "prod-preview")
	assert.NoError(t, err)
	assert.NotNil(t, collabList)
	assert.Len(t, collabList, 2)
	assert.Contains(t, collabList[0].Attributes.Username, "bose")
}

func TestCreateSpace(t *testing.T) {
	// NOTE: This test only runs locally with specific test data.

	// generated by
	// curl --request POST   --url http://localhost:8089/api/token --data "grant_type=client_credentials&client_id=d3170241-97cb-43e3-acee-41355ecc5edb&client_secret=secret"

	localServiceAccountToken := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjlNTG5WaWFSa2hWajFHVDlrcFdVa3dISXdVRC13WmZVeFItM0Nwa0UtWHMiLCJ0eXAiOiJKV1QifQ.eyJpYXQiOjE1MjgzNjc0MzksImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3QiLCJqdGkiOiJmMTc2MjYwYi0yOGIwLTRmMjYtOWU5Yy1lZWVlMmI2YmQ2YmEiLCJzY29wZXMiOlsidW1hX3Byb3RlY3Rpb24iXSwic2VydmljZV9hY2NvdW50bmFtZSI6InNwYWNlLW1pZ3JhdGlvbiIsInN1YiI6ImQzMTcwMjQxLTk3Y2ItNDNlMy1hY2VlLTQxMzU1ZWNjNWVkYiJ9.a-sK6ZeepYmPuOhOYemkaVOr__IzGkKDemPNKZdjnWCiUliz1hMtCKIDQPQK7w5HA2aGEa8Yt12aGtcZC5MfDyq46xNjGDL00neq0Sc3PKA7BjrySuP_VJHEHBvv5sildwb9m_nsjaKryE1JwN86LnHNystpJGhGDQhKwx-CXMwNt4VNTcwVqQ1ikQuOi_Bu6VkltHdEGv3uZqbQvt_4T-goaglxFmLyyS9bdtI3NyYrTPcdJy34h28VlR5b6c1GqcQcQV2Ee4K4uUdFVkThYscgKViNBjfJ86s3pOk_wMYVVJEqF06oAQP2z1rlWwKYvfRTbOFCpfllYCKr_NVDYg"

	// generated by /api/login/generate
	creatorID := "a3fc78c6-eab9-4f50-84a0-e74fc3ced3e7"

	// login locally using 2 different RHD accounts.
	collab := "837f2447-2e42-4db9-9f32-817d4866178a"  // user shbose
	collab2 := "d689af02-329e-492d-97c3-98f281ff149d" // user sbose78

	spaceID := uuid.NewV4()
	created, err := migrate.CreateSpace(spaceID.String(), creatorID, localServiceAccountToken, "localhost")
	require.NoError(t, err)
	require.True(t, created)

	userList := []*migrate.Data{
		&migrate.Data{
			Attributes: migrate.Attributes{IdentityID: collab},
		},
		&migrate.Data{
			Attributes: migrate.Attributes{IdentityID: collab2},
		},
	}

	err = migrate.AddUsersToSpace(userList, spaceID.String(), creatorID, localServiceAccountToken, "localhost")
	require.NoError(t, err)

}

// TODO: Add a test which does a full migration against local API servers
