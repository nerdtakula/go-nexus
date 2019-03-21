package nexus

// Repository object
type Repository struct {
	Name   string `json:"name"`
	Format string `json:"format"`
	Type   string `json:"type"`
	URL    string `json:"url"`
}

// Repositories list
func (c Client) Repositories() ([]Repository, error) {
	var result []Repository
	err := c.makeRequest("GET", "/repositories", nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Repository lookup
func (c Client) Repository(repositoryID string) (*Repository, error) {
	repos, err := c.Repositories()
	if err != nil {
		return nil, err
	}

	for _, repo := range repos {
		if repo.Name == repositoryID {
			return &repo, nil
		}
	}
	return nil, ErrNotFound
}
