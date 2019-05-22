package main

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	// endpont to redirect shortlinks
	r.GET("/:slug", func(c *gin.Context) {
		slug := c.Param("slug")

		val, err := redisClient.Get(slug).Result()

		if err == redis.Nil {
			c.JSON(404, gin.H{
				"status": "nothing found",
			})
		} else if err != nil {
			panic(err)
		} else {
			storeRequest(redisClient)
			c.Redirect(301, val)
		}
	})

	// endpoint to create shortlink
	// // message := c.PostForm("url")
	// // nick := c.DefaultPostForm("nick", "anonymous")
	r.POST("/create/:url", func(c *gin.Context) {

		// loop to make a unique key
		key := ""
		exists := int64(0)
		success := int64(1)
		for exists == success {
			key := RandASCIIBytes(4)
			exists = redisClient.Exists(string(key)).Val()
		}

		url := c.Param("url")
		err := redisClient.Set(key, url, 0).Err()
		if err != nil {
			panic(err)
		}

		fullUrl := fmt.Sprintf("https://myweb.com/%s", key)

		// see if key exists in redis
		// save
		//
		c.JSON(202, gin.H{
			"status":     "created",
			"shorturl":   fullUrl,
			"accessToke": "randomAccessToken",
		})
	})

	// endpoint to get stats
	r.POST("/read/:access_token", func(c *gin.Context) {
		// token := c.Param("access_token")
		// url := c.Param("shorturl")

		// gets token from redis
		// checks to confirm shorturl is the same as in hash
		// returns hash data
		c.JSON(200, gin.H{
			"status": "ok",
			"data":   "data",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandASCIIBytes(n int) []byte {
	output := make([]byte, n)

	// We will take n bytes, one byte for each character of output.
	randomness := make([]byte, n)

	// read all random
	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}

	l := len(letterBytes)
	// fill output
	for pos := range output {
		// get random item
		random := uint8(randomness[pos])

		// random % 64
		randomPos := random % uint8(l)

		// put into output
		output[pos] = letterBytes[randomPos]
	}

	return output
}

func storeRequest(client *redis.Client) {
	pipe := client.Pipeline()
	// increment total views
	pipe.Incr("key")
	// store ip in unique set
	pipe.SAdd("key", "ip_value")
	// store countries
	pipe.HIncrBy("key", "field", 1)
	// store referral
	pipe.HIncrBy("key", "field", 1)

	_, err := pipe.Exec()
	if err != nil {
		panic(err)
	}
}

func getData(client *redis.Client) map[string]string {
	pipe := client.Pipeline()

	// increment total views
	keyViews := pipe.Get("key")
	// store ip in unique set
	keyUniqueViews := pipe.SCard("key")
	// get countires
	keyCountries := pipe.HGetAll("key")
	// get referral
	keyReferals := pipe.HGetAll("key")
	res, _ := keyUniqueViews.Result()

	return results
}
