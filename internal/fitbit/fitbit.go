package fitbit

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"
)

const (
	clientIdEnv     string = "FITBIT_CLIENT_ID"
	clientSecretEnv string = "FITBIT_CLIENT_SECRET"
	redirectUrlEnv  string = "FITBIT_REDIRECT_URL"
)

const (
	profileScope  string = "profile"
	activityScope string = "activity"
)

var requestedScopes []string = []string{profileScope, activityScope}

func NewFitbitOauthConfig() (*oauth2.Config, error) {
	clientId := os.Getenv(clientIdEnv)
	if strings.TrimSpace(clientId) == "" {
		return nil, fmt.Errorf("Missing environment variable: %s", clientIdEnv)
	}
	clientSecret := os.Getenv(clientSecretEnv)
	if strings.TrimSpace(clientSecret) == "" {
		return nil, fmt.Errorf("Missing environment variable: %s", clientSecretEnv)
	}
	redirectUrl := os.Getenv(redirectUrlEnv)
	if strings.TrimSpace(redirectUrl) == "" {
		return nil, fmt.Errorf("Missing environment variable: %s", redirectUrlEnv)
	}

	return &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
		Scopes:       requestedScopes,
		Endpoint:     fitbit.Endpoint,
	}, nil
}
