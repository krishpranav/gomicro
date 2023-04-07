package tool

import (
	"context"
	"fmt"
	"time"
)

type JobHandler = func([]byte) error

type JobResponseHandler = func(error)

type JobResponse struct {
	Error     error
	Completed bool
}

type Job struct {
	Handler    JobHandler
	Args       []byte
	DidSucceed JobResponseHandler
	DidFail    JobResponseHandler
	DidTimeOut JobResponseHandler
	Timeout    time.Duration
	RetryCount int
}

func RunJob(job Job) {
	ctx := context.Background()
	result := JobResponse{}
	if job.Timeout <= 0 {
		job.Timeout = 10 * time.Second
	}
	if job.RetryCount < 0 {
		job.RetryCount = 0
	}
	for i := 0; i <= job.RetryCount; i++ {
		context.WithTimeout(context.Background(), job.Timeout)
		result = runJobWithContext(ctx, job, job.Timeout)
		if result.Completed {
			jobDidComplete(job, result)
			return
		}
		if i >= job.RetryCount {
			runJobResponseHandler(job.DidTimeOut, result.Error)
		}
	}
}

func runJobWithContext(ctx context.Context, job Job, timeout time.Duration) JobResponse {
	channel := make(chan JobResponse)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	go func() {
		response := job.Handler(job.Args)
		channel <- JobResponse{Error: response, Completed: true}
	}()
	select {
	case <-ctx.Done():
		return JobResponse{Error: fmt.Errorf("The job could not be run in the given time slot (%s)", timeout)}
	case response := <-channel:
		return response
	}
}

func jobDidComplete(job Job, response JobResponse) {
	if err := response.Error; err != nil {
		runJobResponseHandler(job.DidFail, err)
		return
	}
	runJobResponseHandler(job.DidSucceed, nil)
}

func runJobResponseHandler(handler JobResponseHandler, err error) {
	if handler != nil {
		handler(err)
	}
}
