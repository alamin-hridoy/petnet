package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/svcutil/mainpkg"
)

func main() {
	cfg := viper.New()
	log := logging.NewLogger()
	log.SetFormatter(&logrus.JSONFormatter{})
	var initTest bool
	flag.BoolVar(&initTest, "init", false, "fail on init")
	flag.Parse()
	svr, err := mainpkg.Setup(cfg, log.WithField("test", "mainpkg"),
		mainpkg.WithPath("/tmp/test.socket"),
		mainpkg.WithInit("inittest", func(ctx context.Context) error {
			log.WithField("init testing", initTest).Info("Init Func Called")
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Second):
				log.Info("Init Func complete")
			}
			if initTest {
				return fmt.Errorf("test init failed")
			}
			return nil
		}),
		mainpkg.WithCleanup("cleanup", func(ctx context.Context) error {
			log.Info("Cleanup Func Called")
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(3 * time.Second):
				log.Info("Cleanup Func complete")
			}
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	svr.Run()
}
