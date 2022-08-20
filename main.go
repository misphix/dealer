package main

import (
	"context"
	"dealer/internal/configmanager"
	"dealer/internal/dao"
	"dealer/internal/handler"
	"dealer/internal/models"
	"dealer/internal/service"
	"encoding/json"

	"dealer/internal/logger"
	"fmt"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
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

	ch, err := newMessageQueue(config.MessageQueue)
	if err != nil {
		panic(err)
	}

	orderDAO := dao.NewOrder()
	dealDAO := dao.NewDeal()
	orderProcessor := service.NewOrderProcessor(ch, config.MessageQueue.QueueName, db, orderDAO)
	dealer := service.NewDealer(db, orderDAO, dealDAO)
	h := handler.NewHandler(orderProcessor)

	if err := startConsumer(ch, config.MessageQueue.QueueName, dealer); err != nil {
		panic(err)
	}

	engine := gin.New()
	handler.RegisterRoutes(engine, h)

	engine.Run(fmt.Sprintf(":%d", config.HTTPServer.Port))
}

func startConsumer(ch *amqp.Channel, name string, dealer service.DealerInterface) error {
	msgs, err := ch.Consume(name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var order *models.Order
			if err := json.Unmarshal(msg.Body, &order); err != nil {
				logger.GetLogger().Error(err.Error())
				continue
			}

			if err := dealer.ProcessOrder(context.Background(), order); err != nil {
				logger.GetLogger().Error(err.Error())
				continue
			}
		}
	}()

	return nil
}
