// Package userip provides functions for extracting a user IP address from a
// request and associating it with a Context.
package backend

import (
	"fmt"
	"net"
	"net/http"

	"golang.org/x/net/context"
)

// GetIP: ip address in the ip:port string
func GetIP(ipport string) (net.IP, error) {
	ip, _, err := net.SplitHostPort(ipport)
	if err != nil {
		return nil, fmt.Errorf("userip: %q cannot split IP IP:port", ipport)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		ips, err := net.LookupIP(ip)
		if err != nil {
			return nil, fmt.Errorf("userip: %q cannot parse IP string", ip)
		}
		userIP = ips[0]
	}
	return userIP, nil
}

// FromRequest extracts the user IP address from req, if present.
func FromRequest(req *http.Request) (net.IP, error) {
	return GetIP(req.RemoteAddr)
}

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

// userIPkey is the context key for the user IP address.  Its value of zero is
// arbitrary.  If this package defined other context keys, they would have
// different integer values.
const userIPKey key = 0

// NewContext returns a new Context carrying userIP.
func NewContext(ctx context.Context, userIP net.IP) context.Context {
	return context.WithValue(ctx, userIPKey, userIP)
}

// FromContext extracts the user IP address from ctx, if present.
func FromContext(ctx context.Context) (net.IP, bool) {
	// ctx.Value returns nil if ctx has no value for the key;
	// the net.IP type assertion returns ok=false for nil.
	userIP, ok := ctx.Value(userIPKey).(net.IP)
	return userIP, ok
}
