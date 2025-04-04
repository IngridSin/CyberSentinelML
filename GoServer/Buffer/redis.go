package Buffer

import (
	"context"
	"encoding/json"
	"fmt"
	"goServer/models"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

const redisQueueKey = "packet_queue"

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	go producer(rdb)
	go consumer(rdb)

	select {}
}

// packet producer
func producer(rdb *redis.Client) {
	for i := 0; i < 10; i++ {
		packet := models.PacketFlow{
			FlowID:        fmt.Sprintf("flow-%d", i),
			SourceIP:      "10.0.0.1",
			DestinationIP: "10.0.0.2",
			Protocol:      6,
			StartTime:     time.Now(),
		}

		data, _ := json.Marshal(packet)
		err := rdb.RPush(ctx, redisQueueKey, data).Err()
		if err != nil {
			log.Printf("Failed to enqueue: %v", err)
		}
		log.Println("Enqueued:", packet.FlowID)
		time.Sleep(500 * time.Millisecond)
	}
}

// Consumer pulls from Redis and processes
func consumer(rdb *redis.Client) {
	for {
		result, err := rdb.BLPop(ctx, 5*time.Second, redisQueueKey).Result()
		if err == redis.Nil {
			continue // timeout, retry
		} else if err != nil {
			log.Printf("Redis error: %v", err)
			continue
		}

		var packet models.PacketFlow
		err = json.Unmarshal([]byte(result[1]), &packet)
		if err != nil {
			log.Printf("Unmarshal error: %v", err)
			continue
		}

		log.Printf("Processing flow from Redis: %+v\n", packet)
	}
}
