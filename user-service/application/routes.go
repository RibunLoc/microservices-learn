package application

import (
	"net/http"

	"github.com/RibunLoc/microservices-learn/user-service/handler"
	repository "github.com/RibunLoc/microservices-learn/user-service/repository/user"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/register", a.loadUserRoutes)
	router.Route("/login", a.loadUserLogin)
	router.Route("/user", a.loadUserChangePassword)

	a.router = router
}

func (a *App) loadUserRoutes(router chi.Router) {

	userHandler := &handler.UserRegister{
		Repo: &repository.RedisMongo{
			Collection: a.mgdb.Collection("users"),
		},
	}

	router.Post("/", userHandler.CreateUserHandler)
}

func (a *App) loadUserLogin(router chi.Router) {
	userHandler := &handler.UserLogin{
		Repo: &repository.RedisMongo{
			Collection: a.mgdb.Collection("users"),
			JwtSecret:  a.config.JwtSecret,
		},
	}

	router.Post("/", userHandler.LoginHandler)
}

func (a *App) loadUserChangePassword(router chi.Router) {
	// khởi tạo các handler cho việc thao tác người dùng
	userChangePassHandler := &handler.UserChangePassword{
		Repo: &repository.RedisMongo{
			Collection: a.mgdb.Collection("users"),
			JwtSecret:  a.config.JwtSecret,
		},
	}
	userUpdateHandler := &handler.UserUpdateInfo{
		Repo: &repository.RedisMongo{
			Collection: a.mgdb.Collection("users"),
			JwtSecret:  a.config.JwtSecret,
		},
	}
	userGetHandler := &handler.UserGetInfo{
		Repo: &repository.RedisMongo{
			Collection: a.mgdb.Collection("users"),
			JwtSecret:  a.config.JwtSecret,
		},
	}

	router.Put("/{id}/change-password", userChangePassHandler.ChangePasswordHandler)
	router.Put("/{id}/update-info", userUpdateHandler.UpdateInfoHandler)
	router.Get("/{id}", userGetHandler.GetInfoHandler) // Thêm route để lấy thông tin người dùng
}
