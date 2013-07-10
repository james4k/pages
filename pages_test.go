package pages

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatic(t *testing.T) {
	g := New("testdata", "testdata/layouts")
	ts := httptest.NewServer(g.Static("index.html"))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	expected, err := ioutil.ReadFile("testdata/indexresult.html")
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(actual, expected) != 0 {
		t.Fatal("result does not match!")
	}
}
