package noid

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/bogem/id3v2"
	"github.com/go-redis/redis/v8"
)

// Songs .
type Songs struct {
	file string
}

// JSON return the json object to export
func (s *Songs) JSON() (map[string]string, error) {
	title, err := s.Title()
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"id":    s.file,
		"file":  "/song/" + s.file,
		"title": title,
		"cover": "/cover/" + s.file,
	}, nil
}

// Title return the song title
func (s *Songs) Title() (string, error) {
	ctx := context.Background()
	rdb := redis.NewClient(rdbConf)
	defer rdb.Close()

	title, err := rdb.HGet(ctx, s.file, "title").Result()

	if err == redis.Nil {
		t, _, _, err := s.setCache(ctx, rdb)

		if err != nil {
			return "", err
		}

		return t, nil
	} else if err != nil {
		return "", err
	}

	return title, nil
}

// Cover return the song cover and the cover mimetype
func (s *Songs) Cover() ([]byte, string, error) {
	ctx := context.Background()
	rdb := redis.NewClient(rdbConf)
	defer rdb.Close()

	field, err := rdb.HMGet(ctx, s.file, "cover", "mime").Result()

	if err == redis.Nil || len(field) != 2 || field[0] == nil || field[1] == nil {
		_, cover, mime, err := s.setCache(ctx, rdb)

		if err != nil {
			return nil, "", err
		}

		return cover, mime, nil
	} else if err != nil {
		return nil, "", err
	}

	cover, ok1 := field[0].(string)
	mime, ok2 := field[1].(string)

	if !ok1 || !ok2 {
		return nil, "", errors.New("Cannot read field value")
	}

	return []byte(cover), mime, nil
}

func (s *Songs) setCache(ctx context.Context, rdb *redis.Client) (string, []byte, string, error) {
	tag, err := id3v2.Open(filepath.Join(root, s.file), id3v2.Options{Parse: true})
	if err != nil {
		return "", nil, "", err
	}
	defer tag.Close()

	title := tag.Title()

	pictures := tag.GetFrames(tag.CommonID("Attached picture"))
	pic, ok := pictures[0].(id3v2.PictureFrame)
	if !ok {
		return "", nil, "", errors.New("Cannot load cover")
	}

	rdb.HSet(ctx, s.file, map[string]interface{}{
		"title": title,
		"cover": pic.Picture,
		"mime":  pic.MimeType,
	})

	return title, pic.Picture, pic.MimeType, nil
}
