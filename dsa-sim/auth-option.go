package main

import (
	"net/http"
)

func (s *Server) getAuthOption(w http.ResponseWriter, r *http.Request) {
	template := s.templates.Lookup("auth-option.html")
	if template == nil {
		http.Error(w, "template doesn't exist", http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, nil); err != nil {
		http.Error(w, "executing template", http.StatusInternalServerError)
		return
	}
}
