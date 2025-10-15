package main

import (
	"fmt"
	"go-advanced/configs"
	"go-advanced/internal/auth"
	"go-advanced/internal/link"
	"go-advanced/internal/stat"
	"go-advanced/internal/user"
	"go-advanced/pkg/db"
	"go-advanced/pkg/event"
	"go-advanced/pkg/middleware"
	"net/http"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	db := db.NewDB(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	//* Reposittories
	linkRepo := link.NewLinkRepossitory(db)
	userRepo := user.NewUserRepository(db)
	statRepo := stat.NewStatRepository(db)

	//* Services
	authService := auth.NewAuthService(userRepo)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepo,
	})

	//* Handlers
	auth.NewAuthHandler(router, &auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, &link.LinkHandlerDeps{
		LinkRepository: linkRepo,
		Config:         conf,
		EventBus:       eventBus,
	})
	stat.NewStatHandler(router, &stat.StatHanddlerDeps{
		StatRepository: statRepo,
		Config:         conf,
	})

	// Click listener
	go statService.AddClick()

	//Middlewares
	stack := middleware.Chain(middleware.CORS, middleware.Logging)
	return stack(router)
}

func main() {
	app := App()

	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	fmt.Println("Sever is listening on port 8081")
	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
