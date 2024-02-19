package content

import (
	"encoding/json"
	"errors"
	"fmt"

	"io/fs"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type Post struct {
	Title   string
	Content string
	Slug    string
}

var posts = map[string]Post{}

func LoadPosts(postsDir string) error {
	cwd, _ := os.Getwd()
	err := filepath.WalkDir(filepath.Join(cwd, postsDir), LoadPost)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to load posts: %v", err))
	}

	return nil
}

func LoadPost(path string, file fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	log.Info().Msgf("Processing: %s", path)

	if !file.IsDir() {
		bytes, _ := os.ReadFile(path)
		if err != nil {
			return err
		}

		var post Post

		err := json.Unmarshal(bytes, &post)
		if err != nil {
			return err
		}

		posts[post.Slug] = post

	}
	return nil
}

func GetPost(slug string) (Post, error) {
	if _, ok := posts[slug]; !ok {
		return Post{}, errors.New("Unable to find post")
	}

	return posts[slug], nil
}

func GetPosts() map[string]Post {
	return posts
}
