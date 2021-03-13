package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	readerSleepDuration = time.Second
	writerSleepDuration = time.Millisecond * 50
	readerCount         = 2
	writerCount         = 2
	dataRandMax         = 32
	dataChannelBuffer   = 10
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func reader(waitGroup *sync.WaitGroup, dataChannel chan int, readerId int) {
	defer waitGroup.Done()

	// Loops until dataChannel gets closed. Also, when channel gets closed, we can still read the buffered data with
	// this for loop
	for data := range dataChannel {
		fmt.Printf("Reader %d: %d\n", readerId, data)
		time.Sleep(readerSleepDuration)
	}
}

func writer(ctx context.Context, waitGroup *sync.WaitGroup, dataChannel chan int, writerId int) {
	defer waitGroup.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("\t\t\t\tExiting Writer %d\n", writerId)
			return
		default:
			data := rand.Intn(dataRandMax)
			fmt.Printf("\t\tWriter %d: %d\n", writerId, data)
			dataChannel <- data
			time.Sleep(writerSleepDuration)
		}
	}
}

func main() {
	ctx, ctxCancel := context.WithCancel(context.Background())

	dataChannel := make(chan int, dataChannelBuffer)
	exitSignalChannel := make(chan os.Signal, 1)

	writerWaitGroup := sync.WaitGroup{}
	func(numberOfWorkers int) {
		for i := 0; i < numberOfWorkers; i++ {
			writerWaitGroup.Add(1)
			go writer(ctx, &writerWaitGroup, dataChannel, i)
		}
	}(writerCount)

	readerWaitGroup := sync.WaitGroup{}
	func(numberOfWorkers int) {
		for i := 0; i < numberOfWorkers; i++ {
			readerWaitGroup.Add(1)
			go reader(&readerWaitGroup, dataChannel, i)
		}
	}(readerCount)

	signal.Notify(exitSignalChannel, syscall.SIGINT, syscall.SIGTERM)

	<-exitSignalChannel
	ctxCancel()
	writerWaitGroup.Wait()
	close(dataChannel)
	readerWaitGroup.Wait()
}
