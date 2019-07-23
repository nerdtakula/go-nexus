package nexus

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Asset object
type Asset struct {
	DownloadURL string            `json:"downloadUrl"`
	Path        string            `json:"path"`
	ID          string            `json:"id"`
	Repository  string            `json:"repository"`
	Format      string            `json:"format"`
	Checksum    map[string]string `json:"checksum"`
}

// AssetGroup object
type AssetGroup struct {
	ID         string  `json:"id"`
	Group      string  `json:"group"`
	Name       string  `json:"name"`
	Version    string  `json:"version"`
	Repository string  `json:"repository"`
	Format     string  `json:"format"`
	Assets     []Asset `json:"assets"`
}

// Assets list via endpoint
func (c Client) Assets(repositoryID, continuationToken string) (assets []Asset, token string, err error) {
	args := map[string]interface{}{
		"repository":        repositoryID,
		"continuationToken": continuationToken,
	}

	if continuationToken == "" {
		delete(args, "continuationToken")
	}

	result := struct {
		Items             []Asset `json:"items"`
		ContinuationToken string  `json:"continuationToken"`
	}{}

	_, err = c.makeRequest("GET", "/assets", args, &result)
	if err != nil {
		return nil, "", errors.Wrap(err, "Assets")
	}
	return result.Items, result.ContinuationToken, nil
}

// Asset lookup via endpoint
func (c Client) Asset(id string) (*Asset, error) {
	if len(strings.TrimSpace(id)) == 0 {
		return nil, fmt.Errorf("asset id can not be empty")
	}

	var asset *Asset

	if _, err := c.makeRequest("GET", fmt.Sprintf("/assets/%s", id), nil, &asset); err != nil {
		return nil, errors.Wrap(err, "Asset")
	}
	return asset, nil
}

// DeleteAsset via endpoint
func (c Client) DeleteAsset(id string) error {
	if len(strings.TrimSpace(id)) == 0 {
		return fmt.Errorf("asset id can not be empty")
	}

	statusCode, err := c.makeRequest("DELETE", fmt.Sprintf("/assets/%s", id), nil, nil)
	switch statusCode {
	case -1:
		// Other error message from request
		return errors.Wrap(err, "Delete Asset")
	case 204:
		// Asset was successfully deleted
		return nil
	case 403:
		// Insufficient permissions to delete asset
		return ErrInsufficientPermissions
	case 404:
		// Asset not found
		return ErrNotFound
	case 422:
		// Malformed ID
		return ErrMalformedID
	}

	// Safety check
	if err != nil {
		return errors.Wrap(err, "Delete Asset")
	}
	return nil
}
