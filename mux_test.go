package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMux(t *testing.T) {
	w := httptest.NewRecorder() // ResponseWriterインタフェースを満たす*ResponseRecorder型の値を受け取れる
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	sut := NewMux()
	sut.ServeHTTP(w, r)
	resp := w.Result() // クライアントが受け取るhttp.Response型の値を取得できる
	t.Cleanup(func() {
		_ = resp.Body.Close()
	})

	if resp.StatusCode != http.StatusOK {
		t.Errorf("NewMux: got status %d, want %d", resp.StatusCode, http.StatusOK)
	}

	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	want := `{"status":"ok"}`
	if string(got) != want {
		t.Errorf("NewMux: got %q, want %q", got, want)
	}
}
