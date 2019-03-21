package nexus

import (
	"testing"
)

const (
	testRepositoryID = "maven-releases"
	testAssetID      = "testAsset"
)

// Test querying a list of assets
func TestAssets(t *testing.T) {
	_, _, err := client.Assets(testRepositoryID, "")
	if err != nil {
		t.Fatal(err)
	}
	// t.Logf("Results: %+v\n", assets)
	// TODO: Test for a known asset in the results
}

// Test lookup of single asset
func TestAsset(t *testing.T) {
	asset, err := client.Asset(testAssetID)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Results: %+v\n", asset)
	// TODO: Check that asset is expected result
}
