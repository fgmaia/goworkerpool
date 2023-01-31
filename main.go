// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"sync"
	"time"
)

type JobInput struct {
	startTime time.Time
	endTime   time.Time
}

type JobResponse struct {
}

type Response struct {
	finalOutput string
}

func startWorker(jobInputChan <-chan JobInput, jobOutputChan chan<- JobResponse, wg *sync.WaitGroup) {
	defer wg.Done()
	for jobInput := range jobInputChan {
		jobOutputChan <- businessFunctionalityJob(jobInput)
	}
}

func businessFunctionality(jobInputChan chan<- JobInput) {
	jobInputs := []JobInput{JobInput{}, JobInput{}}
	for _, jobInput := range jobInputs {
		jobInputChan <- jobInput
	}
	close(jobInputChan)
}

func combineResponses(jobResponses []JobResponse) Response {
	return Response{finalOutput: "Well done!!"}
}

func businessFunctionalityJob(jobInput JobInput) JobResponse {
	fmt.Println("Executing Job..")
	return JobResponse{}
}

func main() {
	num_of_workers := 3
	jobsChan := make(chan JobInput, 10)
	resultsChan := make(chan JobResponse, 10)
	wg := sync.WaitGroup{}

	for i := 0; i < num_of_workers; i++ {
		wg.Add(1)
		go startWorker(jobsChan, resultsChan, &wg)
	}

	go businessFunctionality(jobsChan)

	var responses []JobResponse
	wgResp := sync.WaitGroup{}
	wgResp.Add(1)
	go func() {
		defer wgResp.Done()
		for resp := range resultsChan {
			responses = append(responses, resp)
		}
	}()

	wg.Wait()
	close(resultsChan)
	wgResp.Wait()

	fmt.Println(combineResponses(responses).finalOutput)
}
