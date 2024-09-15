package handler

import (
	"context"
	"errors"
	"github.com/stutkhd-0709/go_todo_app/entity"
	"github.com/stutkhd-0709/go_todo_app/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListTask(t *testing.T) {
	t.Parallel()

	type want struct {
		status  int
		resFile string
	}

	tests := map[string]struct {
		tasks []*entity.Task
		want  want
	}{
		"ok": {
			tasks: []*entity.Task{
				{
					ID:     1,
					Title:  "test1",
					Status: "todo",
				},
				{
					ID:     2,
					Title:  "test2",
					Status: "done",
				},
			},
			want: want{
				status:  http.StatusOK,
				resFile: "testdata/list_task/ok_res.json.golden",
			},
		},
		"empty": {
			tasks: []*entity.Task{},
			want: want{
				status:  http.StatusOK,
				resFile: "testdata/list_task/empty_res.json.golden",
			},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/tasks", nil)

			moq := &ListTasksServiceMock{ListTasksFunc: func(ctx context.Context) (entity.Tasks, error) {
				if tc.tasks != nil {
					return tc.tasks, nil
				}
				return nil, errors.New("mock error")
			}}

			sut := ListTask{
				Service: moq,
			}
			sut.ServeHTTP(w, r)
			resp := w.Result()
			testutil.AssertResponse(t, resp, tc.want.status, testutil.LoadFile(t, tc.want.resFile))
		})
	}
}
