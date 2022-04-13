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

	// TODO 환경변수에서 값을 가져오는 등의 보안적 조치를 취할 것 !
	secret := "secret"
	autenticationMiddelware := auth.NewAuthentication(secret)

	// TODO 의존성 주입 엮은 UserHandler 사용하도록 변경할 것!
	r.Route("/users", func(r chi.Router) {
		// 등록
		r.Post("/", users.UserHandler{}.SignUp)

		// 등록한 유저 대상으로 토큰 발급
		r.Post("/token", users.UserHandler{}.SignIn)

		// 토큰 인증 테스트
		r.Route("/who_am_i", func(r chi.Router) {
			r.Use(autenticationMiddelware.StripTokenMiddleware)
			r.Get("/", users.UserHandler{}.GetSelfUser)
		})
	})

	err := http.ListenAndServe(":3333", r)
	if err != nil {
		panic(err)
	}
}
