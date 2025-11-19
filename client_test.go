package freeproxy

import (
"encoding/json"
"net/http"
"net/http/httptest"
"testing"
)

func TestNew(t *testing.T) {
apiKey := "test-api-key"

client := New(apiKey)

if client.apiKey != apiKey {
t.Errorf("expected apiKey %q, got %q", apiKey, client.apiKey)
}
if client.httpClient == nil {
t.Error("expected default httpClient to be set, got nil")
}
if client.proxiesURL != ProxiesURL {
t.Errorf("expected proxiesURL %q, got %q", ProxiesURL, client.proxiesURL)
}
}

func TestNewWith(t *testing.T) {
apiKey := "test-api-key"
httpClient := &http.Client{}

client := NewWith(apiKey, httpClient)

if client.apiKey != apiKey {
t.Errorf("expected apiKey %q, got %q", apiKey, client.apiKey)
}
if client.httpClient != httpClient {
t.Error("expected custom httpClient to be set")
}
if client.proxiesURL != ProxiesURL {
t.Errorf("expected proxiesURL %q, got %q", ProxiesURL, client.proxiesURL)
}
}

func TestNewWithNilHTTPClient(t *testing.T) {
apiKey := "test-api-key"

client := NewWith(apiKey, nil)

if client.httpClient == nil {
t.Error("expected default httpClient to be set, got nil")
}
}

func TestQueryWithOptions(t *testing.T) {
mockResponse := []Proxy{
{
ID:           "1",
Protocol:     "socks5",
IP:           "192.168.1.1",
Port:         1080,
User:         "user1",
Passwd:       "pass1",
CountryCode:  "US",
Region:       "New York",
AsnNumber:    "AS1234",
AsnName:      "Test ASN",
Anonymity:    "Elite",
Uptime:       99,
ResponseTime: 0.5,
LastAliveAt:  "2025-11-18T10:00:00Z",
ProxyUrl:     "socks5://user1:pass1@192.168.1.1:1080",
Https:        true,
Google:       true,
},
}

server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
if r.Header.Get("Authorization") != "Bearer test-api-key" {
t.Error("expected Authorization header with Bearer token")
}

query := r.URL.Query()
if query.Get("country") != "US" {
t.Errorf("expected country query param US, got %s", query.Get("country"))
}
if query.Get("protocol") != "socks5" {
t.Errorf("expected protocol query param socks5, got %s", query.Get("protocol"))
}
if query.Get("page") != "1" {
t.Errorf("expected page query param 1, got %s", query.Get("page"))
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(mockResponse)
}))
defer server.Close()

client := New("test-api-key")
client.proxiesURL = server.URL

response, err := client.Query(
WithCountry("US"),
WithProtocol("socks5"),
WithPage(1),
)

if err != nil {
t.Fatalf("expected no error, got %v", err)
}
if response == nil {
t.Fatal("expected response, got nil")
}
if len(response) != 1 {
t.Errorf("expected 1 proxy, got %d", len(response))
}
if response[0].ID != "1" {
t.Errorf("expected proxy ID 1, got %s", response[0].ID)
}
}

func TestQueryNoOptions(t *testing.T) {
mockResponse := []Proxy{}

server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(mockResponse)
}))
defer server.Close()

client := New("test-api-key")
client.proxiesURL = server.URL

response, err := client.Query()

if err != nil {
t.Fatalf("expected no error, got %v", err)
}
if response == nil {
t.Fatal("expected response, got nil")
}
if len(response) != 0 {
t.Errorf("expected 0 proxies, got %d", len(response))
}
}

func TestQueryError(t *testing.T) {
errorResponse := map[string]string{
"error": "INVALID_PARAMETER",
}

server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
w.WriteHeader(http.StatusBadRequest)
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(errorResponse)
}))
defer server.Close()

client := New("test-api-key")
client.proxiesURL = server.URL

response, err := client.Query()

if err == nil {
t.Fatal("expected error, got nil")
}
if response != nil {
t.Errorf("expected nil response, got %v", response)
}

apiErr, ok := err.(*Error)
if !ok {
t.Errorf("expected Error, got %T", err)
}
if apiErr.Message != "INVALID_PARAMETER" {
t.Errorf("expected error message 'INVALID_PARAMETER', got %q", apiErr.Message)
}
}

func TestQueryHTTPError(t *testing.T) {
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
w.WriteHeader(http.StatusInternalServerError)
w.Write([]byte("Internal Server Error"))
}))
defer server.Close()

client := New("test-api-key")
client.proxiesURL = server.URL

response, err := client.Query()

if err == nil {
t.Fatal("expected error, got nil")
}
if response != nil {
t.Errorf("expected nil response, got %v", response)
}

apiErr, ok := err.(*Error)
if !ok {
t.Errorf("expected Error, got %T", err)
}
if apiErr.Message == "" {
t.Error("expected error message, got empty string")
}
}

func TestWithCountry(t *testing.T) {
params := &QueryParams{}
opt := WithCountry("US")
opt(params)
if params.Country == nil || *params.Country != "US" {
t.Errorf("expected country to be 'US', got %v", params.Country)
}
}

func TestWithProtocol(t *testing.T) {
params := &QueryParams{}
opt := WithProtocol("https")
opt(params)
if params.Protocol == nil || *params.Protocol != "https" {
t.Errorf("expected protocol to be 'https', got %v", params.Protocol)
}
}

func TestWithPage(t *testing.T) {
	params := &QueryParams{}
	opt := WithPage(2)
	opt(params)
	if params.Page == nil || *params.Page != 2 {
		t.Errorf("expected page to be 2, got %v", params.Page)
	}
}

func TestQueryCountry(t *testing.T) {
	mockResponse := []Proxy{
		{
			ID:          "1",
			Protocol:    "http",
			IP:          "192.168.1.1",
			Port:        8080,
			CountryCode: "US",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("country") != "US" {
			t.Errorf("expected country query param US, got %s", r.URL.Query().Get("country"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := New("test-api-key")
	client.proxiesURL = server.URL

	response, err := client.QueryCountry("US")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(response) != 1 {
		t.Errorf("expected 1 proxy, got %d", len(response))
	}
}

func TestQueryProtocol(t *testing.T) {
	mockResponse := []Proxy{
		{
			ID:       "1",
			Protocol: "https",
			IP:       "192.168.1.1",
			Port:     443,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("protocol") != "https" {
			t.Errorf("expected protocol query param https, got %s", r.URL.Query().Get("protocol"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := New("test-api-key")
	client.proxiesURL = server.URL

	response, err := client.QueryProtocol("https")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(response) != 1 {
		t.Errorf("expected 1 proxy, got %d", len(response))
	}
}

func TestQueryPage(t *testing.T) {
	mockResponse := []Proxy{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "2" {
			t.Errorf("expected page query param 2, got %s", r.URL.Query().Get("page"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := New("test-api-key")
	client.proxiesURL = server.URL

	response, err := client.QueryPage(2)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(response) != 0 {
		t.Errorf("expected 0 proxies, got %d", len(response))
	}
}