package handlers

import (
	"net/http"
	"path"
	"strings"

	"github.com/ayayaakasvin/oneflick-ticket/internal/models/response"
)

func (h *Handlers) ServeImages() http.HandlerFunc {
	prefixOfHandler := "/images/"

	return func(w http.ResponseWriter, r *http.Request) {
		imageURL := strings.TrimPrefix(r.URL.Path, prefixOfHandler)
		if imageURL == "" {
			response.SendErrorJson(w, http.StatusBadRequest, "image url is missing after pattern")
			return
		}

		imageURL = path.Clean(imageURL)
		http.ServeFileFS(w, r, h.lfs, imageURL)
	}
}