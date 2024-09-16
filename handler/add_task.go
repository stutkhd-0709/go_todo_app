package handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/stutkhd-0709/go_todo_app/entity"
	"net/http"
)

// AddTask はhttp.HandleFunc型を満たすServeHTTPメソッドを実装している
type AddTask struct {
	Service   AddTaskService
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
	if err := at.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrorResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	t, err := at.Service.AddTask(ctx, b.Title)
	if err != nil {
		RespondJSON(ctx, w, &ErrorResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	rsp := struct {
		ID entity.TaskID `json:"id"`
	}{ID: t.ID}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
