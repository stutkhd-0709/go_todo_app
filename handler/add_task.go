package handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/stutkhd-0709/go_todo_app/entity"
	"github.com/stutkhd-0709/go_todo_app/store"
	"net/http"
	"time"
)

// AddTask はhttp.HandleFunc型を満たすServeHTTPメソッドを実装している
type AddTask struct {
	Store     *store.TaskStore
	Validator *validator.Validate
}

func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// validateする内容を定義
	var b struct {
		Title string `json:"title" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrorResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	// bで定義したvalidateに合致してるか確認 -> titleが必ず入ってるかどうか
	err := at.Validator.Struct(b)
	if err != nil {
		RespondJSON(ctx, w, &ErrorResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	t := &entity.Task{
		Title:     b.Title,
		Status:    entity.TaskStatusTodo,
		CreatedAt: time.Now(),
	}
	id, err := at.Store.Add(t)
	if err != nil {
		RespondJSON(ctx, w, &ErrorResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	rsp := struct {
		ID int `json:"id"`
	}{ID: int(id)}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
