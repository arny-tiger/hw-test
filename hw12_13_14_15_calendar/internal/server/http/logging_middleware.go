package internalhttp

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/logger"
)

type responseCodeWriter struct {
	http.ResponseWriter
	responseCode int
}

func (w *responseCodeWriter) WriteHeader(statusCode int) {
	w.responseCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func loggingMiddleware(next http.Handler, logger logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rcw := responseCodeWriter{w, http.StatusOK}
		next.ServeHTTP(&rcw, r)

		host := r.RemoteAddr
		currentTime := time.Now().Format("02/Jan/2006:15:04:05")
		httpInfo := r.Method + " " + r.URL.Path + " " + r.Proto + " " + strconv.Itoa(rcw.responseCode)
		latency, ok := r.Context().Value("latency").(time.Duration)
		if !ok {
			fmt.Print("Not found")
			return
		}
		userAgent := r.Header.Get("User-Agent")
		logStr := host + " " + currentTime + " " + httpInfo + " " + latency.String() + " " + userAgent
		logger.Info(logStr)
	})
}
