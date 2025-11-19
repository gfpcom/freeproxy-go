package freeproxy

// QueryParams holds the query parameters for Query request
type QueryParams struct {
	Country  *string
	Protocol *string
	Page     *int
}

// QueryOption is a function that modifies the query parameters
type QueryOption func(*QueryParams)

// WithCountry sets the country filter
func WithCountry(country string) QueryOption {
	return func(qp *QueryParams) {
		qp.Country = &country
	}
}

// WithProtocol sets the protocol filter
func WithProtocol(protocol string) QueryOption {
	return func(qp *QueryParams) {
		qp.Protocol = &protocol
	}
}

// WithPage sets the page number
func WithPage(page int) QueryOption {
	return func(qp *QueryParams) {
		qp.Page = &page
	}
}
