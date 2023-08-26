package store

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"os"
)

var RDB *redis.Client

type Config struct {
	Addr     string `json:"Addr"`
	Password string `json:"Password"`
	DB       int    `json:DB`
}

func InitRedisClient() (err error) {
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		return err
	}
	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return err
	}
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = RDB.Ping().Result()
	if err != nil {
		return err
	}
	fmt.Println("Redis client initiated")
	return nil
}
