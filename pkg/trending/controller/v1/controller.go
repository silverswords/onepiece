package v1

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	influxdb "github.com/influxdata/influxdb-client-go/v2"

	"github.com/silverswords/onepiece/pkg/trending/model/v1"
)

type Controller struct {
	client influxdb.Client
}

func New(client influxdb.Client) *Controller {
	return &Controller{
		client: client,
	}
}

func (c *Controller) Create() error {
	return model.Create(c.client)
}

func (c *Controller) Update() error {
	return nil
}

func (c *Controller) Register(router gin.IRouter) {
	router.POST("/daily/save", c.saveDaily)
}

func (c *Controller) saveDaily(ctx *gin.Context) {
	var req struct {
		Date      time.Time        `json:"date,omitempty"`
		DailyData []*model.Project `json:"daily,omitempty"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		log.Printf("[trending] save daily, binding error: %s\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	if err := model.SaveDailyTrending(c.client, req.Date, req.DailyData); err != nil {
		log.Printf("[trending] save daily, db error: %s\n", err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
