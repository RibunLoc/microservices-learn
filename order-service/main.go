package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/RibunLoc/microservices-learn/application"
)

func main() {
	app := application.New(application.LoadConfig()) // Khởi tạo cấu hình cho server

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt) // kiểm tra ngắt đột ngột
	// Hàm cancel này được trả về từ mã trên, đảm bảo rằng cancel() sẽ được gọi khi main() kết thúc
	// mục đích là để giải phóng tài nguyên liên quan đến context, dọn dẹp goroutine
	defer cancel()

	err := app.Start(ctx) // Chạy Server
	if err != nil {
		fmt.Println("failed to start app:", err)
	}

	cancel()
}
