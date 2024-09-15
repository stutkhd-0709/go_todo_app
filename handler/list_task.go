package handler

import (
	"github.com/stutkhd-0709/go_todo_app/entity"
	"net/http"
)

type ListTask struct {
	Service ListTasksService
}

type task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks, err := lt.Service.ListTasks(ctx)
	if err != nil {
		RespondJSON(ctx, w, &ErrorResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
	}
	rsp := []task{} // varにするとemptyの時のレスポンスがgoldenファイルのものと合わなくなる
	for _, t := range tasks {
		// entityのtask型からfieldを絞り込んだ新しいtask型を定義してる
		rsp = append(rsp, task{
			ID:     t.ID,
			Title:  t.Title,
			Status: t.Status,
		})
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)

}
