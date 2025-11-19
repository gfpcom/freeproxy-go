package freeproxy

// Proxy represents a single proxy entry from the GetFreeProxy API
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
