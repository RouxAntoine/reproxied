package reproxied_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/RouxAntoine/reproxied"
	"github.com/RouxAntoine/reproxied/internal/logging"
)

type HTTPMocking struct {
	executedRequest []*http.Request
}

func (mock *HTTPMocking) RoundTrip(req *http.Request) (*http.Response, error) {
	mock.executedRequest = append(mock.executedRequest, req)
	return &http.Response{Body: io.NopCloser(strings.NewReader("")), StatusCode: http.StatusOK}, nil
}

func TestShouldChangeHost(t *testing.T) {
	httpMock := &HTTPMocking{
		executedRequest: nil,
	}
	cfg := reproxied.CreateConfig()
	cfg.Proxy = "http://proxy:3128"
	cfg.TargetHost = "https://target.com"
	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := reproxied.NewWithRoundTripperAndWriter(ctx, next, cfg, "reProxied", httpMock, os.Stdout)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://internal.url/", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	if httpMock.executedRequest[0].Host != "target.com" {
		t.Errorf("expected request host to be updated to \"target.com\" but was actually: %v", httpMock.executedRequest[0].Host)
	}
}

func TestShouldKeepHostHeader(t *testing.T) {
	httpMock := &HTTPMocking{
		executedRequest: nil,
	}

	cfg := reproxied.CreateConfig()
	cfg.Proxy = "http://proxy:3128"
	cfg.TargetHost = "https://target.com"
	cfg.KeepHostHeader = true

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := reproxied.NewWithRoundTripperAndWriter(ctx, next, cfg, "reProxied", httpMock, os.Stdout)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://internal.url/", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	if httpMock.executedRequest[0].URL.Host != "target.com" {
		t.Errorf("expected request host to be updated to \"target.com\" but was actually: %v", httpMock.executedRequest[0].URL.Host)
	}
	if httpMock.executedRequest[0].Host != "internal.url" {
		t.Errorf("expected request header host to be kept to \"internal.url\" but was actually: %v", httpMock.executedRequest[0].Host)
	}
}

func TestShouldCustomizeLogLevel(t *testing.T) {
	httpMock := &HTTPMocking{
		executedRequest: nil,
	}

	cfg := reproxied.CreateConfig()
	cfg.Proxy = "http://proxy:3128"
	cfg.TargetHost = "https://target.com"
	cfg.KeepHostHeader = true
	cfg.LogLevel = logging.Levels.DEBUG

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	var byteBuffer bytes.Buffer

	_, err := reproxied.NewWithRoundTripperAndWriter(ctx, next, cfg, "reProxied", httpMock, &byteBuffer)
	if err != nil {
		t.Fatal(err)
	}

	_, err = http.NewRequestWithContext(ctx, http.MethodGet, "https://internal.url/", nil)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(byteBuffer.String(), "create logger with level 1") || !strings.Contains(byteBuffer.String(), "[DEBUG]") {
		t.Errorf("Expect logger at level DEBUG and log some message")
	}
}

func TestShouldParseConfig(t *testing.T) {
	data := `
		{
		    "proxy": "http://proxy:3128",
			"targetHost": "https://example.com",
			"keepHostHeader": true,
			"logLevel": 1
        }`

	result := reproxied.CreateConfig()
	err := json.Unmarshal([]byte(data), result)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	if result.KeepHostHeader != true {
		t.Errorf("bad KeepHostHeader value : %v", result.KeepHostHeader)
	}
	if result.Proxy != "http://proxy:3128" {
		t.Errorf("bad Proxy value : %v", result.Proxy)
	}
	if result.TargetHost != "https://example.com" {
		t.Errorf("bad TargetHost value : %v", result.TargetHost)
	}
	if result.LogLevel != logging.Levels.DEBUG {
		t.Errorf("bad LogLevel value : %v", result.LogLevel)
	}
}
