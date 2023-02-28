package main

import (
	"context"

	kafka_consumer_service "fan-out/internal/services/consumer/kafka"
	result_processor_service "fan-out/internal/services/result-processor"
	worker_service "fan-out/internal/services/worker"
)

func main() {
	_, cnlFn := context.WithCancel(context.Background())

	kafkaConsumer := kafka_consumer_service.New()
	worker := worker_service.New()
	resultProcessor := result_processor_service.New()

	go func() {
		worker.StartConsuming()
	}()

	go func() {
		resultProcessor.StartConsuming()
	}()

	err := kafkaConsumer.StartConsuming()
	if err != nil {
		cnlFn()
		panic("consuming")
	}
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
