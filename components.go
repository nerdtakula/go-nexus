package nexus

import "github.com/pkg/errors"

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
func (c Client) Components() (components []Component, token string, err error) {
	return nil, "", nil
}

// UploadComponent to nexus
func (c Client) UploadComponent(repositoryID string, parameters UploadParameters) (*Component, error) {
	// Check Repo exists
	repo, err := c.Repository(repositoryID)
	if err != nil {
		return nil, errors.Wrap(err, "UploadComponent")
	}

	switch repo.Format {
	case "maven2":
		return c.uploadMaven2Component(repositoryID, parameters)
	case "raw":
		return c.uploadRawComponent(repositoryID, parameters)
	case "pypi":
		return c.uploadPyPiComponent(repositoryID, parameters)
	case "rubygems":
		return c.uploadRubyGemsComponent(repositoryID, parameters)
	case "nuget":
		return c.uploadNugetComponent(repositoryID, parameters)
	case "npm":
		return c.uploadNPMComponent(repositoryID, parameters)
	}
	return nil, errors.Wrap(ErrUnknownRepoFormat, "UploadComponent")
}

// Component single lookup
func (c Client) Component(id string) (*Component, error) {
	return nil, nil
}

// DeleteComponent from nexus
func (c Client) DeleteComponent(id string) error {
	return nil
}
