package server

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func newTestRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func Test_application_ping(t *testing.T) {
	// enable parallel test
	t.Parallel()

	app := newTestApplication(t)

	// We then use the httptest.NewServer() function to create a new test
	// server, passing in the value returned by our app.routes() method as the
	// handler for the server. This starts up a HTTPS server which listens on a
	// randomly-chosen port of your local machine for the duration of the test.
	// Notice that we defer a call to ts.Close() to shutdown the server when
	// the test finishes.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{
			name:     "OK",
			urlPath:  "/ping",
			wantCode: http.StatusOK,
			wantBody: []byte("OK"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// The network address that the test server is listening on is contained
			// in the ts.URL field. We can use this along with the ts.Client().Get()
			// method to make a GET /ping request against the test server. This
			// returns a http.Response struct containing the response.
			code, _, body := ts.get(t, tt.urlPath, tt.wantBody != nil)

			// We can then examine the http.Response to check that the status code
			// written by the ping handler was 200.
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if tt.wantBody != nil && !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}

func Test_application_showSnippet(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked
	// dependencies.
	app := newTestApplication(t)
	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Set up some table-driven tests to check the responses sent by our
	// application for different URLs.
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath, tt.wantBody != nil)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if tt.wantBody != nil && !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}
