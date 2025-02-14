// Package reproxied is a plugin
package reproxied

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/RouxAntoine/reproxied/internal/logging"
)

// Config the plugin configuration.
type Config struct {
	Proxy          string        `json:"proxy"`
	TargetHost     string        `json:"targetHost"`
	KeepHostHeader bool          `json:"keepHostHeader"`
	LogLevel       logging.Level `json:"logLevel,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		LogLevel:       logging.Levels.INFO,
		KeepHostHeader: false,
		Proxy:          "",
		TargetHost:     "",
	}
}

// reProxied a Traefik plugin.
type reProxied struct {
	next   http.Handler
	proxy  *httputil.ReverseProxy
	logger logging.Logger
}

// New creates a new reProxied plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	proxyURL, err := url.Parse(config.Proxy)
	if err != nil {
		return nil, fmt.Errorf("unable to parse proxy url: %w", err)
	}

	transportWithHTTPProxy := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	return NewWithRoundTripperAndWriter(ctx, next, config, name, transportWithHTTPProxy, os.Stdout)
}

// NewWithRoundTripperAndWriter creates a new reProxied plugin.
func NewWithRoundTripperAndWriter(_ context.Context, next http.Handler, config *Config, name string, transport http.RoundTripper, loggingWriter logging.Writer) (http.Handler, error) {
	logger := logging.NewReProxiedLoggerWithLevel(name, loggingWriter, config.LogLevel)
	logger.Debug("plugin called with configuration %+v", config)
	logger.Debug("create logger with level %+v", config.LogLevel)

	targetHostURL, err := url.Parse(config.TargetHost)
	if err != nil {
		return nil, fmt.Errorf("unable to parse target host url: %w", err)
	}

	proxyErrorLogger := log.New(os.Stdout, fmt.Sprintf("%s[reproxied]%s - ", logging.Color.CYAN, logging.Color.CLEAR), log.Ldate|log.Ltime|log.LUTC|log.Llongfile)
	return &reProxied{
		next:   next,
		logger: logger,
		proxy: &httputil.ReverseProxy{
			Transport:      transport,
			ErrorLog:       proxyErrorLogger,
			Director:       createProxyRequest(targetHostURL, config, logger),
			ModifyResponse: modifyResponseLogger(logger),
		},
	}, nil
}

// createProxyRequest build new request send through HTTP Proxy.
func createProxyRequest(targetHostURL *url.URL, config *Config, logger logging.Logger) func(request *http.Request) {
	return func(request *http.Request) {
		request.URL.Scheme = targetHostURL.Scheme
		request.URL.User = targetHostURL.User
		request.URL.Host = targetHostURL.Host
		if !config.KeepHostHeader {
			request.Host = targetHostURL.Host
		}
		logger.Info("proxied req : %+v", request)
	}
}

// modifyResponseLogger logging response callback.
func modifyResponseLogger(logger logging.Logger) func(response *http.Response) error {
	return func(response *http.Response) error {
		logger.Debug("resp : %+v", response)
		return nil
	}
}

// ServeHTTP doing reverse call.
func (c *reProxied) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c.logger.Debug("original req : %+v", req)
	c.proxy.ServeHTTP(rw, req)
}
