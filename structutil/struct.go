package structutil

import (
	"encoding/json"
	"io"
)

// decodes JSON from io.Reader (e.g., http.Request.Body)
func DecodeJSON(body io.Reader, target any) error {
	data, err := io.ReadAll(body)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, target); err != nil {
		return err
	}

	return nil
}
