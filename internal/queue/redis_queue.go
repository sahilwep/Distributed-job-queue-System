package queue

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisQueue struct {
	client *redis.Client
	queue  string
	dlq    string // used for dead letter queue
}

func NewRedisQueue() *RedisQueue {

	addr := os.Getenv("REDIS_ADDR")

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisQueue{
		client: rdb,
		queue:  "job_queue",
		dlq:    "dead_jobs_queue",
	}
}

func (q *RedisQueue) Push(jobID string) error {
	return q.client.LPush(ctx, q.queue, jobID).Err()
}

func (q *RedisQueue) Pop() (string, error) {
	result, err := q.client.RPop(ctx, q.queue).Result()
	if err != nil {
		return "", err
	}

	return result, nil
}

// DLQ push for failed entries:
func (q *RedisQueue) PushDead(jobID string) error {
	return q.client.LPush(ctx, q.dlq, jobID).Err()
}
