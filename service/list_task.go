package service

import (
	"context"
	"fmt"
	"github.com/stutkhd-0709/go_todo_app/entity"
	"github.com/stutkhd-0709/go_todo_app/store"
)

type ListTask struct {
	DB   store.Queryer
	Repo TaskLister
}

func (l *ListTask) ListTasks(ctx context.Context) (entity.Tasks, error) {
	ts, err := l.Repo.ListTasks(ctx, l.DB)
	if err != nil {
		return nil, fmt.Errorf("ListTasks: %w", err)
	}
	return ts, nil
}
