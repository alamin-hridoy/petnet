package mw

import (
	"net/http"
	"reflect"
	"strconv"

	pmpb "brank.as/rbac/gunk/v1/permissions"

	pfpb "brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	"github.com/spf13/viper"
)

func PetnetAdmin(config *viper.Viper) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logging.FromContext(r.Context()).WithField("middleware", "PetnetAdmin")

			if config.GetString("runtime.environment") == "localdev" {
				h.ServeHTTP(w, r)
				return
			}

			ot := mw.GetOrgType(r.Context())
			if ot == "" {
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
			oti, err := strconv.Atoi(ot)
			if err != nil {
				logging.WithError(err, log).Error("converting string to int")
				http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
				return
			}
			if pfpb.OrgType(oti) != pfpb.OrgType_PetNet {
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func ValidatePermission(config *viper.Viper, cl pmpb.ValidationServiceClient) func(h http.HandlerFunc, res, act string) http.Handler {
	return func(h http.HandlerFunc, res, act string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logging.FromContext(r.Context()).WithField("middleware", "ValidatePermission")
			ctx := r.Context()

			if config.GetString("runtime.environment") == "localdev" {
				h.ServeHTTP(w, r)
				return
			}

			if !mw.IsPetnetOwner(ctx) {
				if _, err := cl.ValidatePermission(ctx, &pmpb.ValidatePermissionRequest{
					ID:       mw.GetUserID(ctx),
					Action:   act,
					Resource: res,
					OrgID:    mw.GetOrgID(ctx),
				}); err != nil {
					logging.WithError(err, log).Error("permissions denied admin permission.")
					http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
					return
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}
