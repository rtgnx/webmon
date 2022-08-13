package webmon

import (
	"fmt"
	"net/http"
	"net/url"
)

type HTTPProbe struct {
	url *url.URL
}

func NewHTTPProbe(target string) (probe *HTTPProbe, err error) {
	probe = new(HTTPProbe)

	probe.url, err = url.Parse(target)
	return
}

func (p *HTTPProbe) Ping() error {
	res, err := http.DefaultClient.Get(p.url.String())
	if err != nil {
		return err
	}

	if (res.StatusCode - 200) > 100 {
		return fmt.Errorf("HTTP %d %s", res.StatusCode, p.url.String())
	}

	return nil
}
