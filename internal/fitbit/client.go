package fitbit

import "fmt"

const (
	ClientIdEnvVarName     string = "FITBIT_CLIENT_ID"
	ClientSecretEnvVarName string = "FITBIT_CLIENT_SECRET"
)

type FitbitTokenPair struct {
	AccessToken  string
	RefreshToken string
}

type FitbitAPICred struct {
	ClientId     string
	ClientSecret string
}

func NewFitbitAPICred(m map[string]string) (*FitbitAPICred, error) {
	creds := new(FitbitAPICred)

	val, ok := m[ClientIdEnvVarName]
	if !ok {
		return nil, fmt.Errorf("Missing config for %s", ClientIdEnvVarName)
	}
	creds.ClientId = val

	val, ok = m[ClientSecretEnvVarName]
	if !ok {
		return nil, fmt.Errorf("Missing config for %s", ClientSecretEnvVarName)
	}
	creds.ClientSecret = val

	return creds, nil
}
