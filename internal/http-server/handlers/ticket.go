package handlers

import (
	"net/http"

	"github.com/ayayaakasvin/oneflick-ticket/internal/domain"
	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/ctx"
	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/helper"
	"github.com/google/uuid"
)

// InsertTicketAfterwards creates a ticket type for an event owned by the authenticated user.
// @Summary      Create ticket
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ticket  body      domain.Ticket  true  "Ticket to create"
// @Success      200     {object} map[string]interface{}
// @Failure      400     {object} map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      404     {object} map[string]string
// @Failure      500     {object} map[string]string
// @Router       /api/ticket [post]
func (h *Handlers) InsertTicketAfterwards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(ctx.CtxUserIDKey).(uint)

		var saveTicketDTO domain.Ticket
		if err := helper.BindJson(r.Body, &saveTicketDTO); err != nil {
			h.logger.Warn("invalid ticket request body", "err", err)
			helper.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		if event, err := h.er.GetEvent(r.Context(), saveTicketDTO.EventUUID); err != nil {
			helper.WriteJSONResponse(w, http.StatusNotFound, map[string]string{"error": "failed to find event"})
			return
		} else if event.OrganizerID != userID {
			helper.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "access denied"})
			return
		}

		saveTicketDTO.TicketUUID = uuid.NewString()

		if err := h.er.CreateTicket(r.Context(), &saveTicketDTO); err != nil {
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to insert ticket"})
			return
		}

		data := make(map[string]any)
		data["ticket_uuid"] = saveTicketDTO.TicketUUID

		helper.WriteJSONResponse(w, http.StatusOK, data)
	}
}

// DeleteTicket deletes a ticket type from an event owned by the authenticated user.
// @Summary      Delete ticket
// @Tags         tickets
// @Produce      json
// @Security     BearerAuth
// @Param        ticket_uuid  query     string  true  "Ticket UUID"
// @Success      200          {object} map[string]interface{}
// @Failure      400          {object} map[string]string
// @Failure      401          {object}  map[string]string
// @Failure      500          {object} map[string]string
// @Router       /api/ticket [delete]
func (h *Handlers) DeleteTicket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(ctx.CtxUserIDKey).(uint)

		ticketUUID := r.URL.Query().Get("ticket_uuid")

		if ticketUUID == "" {
			helper.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "ticket_uuid is missing"})
			return
		}

		if ticket, err := h.er.GetTicket(r.Context(), ticketUUID); err != nil {
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to find ticket"})
			return
		} else {
			if event, err := h.er.GetEvent(r.Context(), ticket.EventUUID); err != nil {
				helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to find event"})
				return
			} else if event.OrganizerID != userID {
				helper.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "access denied"})
				return
			}
		}

		if err := h.er.DeleteTicket(r.Context(), ticketUUID); err != nil {
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete ticket"})
			return
		}

		helper.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "success"})
	}
}

// GetTicket returns a ticket type by UUID.
// @Summary      Get ticket
// @Tags         tickets
// @Produce      json
// @Param        ticket_uuid  query     string  true  "Ticket UUID"
// @Success      200          {object} map[string]interface{}
// @Failure      400          {object} map[string]string
// @Failure      500          {object} map[string]string
// @Router       /api/ticket [get]
func (h *Handlers) GetTicket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ticketUUID := r.URL.Query().Get("ticket_uuid")

		if ticketUUID == "" {
			helper.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "ticket_uuid is missing"})
			return
		}

		data := make(map[string]any)

		if ticket, err := h.er.GetTicket(r.Context(), ticketUUID); err != nil {
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to find ticket"})
			return
		} else {
			data["ticket"] = ticket
		}

		helper.WriteJSONResponse(w, http.StatusOK, data)
	}
}
