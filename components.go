package nexus

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Component object
type Component struct {
	ID         string  `json:"id"`
	Repository string  `json:"repository"`
	Format     string  `json:"format"`
	Group      string  `json:"group"`
	Name       string  `json:"name"`
	Version    string  `json:"version"`
	Assets     []Asset `json:"assets"`
}

// Components list
func (c Client) Components(repository string) (components []Component, continuationToken string, err error) {
	if len(strings.TrimSpace(repository)) == 0 {
		return nil, "", fmt.Errorf("repository can not be empty")
	}

	args := map[string]interface{}{
		"repository":        repository,
		"continuationToken": continuationToken,
	}

	if continuationToken == "" {
		delete(args, "continuationToken")
	}

	result := struct {
		Items             []Component `json:"items"`
		ContinuationToken string      `json:"continuationToken"`
	}{}

	_, err = c.makeRequest("GET", "/components", args, &result)
	if err != nil {
		return nil, "", errors.Wrap(err, "Components")
	}
	return result.Items, result.ContinuationToken, nil
}

// UploadComponent to upload single component
func (c Client) UploadComponent(repository string, parameters UploadParameters) (*Component, error) {
	if len(strings.TrimSpace(repository)) == 0 {
		return nil, fmt.Errorf("repository can not be empty")
	}

	// Check Repo exists
	repo, err := c.Repository(repository)
	if err != nil {
		return nil, errors.Wrap(err, "UploadComponent")
	}

	switch repo.Format {
	case "maven2":
		return c.uploadMaven2Component(repository, parameters)
	case "raw":
		return c.uploadRawComponent(repository, parameters)
	case "pypi":
		return c.uploadPyPiComponent(repository, parameters)
	case "rubygems":
		return c.uploadRubyGemsComponent(repository, parameters)
	case "nuget":
		return c.uploadNugetComponent(repository, parameters)
	case "npm":
		return c.uploadNPMComponent(repository, parameters)
	}
	return nil, errors.Wrap(ErrUnknownRepoFormat, "UploadComponent")
}

// Component single lookup
func (c Client) Component(id string) (*Component, error) {
	if len(strings.TrimSpace(id)) == 0 {
		return nil, fmt.Errorf("component id can not be empty")
	}

	var component *Component

	if _, err := c.makeRequest("GET", fmt.Sprintf("/components/%s", id), nil, &component); err != nil {
		return nil, errors.Wrap(err, "Component")
	}
	return component, nil
}

// DeleteComponent from nexus
func (c Client) DeleteComponent(id string) error {
	if len(strings.TrimSpace(id)) == 0 {
		return fmt.Errorf("component id can not be empty")
	}

	statusCode, err := c.makeRequest("DELETE", fmt.Sprintf("/components/%s", id), nil, nil)
	switch statusCode {
	case -1:
		// Other error message from request
		return errors.Wrap(err, "Delete Component")
	case 204:
		// Component was successfully deleted
		return nil
	case 403:
		// Insufficient permissions to delete component
		return ErrInsufficientPermissions
	case 404:
		// Component not found
		return ErrNotFound
	case 422:
		// Malformed ID
		return ErrMalformedID
	}

	// Safety check
	if err != nil {
		return errors.Wrap(err, "Delete Component")
	}
	return nil
}
