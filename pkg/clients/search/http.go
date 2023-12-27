package search

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/buildsafedev/bsf/pkg/version"
)

// Client is the means of connecting to the Search API service
type Client struct {
	BaseURL          *url.URL
	UserAgent        string
	APIKey           string
	LastJSONResponse string

	httpClient *http.Client
}

// Component is a struct to define a User-Agent from a client
type Component struct {
	ID, Name, Version string
}

// HTTPError is the error returned when the API fails with an HTTP error
type HTTPError struct {
	Code   int
	Status string
	Reason string
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("%d: %s, %s", e.Code, e.Status, e.Reason)
}

// NewClientWithURL initializes a Client with a specific API URL
func NewClientWithURL(apiurl string) (*Client, error) {
	parsedURL, err := url.Parse(apiurl)
	if err != nil {
		return nil, err
	}

	var httpTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	client := &Client{
		BaseURL:   parsedURL,
		UserAgent: "bsf/" + version.GetVersion(),
		httpClient: &http.Client{
			Transport: httpTransport,
		},
	}
	return client, nil
}

// SendGetRequest sends a correctly authenticated get request to the API server
func (c *Client) SendGetRequest(requestURL string) ([]byte, error) {
	u := c.prepareClientURL(requestURL)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	return c.sendRequest(req)
}

func (c *Client) prepareClientURL(requestURL string) *url.URL {
	u, _ := url.Parse(c.BaseURL.String() + requestURL)
	return u
}

func (c *Client) sendRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", c.APIKey))

	c.httpClient.Transport = &http.Transport{
		DisableCompression: false,
	}

	if req.Method == "GET" || req.Method == "DELETE" {
		param := req.URL.Query()
		req.URL.RawQuery = param.Encode()
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	c.LastJSONResponse = string(body)

	if resp.StatusCode >= 300 {
		return nil, HTTPError{Code: resp.StatusCode, Status: resp.Status, Reason: string(body)}
	}

	return body, err
}
