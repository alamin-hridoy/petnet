package storage

import (
	"context"
	"errors"

	gcs "cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

type ResourceType string

const (
	UnknownResource ResourceType = "unknown"
	Status          ResourceType = "status"
	EnableMFA       ResourceType = "enable-mfa"
	RolePermission  ResourceType = "role-permission"
	RoleDelete      ResourceType = "role-delete"
	UserDelete      ResourceType = "user-delete"
	InviteUser      ResourceType = "invite-user"
	CreateRole      ResourceType = "create-role"
	ChangePassword  ResourceType = "change-password"
	InviteProvider  ResourceType = "invite-provider"
)

// MFANotFound is returned when mfa is either disabled or doesn't exist.
var MFANotFound = errors.New("mfa not found")

func NewOnboardingGCSStorage(logger logrus.FieldLogger, config *viper.Viper) (GCS, error) {
	var gcsClient *gcs.Client

	bucketName := config.GetString("gcs.bucketName")
	if bucketName == "" {
		logger.Warn("GCS bucket name is not available, using fake storage")
		return NewFakeStorage(), nil
	}

	var opts []option.ClientOption
	if credentialsJSON := config.GetString("gcs.credentialsJSON"); credentialsJSON != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(credentialsJSON)))
	}
	if len(opts) == 0 {
		logger.Warn("GCS authentication configuration for logo storage is not available, using fake storage")
		return NewFakeStorage(), nil
	}

	gcsClient, err := gcs.NewClient(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	return NewGoogleCloudStorage(gcsClient, bucketName), nil
}
