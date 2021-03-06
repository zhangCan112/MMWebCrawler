package downloader

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/* 测试辅助函数 */
func expect(t *testing.T, a, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

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

	res, err := get("%%%\aasdadasdasdass")

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

	var api = ts.URL
	doc, err := Download(api)

	if err != nil {
		t.Errorf("error should be nil:%s", err)
	}
	refute(t, doc, nil)
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

	var api = ts.URL
	//URL错误测试

	doc, err := Download(api + "failed")

	expect(t, reflect.ValueOf(err).IsNil(), false)
	expect(t, reflect.ValueOf(doc).IsNil(), true)

	//URL正确但返回状态码错误
	doc, err = Download(api)

	expect(t, reflect.ValueOf(err).IsNil(), false)
	expect(t, reflect.ValueOf(doc).IsNil(), true)
}
