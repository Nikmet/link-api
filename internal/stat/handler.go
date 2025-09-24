package stat

import (
	"go-advanced/configs"
	"go-advanced/pkg/middleware"
	"go-advanced/pkg/response"
	"net/http"
	"time"
)

const (
	GROUP_BY_DAY   = "day"
	GROUP_BY_MONTH = "month"
)

type StatHanddler struct {
	StatRepository *StatRepository
}

type StatHanddlerDeps struct {
	StatRepository *StatRepository
	Config         *configs.Config
}

func NewStatHandler(router *http.ServeMux, deps *StatHanddlerDeps) {
	handler := &StatHanddler{
		StatRepository: deps.StatRepository,
	}
	router.Handle("GET /stat", middleware.IsAuthed(handler.GetStat(), deps.Config))
}

func (h *StatHanddler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryData := r.URL.Query()
		from, err := time.Parse("2006-01-02", queryData.Get("from"))

		if err != nil {
			http.Error(w, "Invalid from param", http.StatusBadRequest)
			return
		}

		to, err := time.Parse("2006-01-02", queryData.Get("to"))

		if err != nil {
			http.Error(w, "Invalid to param", http.StatusBadRequest)
			return
		}

		by := queryData.Get("by")

		if by != GROUP_BY_DAY && by != GROUP_BY_MONTH {
			http.Error(w, "Invalid by param", http.StatusBadRequest)
			return
		}

		stats := h.StatRepository.GetStats(by, from, to)
		response.SendJSON(w, 200, stats)
	}
}
