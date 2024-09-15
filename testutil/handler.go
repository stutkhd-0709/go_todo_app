package testutil

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"io"
	"net/http"
	"os"
	"testing"
)

func AssertJSON(t *testing.T, want, got []byte) {
	t.Helper()

	var jw, jg any
	// Unmarshalした結果を&jwに書き込む
	if err := json.Unmarshal(want, &jw); err != nil {
		t.Fatalf("cannot unmarshal want %q: json: %v", want, err)
	}
	if err := json.Unmarshal(got, &jg); err != nil {
		t.Fatalf("cannot unmarshal got %q: json: %v", got, err)
	}
	if diff := cmp.Diff(jw, jg); diff != "" {
		t.Errorf("json mismatch (-want +got):\n%s", diff)
	}
}

func AssertResponse(t *testing.T, got *http.Response, status int, body []byte) {
	t.Helper()
	t.Cleanup(func() {
		_ = got.Body.Close()
	})
	gb, err := io.ReadAll(got.Body)
	if err != nil {
		t.Fatalf("cannot read response body: %v", err)
	}
	if got.StatusCode != status {
		t.Fatalf("got status %d, want %d", got.StatusCode, status)
	}

	if len(gb) == 0 && len(body) == 0 {
		// 期待としても実体としてもレスポンスがないので
		// AssertJsonを呼ばなくていい
		return
	}
	AssertJSON(t, body, gb)
}

// LoadFile はゴールデンファイルテスト用
func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	bt, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read %q: %v", path, err)
	}
	return bt
}
