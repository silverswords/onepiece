package main

import (
	"log"

	"github.com/gin-gonic/gin"
	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/silverswords/onepiece/pkg/register"
	trendingV1 "github.com/silverswords/onepiece/pkg/trending/controller/v1"
)

const (
	influxAddr  = "http://localhost:8086"
	influxToken = "EFm9R2pGsgh1E7JHHBAnCsxp1EcjyepOytj1PUqyMkKqByulwAxfKfbvIRM0IOg-dg_SyNeODPcqugCTB48fQw=="
	version     = 1

	listenAddr = "0.0.0.0:8080"
)

const (
	trendingServiceName = "trending"
)

func main() {
	client := influxdb.NewClient(influxAddr, influxToken)
	defer client.Close()

	engine := gin.Default()

	register.Register(1, trendingServiceName, trendingV1.New(client))

	err := register.Init(version, engine)
	if err != nil {
		log.Fatal(err)
	}

	if err := engine.Run(listenAddr); err != nil {
		log.Fatal(err)
	}
}
