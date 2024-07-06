package main

import (
	"log"

	_ "github.com/lib/pq"

	svcpostgres "brank.as/rbac/serviceutil/storage/postgres"
)

func main() {
	if err := svcpostgres.Migrate(); err != nil {
		log.Fatal(err)
	}
}
