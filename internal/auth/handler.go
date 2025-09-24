package auth

import (
	"go-advanced/configs"
	"go-advanced/pkg/jwt"
	"go-advanced/pkg/request"
	"go-advanced/pkg/response"
	"net/http"
)

type AuthHandler struct {
	*configs.Config
	*AuthService
}

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps *AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
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
		email, err := h.AuthService.Login(body.Email, body.Password)

		if err != nil {
			response.SendJSON(w, 400, err.Error())
			return
		}

		jwtStr := jwt.NewJWT(h.Auth.Secret)
		token, err := jwtStr.Create(jwt.JWTData{
			Email: email,
		})

		if err != nil {
			response.SendJSON(w, 500, err.Error())
			return
		}

		resp := LoginResponse{
			Token: token,
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
		email, err := h.AuthService.Register(body.Email, body.Password, body.Name)

		if err != nil {
			response.SendJSON(w, 400, err.Error())
			return
		}

		jwtStr := jwt.NewJWT(h.Auth.Secret)
		token, err := jwtStr.Create(jwt.JWTData{
			Email: email,
		})

		if err != nil {
			response.SendJSON(w, 500, err.Error())
			return
		}

		resp := LoginResponse{
			Token: token,
		}

		response.SendJSON(w, http.StatusOK, resp)
	}
}
