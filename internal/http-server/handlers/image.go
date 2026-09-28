package handlers

import (
	"net/http"
	"path"
	"strings"

	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/helper"
)

// ServeImages returns an image stored by the application.
// @Summary      Get stored image
// @Tags         images
// @Produce      image/png,image/jpeg,image/webp
// @Security     BearerAuth
// @Param        image  path      string  true  "Image filename"
// @Success      200    {file}    file
// @Failure      400    {object} map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      404    {object} map[string]string
// @Router       /api/images/{image} [get]
func (h *Handlers) ServeImages() http.HandlerFunc {
	prefixOfHandler := "/images/"

	return func(w http.ResponseWriter, r *http.Request) {
		imageURL := strings.TrimPrefix(r.URL.Path, prefixOfHandler)
		if imageURL == "" {
			helper.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "image url is missing after pattern"})
			return
		}

		imageURL = path.Clean(imageURL)
		http.ServeFileFS(w, r, h.lfs, imageURL)
	}
}
