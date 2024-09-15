package logger

import (
	"context"
	"customers_kuber/closer"
	"fmt"
	"sync"
)

var logKafkaWorkerInstance LogKafkaWorker

// logKafkaWorker создает пул воркеров, которые забирают сообщения из канала для логирования
// и передают их в kafkaLogProducer, количество воркеров регулируется через конфиг
type logKafkaWorker struct {
	wg               *sync.WaitGroup
	cancelFunc       context.CancelFunc
	logChannel       chan string
	kafkaLogProducer LogProducer
}

type LogKafkaWorker interface {
	startWorkers(context.Context)
	stopWorkers()
	GetLogChannel() chan string
}

// GetLogKafkaWorker принимает канал, в который логгер будет отправлять логи для записи в kafka
// создает LogKafkaWorker, который будет забирать логи из канала и передавать их в kafkaLogProducer
// если интерфейс уже создан, вернет существующий
func GetLogKafkaWorker(logChannel chan string) LogKafkaWorker {
	if logKafkaWorkerInstance != nil {
		return logKafkaWorkerInstance
	}
	producer, err := GetLogProducer()
	if err != nil {
		fmt.Printf("пока лог через фмт, не получилось создать лог продюсер %s", err)
	}
	worker := logKafkaWorker{
		wg:               new(sync.WaitGroup),
		logChannel:       logChannel,
		kafkaLogProducer: producer,
	}
	worker.startWorkers(context.Background())
	closer.CloseFunctions = append(closer.CloseFunctions, func() { worker.stopWorkers() })
	return &worker
}

func (worker *logKafkaWorker) startWorkers(ctx context.Context) {
	ctx, cancelFunc := context.WithCancel(ctx)
	worker.cancelFunc = cancelFunc

	//todo вынести количество воркеров в конфиг
	for i := 0; i <= 15; i++ {
		worker.wg.Add(1)
		go worker.spawnWorker(ctx)
	}
}

func (worker *logKafkaWorker) stopWorkers() {
	worker.cancelFunc()
	worker.wg.Wait()
	fmt.Println("workers done")
}

func (worker *logKafkaWorker) GetLogChannel() chan string {
	return worker.logChannel
}

func (worker *logKafkaWorker) spawnWorker(ctx context.Context) {
	defer worker.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case value := <-worker.logChannel:
			worker.kafkaLogProducer.ProduceLogToKafka([]byte(value))
		}
	}
}
