package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := &WrapperWriter{ResponseWriter: w, StatusCode: 200}

		next.ServeHTTP(wrapper, r)

		// Цвета
		methodColor := color.New(color.FgCyan).SprintFunc()
		statusColor := color.New(color.FgGreen).SprintFunc()
		if wrapper.StatusCode >= 400 && wrapper.StatusCode < 500 {
			statusColor = color.New(color.FgYellow).SprintFunc()
		} else if wrapper.StatusCode >= 500 {
			statusColor = color.New(color.FgRed).SprintFunc()
		}

		log.Printf(
			"| %-6s | %-40s | %8v | %s",
			methodColor(r.Method),
			r.URL.Path,
			time.Since(start).Truncate(time.Millisecond),
			statusColor(wrapper.StatusCode),
		)
	})
}
