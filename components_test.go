package nexus

import (
	"io/ioutil"
	"testing"
)

func TestComponents(t *testing.T) {
	components, _, err := client.Components("default")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Results: %+v\n", components)
	// TODO: Test for a known asset in the results
}

func TestUploadMaven2Component(t *testing.T) {
	assetPath := "/tmp/test_asset.txt"

	// Write temp file
	err := ioutil.WriteFile(assetPath, []byte("hello\ngo\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write file, %s", err)
	}

	// Params
	generatePOM := false
	params := UploadParameters{
		Maven2GroupID:         "com.example.test",
		Maven2ArtifactID:      "test",
		Maven2Version:         "0.0.1",
		Maven2GeneratePOM:     &generatePOM,
		Maven2Asset1:          assetPath,
		Maven2Asset1Extension: "txt",
	}

	// Upload file
	_, err = client.UploadComponent("maven-releases", params)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUploadRawComponent(t *testing.T) {
	assetPath := "/tmp/test_asset.txt"

	// Write temp file
	err := ioutil.WriteFile(assetPath, []byte("hello\ngo\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write file, %s", err)
	}

	// Params
	params := UploadParameters{
		RawAsset1:         assetPath,
		RawAsset1Filename: "test_asset.txt",
		RawDirectory:      "/com/example/test",
	}

	// Upload file
	_, err = client.UploadComponent("raw-repo", params)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUploadPyPiComponent(t *testing.T) { t.Skip("Not Implemented") }

func TestUploadRubyGemComponent(t *testing.T) { t.Skip("Not Implemented") }

func TestUploadNugetComponent(t *testing.T) { t.Skip("Not Implemented") }

func TestUploadNPMComponent(t *testing.T) { t.Skip("Not Implemented") }

func TestComponent(t *testing.T) { t.Skip("Not Implemented") }

func TestDeleteComponent(t *testing.T) { t.Skip("Not Implemented") }
