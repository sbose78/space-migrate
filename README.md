
```
$ go install
$ BINARY_PATH/space-migrate -key keypath/keypath2 -users file.txt -session somestate -env preview
```
### What 
Migrate all spaces in prod-preview and prod, from keycloak-managed resources to native Auth service managed resources.

### Why

Current spaces use keycloak authorization resources. We are migrating away from keycloak and all spaces should use the native Auth service spaces.

### How

For each identity,
1. Get list of spaces owned by Identity using the `api.openshift.io/api/namedspaces` API
2. Create the corresponding space resource for the new authz model using the `POST auth.openshift.io/api/spaces/....` API using a service account token.
3. Get the 'old' collaborators using the `api.openshift.io/api/spaces/spaceID/collaborators` API
4. Use the `Add assignee API` to add the 'old' collaborators as `contributor`s  to the space resource created in (2)

All steps in the script are re-runnable. Any conflicts arising from a pre-existing migration is handled.

#### Alternate Path - Space creation fails in step (2): 

This can happen if space resource is already present. In that case, confirm that the space resource exists using the `GET /api/resources/resourceID/roles` API which returns all assigned roles in a resource.

If present, calculate the diff between list of 'old' space collaborators and 'new' resource assignees. Then, assign the diff using the `Add assignee API`

### Validation of migration

As a last step in the migration, we validate if 
`List of collaborators returned by Collabortors API` match `List of Assignees returned by GET /api/resources/resourceID/roles API`


