package webcrawler

import (
	"net/http"
	"net/http/httptest"
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

func Test_Downloader_get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(testDoc))
		expect(t, r.Method, "GET")
	}))
	defer ts.Close()

	var api = ts.URL
	res, err := get(api)

	expect(t, err, nil)
	refute(t, res, nil)

	defer res.Body.Close()

	buf := make([]byte, len(testDoc))
	count, err := res.Body.Read(buf)

	expect(t, err, nil)
	expect(t, count, len(testDoc))

}

func Test_DownloadSuccess(t *testing.T) {
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
}
