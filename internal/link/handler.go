package link

import (
	"fmt"
	"go-advanced/pkg/middleware"
	"go-advanced/pkg/request"
	"go-advanced/pkg/response"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type LinkHandler struct {
	LinkRepository *LinkRepository
}

type LinkHandlerDeps struct {
	LinkRepository *LinkRepository
}

func NewLinkHandler(router *http.ServeMux, deps *LinkHandlerDeps) {
	handler := &LinkHandler{
		LinkRepository: deps.LinkRepository,
	}
	router.HandleFunc("POST /link", handler.Create())
	router.Handle("PATCH /link/{id}", middleware.IsAuthed(handler.Update()))
	router.HandleFunc("DELETE /link/{id}", handler.Delete())
	router.HandleFunc("GET /{hash}", handler.GoTo())
}

func (h *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[LinkCreateRequest](&w, r)

		if err != nil {
			return
		}
		link := NewLink(body.Url)
		for {
			findedLink, _ := h.LinkRepository.GetByHash(link.Hash)

			if findedLink == nil {
				break
			}
			link.GenerateHash()
		}

		createdLink, err := h.LinkRepository.Create(link)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response.SendJSON(w, 201, createdLink)
	}
}

func (h *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[LinkUpdateRequest](&w, r)

		if err != nil {
			return
		}
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		link, err := h.LinkRepository.Update(&Link{
			Model: &gorm.Model{ID: uint(id)},
			Url:   body.Url,
			Hash:  body.Hash,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response.SendJSON(w, 201, link)
	}
}

func (h *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = h.LinkRepository.GetById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		err = h.LinkRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response.SendJSON(w, 200, fmt.Sprintf("Delete link by %d id", id))
	}
}

func (h *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		link, err := h.LinkRepository.GetByHash(hash)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect)
	}
}
