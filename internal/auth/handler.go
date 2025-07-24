package auth

import (
	"fmt"
	"go-advanced/configs"
	"go-advanced/pkg/request"
	"go-advanced/pkg/response"
	"net/http"
)

type AuthHandler struct {
	*configs.Config
}

type AuthHandlerDeps struct {
	*configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps *AuthHandlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[LoginRequest](&w, r)

		if err != nil {
			response.SendJSON(w, 400, err.Error())
			return
		}
		fmt.Println(body)
		resp := LoginResponse{
			Token: "777",
		}
		response.SendJSON(w, http.StatusOK, resp)
	}
}

func (h *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[RegisterRequest](&w, r)
 		if err != nil {
			response.SendJSON(w, 400, err.Error())
			return
		}
		fmt.Println(body)

		response.SendJSON(w, http.StatusOK, "Register is succes")
	}
}
