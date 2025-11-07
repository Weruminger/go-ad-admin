package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Weruminger/go-ad-admin/internal/errs"
)

type httpErr struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	status := http.StatusInternalServerError
	code := errs.Internal
	msg := "Ein unerwarteter Fehler ist aufgetreten."

	var e *errs.E
	if errors.As(err, &e) {
		switch e.Code {
		case errs.InvalidInput:
			status = http.StatusUnprocessableEntity
			code = e.Code
			msg = "Eingabe ungültig."
		case errs.NotFound:
			status = http.StatusNotFound
			code = e.Code
			msg = "Nicht gefunden."
		case errs.Timeout:
			status = http.StatusServiceUnavailable
			code = e.Code
			msg = "Timeout / Dienst nicht erreichbar."
		case errs.Unavailable:
			status = http.StatusServiceUnavailable
			code = e.Code
			msg = "Dienst vorübergehend nicht verfügbar."
		default:
			code = e.Code
		}
	}
	if rid := r.Header.Get("X-Request-ID"); rid != "" {
		w.Header().Set("X-Request-ID", rid)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(httpErr{Code: string(code), Message: msg, RequestID: reqIDFrom(r)})
}
