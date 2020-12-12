package noid

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

var ctx = context.Background()
var rdb *redis.Client = nil
var root string = ""

// New create a new Noid server
func New(musicPath string) {
	root = musicPath
	rdb = getRedis()

	app := fiber.New()
	app.Use(compress.New())

	app.Static("/", "./public", fiber.Static{
		ByteRange: true,
	})
	app.Static("/song", root, fiber.Static{
		ByteRange: true,
	})
	app.Get("/songs", songsHandler)
	app.Get("/cover/:file", coverHandler)

	app.Listen(":3000")
}
