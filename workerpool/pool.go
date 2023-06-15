package workerpool

import (
	"fmt"
	"sync"
	"time"
)

// job
type Job struct {
	Url string
}

// result
type Result struct {
	Url          string
	StatusCode   int
	ResponseTime time.Duration
	Error        error
}

// return
func (r Result) Info() string {
	if r.Error != nil {
		return fmt.Sprintf("Error! Fetched %s. Error message: ", r.Url) + r.Error.Error() + "\n"
	}

	return fmt.Sprintf("Fetched %s. Status: %d. Time needed: %s\n", r.Url, r.StatusCode, r.ResponseTime)
}

type Pool struct {
	worker       *worker
	workersCount int

	results chan Result
	Jobs    chan Job

	wg      *sync.WaitGroup
	stopped bool
}

func NewPool(workersCount int, results chan Result, timeout time.Duration) *Pool {
	return &Pool{
		worker:       newWorker(timeout),
		workersCount: workersCount,

		results: results,
		Jobs:    make(chan Job),

		wg:      &sync.WaitGroup{},
		stopped: false,
	}
}

// initializing workers
func (p *Pool) Init() {
	for i := 0; i < p.workersCount; i++ {
		go p.InitWorker()
	}
}

func (p *Pool) InitWorker() {
	for j := range p.Jobs {
		p.results <- p.worker.process(j.Url)
		p.wg.Done()
	}
}

// push jobs
func (p *Pool) Push(job Job) {
	if p.stopped { // if program is stopped
		return
	}

	p.wg.Add(1)
	p.Jobs <- job
}

// stop pool
func (p *Pool) Stop() {
	p.stopped = true
	close(p.Jobs)
	p.wg.Wait()
}
