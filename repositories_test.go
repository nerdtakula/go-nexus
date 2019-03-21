package nexus

import "testing"

func TestRepositories(t *testing.T) {
	repos, err := client.Repositories()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Results: %+v\n", repos)
	// TODO: check for a known repo
}
