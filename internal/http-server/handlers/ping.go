package handlers

import "net/http"

// PingHandler handler
// @Summary      Ping the service
// @Description  Writes "pong" in response
// @Tags         ping
// @Produce      text/plain
// @Success      200  {string}  string  "pong"
// @Router       /ping [get]
func (h *Handlers) PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}
}
