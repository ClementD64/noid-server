package noid

import (
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func loadRedis() {
	url, ok := os.LookupEnv("NOID_REDIS")

	if !ok {
		url = "redis://localhost:6379/0"
	}

	conf, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}
	rdbConf = conf
}

func sendError(c *fiber.Ctx, err error) error {
	log.Print(err)
	c.Status(500)

	return c.JSON(fiber.Map{
		"status": 500,
		"error":  err.Error(),
	})
}
