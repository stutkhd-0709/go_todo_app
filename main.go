package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"os"
)

//func main() {
//	err := http.ListenAndServe(":18080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		fmt.Fprintf(w, "Hello World")
//	}))
//	if err != nil {
//		fmt.Println("failed to terminate server: %w", err)
//		os.Exit(1)
//	}
//}

func main() {
	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}
	p := os.Args[1]
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen port %s: %v", p, err)
	}
	if err := run(context.Background(), l); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, l net.Listener) error {
	// *http.Serverを使うことで、サーバーのタイムアウト時間など柔軟に変えられるため、ListenAndServe関数より実用的
	// 関数 -> パッケージが提供してる関数
	// メソッド -> 構造体に紐づいている関数
	s := &http.Server{
		// 引数で受け取ったnet.Listenerを利用するため、
		// Addrフィールドは指定しない
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}

	eg, ctx := errgroup.WithContext(ctx)
	// 別ゴルーチンでHTTPサーバーを起動
	// http.ErrServerClosed は http.Server.Shutdown() が正常に終了したことを表すため以上ではない
	eg.Go(func() error {
		// ListenAndServe -> *http.ServerでAddrを指定した場合に使う
		// Serve -> 引数にListenしたいポート情報を渡す
		if err := s.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to listen and serve: %v", err)
			return err
		}
		return nil
	})

	// チャネルからの通知(終了通知)を待機する
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown server: %v", err)
	}
	// Goメソッドで起動した別ゴルーチンの終了をまつ
	return eg.Wait()
}
