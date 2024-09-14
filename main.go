package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/stutkhd-0709/go_todo_app/config"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// シグナルの受信を検知できるようにする -> graceful shutdownができるようになる
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	// *http.Serverを使うことで、サーバーのタイムアウト時間など柔軟に変えられるため、ListenAndServe関数より実用的
	// 関数 -> パッケージが提供してる関数
	// メソッド -> 構造体に紐づいている関数
	s := &http.Server{
		// 引数で受け取ったnet.Listenerを利用するため、
		// Addrフィールドは指定しない
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
			// コマンドラインで実験するため
			time.Sleep(5 * time.Second)
		}),
	}

	eg, ctx := errgroup.WithContext(ctx)
	// 別ゴルーチンでHTTPサーバーを起動
	// http.ErrServerClosed は http.Server.Shutdown() が正常に終了したことを表すため異常ではない
	eg.Go(func() error {
		// ListenAndServe -> *http.ServerでAddrを指定した場合に使う
		// Serve -> 引数にListenしたいポート情報を渡す
		if err := s.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to listen and serve: %v", err)
			return err
		}
		return nil
	})

	/*
		1. <-ctx.Doneの次でhttp.Server.Shutdownメソッドが実行
		2. 別ゴルーチンで実行してたhttp.Server.Serveが停止
		3. 別ゴルーチンで実行していた無名関数 func()error が終了
		4. run関数の最後で別ゴルーチンが終了するのを待機していたeg.waitメソッドが終了
		5. 別ゴルーチンで実行していた無名関数 func() error の戻り値がrunの戻り値になる
	*/

	// チャネルからの通知(終了通知)を待機する
	// 今回の場合だと、NotifyContextから受け取るSIGTERMなど
	<-ctx.Done()

	// ShutdownメソッドはGraceful Shutdownする
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown server: %v", err)
	}
	// Goメソッドで起動した別ゴルーチンの終了をまつ
	return eg.Wait()
}
