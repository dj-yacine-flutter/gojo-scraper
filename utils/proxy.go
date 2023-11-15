package utils

import (
	"fmt"
	"net/http"
	"net/url"
)

// RoundTripperFunc is a type that implements the http.RoundTripper interface with a function.
type RoundTripperFunc func(req *http.Request) (*http.Response, error)

// RoundTrip implements the RoundTripper interface for RoundTripperFunc.
func (rt RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}

// SwitchableTransport is an http.RoundTripper that switches between proxies for each request.
type SwitchableTransport struct {
	proxies []*url.URL
	index   int
}

// NewSwitchableTransport creates a new SwitchableTransport with the provided proxies.
func NewSwitchableTransport(proxies []*url.URL) *SwitchableTransport {
	return &SwitchableTransport{
		proxies: proxies,
		index:   0,
	}
}

// RoundTrip implements the RoundTripper interface for SwitchableTransport.
// RoundTrip executes a single HTTP transaction, switching between proxies for each request.
func (t *SwitchableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log the current proxy being used
	fmt.Printf("Using proxy: %s\n", t.proxies[t.index])

	// Set the current proxy for this request
	req.URL.Scheme = t.proxies[t.index].Scheme
	req.URL.Host = t.proxies[t.index].Host

	// Switch to the next proxy for the next request
	t.index = (t.index + 1) % len(t.proxies)

	// Create a default transport for making the actual request
	defaultTransport := http.DefaultTransport.(*http.Transport)

	// Make the request using the default transport
	return defaultTransport.RoundTrip(req)
}
