package nexus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

var (
	// ErrNotFound when searching for something returns no results
	ErrNotFound = errors.New("not found")
	// ErrUnknownRepoFormat when we don't know what the format of a repo is
	ErrUnknownRepoFormat = errors.New("can't handle unknown repo format")
	// ErrMissingFiles when we expect files but don't find any
	ErrMissingFiles = errors.New("expecting files, but none were found")
)

// Client hander for making REST API calls
type Client struct {
	uri      *url.URL
	username string
	password string
}

// New Client handler
func New(nexusRestURL string) (Client, error) {
	u, err := url.Parse(nexusRestURL)
	if err != nil {
		return Client{}, err
	}

	return Client{
		uri: u,
	}, nil
}

func (c Client) SetBasicAuth(username, password string) Client {
	return Client{
		uri:      c.uri,
		username: username,
		password: password,
	}
}

// Address returns the address string
func (c Client) Address() string { return c.uri.Host }

func (c Client) url() string {
	return fmt.Sprintf("%s://%s%s", c.uri.Scheme, c.uri.Host, c.uri.Path)
}

func (c Client) makeRequest(method, endpoint string, args map[string]interface{}, result interface{}) error {
	url := c.url() + endpoint
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	for key, value := range args {
		q.Add(key, fmt.Sprintf("%v", value))
	}
	req.URL.RawQuery = q.Encode()

	httpClient := http.Client{
		Timeout: time.Second * 5,
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// fmt.Printf("makeRequest: URL     (%s): %s\n", endpoint, req.URL.String())
	// fmt.Printf("makeRequest: Status  (%s): %s\n", endpoint, res.Status)
	// fmt.Printf("makeRequest: Headers (%s): %s\n", endpoint, req.Header.Get("Content-Type"))
	// fmt.Printf("makeRequest: Body    (%s): %s\n", endpoint, body)

	fmt.Printf("Client:makeRequest -> [%s] --> [%s] \n%s\n\n", endpoint, args, body)
	return json.Unmarshal(body, result)
}

func (c Client) makeMultiPartRequest(method, endpoint string, args map[string]interface{}, headers map[string]string, body *bytes.Buffer, result interface{}) error {
	if c.username == "" {
		return fmt.Errorf("missing user authentication for upload")
	}

	url := c.url() + endpoint
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Accept", "application/json")

	req.SetBasicAuth(c.username, c.password)

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	// req.Header.Set("Content-Type", "multipart/form-data")

	q := req.URL.Query()
	for key, value := range args {
		q.Add(key, fmt.Sprintf("%v", value))
	}
	req.URL.RawQuery = q.Encode()

	// log.Printf("url: %s", req.URL)
	// log.Printf("message: %s", body)

	httpClient := http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "makeMultiPartRequest")
	}
	defer res.Body.Close()
	log.Printf("Upload resp: %#v\n", res)

	rbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "makeMultiPartRequest")
	}

	// log.Printf("response: %s", rbody)
	if result == nil {
		return nil
	}
	return json.Unmarshal(rbody, result)
}

// Ping is used to test we can connect to the service
func (c Client) Ping() error {
	var result map[string]interface{}
	return c.makeRequest("GET", "/read-only", nil, &result)
}
