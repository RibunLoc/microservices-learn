package application

import (
	"os"
	"strconv"
)

type Config struct {
	RedisAddress string // địa chỉ redis server
	ServerPort   uint16 // cổng lắng nghe của backend
}

func LoadConfig() Config {
	// tạo mơi config với cấu hình mặc định
	cfg := Config{
		RedisAddress: "localhost:6379",
		ServerPort:   3000,
	}

	// Kiểm tra biến môi trường với REDIS_ADDR có tồn tại hay không
	if redisAddr, exist := os.LookupEnv("REDIS_ADDR"); exist {
		cfg.RedisAddress = redisAddr

	}

	// Kiểm tra biến môi trường SERVER_PORT
	if serverPort, exist := os.LookupEnv("SERVER_PORT"); exist {
		// Chuyển kiểu dạng số về chuỗi (string)
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint16(port)
		}
	}

	return cfg // trả về cấu hình config đã thiết lập
}
