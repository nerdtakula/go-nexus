package nexus

import (
	"encoding/json"
	"fmt"
	"strings"
)

func structToMap(in interface{}, ignoreEmpty bool) (map[string]interface{}, error) {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(in)
	json.Unmarshal(inrec, &inInterface)
	if !ignoreEmpty {
		return inInterface, nil
	}

	outInterface := make(map[string]interface{})
	for key, value := range inInterface {
		if fmt.Sprintf("%v", value) != "" {
			outInterface[key] = value
		}
	}
	return outInterface, nil
}

// SearchParameters available when preforming a search
type SearchParameters struct {
	ContinuationToken   string `json:"continuationToken"`
	Query               string `json:"q"`
	Format              string `json:"format"`
	Repository          string `json:"repository"`
	Group               string `json:"group"`
	Name                string `json:"name"`
	Version             string `json:"version"`
	MD5                 string `json:"md5"`
	SHA1                string `json:"sha1"`
	SHA256              string `json:"sha256"`
	SHA512              string `json:"sha512"`
	MavenGroupID        string `json:"maven.groupId"`
	MavenArtifactID     string `json:"maven.artifactId"`
	MavenBaseVersion    string `json:"maven.baseVersion"`
	MavenExtension      string `json:"maven.extension"`
	MavenClassifier     string `json:"maven.classifier"`
	NugetID             string `json:"nuget.id"`
	NugetTags           string `json:"nuget.tags"`
	NPMScope            string `json:"npm.scope"`
	DockerImageName     string `json:"docker.imageName"`
	DockerImageTag      string `json:"docker.imageTag"`
	DockerLayerID       string `json:"docker.layerId"`
	DockerContentDigest string `json:"docker.contentDigest"`
	PyPiClassifiers     string `json:"pypi.classifiers"`
	PyPiDescription     string `json:"pypi.description"`
	PyPiKeywords        string `json:"pypi.keywords"`
	PyPiSummary         string `json:"pypi.summary"`
	RubyGemsDescription string `json:"rubygems.description"`
	RubyGemsPlatform    string `json:"rubygems.platform"`
	RubyGemsSummary     string `json:"rubygems.summary"`
}

func searchEscapeVersion(in string) string {
	return strings.Replace(in, ":", "\\:", -1)
}

// SearchComponents via end point
func (c Client) SearchComponents(parameters SearchParameters) ([]Component, string, error) {
	args, _ := structToMap(parameters, true)
	// args["version"] = searchEscapeVersion(fmt.Sprintf("%v", args["version"]))

	result := struct {
		Items             []Component `json:"items"`
		ContinuationToken string      `json:"continuationToken"`
	}{}
	if _, err := c.makeRequest("GET", "/search", args, &result); err != nil {
		return nil, "", err
	}
	return result.Items, result.ContinuationToken, nil
}

// SearchAssets via end point
func (c Client) SearchAssets(parameters SearchParameters) ([]Asset, string, error) {
	args, _ := structToMap(parameters, true)
	// args["version"] = searchEscapeVersion(fmt.Sprintf("%v", args["version"]))

	result := struct {
		Items             []Asset `json:"items"`
		ContinuationToken string  `json:"continuationToken"`
	}{}
	if _, err := c.makeRequest("GET", "/search/assets", args, &result); err != nil {
		return nil, "", err
	}
	return result.Items, result.ContinuationToken, nil
}
