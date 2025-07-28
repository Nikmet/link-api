package main

import (
	"fmt"
	"go-advanced/configs"
	"go-advanced/internal/auth"
	"go-advanced/internal/link"
	"go-advanced/pkg/db"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDB(conf)
	router := http.NewServeMux()

	//* Reposittories
	linkRepo := link.NewLinkRepossitory(db)

	//* Handlers
	auth.NewAuthHandler(router, &auth.AuthHandlerDeps{
		Config: conf,
	})
	link.NewLinkHandler(router, &link.LinkHandlerDeps{
		LinkRepository: linkRepo,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Sever is listening on port 8081")
	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}

}
