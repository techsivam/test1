package main

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

type Data struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

var rh *rejson.Handler

func init() {
	var rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	rh = rejson.NewReJSONHandler()
	rh.SetGoRedisClient(rdb)
}

func GetRedis(c *gin.Context) {

	tenant := c.Param("tenant")
	key := c.Param("key")

	// Construct the Redis key using the tenant ID and key
	redisKey := tenant + ":" + key

	// Retrieve the JSON value from Redis as bytes
	jsonBytes, err := rh.JSONGet(redisKey, ".")
	if err != nil {
		c.JSON(404, gin.H{"error": "Key not found"})
		return
	}

	// Perform a type assertion to convert the bytes to []byte
	bytes, ok := jsonBytes.([]byte)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid data type"})
		return
	}

	// Unmarshal the JSON bytes into a map[string]interface{}
	var jsonData map[string]interface{}
	err = json.Unmarshal(bytes, &jsonData)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error decoding JSON"})
		return
	}
	fmt.Println("JSON READ: ", jsonData)
	// Return the JSON data as a response
	c.JSON(200, jsonData)
}

func PutRedis(c *gin.Context) {
	tenant := c.Param("tenant")
	var data Data

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	redisKey := tenant + ":" + data.Key
	_, err := rh.JSONSet(redisKey, ".", data.Value)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("JSON WRITE: ", data)
	c.JSON(200, gin.H{"status": "success"})
}

func main() {
	router := gin.Default()
	router.GET("/:tenant/:key", GetRedis)
	router.POST("/:tenant", PutRedis)

	router.Run(":8082")
}
