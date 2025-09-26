package handlers

import (
	"net/http"

	"github.com/ayayaakasvin/oneflick-ticket/internal/models/response"
)

func (h *Handlers) GetAllCategories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := response.NewData()

		if categories, err := h.eventRepo.GetAllCategories(r.Context()); err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "failed to fetch categories")
			return
		} else {
			data["categories"] = categories
		}

		response.SendSuccessJson(w, http.StatusOK, data)
	}
}