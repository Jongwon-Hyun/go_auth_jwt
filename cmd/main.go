package main

import (
	"authentication/auth"
	"authentication/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// TODO 환경변수에서 값을 얻어오는 등의 보안적인 조치를 취할 것
	secret := "secret"
	authentication := auth.NewAuthentication(secret)
	userRepository := users.NewUserRepository()
	userService := users.NewUserService(userRepository, secret)
	userHandler := users.NewUserHandler(userService)

	r.Route("/users", func(r chi.Router) {
		// 등록
		r.Post("/", userHandler.SignUp)

		// 등록한 유저 대상으로 토큰 발급
		r.Post("/token", userHandler.SignIn)

		// 토큰 인증 테스트
		r.Route("/who_am_i", func(r chi.Router) {
			r.Use(authentication.StripTokenMiddleware)
			r.Get("/", userHandler.GetSelfUser)
		})
	})

	err := http.ListenAndServe(":3333", r)
	if err != nil {
		panic(err)
	}
}
