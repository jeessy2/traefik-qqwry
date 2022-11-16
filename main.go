// Package plugindemo a demo plugin.
package traefik_qqwry

import (
	"context"
	"net"
	"net/http"
	"time"

	cache "github.com/patrickmn/go-cache"
)

const (
	// RealIPHeader real ip header.
	RealIPHeader       = "X-Real-IP"
	DefaultCacheExpire = 30 * time.Minute
	DefaultCachePurge  = 2 * time.Hour
)

// Headers part of the configuration
type Headers struct {
	City string `json:"city"`
	ISP  string `json:"isp"`
}

// IpResult Ip result.
type IpResult struct {
	City string
	ISP  string
}

// Config the plugin configuration.
type Config struct {
	DBPath  string   `json:"dbPath,omitempty"`
	Headers *Headers `json:"headers"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		DBPath:  "qqwry.dat",
		Headers: &Headers{City: "X-Cz-City", ISP: "X-Cz-Isp"},
	}
}

// TraefikQQWry a Demo plugin.
type TraefikQQWry struct {
	next    http.Handler
	name    string
	headers *Headers
	cache   *cache.Cache
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	LoadFile(config.DBPath)
	return &TraefikQQWry{
		next:    next,
		name:    name,
		headers: config.Headers,
		cache:   cache.New(DefaultCacheExpire, DefaultCachePurge),
	}, nil
}

func (a *TraefikQQWry) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	ipStr := req.Header.Get(RealIPHeader)
	if ipStr == "" {
		ipStr = req.RemoteAddr
		tmp, _, err := net.SplitHostPort(ipStr)
		if err == nil {
			ipStr = tmp
		}
	}

	var (
		result *IpResult
	)

	if c, found := a.cache.Get(ipStr); found {
		result = c.(*IpResult)
	} else {
		city, isp, err := QueryIP(ipStr)
		if err == nil {
			result = &IpResult{City: city, ISP: isp}
			a.cache.Set(ipStr, result, cache.DefaultExpiration)
		}
	}

	a.addHeaders(req, result)

	a.next.ServeHTTP(rw, req)
}

func (a *TraefikQQWry) addHeaders(req *http.Request, result *IpResult) {
	if result != nil {
		req.Header.Add(a.headers.City, result.City)
		req.Header.Add(a.headers.ISP, result.ISP)
	} else {
		req.Header.Add(a.headers.City, "NotFound")
		req.Header.Add(a.headers.ISP, "NotFound")
	}

}
