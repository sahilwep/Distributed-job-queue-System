# Distributed Job Queue System

A Distributed background job processing system build with **Go, Gin, Redis, and PostgreSQL** where tasks are queued and executed asynchronously by worker services. Redis acts as the job queue, PostgreSQL stores persistent job metadata, and worker pools process jobs concurrently. This design keeps APIs responsive while ensuring scalability, reliability, and fault tolerance.

## Key Features:
### Scalability
- The system is designed to be horizontally scalable. Worker pools process jobs concurrently using goroutines, and additional worker instances can be added to handle increased workload.

### Decoupled Components
- API, queue, database, and workers are separated, enabling independent scaling and maintenance. Redis acts as the job broker while PostgreSQL maintains persistent job metadata.

### Resilient Job Processing
- Jobs are stored in PostgreSQL, ensuring persistence even if Redis or workers restart. The system can recover job states and continue processing without data loss.

### Retry & Dead Letter Queue
- Failed jobs are automatically retried up to a configurable limit. If retries are exhausted, the job is moved to a Dead Letter Queue for further inspection and debugging.

### Concurrent  Workers Execution
- Workers process jobs using Go goroutines, allowing multiple tasks to execute simultaneously while maintaining efficient resource utilization.


## HLD:
![HLD](assets/hld.png )

## Architecture

```
                +-------------+
                |   Client    |
                +-------------+
                       |
                       v
                +-------------+
                |   API       |
                |   (Gin)     |
                +-------------+
                       |
                       v
                +-------------+
                | PostgreSQL  |
                | (Job State) |
                +-------------+
                       |
                       v
                +-------------+
                |   Redis     |
                |   Queue     |
                +-------------+
                       |
                       v
                +-------------+
                | Worker Pool |
                | (Goroutines)|
                +-------------+
                       |
                       v
                Job Processing
```

## Running the Project

```sh
$ git clone 
$ cd distributed_job_queue_system
$ docker compose up --build
```

This will start:

* API Server
* Worker Service
* PostgreSQL
* Redis

## Test the API

### Health Check

```
GET http://localhost:8080/health
```

Response:

```
{
  "status": "ok"
}
```


### Create Job

```
POST http://localhost:8080/jobs
```

Example request:

```
curl -X POST http://localhost:8080/jobs \
-H "Content-Type: application/json" \
-d '{
"type": "email",
"payload": {
"to": "sahilwep@gmail.com"
}
}'
```

Example response:

```
{
  "job_id": "f490adfe-b5a6-41c3-b840-54ae0753a4d6",
  "job_status": "queued"
}
```


### Check Job Status

```
GET http://localhost:8080/jobs/{job_id}
```

Example:

```
curl http://localhost:8080/jobs/f490adfe-b5a6-41c3-b840-54ae0753a4d6
```

Response:

```
{
  "id": "f490adfe-b5a6-41c3-b840-54ae0753a4d6",
  "type": "email",
  "status": "completed"
}
```


## Job Lifecycle

```
queued
   ↓
processing
   ↓
completed
```

If a job fails:

```
failed → retry → retry → dead_letter_queue
```

## Containers Inspection:
### Redis:
```sh
$ docker exec -it job_queue_redis redis-cli

127.0.0.1:6379> LRANGE job_queue 0 -1
(empty array)
127.0.0.1:6379> LRANGE dead_jobs_queue 0 -1
1) "8c0213b6-d463-4add-981e-e9605c6af612"
2) "a6de4539-f1bb-4d80-974f-b6f8d2643948"
3) "b2f3851a-0110-4bcf-a436-c5e48224c306"
127.0.0.1:6379> 
```

### postgres:
```sh
$ docker exec -it job_queue_postgres psql -U admin -d job_queue

----\DB# SELECT * FROM jobs;
----\DB# \q
```


# 🧑‍💻 Thanks Giving 🌟


