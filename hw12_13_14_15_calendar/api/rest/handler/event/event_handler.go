package event

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/arny_tiger/hw-test/hw12_13_14_15_calendar/internal/storage/entity"
)

type Handler struct {
	storage storage.Storage
	logger  logger.Logger
}

func NewHandler(storage storage.Storage, logger logger.Logger) Handler {
	return Handler{
		storage: storage,
		logger:  logger,
	}
}

func (h Handler) GetDispatcher(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getEvents(w, r)
	case http.MethodPost:
		h.createEvent(w, r)
	case http.MethodPut:
		h.updateEvent(w, r)
	case http.MethodDelete:
		h.deleteEvent(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h Handler) createEvent(w http.ResponseWriter, r *http.Request) {
	var event entity.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if event.ID == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err != nil {
		h.respondWithInternalError(w, "CreateEvent failed", err)
		return
	}
	id, err := h.storage.CreateEvent(event)
	if err != nil {
		h.respondWithInternalError(w, "CreateEvent failed", err)
		return
	}
	idJSON, err := json.Marshal(id)
	if err != nil {
		h.respondWithInternalError(w, "CreateEvent failed", err)
		return
	}
	h.respondWithJSON(w, idJSON)
}

func (h Handler) getEvents(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit, err := getIntVal("limit", query)
	if err != nil || limit == 0 {
		limit = 20
	}
	offset, err := getIntVal("offset", query)
	if err != nil {
		offset = 0
	}
	dbEvents, err := h.storage.GetEvents(limit, offset)
	if err != nil {
		h.respondWithInternalError(w, "GetEvents failed", err)
		return
	}
	h.respondWithJSON(w, dbEvents)
}

func (h Handler) updateEvent(w http.ResponseWriter, r *http.Request) {
	var event entity.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if event.ID == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err != nil {
		h.respondWithInternalError(w, "UpdateEvent failed", err)
		return
	}
	_, err = h.storage.UpdateEvent(event)
	if err != nil {
		h.respondWithInternalError(w, "UpdateEvent failed", err)
		return
	}
}

func (h Handler) deleteEvent(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id, err := getIntVal("id", query)
	if err != nil {
		h.respondWithInternalError(w, "DeleteEvent failed", err)
		return
	}
	if id == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = h.storage.DeleteEvent(id)
	if err != nil {
		h.respondWithInternalError(w, "DeleteEvent failed", err)
		return
	}
}

func (h Handler) respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		err = fmt.Errorf("encoding failed: %w", err)
		h.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h Handler) respondWithInternalError(w http.ResponseWriter, msg string, err error) {
	respErr := fmt.Errorf(msg+": %w", err)
	h.logger.Error(respErr.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func getIntVal(param string, q url.Values) (int, error) {
	if q.Get(param) != "" {
		var err error
		var res int

		res, err = strconv.Atoi(q.Get(param))
		if err != nil {
			return 0, err
		}
		return res, nil
	}
	return 0, nil
}
