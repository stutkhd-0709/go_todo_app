package service

import (
	"context"
	"fmt"
	"github.com/stutkhd-0709/go_todo_app/entity"
	"github.com/stutkhd-0709/go_todo_app/store"
)

type AddTask struct {
	DB   store.Execer
	Repo TaskAdder // インタフェースをDIして特定のstoreパッケージに依存しないようになっている
}

func (a *AddTask) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	t := &entity.Task{
		Title:  title,
		Status: entity.TaskStatusTodo,
	}
	err := a.Repo.AddTask(ctx, a.DB, t)
	if err != nil {
		return nil, fmt.Errorf("error adding task: %w", err)
	}
	return t, nil
}
