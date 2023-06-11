package httpjson

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// WriteResponse writes the given data as the JSON response.
func WriteResponse(w http.ResponseWriter, code int, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "JSON marshal")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
	return nil
}
