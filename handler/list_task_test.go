package handler

import (
	"github.com/stutkhd-0709/go_todo_app/entity"
	"github.com/stutkhd-0709/go_todo_app/store"
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
		tasks map[entity.TaskID]*entity.Task
		want  want
	}{
		"ok": {
			tasks: map[entity.TaskID]*entity.Task{
				1: {
					ID:     1,
					Title:  "test1",
					Status: "todo",
				},
				2: {
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
			tasks: map[entity.TaskID]*entity.Task{},
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

			sut := ListTask{
				Store: &store.TaskStore{
					Tasks: tc.tasks,
				},
			}
			sut.ServeHTTP(w, r)
			resp := w.Result()
			testutil.AssertResponse(t, resp, tc.want.status, testutil.LoadFile(t, tc.want.resFile))
		})
	}
}
