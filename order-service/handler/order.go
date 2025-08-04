package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/RibunLoc/microservices-learn/model"
	"github.com/RibunLoc/microservices-learn/repository/order"
	"github.com/RibunLoc/microservices-learn/util"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Order là một HTTP handler chứa tham chiếu đến RedisRepoo để thao tác dữ liệu
type Order struct {
	Repo *order.RedisRepo
}

// Create là HTTP handler để tạo một đơn hàng mới (POST /orders)
func (h *Order) Create(w http.ResponseWriter, r *http.Request) {
	// Định nghĩa struct tạm thời đề nhận dữ liệu JSON từ client gửi lên
	var body struct {
		CustomerID  uuid.UUID        `json:"customer_id"` // ID của khách hàng
		LineItems   []model.LineItem `json:"line_items"`  // Danh sách các mặt hàng trong đơn
		OrderStatus string           `json:"order_status"`
	}

	// Giải mã (decode) dữ liệu JSON từ body request vào struct `body`
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest) // Nếu lỗi, trả về 400 Bad Request
		return
	}

	// Lấy thời gian thực
	time_zone := time.Now().UTC()
	now := time_zone

	// Tạo struct Order từ dữ liệu nhận được
	order := model.Order{
		OrderID:     rand.Uint64(), // Tạo ID ngẫu nhiên cho đơn hàng
		CustomerID:  body.CustomerID,
		LineItems:   body.LineItems,
		OrderStatus: body.OrderStatus,
		CreateAt:    &now,
	}

	// Gọi Repo để chèn đơn hàng vào Redis
	err := h.Repo.Insert(r.Context(), order)
	if err != nil {
		fmt.Println("failed to insert: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Chuyển order thành JSON để trả về cho client
	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println("failed to insert: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Ghi JSON và response
	w.Write(res)
	w.WriteHeader(http.StatusCreated) // Trả về status 201 created

}

// List là HTTP handler để liệt kê tất cả mặt hàng của User đó (GET /)
func (h *Order) List(w http.ResponseWriter, r *http.Request) {
	// Lấy giá trị cursor từ query string (?cursor=...), nếu không có thì mặc định là 0
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}

	// Chuyển đổi cursor từ string -> uint64 để phục vụ truy vấn dữ liệu
	const decimal = 10
	const bitsize = 64
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitsize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Số lượng item mỗi trang, có thể thay đổi tùy thiết kế
	const size = 50

	// Truy vấn dữ liệu từ Redis(hoặc DB) thông qua Repo
	res, err := h.Repo.FindAll(r.Context(), order.FindAllPage{
		Offset: cursor,
		Size:   size,
	})
	if err != nil {
		fmt.Println("failed to find all: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Cấu trúc response trả về client
	var response struct {
		Items []model.Order `json:"items"`          // Danh sách đơn hàng
		Next  uint64        `json:"next,omitempty"` // Cursor tiếp theo, dùng để lấy trang kế tiếp
	}
	response.Items = res.Orders
	response.Next = res.Cursor

	// Chuyển response thành JSON
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed to marshal: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Gửi JSON về client
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

// HTTP handler dùng để lấy chi tiết một đơn hàng theo ID.
func (h *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	// Lấy tham số "id" từ URL path, ví dụ: /order/123 -> id = 123
	idParam := chi.URLParam(r, "id")

	const base = 10    // hệ cơ số dùng để parse
	const bitSize = 64 // số bit của uint64

	// Chuyển id từ chuỗi sang uint64
	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Gọi hàm repo để tìm đơn hàng theo ID
	o, err := h.Repo.FindByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		// Nếu không tìm thấy đơn hàng, trả về lỗi 400
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Encode struct đơn hàng thành JSON và ghi vào response
	if err := json.NewEncoder(w).Encode(o); err != nil {
		fmt.Println("failed to marshal: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// xử lý yêu cầu HTTP để cập nhật trạng thái đơn hàng theo ID
func (h *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	// Định nghĩa struct để prase phần thanh JSON có chứa trường "status"
	var body struct {
		Status string `json:"status"`
	}

	// Giải mã JSON từ body của request
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	// Lấy ID từ URL param và chuyển sang kiểu uint64
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Tìm đơn hàng theo ID trong repository
	theOrder, err := h.Repo.FindByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Khai báo giá trị trạng thái hợp lệ và thời gian hiện tại
	const completedStatus = "completed"
	const shippedStatus = "shipped"
	now := util.CustomTime(time.Now())

	// xử lsy cập nhật trạng thái đơn hàng theo logic nghiệp vụ
	switch body.Status {
	case shippedStatus:
		// không cho phép shipped nếu đã completed
		if theOrder.CompletedAt != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.ShippedAt = &now

	case completedStatus:
		// Chỉ cho phép completed nếu đã shipped và chưa completed
		if theOrder.CompletedAt != nil || theOrder.ShippedAt == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.CompletedAt = &now
	default:
		// Trạng thái không hợp lệ
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Gọi repositoy để cập nhật đơn hàng
	err = h.Repo.Update(r.Context(), theOrder)
	if err != nil {
		fmt.Println("failed to insert: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Trả về đơn hàng đã cập nhật dưới dạng JSON
	if err := json.NewEncoder(w).Encode(theOrder); err != nil {
		fmt.Println("Failed to Marshal: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// HTTP handler để xóa đơn hàng theo ID
func (h *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[command] DeletedByID called")

	// Lấy ID từ URL
	idParam := chi.URLParam(r, "id")
	fmt.Println("[command] URL param id: ", idParam)

	const base = 10
	const bitSize = 64

	// Prase ID sang uint64
	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("[command] Prased order ID: ", orderID)

	// Gọi repository để xóa theo ID
	err = h.Repo.DeleteByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by ID: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("[command] Sucessfully deleted order ID: ", orderID)
	w.WriteHeader(http.StatusNoContent) // 204 - xóa thành công, khoogn trả body
}
