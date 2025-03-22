package transportlayer

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/YankovskiyVS/eventProject/events/internal/models"
	"github.com/gorilla/mux"
)

func (s *APIServer) handleGetEvent(w http.ResponseWriter, r *http.Request) error {
	//Check the role of the user from request header
	role := r.Header.Get("X-User-Role")
	if role != "user" {
		return errors.New("status forbidden")
	}

	//Get ID from URI
	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["id"])
	if err != nil || eventID < 1 {
		return errors.New("invalid id format")
	}

	//Connect handler with the method to DB
	event, err := s.event.GetEvent(uint(eventID))
	if err != nil {
		return errors.New("internal server error")
	}

	return WriteJSON(w, http.StatusOK, event)
}

func (s *APIServer) handleListEvents(w http.ResponseWriter, r *http.Request) error {
	//Check the role of the user from request header
	role := r.Header.Get("X-User-Role")
	if role != "user" {
		return errors.New("status forbidden")
	}

	//Get URI query
	query := r.URL.Query()

	//Date filter
	var dateTo *time.Time
	if dateStr := query.Get("dateTo"); dateStr != "" {
		parsedDate, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return WriteJSON(w, http.StatusBadRequest, "Invalid date format, use RFC3339")
		}
		dateTo = &parsedDate
	}

	//Parse pagination with defaults
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}

	itemsCount, _ := strconv.Atoi(query.Get("items_count"))
	if itemsCount < 1 {
		itemsCount = 20
	} else if itemsCount > 100 {
		itemsCount = 100
	}

	//Connect handler with the method to DB
	events, total, err := s.event.ListEvents(dateTo, page, itemsCount)
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, "Could not retrieve events")
	}

	//Calculate total pages
	totalPages := total / itemsCount
	if total%itemsCount != 0 {
		totalPages++
	}

	response := models.EventListResponse{
		Data:       events,
		Page:       page,
		ItemsCount: itemsCount,
		TotalItems: total,
		TotalPages: totalPages,
	}

	return WriteJSON(w, http.StatusOK, response)
}
