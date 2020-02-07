package odsrestclient

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

//Client : Client to consume ODS Rest API
type Client struct {
	baseURL    *url.URL
	userAgent  string
	httpClient *http.Client
}

// ODSParameters interface to convert parameters into url.Values type
type ODSParameters interface {
	Values() *url.Values
}

//NewClient create a ODS Client
func NewClient(httpClient *http.Client, baseURL *url.URL) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		httpClient: httpClient,
		userAgent:  "golang ods rest client",
		baseURL:    baseURL,
	}
}

//DatasetSearch Query the dataset DB
func (c Client) DatasetSearch(p DatasetSearchParameters) (*Catalog, error) {

	catalog := &Catalog{}
	if err := c.decode("/api/datasets/1.0/search", p, catalog); err != nil {
		return nil, err
	}

	log.Println(catalog)

	return catalog, nil
}

//DownloadRecords download dataset records into a file
func (c *Client) DownloadRecords(p RecordsDownloadParameters, file io.Writer) error {
	if err := c.download("/api/records/1.0/download", p, file); err != nil {
		return err
	}

	return nil
}

// Get : Do a Http get Request on a endpoint
func (c *Client) decode(endpoint string, p ODSParameters, result interface{}) error {
	resp, err := c.exec(endpoint, p)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}

	return nil
}

func (c *Client) download(endpoint string, p ODSParameters, writer io.Writer) error {
	resp, err := c.exec(endpoint, p)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if _, err = io.Copy(writer, resp.Body); err != nil {
		return err
	}

	return nil
}

func (c *Client) exec(endpoint string, p ODSParameters) (*http.Response, error) {
	req, error := c.buildRequest(endpoint, p)

	if error != nil {
		return nil, error
	}

	return c.doRequest(req)
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
		return nil, err
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

	values, ok := headers[http.CanonicalHeaderKey(key)]

	if ok == false {
		return 0
	}

	limit, err := strconv.ParseInt(values[0], 10, 16)

	if err != nil {
		panic(err)
	}

	return uint16(limit)
}
