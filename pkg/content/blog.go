package content

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Post struct {
	Title       string        `yaml:"title"`
	Content     template.HTML `yaml:"content"`
	Slug        string        `yaml:"slug"`
	Description string        `yaml:"description"`
	Image       string        `yaml:"image"`
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

	if !file.IsDir() {
		bytes, _ := os.ReadFile(path)
		if err != nil {
			return err
		}

		var post Post

		err := yaml.Unmarshal(bytes, &post)
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
