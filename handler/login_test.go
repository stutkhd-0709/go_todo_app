package handler

import (
	"bytes"
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/stutkhd-0709/go_todo_app/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin_ServeHTTP(t *testing.T) {
	type moq struct {
		token string
		err   error
	}
	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		reqFile string
		moq     moq
		want    want
	}{
		"ok": {
			reqFile: "testdata/login/ok_req.json.golden",
			moq: moq{
				token: "from_moq",
			},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/login/ok_res.json.golden",
			},
		},
		"barRequest": {
			reqFile: "testdata/login/bad_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/login/bad_res.json.golden",
			},
		},
		"internal_server_error": {
			reqFile: "testdata/login/ok_req.json.golden",
			moq: moq{
				err: errors.New("error from mock"),
			},
			want: want{
				status:  http.StatusInternalServerError,
				rspFile: "testdata/login/internal_server_error_res.json.golden",
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(testutil.LoadFile(t, tt.reqFile)))

			moq := &LoginServiceMock{}
			moq.LoginFunc = func(ctx context.Context, name string, pw string) (string, error) {
				return tt.moq.token, tt.moq.err
			}
			sut := Login{
				Service:  moq,
				Validate: validator.New(),
			}
			sut.ServeHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(t, resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})
	}
}
