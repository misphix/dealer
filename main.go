package main

import (
	"dealer/internal/configmanager"
	"dealer/internal/dao"
	"dealer/internal/handler"
	"dealer/internal/service"

	"dealer/internal/logger"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := configmanager.Get()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	logger.SetLevel(config.Logger.Level)
	l := logger.GetLogger()
	defer l.Sync()

	db, err := newMySQL(config.Database)
	if err != nil {
		panic(err)
	}

	dealer := service.NewDealer(db, dao.NewOrder(), dao.NewDeal())
	h := handler.NewHandler(dealer)
	engine := gin.New()
	handler.RegisterRoutes(engine, h)

	engine.Run(fmt.Sprintf(":%d", config.HTTPServer.Port))
}
