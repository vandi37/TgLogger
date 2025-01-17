package api

import (
	"encoding/json"
	"net/http"
)

// A response
type Response struct {
	Ok          bool   `json:"ok"`
	StatusCode  int    `json:"status_code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

type Request struct {
	Id            int64  `json:"id"`
	Text          string `json:"text"`
	Mode          string `json:"mode"`
	Notifications bool   `json:"notifications"`
	WebPreview    bool   `json:"web-preview"`
}

// Sends the response
func (r Response) Send(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(r)
}

// Sends an error
func SendError(w http.ResponseWriter, code int, err error) error {
	w.WriteHeader(code)
	return Response{
		Ok:          false,
		StatusCode:  code,
		Message:     http.StatusText(code),
		Description: err.Error(),
	}.Send(w)
}

// Sends a response
func Send(w http.ResponseWriter, description string) error {
	return Response{
		Ok:          true,
		StatusCode:  http.StatusOK,
		Message:     http.StatusText(http.StatusOK),
		Description: description,
	}.Send(w)
}
