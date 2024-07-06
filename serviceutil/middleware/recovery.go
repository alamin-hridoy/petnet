package middleware

import (
	"runtime/debug"

	"brank.as/petnet/serviceutil/logging"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type recoveryHandler struct {
	logger *logrus.Entry
}

func newRecoveryHandler(l *logrus.Entry) *recoveryHandler {
	return &recoveryHandler{
		logger: l,
	}
}

func (r *recoveryHandler) recover(p interface{}) error {
	var entry *logrus.Entry
	if err, ok := p.(error); ok {
		// runtime panic
		entry = logging.WithError(err, r.logger)
	} else {
		// explicit panic call
		entry = r.logger.WithFields(logrus.Fields{
			"panic": p,
		})
	}
	entry.WithFields(logrus.Fields{
		"stacktrace": string(debug.Stack()),
	}).Error("recovered from panic")

	return status.Errorf(codes.Internal, "internal error occurred")
}
