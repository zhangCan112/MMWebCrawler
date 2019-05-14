package webcrawler

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var (
	testDoc = `<html>
		<body>

		<h1>My First Heading</h1>

		<p>My first paragraph.</p>

		</body>
		</html>`
)

func Test_randomUserAgent(t *testing.T) {
	ua1 := randomUserAgent()
	ua2 := randomUserAgent()
	ua3 := randomUserAgent()
	ua4 := randomUserAgent()

	//随机不能连续几个都一样
	expect(t, ((ua1 == ua2) && (ua2 == ua3) && (ua3 == ua4)), false)
}

func Test_Downloader_get_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(testDoc))
		expect(t, r.Method, "GET")
	}))
	defer ts.Close()

	// 正确请求的情况
	var api = ts.URL
	res, err := get(api)

	expect(t, err, nil)
	refute(t, res, nil)

	defer res.Body.Close()

	buf := make([]byte, len([]byte(testDoc)))
	_, err = res.Body.Read(buf)

	expect(t, err.Error(), "EOF")
	expect(t, string(buf), testDoc)
}

func Test_Downloader_get_Failed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(testDoc))
		expect(t, r.Method, "GET")
	}))
	defer ts.Close()

	res, err := get("aasdadasdasdass")

	refute(t, reflect.ValueOf(err).IsNil(), true)
	expect(t, reflect.ValueOf(res).IsNil(), true)
}

func Test_Downloader_DownloadSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(testDoc))
		expect(t, r.Method, "GET")
	}))
	defer ts.Close()

	docReceiver := RunDownloader()
	var api = ts.URL
	Download(api)

	resultFunc := <-docReceiver
	doc, url, err := resultFunc()

	if err != nil {
		t.Errorf("error should be nil:%s", err)
	}
	refute(t, doc, nil)
	expect(t, url, api)
	h1Text := doc.Find("h1").Text()
	expect(t, h1Text, "My First Heading")
}

func Test_Downloader_DownloadFailed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(testDoc))
		expect(t, r.Method, "GET")
	}))
	defer ts.Close()

	docReceiver := RunDownloader()
	var api = ts.URL
	Download(api + "failed")

	resultFunc := <-docReceiver
	doc, _, err := resultFunc()

	expect(t, reflect.ValueOf(err).IsNil(), false)
	expect(t, reflect.ValueOf(doc).IsNil(), true)
}
