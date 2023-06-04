package queue

import (
	"context"
	"encoding/json"
	"github.com/wagslane/go-rabbitmq"
	"soa_hw_4/internal/tasks"
)

type TaskPublisher struct {
	p *rabbitmq.Publisher
}

func NewTaskPublisher(conn *rabbitmq.Conn) (*TaskPublisher, error) {
	p, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("events"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, err
	}
	return &TaskPublisher{p}, nil
}

func (tp *TaskPublisher) Publish(task tasks.Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	err = tp.p.PublishWithContext(
		context.Background(),
		data,
		[]string{"ROUTE"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsExchange("tasks"),
	)
	return err
}
