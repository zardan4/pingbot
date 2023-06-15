package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zardan4/pingbot/workerpool"
)

func main() {
	// configuration
	const (
		WORKERS_COUNT = 3
		TIMEOUT       = time.Second * 3
		INTERVAL      = time.Second * 5
	)
	results := make(chan workerpool.Result) // result channel

	// define all links
	links := []string{
		"https://www.udemy.com/",
		"https://github.com/",
		"https://goihab.com/",
		"https://chess.com/",
	}

	pool := workerpool.NewPool(WORKERS_COUNT, results, TIMEOUT)
	pool.Init()

	// generate/read results
	go generateJobs(INTERVAL, links, pool)
	go processResults(results)

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT) // записуємо сигнал, що закінчилася програма

	<-quit // чекаємо, коли поступить інфа про завершення програми. так як це канал, то читання з нього блокує процес проограм

	pool.Stop()

}

// генеримо Jobs раз на INTERVAL секунд
func generateJobs(interval time.Duration, links []string, pool *workerpool.Pool) {
	for {
		for _, l := range links {
			pool.Push(workerpool.Job{Url: l})
		}
		time.Sleep(interval)
	}
}

// виводимо результати
func processResults(results <-chan workerpool.Result) {
	for r := range results {
		fmt.Print(r.Info())
	}
}
