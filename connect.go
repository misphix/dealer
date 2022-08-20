package main

import (
	"dealer/internal/configmanager"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func newMySQL(config configmanager.DatabaseConfig) (*gorm.DB, error) {
	return gorm.Open(mysql.New(mysql.Config{
		DSN:                       config.DSN,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})
}

func newMessageQueue(config configmanager.MessageQueueConfig) (*amqp.Channel, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if _, err := ch.QueueDeclare(config.QueueName, false, false, false, false, nil); err != nil {
		return nil, err
	}

	return ch, nil
}
