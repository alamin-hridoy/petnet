package handler

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/h2non/filetype"
	"github.com/kenshaw/goji"

	"brank.as/petnet/cms/storage"
	fpb "brank.as/petnet/gunk/dsa/v2/file"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Server) getViewGCSFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	sess, _ := s.sess.Get(r, sessionCookieName)
	if s.gcs.IsMock() {
		b, err := base64.StdEncoding.DecodeString(storage.TestB64Image)
		if err != nil {
			log.Error("decoding test base64 image")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(200)
		w.Write(b)
		return
	}
	id := goji.Param(r, "id")
	if id == "" {
		log.Error("missing file id param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	queryParams := r.URL.Query()
	downloadable, err := url.PathUnescape(queryParams.Get("downloadable"))
	if err != nil {
		log.Error("unable to downloadable param")
	}
	oid, err := s.getOrgIdFromSessionUser(ctx, sess)
	if err != nil {
		log.Error("missing org id")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	uid, err := s.getUsrIdFromSessionUser(ctx, sess)
	if err != nil {
		log.Error("missing user id")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if !s.IsPetnetOwner(ctx, oid, uid) {
		res, err := s.pf.ListFiles(ctx, &fpb.ListFilesRequest{
			OrgID: oid,
		})
		if err != nil {
			logging.WithError(err, log).Error("getting files")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
		}
		var urls []string
		for _, f := range res.FileUploads {
			urls = append(urls, f.FileNames...)
		}
		var isOwner bool
		for _, u := range urls {
			if u == id {
				isOwner = true
				break
			}
		}
		if !isOwner {
			logging.WithError(err, log).Info("is not owner of file")
			http.Redirect(w, r, errorPath, http.StatusNotFound)
			return
		}
	}
	obj, err := s.gcs.Get(ctx, id)
	if err != nil {
		logging.WithError(err, log).Error("getting file from gcs")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	objAttrs, err := obj.Attrs(ctx)
	if err != nil {
		logging.WithError(err, log).Error("getting attributes for file")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	o := obj.ReadCompressed(true)
	rc, err := o.NewReader(ctx)
	if err != nil {
		logging.WithError(err, log).Error("creating new reader")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	defer rc.Close()
	buff, err := ioutil.ReadAll(rc)
	if err != nil {
		logging.WithError(err, log).Error("reading file in buffer")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	contentType := objAttrs.ContentType
	kind, err := filetype.Match(buff)
	if err == nil {
		contentType = kind.MIME.Value
	} else {
		logging.WithError(err, log).Warn("get content type error")
	}
	w.Header().Set("Content-Type", contentType)
	if downloadable == "true" {
		w.Header().Set("Content-Disposition", "attachment")
	}
	w.WriteHeader(200)
	_, err = w.Write(buff)
	if err != nil {
		logging.WithError(err, log).Error("writing file error")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
	}
}
