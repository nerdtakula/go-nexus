package nexus

import "testing"

func TestSearchComponents(t *testing.T) {
	params := SearchParameters{
		Query:   "TestApplication",
		Version: "2.7.0",
	}
	results, _, err := client.SearchComponents(params)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Results: %+v\n", results)
}

func TestSearchAssets(t *testing.T) { t.Skip("Not Implemented") }
