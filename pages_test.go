package pages

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatic(t *testing.T) {
	g := Group{
		Dir:        "testdata",
		LayoutsDir: "testdata/layouts",
	}
	ts := httptest.NewServer(g.Handler("index.html", nil))
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
