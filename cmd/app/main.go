package main

import (
	"NureUvarenkoAnton/unik_go_lb_4/internal/core"
	"NureUvarenkoAnton/unik_go_lb_4/internal/db"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg/jwt"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg/server"
	"NureUvarenkoAnton/unik_go_lb_4/internal/pkg/translate"
	"NureUvarenkoAnton/unik_go_lb_4/internal/service"
	"NureUvarenkoAnton/unik_go_lb_4/internal/transport"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/olahol/melody"
)

func main() {
	godotenv.Load()

	db := db.Connect()

	repo := core.New(db)

	jwtHandler := jwt.NewJWT(os.Getenv("JWT_SECRET"))

	tranlsator := translate.NewTranlator(os.Getenv("DEEPL_HOST"), os.Getenv("DEEPL_API_KEY"))

	service := service.NewService(
		repo,
		*jwtHandler,
		repo,
		repo,
		repo,
		repo,
		tranlsator,
		tranlsator,
	)

	m := melody.New()

	handler := transport.NewHandler(
		service.AuthService,
		service.ProfileService,
		m,
		service.UsersService,
		service.WalkService,
		service.RatingService,
	)

	s := server.New(handler, *jwtHandler, m)

	fmt.Println("starting server...")

	log.Fatal(s.ListenAndServe())
}
