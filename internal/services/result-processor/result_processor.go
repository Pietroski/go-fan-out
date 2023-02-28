package result_processor_service

import (
	"context"
	"fmt"
	"log"
	"time"

	task_model "fan-out/internal/models/task"
)

type (
	ResultProcessor struct {
		limiter       chan struct{}
		resultChannel chan task_model.Task
	}
)

func New(
	limit int,
	resultChannel chan task_model.Task,
) *ResultProcessor {
	return &ResultProcessor{
		limiter:       make(chan struct{}, limit),
		resultChannel: resultChannel,
	}
}

func (s *ResultProcessor) StartConsuming() {
	for {
		select {
		case task := <-s.resultChannel:
			s.limiter <- struct{}{}
			go s.ConsumeInMemory(task)
		}
	}
}

func (s *ResultProcessor) ConsumeInMemory(task task_model.Task) {
	defer func() {
		<-s.limiter
	}()
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()

	if task.Err != nil {
		err := s.NotifyTaskError(ctx, task.Err)
		if err != nil {
			err = fmt.Errorf("failed to notidy task error on ResultProcessor: %v", err)
			log.Println(err)

			return
		}
	}

	err := s.NotifyTaskSuccess(ctx, task.Msg)
	if err != nil {
		err = fmt.Errorf("failed to notidy task success on ResultProcessor: %v", err)
		log.Println(err)
	}
}

func (s *ResultProcessor) NotifyTaskError(ctx context.Context, taskErr error) error {
	return nil
}

func (s *ResultProcessor) NotifyTaskSuccess(ctx context.Context, taskMsg interface{}) error {
	return nil
}
