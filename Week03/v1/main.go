package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/sullfu/Go-000/Week03/errgroup"
)

func Go(task func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				log.Printf("errgroup: panic recovered: %s\n%s\n", r, buf)
			}
		}()
		task()
	}()
}

func server(ctx context.Context, addr string) error {
	// 参数校验
	srv := http.Server{
		Addr: addr,
	}

	srv.RegisterOnShutdown(func() {
		log.Printf("http server %s shutdown...\n", addr)
	})

	Go(func() {
		<-ctx.Done()
		_ = srv.Shutdown(ctx)
	})

	return srv.ListenAndServe()
	//return errors.New("new error")
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)

	Go(func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		// 使用户 errgroup 的 Cancel() 取消
		g.Cancel()
	})

	g.Go(func() error {
		return server(ctx, ":8281")
	})
	g.Go(func() error {
		return server(ctx, ":8282")
	})
	g.Go(func() error {
		return server(ctx, ":8283")
	})

	if err := g.Wait(); err != nil {
		log.Println("task error: ", err)
	}
}
