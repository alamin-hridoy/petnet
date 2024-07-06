package keto

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func testConfig(t *testing.T) *viper.Viper {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("../../env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		t.Fatalf("error loading configuration: %v", err)
	}
	return config
}

func TestSvc_CreatePermission(t *testing.T) {
	baseURL := os.Getenv("KETO_URL")
	if baseURL == "" {
		t.Skip("missing env 'KETO_URL'")
	}
	t.Parallel()
	s := New(baseURL)

	tests := []struct {
		name        string
		p           Permission
		wantErr     bool
		wantEmptyID bool
	}{
		{
			name: "AllowCRUDPermission",
			p: Permission{
				Description: "Allow CRUD permission",
				Environment: "development",
				Allow:       true,
				Actions:     []string{"delete", "create", "read", "modify"},
				Resources:   []string{"test-resource"},
				Groups:      []string{"test-group"},
			},
		},
		{
			name:        "EmptyPermission",
			wantEmptyID: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			got, err := s.CreatePermission(ctx, test.p)
			if (err != nil) != test.wantErr {
				t.Fatalf("Svc.CreatePermission() error = %v, wantErr %v", err, test.wantErr)
			}
			t.Cleanup(func() { s.DeletePermission(ctx, got) })
			if got == "" && !test.wantEmptyID {
				t.Fatalf("got empty ID")
			}
		})
	}
}
