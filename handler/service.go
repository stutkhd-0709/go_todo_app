package handler

import (
	"context"
	"github.com/stutkhd-0709/go_todo_app/entity"
)

/*
インタフェースを定義する理由
1. 特定のパッケージからの参照を取り除いて疎なパッケージ構成にするため
2. インタフェースを介して特定の型に依存しないことで、モックに処理を入れ替えたテストを行うため
*/

//go:generate go run github.com/matryer/moq -out moq_test.go . ListTasksService AddTaskService RegisterUserService
type ListTasksService interface {
	ListTasks(ctx context.Context) (entity.Tasks, error)
}

type AddTaskService interface {
	AddTask(ctx context.Context, title string) (*entity.Task, error)
}

type RegisterUserService interface {
	RegisterUser(ctx context.Context, name, password, role string) (*entity.User, error)
}
