package main

import (
	"log"

	svcpostgres "brank.as/petnet/serviceutil/storage/postgres"
	_ "github.com/lib/pq"
)

func main() {
	if err := svcpostgres.Migrate(); err != nil {
		log.Fatal(err)
	}
}
