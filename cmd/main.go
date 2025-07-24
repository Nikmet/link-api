package main

import (
	"fmt"
	"go-advanced/configs"
	"go-advanced/internal/auth"
	"go-advanced/pkg/db"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	_ = db.NewDB(conf)
	router := http.NewServeMux()
	auth.NewAuthHandler(router, &auth.AuthHandlerDeps{
		Config: conf,
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
