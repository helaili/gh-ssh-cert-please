package utils

import (
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/google/go-github/v60/github"
)

func GetUSer() (*github.User, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, err
	}
	// Get the current user's login
	profileResponse := github.User{}
	err = client.Get("user", &profileResponse)
	if err != nil {
		return nil, err
	}

	return &profileResponse, nil
}
