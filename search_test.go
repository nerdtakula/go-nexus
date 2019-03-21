package nexus

import "testing"

func TestSearchComponents(t *testing.T) {
	params := SearchParameters{
		Query:   "VizFlow",
		Version: "BE:1.3.8-UI:2.0.12",
	}
	results, _, err := client.SearchComponents(params)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Results: %+v\n", results)
}

func TestSearchAssets(t *testing.T) { t.Skip("Not Implemented") }
