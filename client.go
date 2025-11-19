package freeproxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	// ProxiesURL is the full URL for getting proxies from the GetFreeProxy API
	ProxiesURL = "https://api.getfreeproxy.com/v1/proxies"
)

// Client represents a GetFreeProxy API client
type Client struct {
	apiKey     string
	httpClient *http.Client
	proxiesURL string // internal field for testing
}

// New creates a new GetFreeProxy API client with default HTTP client
// apiKey: Your API key from GetFreeProxy
func New(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
		proxiesURL: ProxiesURL,
	}
}

// NewWith creates a new GetFreeProxy API client with custom HTTP client
// apiKey: Your API key from GetFreeProxy
// httpClient: Custom HTTP client (configure timeouts and other settings before passing)
func NewWith(apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		apiKey:     apiKey,
		httpClient: httpClient,
		proxiesURL: ProxiesURL,
	}
}

// buildURL constructs the request URL with query parameters
func (c *Client) buildURL(params *QueryParams) string {
	u, _ := url.Parse(c.proxiesURL)
	q := u.Query()
	if params.Country != nil {
		q.Add("country", *params.Country)
	}
	if params.Protocol != nil {
		q.Add("protocol", *params.Protocol)
	}
	if params.Page != nil {
		q.Add("page", strconv.Itoa(*params.Page))
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// handleErrorResponse parses error responses from the API
func handleErrorResponse(statusCode int, body []byte) error {
	var errorResp map[string]string
	if err := json.Unmarshal(body, &errorResp); err != nil {
		return &Error{
			Message: fmt.Sprintf("http error: status code %d", statusCode),
		}
	}
	if errMsg, ok := errorResp["error"]; ok {
		return &Error{Message: errMsg}
	}
	return &Error{
		Message: fmt.Sprintf("http error: status code %d", statusCode),
	}
}

// Query retrieves a list of proxies from the API with optional filters
// options: Variable number of QueryOption functions to set filters
// Returns a slice of Proxy or an error
func (c *Client) Query(options ...QueryOption) ([]Proxy, error) {
	params := &QueryParams{}
	for _, opt := range options {
		opt(params)
	}

	// Build request
	url := c.buildURL(params)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Handle errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, handleErrorResponse(resp.StatusCode, body)
	}

	// Parse success response
	var proxies []Proxy
	if err := json.Unmarshal(body, &proxies); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return proxies, nil
}

// QueryCountry retrieves proxies filtered by country
func (c *Client) QueryCountry(country string) ([]Proxy, error) {
	return c.Query(WithCountry(country))
}

// QueryProtocol retrieves proxies filtered by protocol
func (c *Client) QueryProtocol(protocol string) ([]Proxy, error) {
	return c.Query(WithProtocol(protocol))
}

// QueryPage retrieves proxies from a specific page
func (c *Client) QueryPage(page int) ([]Proxy, error) {
	return c.Query(WithPage(page))
}
