package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/ctx"
	"github.com/ayayaakasvin/oneflick-ticket/internal/lib/bindjson"
	"github.com/ayayaakasvin/oneflick-ticket/internal/lib/validinput"
	"github.com/ayayaakasvin/oneflick-ticket/internal/models"
	"github.com/ayayaakasvin/oneflick-ticket/internal/models/response"
)

const maxSizeForFile = 10 << 20 // max image size
const (
	PNG  = "image/png"
	JPEG = "image/jpeg"
	WEBP = "image/webp"
)

var validImageMimeTypes map[string]string = map[string]string{
	PNG: ".png",
	JPEG: ".jpeg",
	WEBP: ".webp",
}

func (h *Handlers) SaveEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(ctx.CtxUserIDKey).(uint)

		var saveEventDTO models.Event
		if err := bindjson.BindJson(r.Body, &saveEventDTO); err != nil {
			response.SendErrorJson(w, http.StatusBadRequest, "failed to decode body: %v", err.Error())
			return
		}

		saveEventDTO.OrganizerID = userID

		if err := validinput.ValidateEventSave(&saveEventDTO); err != nil {
			response.SendErrorJson(w, http.StatusBadRequest, err.Error())
			return
		}

		newEventUUID, err := h.eventRepo.InsertEventObjectToDatabase(r.Context(), &saveEventDTO)
		if err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "failed to insert event record")
			return
		}

		data := response.NewData()
		data["event_uuid"] = newEventUUID

		response.SendSuccessJson(w, http.StatusCreated, data)
	}
}

func (h *Handlers) GetEventByUUID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventUUID := r.URL.Query().Get("event_uuid")
		if eventUUID == "" {
			response.SendErrorJson(w, http.StatusBadRequest, "event_uuid is not specified")
			return
		}

		data := response.NewData()
		if event, err := h.eventRepo.GetEventByUUID(r.Context(), eventUUID); err != nil {
			response.SendErrorJson(w, http.StatusNotFound, "failed to find event")
			return
		} else {
			data["event"] = event
		}

		response.SendSuccessJson(w, http.StatusOK, data)
	}
}

func (h *Handlers) GetAllEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := response.NewData()

		if events, err := h.eventRepo.GetAllEvents(r.Context()); err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "failed to find events")
			return
		} else {
			data["events"] = events
		}

		response.SendSuccessJson(w, http.StatusOK, data)
	}
}

func (h *Handlers) GetEventsByCategoryID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categoryIDString := r.URL.Query().Get("category_id")
		if categoryIDString == "" {
			response.SendErrorJson(w, http.StatusBadRequest, "category_id is missing")
			return
		}

		value, err := strconv.ParseUint(categoryIDString, 10, 64) // Base 10, 64-bit size
		if err != nil {
			response.SendErrorJson(w, http.StatusBadRequest, "invalid category_id value")
			return
		}
		categoryID := uint(value)

		data := response.NewData()
		if events, err := h.eventRepo.GetEventsByCategoryID(r.Context(), categoryID); err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "failed to get events")
			return
		} else if len(events) == 0 {
			response.SendErrorJson(w, http.StatusNotFound, "no events found")
			return
		} else {
			data["events"] = events
		}

		response.SendSuccessJson(w, http.StatusOK, data)
	}
}

func (h *Handlers) UpdateEventImageURLUsingExternalSource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(ctx.CtxUserIDKey).(uint)

		eventUUID := r.URL.Query().Get("event_uuid")
		if eventUUID == "" {
			response.SendErrorJson(w, http.StatusBadRequest, "event_uuid is missing")
			return
		}

		if event, err := h.eventRepo.GetEventByUUID(r.Context(), eventUUID); err != nil {
			response.SendErrorJson(w, http.StatusNotFound, "failed to find event")
			return
		} else if event.OrganizerID != userID {
			response.SendErrorJson(w, http.StatusUnauthorized, "access denied")
			return
		}

		imageURL := r.URL.Query().Get("image_url")
		if imageURL == "" {
			response.SendErrorJson(w, http.StatusBadRequest, "image_url is missing")
			return
		}

		if err := h.eventRepo.UpdateEventImageURL(r.Context(), eventUUID, imageURL); err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "failed to update event image")
			return
		}

		response.SendSuccessJson(w, http.StatusOK, nil)
	}
}

func (h *Handlers) UpdateEventImageURLByUploading() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(ctx.CtxUserIDKey).(uint)

		eventUUID := r.URL.Query().Get("event_uuid")
		if eventUUID == "" {
			response.SendErrorJson(w, http.StatusBadRequest, "event_uuid is missing")
			return
		}

		if event, err := h.eventRepo.GetEventByUUID(r.Context(), eventUUID); err != nil {
			response.SendErrorJson(w, http.StatusNotFound, "failed to find event")
			return
		} else if event.OrganizerID != userID {
			response.SendErrorJson(w, http.StatusUnauthorized, "access denied")
			return
		}

		img, mimeType, err := parseImageFromRequest(r)
		if err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "failed to parse and read image")
		}

		imageURL, err := h.lfs.SaveImage(img, fmt.Sprintf("%s.%s", eventUUID, mimeType))
		if err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "failed to save image")
			return
		}

		if err := h.eventRepo.UpdateEventImageURL(r.Context(), eventUUID, imageURL); err != nil {
		  
		}

		response.SendSuccessJson(w, http.StatusOK, nil)
	}
}

func (h *Handlers) DeleteEventByUUID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(ctx.CtxUserIDKey).(uint)

		eventUUID := r.URL.Query().Get("event_uuid")
		if eventUUID == "" {
			response.SendErrorJson(w, http.StatusBadRequest, "event_uuid is missing")
			return
		}

		if event, err := h.eventRepo.GetEventByUUID(r.Context(), eventUUID); err != nil {
			response.SendErrorJson(w, http.StatusNotFound, "failed to find event")
			return
		} else if event.OrganizerID != userID {
			response.SendErrorJson(w, http.StatusUnauthorized, "access denied")
			return
		}

		if err := h.eventRepo.DeleteEventByUUID(r.Context(), eventUUID); err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "failed to delete event")
			return
		}

		response.SendSuccessJson(w, http.StatusOK, nil)
	}
}

func parseImageFromRequest(r *http.Request) (io.Reader, string, error) {
	err := r.ParseMultipartForm(maxSizeForFile)
	if err != nil {
		return nil, "", err
	}

	img, _, err := r.FormFile("image")
	if err != nil {
		return nil, "", err
	}

	mimeType, err := getMimeType(img)
	if err != nil {
		return nil, "", err
	}

	if !checkForValidImageType(mimeType) {
		return nil, "", fmt.Errorf("invalid mime type: %s", mimeType)
	}

	return img, mimeType, nil
}

func getMimeType(file multipart.File) (string, error) {
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil {
		return "", err
	}

	// reset reader to beginning
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	return http.DetectContentType(buf[:n]), nil
}

func checkForValidImageType(mimeType string) bool {
	_, ok := validImageMimeTypes[mimeType]
	return ok
}
