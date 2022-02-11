package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"vkparser.com/internal/jsonlog"
)


func newTestApplication(t *testing.T) *application {
	return &application{
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}
}

type testServer struct {
	*httptest.Server
}


func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
		if err != nil {
		t.Fatal(err)
		}
	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
		if err != nil {
		t.Fatal(err)
		}
	return rs.StatusCode, rs.Header, body
}