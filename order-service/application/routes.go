package application

import (
	"net/http"

	"github.com/RibunLoc/microservices-learn/handler"
	"github.com/RibunLoc/microservices-learn/repository/order"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Dùng để khởi tạo và cấu hình các routes chính cho ứng dụng
func (a *App) loadRoutes() {
	// Tạo một router mới từ thư viện chi, dùng để định nghĩa các endpoint API
	router := chi.NewRouter()

	// Ghi log cho tất cả request - ghi lại method, URL, thời gian xử lý
	router.Use(middleware.Logger)

	// Định nghĩa endpoint "/" kiểm tra app đang chạy
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Gắn nhóm route con /orders vào router, bằng cách gọi hàm a.loadOrderRoutes
	router.Route("/orders", a.loadOrderRoutes)

	// Gắn router đã cấu hình vào App, khi khởi động server sẽ dùng đến nó
	a.router = router
}

// định nghĩa các route con bên trong /orders
func (a *App) loadOrderRoutes(router chi.Router) {
	/*
		Tạo một handler xử lý đơn hàng, gắn với một repo truy xuất dữ liệu
		 - handler xử lý http
		 - repo truy xuất redis
	*/
	orderHandler := &handler.Order{
		Repo: &order.RedisRepo{
			Client: a.rdb,
		},
	}

	router.Post("/", orderHandler.Create)           // Tạo mới một đơn hàng
	router.Get("/", orderHandler.List)              // Trả về danh sách tất cả các đơn hàng
	router.Get("/{id}", orderHandler.GetByID)       // Trả về đơn hàng theo id
	router.Put("/{id}", orderHandler.UpdateByID)    // Cập nhật đơn hàng theo id
	router.Delete("/{id}", orderHandler.DeleteByID) // Xóa đơn hàng theo id
}
