package application

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddress string // địa chỉ redis server
	Username     string // tên user login redis
	Password     string // mật khẩu login
	ServerPort   uint16 // cổng lắng nghe của backend
}

func LoadConfig() Config {
	_ = godotenv.Load()
	// tạo mơi config với cấu hình mặc định
	cfg := Config{
		RedisAddress: "",
		Username:     "",
		Password:     "",
		ServerPort:   3000,
	}

	// Kiểm tra biến môi trường với REDIS_ADDR có tồn tại hay không
	if redisAddr, exist := os.LookupEnv("REDIS_ADDR"); exist {
		cfg.RedisAddress = redisAddr
	}

	// Kiểm tra biến môi trường với REDIS_USERNAME
	if redisUser, exits := os.LookupEnv("REDIS_USERNAME"); exits {
		cfg.Username = redisUser
	}

	// Kiểm tra biến môi trường với REDIS_PASSWORD
	if redisPass, exist := os.LookupEnv("REDIS_PASSWORD"); exist {
		cfg.Password = redisPass
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
