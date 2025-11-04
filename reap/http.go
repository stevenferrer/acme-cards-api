package reap

import (
	"net"
	"net/http"
	"time"
)

// See https://goperf.dev/02-networking/efficient-net-use/#transport-tuning-when-defaults-arent-enough
func newHTTPClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConnsPerHost: 100,
		MaxIdleConns:        10,
		IdleConnTimeout:     time.Minute,
		Dial: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}
}
