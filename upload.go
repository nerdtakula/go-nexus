package nexus

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// UploadParameters for uploading files to nexus
type UploadParameters struct {
	RubyGemsAsset          string `json:"rubygems.asset"`
	NugetAsset             string `json:"nuget.asset"`
	RawDirectory           string `json:"raw.directory"`
	RawAsset1              string `json:"raw.asset1"`
	RawAsset1Filename      string `json:"raw.asset1.filename"`
	RawAsset2              string `json:"raw.asset2"`
	RawAsset2Filename      string `json:"raw.asset2.filename"`
	RawAsset3              string `json:"raw.asset3"`
	RawAsset3Filename      string `json:"raw.asset3.filename"`
	PyPiAsset              string `json:"pypi.asset"`
	NPMAsset               string `json:"npm.asset"`
	Maven2GroupID          string `json:"maven2.groupId"`
	Maven2ArtifactID       string `json:"maven2.artifactId"`
	Maven2Version          string `json:"maven2.version"`
	Maven2GeneratePOM      *bool  `json:"maven2.generate-pom"`
	Maven2Packaging        string `json:"maven2.packaging"`
	Maven2Asset1           string `json:"maven2.asset1"`
	Maven2Asset1Classifier string `json:"maven2.asset1.classifier"`
	Maven2Asset1Extension  string `json:"maven2.asset1.extension"`
	Maven2Asset2           string `json:"maven2.asset2"`
	Maven2Asset2Classifier string `json:"maven2.asset2.classifier"`
	Maven2Asset2Extension  string `json:"maven2.asset2.extension"`
	Maven2Asset3           string `json:"maven2.asset3"`
	Maven2Asset3Classifier string `json:"maven2.asset3.classifier"`
	Maven2Asset3Extension  string `json:"maven2.asset3.extension"`
}

func parseFileUpload(w *multipart.Writer, key, filename string) (io.Writer, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	part, err := w.CreateFormFile(key, filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	return part, nil
}

type maven2File struct {
	Label      string
	AssetPath  string
	Classifier string
	Extension  string
}

func (c Client) uploadMaven2Component(rID string, p UploadParameters) (*Component, error) {
	// Validate parameters
	// Starting with the files
	files := make([]maven2File, 0)
	if p.Maven2Asset1 != "" {
		files = append(files, maven2File{
			Label:      "maven2.asset1",
			AssetPath:  p.Maven2Asset1,
			Classifier: p.Maven2Asset1Classifier,
			Extension:  p.Maven2Asset1Extension,
		})
	}
	if p.Maven2Asset2 != "" {
		files = append(files, maven2File{
			Label:      "maven2.asset2",
			AssetPath:  p.Maven2Asset2,
			Classifier: p.Maven2Asset2Classifier,
			Extension:  p.Maven2Asset2Extension,
		})
	}
	if p.Maven2Asset3 != "" {
		files = append(files, maven2File{
			Label:      "maven2.asset3",
			AssetPath:  p.Maven2Asset3,
			Classifier: p.Maven2Asset3Classifier,
			Extension:  p.Maven2Asset3Extension,
		})
	}

	if len(files) == 0 {
		return nil, errors.Wrap(ErrMissingFiles, "uploadMaven2Component")
	}

	pomSupplied := false
	for _, f := range files {
		if strings.ToLower(f.Extension) == "pom" {
			pomSupplied = true
		}

		if f.Extension == "" {
			return nil, fmt.Errorf("uploadMaven2Component: missing extension for asset '%s'", f.AssetPath)
		}
	}

	if !pomSupplied {
		if p.Maven2GroupID == "" {
			return nil, fmt.Errorf("uploadMaven2Component: missing group id")
		}
		if p.Maven2ArtifactID == "" {
			return nil, fmt.Errorf("uploadMaven2Component: missing artifact id")
		}
		if p.Maven2Version == "" {
			return nil, fmt.Errorf("uploadMaven2Component: missing version")
		}
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, f := range files {
		if _, err := parseFileUpload(writer, f.Label, f.AssetPath); err != nil {
			return nil, errors.Wrap(err, "uploadMaven2Component - parseFileUpload")
		}
		_ = writer.WriteField(fmt.Sprintf("%s.classifier", f.Label), f.Classifier)
		_ = writer.WriteField(fmt.Sprintf("%s.extension", f.Label), f.Extension)
	}

	// ---
	if p.Maven2GroupID != "" {
		_ = writer.WriteField("maven2.groupId", p.Maven2GroupID)
	}
	if p.Maven2ArtifactID != "" {
		_ = writer.WriteField("maven2.artifactId", p.Maven2ArtifactID)
	}
	if p.Maven2Version != "" {
		_ = writer.WriteField("maven2.version", p.Maven2Version)
	}
	if p.Maven2GeneratePOM != nil {
		_ = writer.WriteField("maven2.generate-pom", strconv.FormatBool(*p.Maven2GeneratePOM))
	}
	if p.Maven2Packaging != "" {
		_ = writer.WriteField("maven2.packaging", p.Maven2Packaging)
	}

	headers := map[string]string{"Content-Type": writer.FormDataContentType()}

	// Close writer so we can pass it to the http request
	if err := writer.Close(); err != nil {
		return nil, errors.Wrap(err, "uploadMaven2Component - close")
	}

	err := c.makeMultiPartRequest("POST", "/components", map[string]interface{}{"repository": rID}, headers, body, nil)
	if err != nil {
		return nil, err
	}

	// Query the artifact
	parameters := SearchParameters{
		MavenGroupID:     p.Maven2GroupID,
		MavenArtifactID:  p.Maven2ArtifactID,
		MavenBaseVersion: p.Maven2Version,
	}
	// log.Printf("Params: %#v\n", parameters)

	cpnts, _, err := c.SearchComponents(parameters)
	if err != nil {
		return nil, err
	}
	if len(cpnts) == 0 {
		return nil, ErrNotFound
	}
	return &cpnts[0], nil
}

type rawFile struct {
	Label        string
	SourceAsset  string
	DestFileName string
}

func (c Client) uploadRawComponent(rID string, p UploadParameters) (*Component, error) {
	// Validate parameters
	// Starting with the files
	files := make([]rawFile, 0)
	if p.RawAsset1 != "" {
		files = append(files, rawFile{
			Label:        "raw.asset1",
			SourceAsset:  p.RawAsset1,
			DestFileName: p.RawAsset1Filename,
		})
	}
	if p.RawAsset2 != "" {
		files = append(files, rawFile{
			Label:        "raw.asset2",
			SourceAsset:  p.RawAsset2,
			DestFileName: p.RawAsset2Filename,
		})
	}
	if p.RawAsset3 != "" {
		files = append(files, rawFile{
			Label:        "raw.asset3",
			SourceAsset:  p.RawAsset3,
			DestFileName: p.RawAsset3Filename,
		})
	}

	if len(files) == 0 {
		return nil, errors.Wrap(ErrMissingFiles, "uploadRawComponent")
	}

	if p.RawDirectory == "" {
		return nil, fmt.Errorf("uploadRawComponent: missing upload directory")
	}

	// Process the information
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, f := range files {
		if _, err := parseFileUpload(writer, f.Label, f.SourceAsset); err != nil {
			return nil, errors.Wrap(err, "uploadRawComponent - parseFileUpload")
		}
		_ = writer.WriteField(fmt.Sprintf("%s.filename", f.Label), f.DestFileName)
	}

	_ = writer.WriteField("raw.directory", p.RawDirectory)

	headers := map[string]string{"Content-Type": writer.FormDataContentType()}

	// Close writer so we can pass it to the http request
	if err := writer.Close(); err != nil {
		return nil, errors.Wrap(err, "uploadRawComponent - close")
	}

	err := c.makeMultiPartRequest("POST", "/components", map[string]interface{}{"repository": rID}, headers, body, nil)
	if err != nil {
		return nil, err
	}

	// Query the artifact
	parameters := SearchParameters{
		Format: "raw",
		Query:  files[0].DestFileName,
	}
	// log.Printf("Params: %#v\n", parameters)

	cpnts, _, err := c.SearchComponents(parameters)
	if err != nil {
		return nil, err
	}
	if len(cpnts) == 0 {
		return nil, ErrNotFound
	}
	return &cpnts[0], nil
}

func (c Client) uploadPyPiComponent(rID string, p UploadParameters) (*Component, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c Client) uploadRubyGemsComponent(rID string, p UploadParameters) (*Component, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c Client) uploadNugetComponent(rID string, p UploadParameters) (*Component, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c Client) uploadNPMComponent(rID string, p UploadParameters) (*Component, error) {
	return nil, fmt.Errorf("not implemented")
}
