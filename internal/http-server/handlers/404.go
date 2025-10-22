package handlers

import "net/http"

const NotFoundFilePath = "./docs/index.html"

func (h *Handlers) NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, NotFoundFilePath)
	}
}