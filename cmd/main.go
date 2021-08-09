package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/silverswords/onepiece/pkg/register"
	trendingV1 "github.com/silverswords/onepiece/pkg/trending/controller/v1"
)

const (
	version = 1

	listenAddr = "0.0.0.0:8080"
)

const (
	trendingServiceName = "trending"
)

func main() {
	db, err := sql.Open("postgres", "host=192.168.0.251 port=5432 user=root password=123456 dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	engine := gin.Default()

	register.Register(1, trendingServiceName, trendingV1.New(db))

	if err := register.Init(version, engine); err != nil {
		log.Fatal(err)
	}

	if err := engine.Run(listenAddr); err != nil {
		log.Fatal(err)
	}
}
