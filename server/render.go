package server

import (
	"encoding/json"
	"net/http"
)

func Json(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	bytes, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}

func Html(w http.ResponseWriter, status int, content string) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, err := w.Write([]byte(content))
	return err
}
