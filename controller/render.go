package controller

import (
	"encoding/json"
	"net/http"
)

func GetRender(w http.ResponseWriter) Rendering {
	return Rendering{
		writer: w,
	}
}

type Rendering struct {
	writer http.ResponseWriter
}

func (r Rendering) JSON(status int, value interface{}) {
	response, _ := json.Marshal(value)
	r.Content(status, "application/json", string(response))
}

func (r Rendering) HTML(status int, html string) {
	r.Content(status, "text/html", html)
}

func (r Rendering) Content(status int, contentType string, content string) {
	r.writer.Header().Set("Content-Type", contentType)
	r.writer.WriteHeader(status)
	r.writer.Write([]byte(content))
}
