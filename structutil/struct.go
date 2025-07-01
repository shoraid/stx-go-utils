package structutil

import (
	"encoding/json"
	"net/http"

	"github.com/shoraid/stx-go-utils/apperror"
)

func BindJSON(r *http.Request, input any) error {
	if r.Body == nil {
		return apperror.Err400InvalidBody
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(input); err != nil {
		return err
	}

	return nil
}
