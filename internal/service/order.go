package service

import (
	"context"
	"dealer/internal/dao"
	"dealer/internal/models"

	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type OrderProcessorInterface interface {
	NewOrder(context.Context, *models.Order) error
	CancelOrder(context.Context, int64) error
}

type OrderProcessor struct {
	queueName string
	ch        *amqp.Channel
	db        *gorm.DB
	orderDAO  dao.OrderInterface
}

var _ OrderProcessorInterface = (*OrderProcessor)(nil)

func NewOrderProcessor(ch *amqp.Channel, queueName string, db *gorm.DB, orderDAO dao.OrderInterface) *OrderProcessor {
	return &OrderProcessor{
		ch:        ch,
		queueName: queueName,
		db:        db,
		orderDAO:  orderDAO,
	}
}

func (p *OrderProcessor) NewOrder(ctx context.Context, order *models.Order) error {
	if err := p.orderDAO.Insert(ctx, p.db, order); err != nil {
		return err
	}

	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return p.ch.PublishWithContext(ctx, "", p.queueName, false, false, amqp.Publishing{ContentType: "application/json", Body: data})
}

func (p *OrderProcessor) CancelOrder(ctx context.Context, orderID int64) error {
	order := &models.Order{ID: orderID, IsCancel: true}
	if err := p.orderDAO.Update(ctx, p.db, order); err != nil {
		return err
	}

	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return p.ch.PublishWithContext(ctx, "", p.queueName, false, false, amqp.Publishing{ContentType: "application/json", Body: data})
}
