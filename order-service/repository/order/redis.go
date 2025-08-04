package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/RibunLoc/microservices-learn/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client // Redis client từ go-redis
}

// Tạo key Redis dạng: "order:123"
func orderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

// Hàm Insert lưu order vào Redis
func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {
	// 1. Mã hóa strut Order thành chuỗi JSON
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	// 2. Tạo key Redis cho đơn hàng (ví dụ: "order:123")
	key := orderIDKey(order.OrderID)

	// 3. Tạo pipeline transaction (Gom nhiều lệnh lại và thực thi cùng một lúc)
	txn := r.Client.TxPipeline()

	// 4. Lưu order vào Redis nếu key chưa tồn tại (SETNX)
	res := txn.SetNX(ctx, key, string(data), 0)

	// 5. Kiểm tra khi lỗi gọi SETNX
	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to set: %w", err)
	}

	// 6. Thêm key order vào tập hợp "orders" (để dễ truy vấn sau này)
	if err := txn.SAdd(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add to orders set: %w", err)
	}

	// 7. Thực thi tất cả lệnh trong pipeline
	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}
	return nil
}

// Dùng để báo lỗi khi không tìm thấy order
var ErrNotExist = errors.New("order does not exist")

// findByID truy vấn Redis để lấy order theo ID
func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Order, error) {
	// 1. Tạo key Redis từ order ID
	key := orderIDKey(id)

	// 2. Thử lấy dữ liệu từ Redis
	value, err := r.Client.Get(ctx, key).Result()

	// 3. Nếu key khong tồn tại trong Redis
	if errors.Is(err, redis.Nil) {
		return model.Order{}, ErrNotExist
	} else if err != nil { // 4. Nếu có lỗi Redis khác (timeout, network, v.v...)
		return model.Order{}, fmt.Errorf("get order: %w", err)
	}

	// 5. nếu khong có lỗi, tiến hành giải mã JSON thành struct Order
	var order model.Order
	err = json.Unmarshal([]byte(value), &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to decode order json: %w", err)
	}

	// 6. Trả về order thành công
	return order, nil
}

// DeleteByID xóa order dựa trên id cung cấp
func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	// 1. Tạo key Redis từ order ID
	key := orderIDKey(id)

	// 2. Tạo pipline transaction để gom lệnh thực thi
	txn := r.Client.TxPipeline()

	// 3. Thêm lệnh xóa key order
	err := txn.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		txn.Discard()         // Hủy luôn pipline
		return os.ErrNotExist // trả về lỗi dữ liệu không tồn tại
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("get order: %w", err)
	}

	// 4. Thêm lệnh xóa key khỏi set "orders" (danh sách order)
	if err := txn.SRem(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to remove from orders set: %w", err)
	}

	// 5. Thực thi pipeline, gồm cả lệnh xóa key và xóa khỏi set
	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

// Cập nhật thông tin đơn hàng vảo Redis nếu key đã tồn tại
func (r *RedisRepo) Update(ctx context.Context, order model.Order) error {
	// 1. Chuyển struct Order thành chuỗi JSON để lưu vào Redis
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	// 2. Thông báo tên key cần thay đổi
	key := orderIDKey(order.OrderID)

	// 3. Dùng lệnh SETXX: chỉ cập nhật giá trị nếu key đã tồi tại trong Redis
	err = r.Client.SetXX(ctx, key, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return os.ErrNotExist
	}

	return nil
}

/*
	Định nghĩa thông tin phân trang: lấy bao nhiêu phần tử (Size),

bắt đầu từ đâu (Offset)
*/
type FindAllPage struct {
	Size   uint64 // Số lượng kết quả muốn lấy mỗi lần
	Offset uint64 // Vị trí bắt đầu (cursor) trong tập hợp
}

// là kết quả trả về khi truy vấn: danh sách đơn hàng + cursor tiếp theo để phân trang
type FindResult struct {
	Orders []model.Order // Danh sách đơn hàng lấy được
	Cursor uint64        // Con trỏ tiếp theo (cursor) để try vấn trang tiếp theo
}

// FindALL thực hiện lấy danh sahcs các đơn hàng từ Redis (phân trang bằng SSCAN)
func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	/*
		1. Dùng Redis SSCAN để lấy các key đơn hàng từ tập hợp "orders"
			- page.Offset: con trỏ bắt đầu (cursor)
			- "*": lấy tất cả key
			- page.Size: số lượng key cần lấy
	*/
	res := r.Client.SScan(ctx, "orders", uint64(page.Offset), "*", int64(page.Size))

	// 2. Lấy danh sách key và cursor tiếp theo từ kết quả SSCAN
	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get order ids: %w", err)
	}

	// 3. Nếu không có key nào được trả về, trả về danh sách rỗng
	if len(keys) == 0 {
		return FindResult{
			Orders: []model.Order{},
		}, nil
	}

	// 4. Dùng MGET để lấy dữ liệu chi tiết (giá trị) của các key cùng lúc
	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get orders: %w", err)
	}

	// 5. Tạo danh sách đơn hàng (orders) để lưu kết quả slice go
	orders := make([]model.Order, len(xs))

	// 6. Lặp qua từng kết quả Redis trả về
	for i, x := range xs {
		x := x.(string) // Ép kiểu kết quả sang string (Redis trả về kiểu interface{})
		var order model.Order

		// 7. Giải mã chuỗi JSON thành struct `model.Order`
		err := json.Unmarshal([]byte(x), &order)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to decode order json: %w", err)
		}
		orders[i] = order // Lưu đơn hàng vào danh sách
	}

	// 8. Trả về danh sách đơn hàng và cursor để phân trang cho lần sau
	return FindResult{
		Orders: orders,
		Cursor: cursor,
	}, nil
}
