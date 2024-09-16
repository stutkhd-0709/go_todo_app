package store

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/stutkhd-0709/go_todo_app/entity"
	"github.com/stutkhd-0709/go_todo_app/testutil"
	"testing"
	"time"
)

func TestKVS_Save(t *testing.T) {
	t.Parallel()

	cli := testutil.OpenRedisForTest(t)

	// sut System Under Test
	sut := &KVS{cli}
	key := "TestKVS_Save"
	uid := entity.UserID(1234)
	ctx := context.Background()
	t.Cleanup(func() {
		cli.Del(ctx, key)
	})
	if err := sut.Save(ctx, key, uid); err != nil {
		t.Errorf("want no error, got %v", err)
	}
}

func TestKVS_Load(t *testing.T) {
	t.Parallel()

	cli := testutil.OpenRedisForTest(t)
	sut := &KVS{cli}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_Load_ok"
		uid := entity.UserID(1234)
		ctx := context.Background()
		// Loadがint64にキャストしてるため、int64にしないと動かない
		cli.Set(ctx, key, int64(uid), 30*time.Minute)

		t.Cleanup(func() {
			cli.Del(ctx, key)
		})
		got, err := sut.Load(ctx, key)
		if err != nil {
			t.Fatalf("want no error, got %v", err)
		}
		if got != uid {
			t.Errorf("want %d, got %d", uid, got)
		}
	})

	t.Run("notFound", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_Save_notFound"
		ctx := context.Background()
		got, err := sut.Load(ctx, key)
		if err == nil || !errors.Is(err, ErrNotFound) {
			t.Errorf("want %v, but got %v(value = %d)", redis.Nil, err, got)
		}
	})
}
