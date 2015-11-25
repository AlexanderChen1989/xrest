package plugs

import (
	"encoding/json"
	"net/http"
)

// WriteJSON write data as json to http.ResponseWriter
func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Add("ContentType", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}
