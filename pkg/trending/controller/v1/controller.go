package v1

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/silverswords/onepiece/pkg/trending/model/v1"
)

type Controller struct {
	db *sql.DB
}

func New(db *sql.DB) *Controller {
	return &Controller{
		db: db,
	}
}

func (c *Controller) Create() error {
	if err := model.CreateSchema(c.db); err != nil {
		return err
	}

	if err := model.CreateDailyTable(c.db); err != nil {
		return err
	}

	if err := model.CreateRepoTable(c.db); err != nil {
		return err
	}

	return nil
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

	tx, err := c.db.Begin()
	if err != nil {
		log.Printf("[trending] save daily, begin tx error: %s\n", err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	for _, project := range req.DailyData {
		id, err := model.TxSelectRepoIDByName(tx, project.Name)
		if err != nil {
			log.Printf("[trending] save daily, select repo id error: %s\n", err)
			ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
			return
		}

		if err := model.TxInsertDailyTrending(tx, req.Date, id, project.Star, project.TodayStar, project.Fork); err != nil {
			log.Printf("[trending] save daily, select repo id error: %s\n", err)
			ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
