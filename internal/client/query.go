package client

import (
	"fmt"
	"net/url"
)

// BuildQuery builds a URL query string from a map, skipping empty values.
func BuildQuery(params map[string]string) string {
	q := url.Values{}
	for k, v := range params {
		if v != "" {
			q.Set(k, v)
		}
	}

	return q.Encode()
}

// Itoa converts an int to its string representation.
func Itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
