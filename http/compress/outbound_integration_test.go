//go:build integration

package compress

import (
	"context"
	"embed"
	"io"
	"io/fs"
	"net"
	"net/http"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

//go:embed testdata/nginx/htdocs
var testData embed.FS

func loadTestData(t *testing.T, filename string) fs.File {
	t.Helper()

	path := filepath.Join("testdata", "nginx", "htdocs", filename)
	f, err := testData.Open(path)
	if err != nil {
		t.Fatalf("error during testdate file opening %s", err.Error())
	}

	return f
}

func defaultClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

func TestNginx(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		URL                   string
		ExpectedTestDataFile  string
		ExpectedContentLength string
		InitClient            func(*http.Client)
	}{
		"gzip-pool": {
			URL:                  "http://127.0.0.1:8888/gzip/",
			ExpectedTestDataFile: "text_1.txt",
			InitClient: func(client *http.Client) {
				client.Transport = RoundTripper(client.Transport, WithCompressionTypeGZIP(true))
			},
		},
		"gzip": {
			URL:                  "http://127.0.0.1:8888/gzip/",
			ExpectedTestDataFile: "text_1.txt",
			InitClient: func(client *http.Client) {
				client.Transport = RoundTripper(client.Transport, WithCompressionTypeGZIP(false))
			},
		},
		"gzip-chunked": {
			URL:                  "http://127.0.0.1:8888/gzip-chunked/",
			ExpectedTestDataFile: "text_2.txt",
			InitClient: func(client *http.Client) {
				client.Transport = RoundTripper(client.Transport, WithCompressionTypeGZIP(true))
			},
		},
		"gzip-pool-chunked": {
			URL:                  "http://127.0.0.1:8888/gzip-chunked/",
			ExpectedTestDataFile: "text_2.txt",
			InitClient: func(client *http.Client) {
				client.Transport = RoundTripper(client.Transport, WithCompressionTypeGZIP(false))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := defaultClient()

			tc.InitClient(client)

			fl := loadTestData(t, tc.ExpectedTestDataFile)
			expected, err := io.ReadAll(fl)
			if err != nil {
				t.Fatalf("error while reading data file: %s", err.Error())
			}

			// for i := range 1500 {
			// 	t.Run(fmt.Sprintf("sub_%d", i), func(t *testing.T) {
			// 		t.Parallel()
			// 		nginxTestCase(t, client, tc.URL, expected)
			// 	})
			// }

			nginxTestCase(t, client, tc.URL, tc.ExpectedContentLength, expected)
		})
	}
}

func nginxTestCase(t *testing.T, client *http.Client, url string, expectedContentLength string, expectedBody []byte) {
	t.Helper()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("error while creating request: %s", err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error while executing request: %s", err.Error())
	}

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error while ReadAll response: %s", err.Error())
	}

	if e := resp.Body.Close(); e != nil {
		t.Fatalf("error while response body closing: %s", err.Error())
	}

	if !reflect.DeepEqual(expectedBody, got) {
		t.Errorf("content mismatch\nexpected: %s\n\n\ngot: %s\n\n\n", string(expectedBody), string(got))
	}
}

func TestGoogle(t *testing.T) {
	client := defaultClient()

	rt := RoundTripper(
		client.Transport,
		WithCompressionTypeDeflate(false),
	)
	client.Transport = rt

	// expectedAcceptEncodingHeader := "gzip, zstd, br"
	// if rt.acceptEncodingHeader != expectedAcceptEncodingHeader {
	// 	t.Fatalf("wring AcceptEncoding header. Expected %s - got %s", expectedAcceptEncodingHeader, rt.acceptEncodingHeader)
	// }
	t.Logf(">>> %s \n", rt.acceptEncodingHeader)

	resp := httpRequestDo(t, client, "https://www.docker.com/blog/faster-multi-platform-builds-dockerfile-cross-compilation-guide/")
	defer resp.Body.Close()
	for k := range resp.Header {
		t.Logf(">>> %s  :  %s\n", k, resp.Header.Get(k))
	}
	t.Logf(">>> %s  :  %s\n", "content-encoding", resp.Header.Get("Content-Encoding"))
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	_ = body

	// t.Logf("\n\n%s\n\n", string(body))
}

func httpRequestDo(t *testing.T, client *http.Client, url string) *http.Response {
	t.Helper()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("error while creating request: %s", err.Error())
	}

	req.Header.Add("Accept", "*")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error while executing request: %s", err.Error())
	}

	return resp
}
