package main

import (
	"fmt"
	"go-advanced/configs"
	"go-advanced/internal/auth"
	"go-advanced/internal/link"
	"go-advanced/internal/user"
	"go-advanced/pkg/db"
	"go-advanced/pkg/middleware"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDB(conf)
	router := http.NewServeMux()

	//* Reposittories
	linkRepo := link.NewLinkRepossitory(db)
	userRepo := user.NewUserRepository(db)

	//* Services
	authService := auth.NewAuthService(userRepo)

	//* Handlers
	auth.NewAuthHandler(router, &auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, &link.LinkHandlerDeps{
		LinkRepository: linkRepo,
	})

	//Middlewares
	stack := middleware.Chain(middleware.CORS, middleware.Logging)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	fmt.Println("Sever is listening on port 8081")
	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}

}
