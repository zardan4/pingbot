package workerpool

import (
	"net/http"
	"time"
)

type worker struct {
	client *http.Client // custom client
}

// initialize new client with custom timeout
func newWorker(timeout time.Duration) *worker {
	return &worker{
		client: &http.Client{Timeout: timeout},
	}
}

func (w *worker) process(url string) Result {
	var result Result // result returned
	result.Url = url  // set url

	now := time.Now() // current time

	response, err := w.client.Get(url) // request

	result.ResponseTime = time.Since(now) // set response time
	if err != nil {                       // if error
		result.Error = err // set error
		return result
	}
	defer response.Body.Close()

	result.StatusCode = response.StatusCode // set status code

	return result
}
