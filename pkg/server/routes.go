package server

import (
	"fmt"
	"os"
	"time"

	"github.com/cabewaldrop/website/pkg/content"
	"github.com/labstack/echo/v4"
)

var baseURL = os.Getenv("BASE_URL")

type HealthCheckResponse struct {
	Status string
}

type OpenGraph struct {
	Title       string
	Description string
	Image       string
}

type Metadata struct {
	Date      string
	Title     string
	OpenGraph OpenGraph
}

type IndexParams struct {
	Meta Metadata
}

type ContentLink struct {
	Title string
	Link  string
}

type RecipeIndexParams struct {
	Meta  Metadata
	Links []ContentLink
}

type RecipeDetailParams struct {
	Meta   Metadata
	Recipe content.Recipe
}

type PostDetailParams struct {
	Meta Metadata
	Post content.Post
}

type BlogIndexParams struct {
	Meta  Metadata
	Links []ContentLink
}

func RegisterRoutes(e *echo.Echo) *echo.Echo {
	now := time.Now().Format("01/02/2006")

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", IndexParams{Meta: Metadata{Date: now, Title: "Cabe Waldrop"}})
	})

	e.File("/favicon.ico", "static/favicon.ico")

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(200, HealthCheckResponse{Status: "ok"})
	})

	e.GET("/recipes", func(c echo.Context) error {
		recipes := content.GetRecipes()
		links := []ContentLink{}
		for _, recipe := range recipes {
			links = append(links,
				ContentLink{Title: recipe.Title, Link: fmt.Sprintf("/recipes/%s", recipe.Slug)})
		}

		return c.Render(200, "recipe-index", RecipeIndexParams{Meta: Metadata{Date: now, Title: "Recipes"}, Links: links})
	})

	e.GET("/recipes/:slug", func(c echo.Context) error {
		slug := c.Param("slug")
		recipe, err := content.GetRecipe(slug)
		if err != nil {
			return &echo.HTTPError{
				Code:    404,
				Message: "Could not find recipe",
			}
		}

		return c.Render(200, "recipe-detail",
			RecipeDetailParams{
				Meta: Metadata{
					Date:  now,
					Title: slug,
					OpenGraph: OpenGraph{
						Title:       recipe.Title,
						Description: recipe.Title,
						Image:       fmt.Sprintf("%s%s", baseURL, recipe.Image),
					},
				},
				Recipe: recipe,
			},
		)
	})

	e.GET("/blog", func(c echo.Context) error {
		posts := content.GetPosts()
		links := []ContentLink{}
		for _, post := range posts {
			links = append(links,
				ContentLink{Title: post.Title, Link: fmt.Sprintf("/blog/%s", post.Slug)})
		}

		return c.Render(200, "blog-index", BlogIndexParams{Meta: Metadata{Date: now, Title: "Blog"}, Links: links})
	})

	e.GET("/blog/:slug", func(c echo.Context) error {
		slug := c.Param("slug")
		post, err := content.GetPost(slug)
		if err != nil {
			return &echo.HTTPError{
				Code:    404,
				Message: "Could not find blog post",
			}
		}

		return c.Render(200, "post", PostDetailParams{Meta: Metadata{Date: now, Title: slug}, Post: post})
	})

	// Used in dev, but actually intercepted and served by CDN in production
	// See [[statics]] in fly.toml
	e.Static("/static", "static")

	return e
}
