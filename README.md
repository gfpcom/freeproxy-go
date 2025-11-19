# FreeProxy Go Client

[![Go Reference](https://pkg.go.dev/badge/github.com/gfpcom/freeproxy-go)](https://pkg.go.dev/github.com/gfpcom/freeproxy-go) [![Release](https://img.shields.io/github/v/release/gfpcom/freeproxy-go?label=release&style=flat-square)](https://github.com/gfpcom/freeproxy-go/releases) [![Build Status](https://img.shields.io/github/actions/workflow/status/gfpcom/freeproxy-go/test.yml?workflow=Test&style=flat-square)](https://github.com/gfpcom/freeproxy-go/actions) [![Coverage](https://img.shields.io/github/actions/workflow/status/gfpcom/freeproxy-go/test.yml?label=coverage&style=flat-square)](https://github.com/gfpcom/freeproxy-go/actions) [![License](https://img.shields.io/github/license/gfpcom/freeproxy-go?style=flat-square)](https://github.com/gfpcom/freeproxy-go/blob/main/LICENSE) [![Go Version](https://img.shields.io/badge/go-1.25-blue?logo=go&style=flat-square)](https://pkg.go.dev/github.com/gfpcom/freeproxy-go)

A lightweight Go client library for the [GetFreeProxy API](https://developer.getfreeproxy.com/docs).

## Features

- Simple and clean API interface
- Full support for proxy filtering (country, protocol, pagination)
- Custom HTTP client configuration
- Comprehensive error handling
- Full test coverage
- No external dependencies (uses only Go standard library)

## Installation

```bash
go get github.com/gfpcom/freeproxy
```

## Requirements

- Go 1.25 or later

## Quick Start

### Get Your API Key

First, obtain your API key from [GetFreeProxy](https://developer.getfreeproxy.com/docs).

### Basic Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/gfpcom/freeproxy"
)

func main() {
	// Initialize client with your API key
	c := freeproxy.New("your-api-key")

	// Get latest 10 proxies
	proxies, err := c.Query()
	if err != nil {
		log.Fatal(err)
	}

	for _, proxy := range proxies {
		fmt.Printf("%s://%s:%d\n", proxy.Protocol, proxy.IP, proxy.Port)
	}
}
```

## API Documentation

### New

```go
func New(apiKey string) *Client
```

Creates a new GetFreeProxy API client with default HTTP client.

- `apiKey`: Your API key from GetFreeProxy (required)

### NewWith

```go
func NewWith(apiKey string, httpClient *http.Client) *Client
```

Creates a new GetFreeProxy API client with custom HTTP client.

- `apiKey`: Your API key from GetFreeProxy (required)
- `httpClient`: Custom HTTP client. Configure timeouts and other settings before passing.


### Query

```go
func (c *Client) Query(options ...QueryOption) ([]Proxy, error)
```

Retrieves a list of proxies from the API with optional filters using the option pattern.

**Parameters:**
- `options`: Variable number of query option functions (optional)

**Returns:**
- Slice of `Proxy` structures
- Error if the request failed

**Example:**

```go
// Get proxies with multiple filters
proxies, err := c.Query(
	freeproxy.WithCountry("US"),
	freeproxy.WithProtocol("http"),
	freeproxy.WithPage(1),
)
if err != nil {
	log.Fatal(err)
}
```

### QueryCountry

```go
func (c *Client) QueryCountry(country string) ([]Proxy, error)
```

Retrieves proxies filtered by country. Convenience method for common use cases.

**Parameters:**
- `country`: Country code (e.g., "US", "GB", "DE")

**Example:**
```go
proxies, err := c.QueryCountry("US")
```

### QueryProtocol

```go
func (c *Client) QueryProtocol(protocol string) ([]Proxy, error)
```

Retrieves proxies filtered by protocol. Convenience method for common use cases.

**Parameters:**
- `protocol`: Protocol type (e.g., "http", "https", "socks5")

**Example:**
```go
proxies, err := c.QueryProtocol("https")
```

### QueryPage

```go
func (c *Client) QueryPage(page int) ([]Proxy, error)
```

Retrieves proxies from a specific page. Convenience method for pagination.

**Parameters:**
- `page`: Page number (default: 1)

**Example:**
```go
proxies, err := c.QueryPage(2)
```

## Response Structure

The API returns a direct JSON array of proxy objects:

```go
type Proxy struct {
    ID           string  `json:"id"`
    Protocol     string  `json:"protocol"`
    IP           string  `json:"ip"`
    Port         int     `json:"port"`
    User         string  `json:"user"`
    Passwd       string  `json:"passwd"`
    CountryCode  string  `json:"countryCode"`
    Region       string  `json:"region"`
    AsnNumber    string  `json:"asnNumber"`
    AsnName      string  `json:"asnName"`
    Anonymity    string  `json:"anonymity"`
    Uptime       int     `json:"uptime"`
    ResponseTime float64 `json:"responseTime"`
    LastAliveAt  string  `json:"lastAliveAt"`
    ProxyUrl     string  `json:"proxyUrl"`
    Https        bool    `json:"https"`
    Google       bool    `json:"google"`
}
```

## Error Handling

The client returns an `Error` for all error cases:

```go
proxies, err := c.Query(freeproxy.WithCountry("US"))
if err != nil {
    apiErr, ok := err.(*freeproxy.Error)
    if ok {
        fmt.Printf("API Error: %s\n", apiErr.Error())
    }
}
```

The error message comes directly from the API's `error` field in the response.

## Examples

See the `examples/` directory for more complete examples.

### Get Proxies by Country

```go
proxies, err := c.QueryCountry("US")
if err != nil {
    log.Fatal(err)
}

for _, proxy := range proxies {
    fmt.Println(proxy.ProxyUrl)
}
```

### Get Proxies by Protocol

```go
proxies, err := c.QueryProtocol("https")
if err != nil {
    log.Fatal(err)
}
```

### Get Proxies from Specific Page

```go
proxies, err := c.QueryPage(2)
if err != nil {
    log.Fatal(err)
}
```

### Combine Multiple Filters

```go
proxies, err := c.Query(
    freeproxy.WithCountry("US"),
    freeproxy.WithProtocol("socks5"),
    freeproxy.WithPage(1),
)
if err != nil {
    log.Fatal(err)
}
```

### Iterate Through Pages

```go
for page := 1; page <= 10; page++ {
    proxies, err := c.QueryPage(page)
    if err != nil {
        log.Fatal(err)
    }
    
    if len(proxies) == 0 {
        break // No more proxies
    }
    
    for _, proxy := range proxies {
        fmt.Println(proxy.ProxyUrl)
    }
}
```

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## License

See the LICENSE file for details.

## Support

For API documentation and support, visit [GetFreeProxy Developer Docs](https://developer.getfreeproxy.com/docs).
