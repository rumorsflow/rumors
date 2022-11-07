package util

import (
	"encoding/json"
	"net/http"
)

type ListResponse struct {
	Data  any   `json:"data"`
	Total int64 `json:"total"`
	Index int64 `json:"index"`
	Size  int64 `json:"size"`
}

func BadRequest(err error) {
	panic(NewHTTPError(http.StatusBadRequest, err))
}

func OK(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, data)
}

func JSON(w http.ResponseWriter, code int, data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(bytes)
}

func Created(w http.ResponseWriter, location string) {
	if location != "" {
		w.Header().Set("Location", location)
	}
	w.WriteHeader(http.StatusCreated)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
