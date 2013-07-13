package pages

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkStatic(b *testing.B) {
	benchStatic(b, false)
}

func BenchmarkStaticPrecached(b *testing.B) {
	benchStatic(b, true)
}

// naive bench just to see what we look like
func benchStatic(b *testing.B, precache bool) {
	g := New("testdata", "testdata/layouts")
	g.SetPrecache(precache)
	ts := httptest.NewServer(g.Static("index.html"))
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(ts.URL)
		if err != nil {
			b.Fatal(err)
		}

		_, err = io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		if err != nil {
			b.Fatal(err)
		}
	}
}
