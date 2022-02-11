package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestHealthCheck(t *testing.T) {

	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	
	code, _, body := ts.get(t, "/v1/healthcheck")


	rs, err := ts.Client().Get(ts.URL + "/v1/healthcheck")
	if err != nil {
	t.Fatal(err)
	}

	if code != http.StatusOK {
	t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()

	var testData struct {
		StatusCode string `json:"status"`
	}

	err = json.Unmarshal(body, &testData)

	if err != nil {
		t.Fatal(err)
	}

	if testData.StatusCode != "available" {
	t.Errorf("want status to equal %q", "available")
	}
}