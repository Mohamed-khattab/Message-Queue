package utils

import (
	"net/url"
)

func IsValidURL(str string) bool {
	parsedURL, err := url.Parse(str)
	if err != nil || parsedURL == nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}
	return true
}
