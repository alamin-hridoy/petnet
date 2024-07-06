package handler

import (
	"encoding/json"
	"net/http"

	fpb "brank.as/petnet/gunk/dsa/v2/file"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/kenshaw/goji"
)

type (
	RemoveFileForm struct {
		Name string
	}
)

func (s *Server) removeFile(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	sess, err := s.sess.Get(r, sessionCookieName)
	if err != nil {
		log.WithError(err).Error("fetching session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	oid, err := s.getOrgIdFromSessionUser(ctx, sess)
	if err != nil {
		log.Error("missing org id")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if oid == "" {
		log.Error("missing org id")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	fileID := goji.Param(r, "id")
	if fileID == "" {
		log.Error("missing file ID query param")
		return
	}
	var form RemoveFileForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if err := validation.ValidateStruct(&form,
		validation.Field(&form.Name, validation.Required),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["Name"] != nil {
				log.Error("File Name is Required")
			}
		}
		return
	}
	if _, err := s.pf.DeleteFileUpload(ctx, &fpb.DeleteFileUploadRequest{
		ID:        fileID,
		FileNames: form.Name,
		OrgID:     oid,
	}); err != nil {
		logging.WithError(err, log).Error(err.Error())
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
	}
	jsn, _ := json.Marshal([]string{})
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsn)
}
