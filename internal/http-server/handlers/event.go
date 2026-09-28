package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/ayayaakasvin/oneflick-ticket/internal/domain"
	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/ctx"
	"github.com/ayayaakasvin/oneflick-ticket/internal/http-server/helper"
	"github.com/ayayaakasvin/oneflick-ticket/internal/lib/validinput"
)

const maxSizeForFile = 10 << 20 // max image size
const (
	PNG                   = "image/png"
	JPEG                  = "image/jpeg"
	WEBP                  = "image/webp"
	trendingUpdateTimeKey = "trending_update_time"
	trendingKey           = "trending_events"
)

var validImageMimeTypes map[string]string = map[string]string{
	PNG:  ".png",
	JPEG: ".jpeg",
	WEBP: ".webp",
}

// SaveEvent creates an event owned by the authenticated user.
// @Summary      Create event
// @Description  Creates an event and assigns the authenticated user as its organizer.
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        event  body      domain.Event  true  "Event to create"
// @Success      201    {object} map[string]interface{}
// @Failure      400    {object} map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      500    {object} map[string]string
// @Router       /api/event [post]
func (h *Handlers) SaveEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(ctx.CtxUserIDKey).(uint)

		var saveEventDTO domain.Event
		if err := helper.BindJson(r.Body, &saveEventDTO); err != nil {
			h.logger.Warn("invalid event request body", "err", err)
			helper.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		saveEventDTO.OrganizerID = userID

		if err := validinput.ValidateEventSave(&saveEventDTO); err != nil {
			helper.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "event data does not meet the required fields or value constraints"})
			return
		}

		newEventUUID, err := h.er.CreateEvent(r.Context(), &saveEventDTO)
		if err != nil {
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to insert event record"})
			return
		}

		data := make(map[string]any)
		data["event_uuid"] = newEventUUID

		helper.WriteJSONResponse(w, http.StatusCreated, data)
	}
}

// GetEventByUUID returns a single event.
// @Summary      Get event
// @Tags         events
// @Produce      json
// @Security     BearerAuth
// @Param        event_uuid  query     string  true  "Event UUID"
// @Success      200         {object} map[string]interface{}
// @Failure      400         {object} map[string]string
// @Failure      401         {object}  map[string]string
// @Failure      404         {object} map[string]string
// @Router       /api/event [get]
func (h *Handlers) GetEventByUUID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventUUID := r.URL.Query().Get("event_uuid")
		if eventUUID == "" {
			helper.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "event_uuid is not specified"})
			return
		}

		data := make(map[string]any)
		if event, err := h.er.GetEvent(r.Context(), eventUUID); err != nil {
			helper.WriteJSONResponse(w, http.StatusNotFound, map[string]string{"error": "failed to find event"})
			return
		} else {
			data["event"] = event
		}

		helper.WriteJSONResponse(w, http.StatusOK, data)
	}
}

// GetAllEvents returns all events.
// @Summary      List events
// @Tags         events
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/event/all [get]
func (h *Handlers) GetAllEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]any)

		if events, err := h.er.GetEvents(r.Context()); err != nil {
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to find events"})
			return
		} else {
			data["events"] = events
		}

		helper.WriteJSONResponse(w, http.StatusOK, data)
	}
}

// GetEventsByCategoryID lists events in a category.
// @Summary      List events by category
// @Tags         events
// @Produce      json
// @Security     BearerAuth
// @Param        category_id  query     integer  true  "Category ID"
// @Success      200          {object} map[string]interface{}
// @Failure      400          {object} map[string]string
// @Failure      401          {object}  map[string]string
// @Failure      404          {object} map[string]string
// @Failure      500          {object} map[string]string
// @Router       /api/event/category [get]
func (h *Handlers) GetEventsByCategoryID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categoryIDString := r.URL.Query().Get("category_id")
		if categoryIDString == "" {
			helper.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "category_id is missing"})
			return
		}

		value, err := strconv.ParseUint(categoryIDString, 10, 64) // Base 10, 64-bit size
		if err != nil {
			helper.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid category_id value"})
			return
		}
		categoryID := uint(value)

		data := make(map[string]any)
		if events, err := h.er.GetEventsByCategory(r.Context(), categoryID); err != nil {
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to get events"})
			return
		} else if len(events) == 0 {
			helper.WriteJSONResponse(w, http.StatusNotFound, map[string]string{"error": "no events found"})
			return
		} else {
			data["events"] = events
		}

		helper.WriteJSONResponse(w, http.StatusOK, data)
	}
}

// GetTop10Events returns the current trending events.
// @Summary      Get trending events
// @Tags         events
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object} map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/event/top-ten [get]
func (h *Handlers) GetTop10Events() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]any)

		var marshalled []byte
		if marshalledAny, err := h.cc.Get(r.Context(), trendingKey); err != nil {
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to get list"})
			return
		} else {
			switch v := marshalledAny.(type) {
			case []byte:
				{
					marshalled = v
				}
			case string:
				{
					marshalled = []byte(v)
				}
			}
		}

		if updatedAtAny, err := h.cc.Get(r.Context(), trendingUpdateTimeKey); err != nil {
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to get update time"})
			return
		} else {
			if updatedAt, ok := updatedAtAny.(string); ok {
				data["updated_at"] = updatedAt
			} else {
				data["updated_at"] = "unknown"
			}
		}

		if trending, err := unmarshalFromValkey([]byte(marshalled)); err != nil {
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "unmarshal error"})
			return
		} else {
			data["trending_events"] = trending
		}

		helper.WriteJSONResponse(w, http.StatusOK, data)
	}
}

// UpdateEventImageURLUsingExternalSource sets an event's image URL.
// @Summary      Update event image URL
// @Tags         events
// @Produce      json
// @Security     BearerAuth
// @Param        event_uuid  query     string  true  "Event UUID"
// @Param        image_url   query     string  true  "Image URL"
// @Success      200         {object} map[string]interface{}
// @Failure      400         {object} map[string]string
// @Failure      401         {object}  map[string]string
// @Failure      404         {object} map[string]string
// @Failure      500         {object} map[string]string
// @Router       /api/event/update/image [post]
func (h *Handlers) UpdateEventImageURLUsingExternalSource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(ctx.CtxUserIDKey).(uint)

		eventUUID := r.URL.Query().Get("event_uuid")
		if eventUUID == "" {
			helper.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "event_uuid is missing"})
			return
		}

		if event, err := h.er.GetEvent(r.Context(), eventUUID); err != nil {
			helper.WriteJSONResponse(w, http.StatusNotFound, map[string]string{"error": "failed to find event"})
			return
		} else if event.OrganizerID != userID {
			helper.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "access denied"})
			return
		}

		imageURL := r.URL.Query().Get("image_url")
		if imageURL == "" {
			helper.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "image_url is missing"})
			return
		}

		if err := h.er.UpdateEventImage(r.Context(), eventUUID, imageURL); err != nil {
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to update event image"})
			return
		}

		helper.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "success"})
	}
}

// UpdateEventImageURLByUploading uploads an event image.
// @Summary      Upload event image
// @Tags         events
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        event_uuid  query     string  true  "Event UUID"
// @Param        image       formData  file    true  "Image file (PNG, JPEG, or WebP; up to 10 MiB)"
// @Success      200         {object} map[string]interface{}
// @Failure      401         {object}  map[string]string
// @Failure      404         {object} map[string]string
// @Failure      500         {object} map[string]string
// @Router       /api/event/update/upload [post]
func (h *Handlers) UpdateEventImageURLByUploading() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(ctx.CtxUserIDKey).(uint)

		eventUUID := r.URL.Query().Get("event_uuid")
		if eventUUID == "" {
			helper.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "event_uuid is missing"})
			return
		}

		if event, err := h.er.GetEvent(r.Context(), eventUUID); err != nil {
			helper.WriteJSONResponse(w, http.StatusNotFound, map[string]string{"error": "failed to find event"})
			return
		} else if event.OrganizerID != userID {
			helper.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "access denied"})
			return
		}

		img, mimeType, err := parseImageFromRequest(r)
		if err != nil {
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to parse and read image"})
		}

		imageURL, err := h.lfs.SaveImage(img, fmt.Sprintf("%s.%s", eventUUID, mimeType))
		if err != nil {
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to save image"})
			return
		}

		if err := h.er.UpdateEventImage(r.Context(), eventUUID, imageURL); err != nil {

		}

		helper.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "success"})
	}
}

// DeleteEventByUUID deletes an event owned by the authenticated user.
// @Summary      Delete event
// @Tags         events
// @Produce      json
// @Security     BearerAuth
// @Param        event_uuid  query     string  true  "Event UUID"
// @Success      200         {object} map[string]interface{}
// @Failure      400         {object} map[string]string
// @Failure      401         {object}  map[string]string
// @Failure      404         {object} map[string]string
// @Failure      500         {object} map[string]string
// @Router       /api/event [delete]
func (h *Handlers) DeleteEventByUUID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(ctx.CtxUserIDKey).(uint)

		eventUUID := r.URL.Query().Get("event_uuid")
		if eventUUID == "" {
			helper.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "event_uuid is missing"})
			return
		}

		if event, err := h.er.GetEvent(r.Context(), eventUUID); err != nil {
			helper.WriteJSONResponse(w, http.StatusNotFound, map[string]string{"error": "failed to find event"})
			return
		} else if event.OrganizerID != userID {
			helper.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "access denied"})
			return
		}

		if err := h.er.DeleteEvent(r.Context(), eventUUID); err != nil {
			helper.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete event"})
			return
		}

		helper.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "success"})
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

func unmarshalFromValkey(marshalled []byte) ([]domain.EventStats, error) {
	var trending []domain.EventStats

	if err := json.Unmarshal(marshalled, &trending); err != nil {
		return nil, err
	}

	return trending, nil
}
