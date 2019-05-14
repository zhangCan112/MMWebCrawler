package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublishWrongResponseStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hellow world!"))
		if r.Method != "GET" {
			t.Errorf("Expected 'GET' request, got '%s'", r.Method)
		}

	}))

	defer ts.Close()
	api := ts.URL
	fmt.Println("url:", api)
	resp, err := http.Get(api)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	fmt.Println("reps:", resp)
}
