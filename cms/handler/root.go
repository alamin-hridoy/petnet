package handler

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"

	"brank.as/petnet/serviceutil/logging"
)

type rootFormParams struct {
	CSRFField template.HTML
	FirstName string
}

func (s *Server) getRoot(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	template := s.lookupTemplate("root.html")
	if template == nil {
		log.Error("unable to load template")
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, &rootFormParams{
		CSRFField: csrf.TemplateField(r),
	}); err != nil {
		log.Infof("error with template execution: %+v", err)
	}
}
