package odsrestclient

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

//NewClient create a ODS Client
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{httpClient: httpClient}
}

//Client : Client to consume ODS Rest API
type Client struct {
	baseURL    *url.URL
	userAgent  string
	httpClient *http.Client
}

func (c *Client) DatasetSearch(p ODSParameters) (*Catalog, error) {

	catalog := &Catalog{}
	if err := c.Get("/api/datasets/1.0/search/", p, catalog); err != nil {
		return nil, err
	}

	return catalog, nil
}

// Get : Do a Http get Request on a endpoint
func (c *Client) Get(endpoint string, p ODSParameters, result interface{}) error {
	req, error := c.buildRequest(endpoint, p)
	if error != nil {
		return error
	}

	resp, error := c.doRequest(req)
	defer resp.Body.Close()

	if error != nil {
		return error
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}

	return nil
}

func (c *Client) buildRequest(endpoint string, p ODSParameters) (*http.Request, error) {

	url := c.baseURL.ResolveReference(
		&url.URL{
			Path:     endpoint,
			RawQuery: p.Values().Encode(),
		},
	)

	return http.NewRequest(
		"GET",
		url.String(),
		nil,
	)
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return resp, err
	}
	if resp.StatusCode == 400 {
		resp.Body.Close()

		return resp, &RateLimitError{
			extractLimit(resp.Header, "X-Ratelimit-Limit"),
			extractLimit(resp.Header, "X-Ratelimit-Remaining"),
			extractLimit(resp.Header, "X-Ratelimit-Reset"),
		}
	}

	return resp, nil
}

func extractLimit(headers map[string][]string, key string) uint16 {
	limit, err := strconv.ParseInt(headers[http.CanonicalHeaderKey(key)][0], 10, 16)

	if err != nil {
		panic(err)
	}

	return uint16(limit)
}
