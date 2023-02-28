package kafka_consumer_service

import (
	"io"
	"log"

	task_model "fan-out/internal/models/task"
)

type (
	KafkaConsumerService struct {
		limiter     chan struct{}
		taskChannel chan task_model.Task
	}
)

func New(
	limit int,
	taskChannel chan task_model.Task,
) *KafkaConsumerService {
	return &KafkaConsumerService{
		limiter:     make(chan struct{}, limit),
		taskChannel: taskChannel,
	}
}

func (s *KafkaConsumerService) StartConsuming() error {
	for {
		msg, err := msgPopper()
		if err != nil {
			if err == io.EOF {
				log.Println("consumer consumer queue's last message")
				break // should we return an error instead?
			}

			log.Println("failed to consume the message")
			continue // it could be a parser error - so it continues
		}

		s.limiter <- struct{}{}
		go s.ConsumeMessage(msg) // consumes message async...
	}

	return nil
}

func (s *KafkaConsumerService) ConsumeMessage(msg interface{}) {
	defer func() {
		<-s.limiter
	}()

	// do any filtering or any verification before publishing the task

	s.taskChannel <- task_model.Task{Msg: msg}
}

func msgPopper() (interface{}, error) {
	return nil, nil
}
