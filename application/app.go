package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
	config Config
}

func New(config Config) *App {
	app := &App{
		rdb: redis.NewClient(&redis.Options{
			Addr: config.RedisAddress,
		}),
		config: config,
	}

	app.loadRoutes()

	return app
}

func (a *App) Start(ctx context.Context) error {
	//Khởi tạo HTTP server với port lấy từ config và gán router làm handler
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.ServerPort),
		Handler: a.router,
	}

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	// defer: trì hoãn thực thi cho dến khi hàm kết thúc
	// Đảm bảo đóng két nối Redis khi hàm Start() kết thúc
	// giúp giải phóng tài nguyên và tránh rò rỉ kết nối.
	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("failed to close redis", err)
		}
	}()

	fmt.Println("Starting server")

	// Tạo channel để nhận lỗi nếu server khởi động thất bại
	ch := make(chan error, 1)

	// chạy server trong goroutine, tránh chặn luồng chính
	// (giống chạy bất đồng bộ, không block)
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failded to start server: %w", err)
		}
		close(ch)
	}()

	// Dùng select để:
	// - Bắt lỗi từ goroutine nếu server chưa crash
	// - Hoặc bắt tín hiệu hủy từ context để tắt server an toàn
	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}

	return nil
}
