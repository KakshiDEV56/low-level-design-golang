/*
Question 3 — System Design (learning simulation)

You're building a backend that accepts jobs quickly and processes them in the background.
Example: user uploads a report, a worker generates a PDF later.

This file is a skeleton. Fill in the TODOs. The comments are hints, not answers.
*/

package main

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"sync"
	"time"
)

// ---------------------------------------------------------
// Core idea
//
// HTTP handler should NOT wait for PDF generation.
// Flow you are simulating:
//
//   client -> SubmitJob() -> enqueue job -> return jobID immediately
//                 |
//                 v
//            job queue (channel)
//                 |
//                 v
//            idle worker takes job -> busy -> process -> idle again
//
// A worker "becoming available" is just: the goroutine finishes one job
// and loops back to wait on the queue.
// ---------------------------------------------------------

type JobStatus string

// global retry mechanism rather than creating it

const (
	StatusQueued     JobStatus = "queued"
	StatusProcessing JobStatus = "processing"
	StatusDone       JobStatus = "done"
	StatusFailed     JobStatus = "failed"
	randomLimit      int       = 100
	StatusCrashed    JobStatus = "Crashed"
)

type Job struct {
	ID     string
	Report *report
	Status JobStatus

	// HINT: you may want CreatedAt, Attempts, Result, Error
}

type (
	report struct{}
	Pdf    struct{}
)

// JobStore is the "database" in this simulation.
// HINT: a map is enough. You will need a mutex if multiple goroutines touch it.
type JobStore struct {
	jobs  map[string]*Job
	retry map[string]*Job
	mu    sync.Mutex
}

func NewJobStore() *JobStore {
	return &JobStore{jobs: make(map[string]*Job), retry: make(map[string]*Job)}
}

func (s *JobStore) Save(job *Job) {
	// TODO: store the job by ID
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.ID] = job
}

func (s *JobStore) Get(id string) (*Job, bool) {
	// TODO: look up a job so the client can poll status

	s.mu.Lock()
	defer s.mu.Unlock()
	_, exist := s.jobs[id]
	// earlier i was returning the direct job pointer
	// that why now i am returning the snapshot pointer which is just the copy of the
	// job pointer
	snapshot := s.jobs[id]
	if exist {
		return snapshot, true
	}

	return nil, false
}

func (s *JobStore) UpdateStatus(id string, status JobStatus) {
	// TODO: queued -> processing -> done/failed
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[id]; ok {
		s.jobs[id].Status = status
	}
}

// ---------------------------------------------------------
// Queue + workers
//
// In Go, a buffered channel is a simple in-memory queue.
//
//   jobs := make(chan *Job, 64)
//
// SubmitJob sends into this channel and returns.
// Workers block on `for job := range jobs` until work arrives.
// ---------------------------------------------------------

type Dispatcher struct {
	jobs       chan *Job
	numWorkers int
	store      *JobStore
}

func NewDispatcher(numWorkers int, queueSize int, store *JobStore) *Dispatcher {
	return &Dispatcher{
		jobs:       make(chan *Job, queueSize),
		numWorkers: numWorkers,
		store:      store,
	}
}

// starting the worker go routine
func (d *Dispatcher) Start() {
	// TODO: start numWorkers goroutines
	for i := 0; i < d.numWorkers; i++ {
		go d.worker(i)
	}

	// HINT: do not wait here. Start() should return immediately.
}

// i dont see the use of the workerId anywhere
func (d *Dispatcher) worker(id int) {
	fmt.Printf("[worker-%d] idle, waiting for a job\n", id)

	// this loop is the whole worker: wait -> busy -> idle -> wait
	for job := range d.jobs {
		fmt.Printf("[worker-%d] took job %s (processing)  jobs-waiting-in-line=%d/%d\n",
			id, job.ID, len(d.jobs), cap(d.jobs))

		// Dequeue already happened: this job is no longer in d.jobs.
		// dont do this it is bad approach use the function updateStatus rather than manually doing this
		// updating the status sutff
		// job.Status = StatusProcessing
		d.store.UpdateStatus(job.ID, StatusProcessing)
		// Crash AFTER dequeue, BEFORE done. Do not mark Failed/Done.
		// Poll will show this job stuck at processing — that is the lost-from-channel case.
		if rand.IntN(3) == 0 {
			fmt.Printf("[worker-%d] CRASH after dequeue job %s — not in queue anymore, status stays processing\n",
				id, job.ID)
			job.Status = StatusCrashed
			d.store.retry[job.ID] = job

			continue
		}
		// if i found anything that is crashing i will use the retry mechansim until i will get successsed
		_, err := processJob(job)
		if err != nil {
			fmt.Printf("[worker-%d] job %s failed: %v (idle again)\n", id, job.ID, err)
			//job.Status = StatusFailed
			d.store.UpdateStatus(job.ID, StatusFailed)
		} else {
			fmt.Printf("[worker-%d] finished job %s (idle again)\n", id, job.ID)
			//job.Status = StatusDone
			d.store.UpdateStatus(job.ID, StatusDone)
		}
	}
}

func (d *Dispatcher) retryJobs() {
	for id, job := range d.store.retry {
		// Crashed jobs AND queue-full rejects both live here.
		// Rejects are still StatusQueued; only skip work that already finished.
		if job.Status == StatusDone {
			delete(d.store.retry, id)
			continue
		}

		was := job.Status
		job.Status = StatusQueued
		select {
		case d.jobs <- job:
			fmt.Printf("[retry] re-queued job %s (was %s)\n", id, was)
			delete(d.store.retry, id)
		default:
			fmt.Printf("[retry] queue full, will try job %s later\n", id)
		}
	}
}

func processJob(job *Job) (*Pdf, error) {
	fmt.Printf("[job %s] generating pdf (busy ~2s)\n", job.ID)
	time.Sleep(2 * time.Second)

	return &Pdf{}, nil
}

// doing some refactor right and see how i am able to fix the issue
// SubmitJob is the HTTP handler equivalent.
// It must return FAST even if every worker is busy.
func (d *Dispatcher) SubmitJob(r *report) (jobID string, err error) {
	// TODO:
	// 1. create a Job with a unique ID and StatusQueued
	// 2. save it in the store FIRST (so it is not lost if we crash after accept)
	//  return the jobID immediately — do NOT call processJob here
	//
	// HINT: if the queue is full, a blocking send waits. For "accept quickly"
	// try a non-blocking send:
	//
	// auto rejection mechanism when the queue is full

	randomNum := rand.IntN(randomLimit)
	for {

		// Check if the map contains the random_num as a key
		if _, exists := d.store.jobs[strconv.Itoa(randomNum)]; !exists {
			// Found a unique one! Mark it as used in the map
			break
		}
		randomNum = rand.IntN(randomLimit)
	}
	// checking if the random_num is unique or not

	UniqueId := strconv.Itoa(randomNum)
	job := &Job{
		ID:     UniqueId,
		Report: r,
		Status: StatusQueued,
	}

	d.store.jobs[UniqueId] = job
	// just going to crash the pogram after the dequeue right

	select {

	case d.jobs <- job:
		fmt.Printf("[queue] job %s enqueued  waiting-in-line=%d/%d\n",
			job.ID, len(d.jobs), cap(d.jobs))

	default:
		fmt.Printf("[queue] FULL waiting-in-line=%d/%d  bouncing job %s\n",
			len(d.jobs), cap(d.jobs), job.ID)
		d.store.retry[job.ID] = job
		// still return the ticket so the client can poll / retry can find it
		return UniqueId, fmt.Errorf("queue full")
	}
	// now i will make sure that the retrying of the failed event took place
	return UniqueId, nil
}

func (d *Dispatcher) GetJob(id string) (*Job, bool) {
	return d.store.Get(id)
}

func main() {
	store := NewJobStore()
	d := NewDispatcher(3 /* workers */, 3 /* queue size */, store)

	fmt.Println("=== 3 workers, queue holds 3 waiting jobs, 10 submit attempts ===")
	fmt.Println("=== [client] is the API  |  [queue] is the line  |  [worker] is background ===")
	fmt.Println()

	d.Start()

	tickets := make([]string, 0, 10)
	n := 10
	// that a better way to do the task right

	for i := 0; i < n; i++ {
		id, err := d.SubmitJob(&report{})
		if id != "" {
			tickets = append(tickets, id)
		}
		if err != nil {
			fmt.Printf("[client] request %d rejected ticket=%s: %v  (saved for retry)\n", i, id, err)
			continue
		}
		fmt.Printf("[client] request %d accepted ticket=%s  (not waiting for the pdf)\n", i, id)
	}
	/// fixing the loop
	poll := func(rounds int, label string) {
		fmt.Println()
		fmt.Println(label)
		for tick := 1; tick <= rounds; tick++ {
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("\n--- poll %d ---\n", tick)
			pending := 0
			for _, id := range tickets {
				job, ok := d.GetJob(id)
				if !ok {
					fmt.Printf("[poll] job %s not in store\n", id)
					pending++
					continue
				}
				if job.Status == StatusDone {
					continue
				}
				pending++
				fmt.Printf("[poll] job %s status=%s\n", id, job.Status)
			}
			if pending == 0 {
				fmt.Println("[poll] all jobs done")
				return
			}
		}
	}

	poll(12, "=== submits are done. client polls while workers finish ===")

	for round := 1; round <= 6; round++ {
		if len(d.store.retry) == 0 {
			break
		}
		fmt.Printf("\n=== retry round %d (%d jobs waiting to retry) ===\n", round, len(d.store.retry))
		d.retryJobs()
		poll(8, "=== client polls after retry ===")
	}

	fmt.Println()
	fmt.Println("=== sim finished ===")
}

// ---------------------------------------------------------
// Ideas to implement next (pick any, in this order)
//
// 1. Happy path - DONE
//    Submit 10 jobs, 3 workers. Watch jobs wait in the queue while
//    workers are busy, then get picked up as workers free themselves.
//
// 2. Job status polling - DONE
//    Client does not wait for the PDF. It gets a jobID, then polls GetJob.
//    That is how async APIs work (202 Accepted + GET /jobs/:id).
//
// 3. Bounded queue - Done this task is somewhat completed i guess
//    If the channel is full, reject the request. This is backpressure.
//    The API stays fast; it does not wait forever on a blocked send.
//
// 4. Crash after dequeue, before done
//    Worker takes a job, "crashes" (return error / panic).
//    Job is no longer in the channel. If you only used a channel,
//    the job is lost. This is why you save to JobStore first.
//    A retry loop can look for StatusProcessing jobs that are too old.
// going to implement the retry loop so that the task will be completed
// this below task has been done by me properly so i can move further
// 5. At-least-once vs exactly-once
//    Retrying a crashed job means processJob may run twice.
//    Prevent duplicate side effects with an idempotency key:
//    if StatusDone already, skip generating the PDF again.
//
// 6. Worker pool you can "see"
//    Instead of ranging on one shared jobs channel, give each worker
//    its own inbox and a shared idle channel:
//
//      idle := make(chan chan *Job, numWorkers)
//      // worker, when free: idle <- myInbox
//      // dispatcher: workerInbox := <-idle; workerInbox <- job
//
//    That models "become available" more explicitly.
//
// 7. Graceful shutdown
//    close(d.jobs) so workers exit the range loop.
//    Use WaitGroup to wait until in-flight jobs finish.
//
// Suggested print format while you debug:
//   [submit] job-3 queued
//   [worker-1] took job-3
//   [worker-1] finished job-3  (now available)
//   [poll] job-3 status=done
// ---------------------------------------------------------
