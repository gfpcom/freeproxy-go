package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gfpcom/freeproxy-go"
)

func main() {
	// Create a custom HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Initialize the client with your API key
	apiKey := "your-api-key-here" // Replace with your actual API key
	c := freeproxy.NewWith(apiKey, httpClient)

	// Example 1: Get all proxies on page 1
	fmt.Println("Example 1: Get all proxies on page 1")
	response, err := c.Query(freeproxy.WithPage(1))
	if err != nil {
		log.Fatalf("Failed to get proxies: %v", err)
	}

	fmt.Printf("Got %d proxies\n", len(response))
	for i, proxy := range response {
		fmt.Printf("Proxy %d: %s://%s:%d (Country: %s)\n",
			i+1, proxy.Protocol, proxy.IP, proxy.Port, proxy.CountryCode)
	}

	// Example 2: Get US proxies with SOCKS5 protocol
	fmt.Println("\nExample 2: Get US proxies with SOCKS5 protocol on page 1")
	response, err = c.Query(
		freeproxy.WithCountry("US"),
		freeproxy.WithProtocol("socks5"),
		freeproxy.WithPage(1),
	)
	if err != nil {
		log.Fatalf("Failed to get proxies: %v", err)
	}

	fmt.Printf("Got %d US SOCKS5 proxies\n", len(response))
	for i, proxy := range response {
		if i < 3 { // Print first 3
			fmt.Printf("Proxy %d: %s (Uptime: %d%%, Response time: %.2fs)\n",
				i+1, proxy.ProxyUrl, proxy.Uptime, proxy.ResponseTime)
		}
	}

	// Example 3: Get proxies from specific country on page 2
	fmt.Println("\nExample 3: Get proxies from GB on page 2")
	response, err = c.Query(
		freeproxy.WithCountry("GB"),
		freeproxy.WithPage(2),
	)
	if err != nil {
		log.Fatalf("Failed to get proxies: %v", err)
	}

	fmt.Printf("Got %d proxies from GB on page 2\n", len(response))

	// Example 4: Error handling - invalid API key
	fmt.Println("\nExample 4: Error handling with invalid API key")
	invalidClient := freeproxy.New("invalid-api-key")
	_, err = invalidClient.Query()
	if err != nil {
		fmt.Printf("Error occurred (as expected): %v\n", err)
	}
}
