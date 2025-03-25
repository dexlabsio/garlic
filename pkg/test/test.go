//go:build unit
// +build unit

package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/dexlabsio/garlic/pkg/logging"
	"github.com/wI2L/jsondiff"
	"go.uber.org/zap"
)

type HttpStatus = int

type APITestCase struct {
	t                *testing.T
	title            string
	handler          http.Handler
	request          *http.Request
	requestCtx       *chi.Context
	requestQuery     url.Values
	expectedStatus   HttpStatus
	expectedResponse []byte
	expectedMessage  string
}

func New(t *testing.T, title string) *APITestCase {
	return &APITestCase{
		t:            t,
		title:        title,
		requestCtx:   chi.NewRouteContext(),
		requestQuery: url.Values{},
	}
}

func (tc *APITestCase) Handler(h func(http.ResponseWriter, *http.Request)) *APITestCase {
	tc.handler = http.HandlerFunc(h)
	return tc
}

func (tc *APITestCase) Post(url string, params ...any) *APITestCase {
	tc.request = tc.newRequest("POST", url, params...)
	return tc
}

func (tc *APITestCase) Get(url string, params ...any) *APITestCase {
	tc.request = tc.newRequest("GET", url, params...)
	return tc
}

func (tc *APITestCase) Put(url string, params ...any) *APITestCase {
	tc.request = tc.newRequest("PUT", url, params...)
	return tc
}

func (tc *APITestCase) Delete(url string, params ...any) *APITestCase {
	tc.request = tc.newRequest("DELETE", url, params...)
	return tc
}

func (tc *APITestCase) Param(key, value string) *APITestCase {
	tc.requestCtx.URLParams.Add(key, value)
	return tc
}

func (tc *APITestCase) Query(key, value string) *APITestCase {
	tc.requestQuery.Add(key, value)
	return tc
}

func (tc *APITestCase) Body(data any) *APITestCase {
	jsonData, err := json.Marshal(data)
	if err != nil {
		tc.t.Fatal("Failed to marshal body data.", err)
	}

	tc.request.Body = io.NopCloser(bytes.NewBuffer(jsonData))
	return tc
}

func (tc *APITestCase) ExpectStatus(status HttpStatus) *APITestCase {
	tc.expectedStatus = status
	return tc
}

func (tc *APITestCase) ExpectResponse(jsonData string) *APITestCase {
	tc.expectedResponse = []byte(jsonData)
	return tc
}

func (tc *APITestCase) ExpectMessage(msg string) *APITestCase {
	tc.expectedMessage = msg
	return tc
}

func (tc *APITestCase) ExpectObject(obj any) *APITestCase {
	jsonData, _ := json.Marshal(obj)
	tc.expectedResponse = jsonData
	return tc
}

func (tc *APITestCase) End() {
	tc.t.Run(tc.title, func(t *testing.T) {
		recorder := httptest.NewRecorder()

		// prepare request before serving
		// add the url parmeters in the request context just like chi.Router does
		// add url query params in the request URL just like chi.Router does
		ctx := context.WithValue(tc.request.Context(), chi.RouteCtxKey, tc.requestCtx)

		// Add the logger to the request context (there's a middleware that does this, and the api layer uses it)
		ctx = context.WithValue(ctx, logging.LoggerKey, zap.NewNop())
		ctx = context.WithValue(ctx, logging.RequestIdKey, "test-request-id")
		ctx = context.WithValue(ctx, logging.SessionIdKey, "test-session-id")
		request := tc.request.WithContext(ctx)

		request.URL.RawQuery = tc.requestQuery.Encode()

		tc.handler.ServeHTTP(recorder, request)

		result := recorder.Result()
		// check status code
		if tc.expectedStatus != result.StatusCode {
			t.Errorf("Expected status code `%d` but got `%d`", tc.expectedStatus, result.StatusCode)
		}

		// check response body
		response, err := io.ReadAll(result.Body)
		if err != nil {
			tc.t.Fatal("Failed to read response body.", err)
		}
		diff := tc.compareJsonResponse(tc.expectedResponse, response)
		if diff != "" {
			tc.t.Errorf("Expected response `%s` but got `%s`. Diff %s", tc.expectedResponse, response, diff)
		}
	})
}

func (tc *APITestCase) newRequest(method string, url string, params ...any) *http.Request {
	endpoint := fmt.Sprintf(url, params...)
	r, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		tc.t.Fatal("Failed to create request endpoint.", err)
	}
	return r
}

func (tc *APITestCase) compareJsonResponse(want []byte, got []byte) string {
	patch, err := jsondiff.CompareJSON(want, got)
	if err != nil {
		tc.t.Fatal("Failed to compare json response.", err)
	}
	if patch == nil {
		return ""
	}
	diff, err := json.Marshal(patch)
	if err != nil {
		tc.t.Fatal("Failed to marshal diff operations.", err)
	}

	return string(diff)
}

func AssertModelJsonFields(t *testing.T, dto any, want string) {
	got, err := json.Marshal(dto)
	if err != nil {
		t.Errorf("Expected 'nil', but got '%s'", err)
	}

	if string(got) != want {
		t.Errorf("Expected '%s', but got '%s'", want, got)
	}
}
