package store

import (
	"context"
	"github.com/stutkhd-0709/go_todo_app/entity"
	"time"
)

func (r *Repository) ListTasks(ctx context.Context, db Queryer) (entity.Tasks, error) {
	tasks := entity.Tasks{}

	sql := `SELECT
			id, title,
			status, created, modified
		FROM task;`

	// SelectContextは複数行クエリした結果を構造体にマッピングして返却してくれる
	if err := db.SelectContext(ctx, &tasks, sql); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) AddTask(ctx context.Context, db Execer, t *entity.Task) error {
	t.Created = time.Now()
	t.Modified = time.Now()
	sql := `INSERT INTO tasks (title, status, created, modified) VALUES (?, ?, ?, ?)`
	result, err := db.ExecContext(ctx, sql, t.Title, t.Status, t.Created, t.Modified)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = entity.TaskID(id)
	return nil
}
