package buffer

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"goServer/config"
	"goServer/models"
	"log"
)

var rdb *redis.Client

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: config.RedisAddr + ":6379",
		DB:   0,
	})
}

func InsertPacket(packet *models.PacketFlow) {
	data, err := json.Marshal(packet)
	if err != nil {
		log.Println("Failed to serialize packet:", err)
		return
	}

	ctx := context.Background()
	err = rdb.RPush(ctx, "packet_queue", data).Err()
	if err != nil {
		log.Println("Failed to push to Redis:", err)
	} else {
		log.Printf("Pushed flow %s to Redis\n", packet.FlowID)
	}
}
