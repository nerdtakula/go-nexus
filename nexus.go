package nexus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrEmptyResults            = errors.New("no results found")                           // when searching for something returns no results
	ErrInsufficientPermissions = errors.New("insufficient permissions to preform action") // when user doesn't have required permissions to preform acction
	ErrMalformedID             = errors.New("the supplied id was malformed")              // when the id supplied didn't match expected format
	ErrMissingFiles            = errors.New("expecting files, but none were found")       // when we expect files but don't find any
	ErrNotFound                = errors.New("not found")                                  // when searching for something returns no results
	ErrRepositoryNotFound      = errors.New("repository not found")                       // when specified repository wasn't found
	ErrUnknownRepoFormat       = errors.New("can't handle unknown repo format")           // when we don't know what the format of a repo is
)

// Client hander for making REST API calls
type Client struct {
	uri      *url.URL
	username string
	password string
}

// New Client handler
func New(nexusRestURL string) (*Client, error) {
	u, err := url.Parse(nexusRestURL)
	if err != nil {
		return nil, err
	}
	// TODO: validate the 'nexusRestURL' input to make sure it's in the correct format
	return &Client{uri: u}, nil
}

// SetBasicAuth to log into nexus
func (c *Client) SetBasicAuth(username, password string) {
	c.username = username
	c.password = password
}

// Address returns the address string
func (c Client) Address() string { return c.uri.Host }

func (c Client) url() string {
	// return fmt.Sprintf("%s://%s%s", c.uri.Scheme, c.uri.Host, c.uri.Path)
	return c.uri.String()
}

func (c Client) makeRequest(method, endpoint string, args map[string]interface{}, result interface{}) (statusCode int, err error) {
	url := c.url() + endpoint
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Accept", "application/json")

	// Provide base auth if provided
	if len(strings.TrimSpace(c.username)) != 0 {
		req.SetBasicAuth(c.username, c.password)
	}

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
		return -1, err
	}
	defer res.Body.Close()
	// fmt.Printf("%s: %+v\n", url, res)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return -1, err
	}
	// fmt.Printf("%s Body (%d): %s\n", url, res.StatusCode, body)
	if result == nil {
		return res.StatusCode, nil
	}
	return res.StatusCode, json.Unmarshal(body, result)
}

func (c Client) makeMultiPartRequest(method, endpoint string, args map[string]interface{}, headers map[string]string, body *bytes.Buffer, result interface{}) (statusCode int, err error) {
	if len(strings.TrimSpace(c.username)) == 0 {
		return -1, fmt.Errorf("missing user authentication for upload")
	}

	url := c.url() + endpoint
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Accept", "application/json")

	req.SetBasicAuth(c.username, c.password)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	q := req.URL.Query()
	for key, value := range args {
		q.Add(key, fmt.Sprintf("%v", value))
	}
	req.URL.RawQuery = q.Encode()

	httpClient := http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		return -1, errors.Wrap(err, "makeMultiPartRequest Do")
	}
	defer res.Body.Close()
	// log.Printf("Upload resp: %#v\n", res)

	rbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return -1, errors.Wrap(err, "makeMultiPartRequest ReadAll")
	}

	if result == nil {
		return res.StatusCode, nil
	}
	return res.StatusCode, json.Unmarshal(rbody, result)
}

// Ping is used to test we can connect to the service
func (c Client) Ping() error {
	statusCode, err := c.makeRequest("GET", "/status", nil, nil)
	switch statusCode {
	case 200:
		return nil
	case 401:
		return errors.Wrap(err, "Unauthorized")
	case 503:
		return fmt.Errorf("Unavailable to service requests")
	}
	return err
}
