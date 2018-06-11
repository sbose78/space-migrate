package migrate

import "fmt"

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
