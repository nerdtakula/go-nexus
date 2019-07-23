package nexus

import (
	"fmt"

	"github.com/pkg/errors"
)

// BlobStore Quota Information
type BlobStore struct {
	IsViolation bool   `json:"isViolation"`
	Message     string `json:"message"`
	Name        string `json:"blobStoreName"`
}

// BlobStore - Get quota status for a given blob store
func (c Client) BlobStore(id string) (*BlobStore, error) {
	result := new(BlobStore)

	if err := c.makeRequest("GET", fmt.Sprintf("/blobstores/%s/quota-status"), nil, &result); err != nil {
		return nil, errors.Wrap(err, "BlobStore")
	}
	return result, nil
}
