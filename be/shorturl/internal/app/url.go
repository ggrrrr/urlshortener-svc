package app

import (
	"fmt"
	"net/url"
)

func parseUrl(fromUrl string) (string, error) {
	u, err := url.Parse(fromUrl)
	if err != nil {
		return "", err
	}

	fmt.Printf("parseUrlparseUrlparseUrl %#v \n", u)
	query := ""

	switch u.Scheme {
	case "http", "https":
	default:
		return "", fmt.Errorf("bad URL schema")
	}

	if u.RawQuery != "" {
		query = fmt.Sprintf("?%s", u.RawQuery)
	}

	return fmt.Sprintf("%s://%s%s%s", u.Scheme, u.Host, u.Path, query), nil
}
