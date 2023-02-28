package worker_service

import (
	"context"
	"time"

	task_model "fan-out/internal/models/task"
)

type (
	Worker struct {
		taskChannel   chan task_model.Task
		resultChannel chan task_model.Task
	}
)

func New() *Worker {
	return &Worker{}
}

func (s *Worker) StartConsuming() {
	for {
		select {
		case task := <-s.taskChannel:
			s.ConsumeInMemory(task)
		}
	}
}

func (s *Worker) ConsumeInMemory(task task_model.Task) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()

	task.Err = s.DoStuff(ctx, task.Msg)

	s.PublishTaskResult(ctx, task)
}

func (s *Worker) PublishTaskResult(ctx context.Context, task task_model.Task) {
	s.resultChannel <- task
}

func (s *Worker) DoStuff(ctx context.Context, msg interface{}) error {
	return nil
}
