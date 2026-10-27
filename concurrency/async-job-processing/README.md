# Asynchronous Job Processing

A small Go learning simulation of how a backend can accept work quickly and process it in the background. The example uses report-to-PDF jobs, a buffered channel as a bounded queue, and goroutines as workers.

## Request and worker flow

```text
client -> SubmitJob -> save job -> enqueue -> return job ID
                                      |
                                      v
                             worker takes the job
                                      |
                         process -> update job status

client -> GetJob(job ID) -> queued / processing / done / failed
```

`SubmitJob` is the HTTP-handler equivalent: it does not wait for PDF generation. The client receives a job ID and can poll `GetJob` for its status. When a worker finishes a job, it loops back to the shared channel and waits for more work; that is how this worker pool becomes available again.

## Concepts demonstrated

- **Background work:** `Dispatcher.Start` starts worker goroutines and returns; workers consume jobs independently of the submitter.
- **Bounded queue and backpressure:** a non-blocking channel send avoids making submission wait for a full queue. A job that cannot be enqueued is retained in the in-memory retry map for a later retry round.
- **Job status polling:** jobs are saved in `JobStore` so callers can look up their progress using the returned ID.
- **Failure and retry simulation:** a worker can randomly crash after taking a job. The simulation records that job for a later retry, illustrating why dequeuing work and completing it are separate failure points.
- **At-least-once processing:** retrying may run a job more than once. Real systems need idempotent job handling to avoid duplicate side effects.

## Run

From this directory:

```sh
go run .
```

The demo starts three workers, submits ten jobs, polls their statuses, then makes a limited number of retry rounds. Processing delays and simulated crashes are intentional, so output and completion timing can vary.

## Scope and limitations

This is a concurrency simulation, not a running HTTP service: `SubmitJob` and `GetJob` stand in for endpoints, but there is no HTTP server or request/response encoding. The queue and store are in memory, so they do not survive a process restart. The retry loop is manually driven by the demo rather than a durable scheduler. These tradeoffs keep the exercise focused on the async request-to-worker flow; they are not production guarantees.