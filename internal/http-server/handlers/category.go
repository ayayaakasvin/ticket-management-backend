package handlers

import (
	"net/http"

	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/helper"
)

// GetAllCategories returns all event categories.
// @Summary      List categories
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/category [get]
func (h *Handlers) GetAllCategories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]any)

		if categories, err := h.er.GetCategories(r.Context()); err != nil {
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch categories"})
			return
		} else {
			data["categories"] = categories
		}

		helper.WriteJSONResponse(w, http.StatusOK, data)
	}
}
