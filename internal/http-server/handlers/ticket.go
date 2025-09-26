package handlers

import (
	"net/http"

	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/ctx"
	"github.com/ayayaakasvin/oneflick-ticket/internal/lib/bindjson"
	"github.com/ayayaakasvin/oneflick-ticket/internal/models"
	"github.com/ayayaakasvin/oneflick-ticket/internal/models/response"
	"github.com/google/uuid"
)

func (h *Handlers) InsertTicketAfterwards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(ctx.CtxUserIDKey).(uint)

		var saveTicketDTO models.Ticket
		if err := bindjson.BindJson(r.Body, &saveTicketDTO); err != nil {
			return
		}

		if event, err := h.eventRepo.GetEventByUUID(r.Context(), saveTicketDTO.EventUUID); err != nil {
			response.SendErrorJson(w, http.StatusNotFound, "failed to find event")
			return
		} else if event.OrganizerID != userID {
			response.SendErrorJson(w, http.StatusUnauthorized, "access denied")
			return
		}

		saveTicketDTO.TicketUUID = uuid.NewString()

		if err := h.eventRepo.InsertTicketAfterwards(r.Context(), &saveTicketDTO); err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "failed to insert ticket")
			return
		}

		data := response.NewData()
		data["ticket_uuid"] = saveTicketDTO.TicketUUID
		
		response.SendSuccessJson(w, http.StatusOK, data)
	}
}

func (h *Handlers) DeleteTicket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(ctx.CtxUserIDKey).(uint)

		ticketUUID := r.URL.Query().Get("ticket_uuid")

		if ticketUUID == "" {
			response.SendErrorJson(w, http.StatusBadRequest, "ticket_uuid is missing")
			return
		}

		if ticket, err := h.eventRepo.GetTicket(r.Context(), ticketUUID); err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "failed to find ticket")
			return 
		} else {
			if event, err := h.eventRepo.GetEventByUUID(r.Context(), ticket.EventUUID); err != nil {
				response.SendErrorJson(w, http.StatusInternalServerError, "failed to find event")
				return
			} else if event.OrganizerID != userID {
				response.SendErrorJson(w, http.StatusUnauthorized, "access denied")
				return 
			}
		}

		if err := h.eventRepo.DeleteTicket(r.Context(), ticketUUID); err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "failed to delete ticket")
			return
		}

		response.SendSuccessJson(w, http.StatusOK, nil)
	}
}

func (h *Handlers) GetTicket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ticketUUID := r.URL.Query().Get("ticket_uuid")

		if ticketUUID == "" {
			response.SendErrorJson(w, http.StatusBadRequest, "ticket_uuid is missing")
			return
		}

		data := response.NewData()

		if ticket, err := h.eventRepo.GetTicket(r.Context(), ticketUUID); err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "failed to find ticket")
			return
		} else {
			data["ticket"] = ticket
		}

		response.SendSuccessJson(w, http.StatusOK, data)
	}
}
