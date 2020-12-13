package noid

import (
	"io/ioutil"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func songsHandler(c *fiber.Ctx) error {
	s := make([]interface{}, 0, 0)

	files, err := ioutil.ReadDir(root)
	if err != nil {
		return sendError(c, err)
	}

	for _, file := range files {
		name := file.Name()

		if !strings.HasSuffix(name, ".mp3") {
			continue
		}

		song := Songs{file.Name()}
		json, err := song.JSON()
		if err != nil {
			return sendError(c, err)
		}

		s = append(s, json)
	}

	c.Set("Access-Control-Allow-Origin", "*")
	return c.JSON(s)
}

func coverHandler(c *fiber.Ctx) error {
	song := Songs{c.Params("file")}
	cover, mime, err := song.Cover()
	if err != nil {
		return sendError(c, err)
	}

	c.Set("Content-Type", mime)
	c.Set("Cache-Control", "public, max-age=2592000")
	return c.Send(cover)
}
