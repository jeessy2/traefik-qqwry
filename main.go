// Package plugindemo a demo plugin.
package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/kayon/iploc"
	cache "github.com/patrickmn/go-cache"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	// RealIPHeader real ip header.
	RealIPHeader       = "X-Real-IP"
	DefaultCacheExpire = 30 * time.Minute
	DefaultCachePurge  = 2 * time.Hour
)

// Headers part of the configuration
type Headers struct {
	Country string `json:"country"`
	Region  string `json:"region"`
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
		Headers: &Headers{},
	}
}

// TraefikQQWry a Demo plugin.
type TraefikQQWry struct {
	next    http.Handler
	headers *Headers
	name    string
	loc     *iploc.Locator
	cache   *cache.Cache
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	loc, err := iploc.Open(config.DBPath)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return &TraefikQQWry{
		headers: config.Headers,
		next:    next,
		name:    name,
		loc:     loc,
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
		detail *iploc.Detail
	)

	if c, found := a.cache.Get(ipStr); found {
		detail = c.(*iploc.Detail)
	} else {
		detail = a.loc.Find(ipStr)
		a.cache.Set(ipStr, detail, cache.DefaultExpiration)
	}

	a.addHeaders(req, &detail.Location)

	a.next.ServeHTTP(rw, req)
}

func (a *TraefikQQWry) addHeaders(req *http.Request, detail *iploc.Location) {
	req.Header.Add(a.headers.Country, gb18030Decode([]byte(detail.Country)))
	req.Header.Add(a.headers.Country, gb18030Decode([]byte(detail.Region)))
}

func gb18030Decode(src []byte) string {
	in := bytes.NewReader(src)
	out := transform.NewReader(in, simplifiedchinese.GB18030.NewDecoder())
	d, _ := ioutil.ReadAll(out)
	return string(d)
}
