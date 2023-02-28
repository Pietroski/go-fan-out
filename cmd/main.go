package main

import (
	"context"
	task_model "fan-out/internal/models/task"
	"log"
	"runtime"

	kafka_consumer_service "fan-out/internal/services/consumer/kafka"
	result_processor_service "fan-out/internal/services/result-processor"
	worker_service "fan-out/internal/services/worker"
)

func main() {
	_, cnlFn := context.WithCancel(context.Background())

	limit := 5 // taken by and env variable
	taskChannel := make(chan task_model.Task, 50)
	resultChannel := make(chan task_model.Task, 50)

	kafkaConsumer := kafka_consumer_service.New(limit, taskChannel)
	worker := worker_service.New(taskChannel, resultChannel)
	resultProcessor := result_processor_service.New(limit, resultChannel)

	workerPoolNumber := 10 // taken by and env variable
	for idx := 0; idx < workerPoolNumber; idx++ {
		go func() {
			worker.StartConsuming()
		}()
	}

	go func() {
		resultProcessor.StartConsuming()
	}()

	log.Println("GOROUTINE COUNT ->", runtime.NumGoroutine())
	err := kafkaConsumer.StartConsuming()
	if err != nil {
		cnlFn()
		log.Println("GOROUTINE COUNT ->", runtime.NumGoroutine())
		panic("consuming")
	}
	log.Println("GOROUTINE COUNT ->", runtime.NumGoroutine())
}

//var (
//	wg = &sync.WaitGroup{}
//)
//
//func main() {
//	kafkaConsumer := kafka_consumer_service.New()
//	worker := worker_service.New()
//	resultProcessor := result_processor_service.New()
//
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		worker.StartConsuming()
//	}()
//
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		resultProcessor.StartConsuming()
//	}()
//
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		_ = kafkaConsumer.StartConsuming()
//	}()
//	wg.Wait()
//}
