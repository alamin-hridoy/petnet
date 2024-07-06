package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/sirupsen/logrus"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/serviceutil/slack"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type recoveryHandler struct {
	logger        *logrus.Entry
	slackHookURL  string
	slackChannel  string // optional
	slackUsername string // optional
	slackIconURL  string // optional
}

func newRecoveryHandler(l *logrus.Entry, slackHookURL, slackChannel, slackUsername, slackIconURL string) *recoveryHandler {
	return &recoveryHandler{
		logger:        l,
		slackHookURL:  slackHookURL,
		slackChannel:  slackChannel,
		slackUsername: slackUsername,
		slackIconURL:  slackIconURL,
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

	stack := string(debug.Stack())

	entry.WithFields(logrus.Fields{
		"stacktrace": stack,
	}).Error("recovered from panic")

	if r.slackHookURL != "" {
		slack.PostToSlackWithRetries(entry, r.slackHookURL, fmt.Sprintf("Unexpected panic: %+v\n\n%s", p, stack),
			r.slackChannel, r.slackUsername, r.slackIconURL)
	}

	return status.Errorf(codes.Internal, "internal error occurred")
}
