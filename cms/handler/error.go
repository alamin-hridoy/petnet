package handler

import (
	"net/http"

	"brank.as/petnet/serviceutil/logging"
)

type errorFormParams struct{}

func (s *Server) handleError(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	template := s.lookupTemplate("error.html")
	if template == nil {
		errMsg := "unable to load template"
		log.Error(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, &errorFormParams{}); err != nil {
		log.Infof("error with template execution: %+v", err)
	}
}
