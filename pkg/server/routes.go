package server

import (
	"fmt"
	"time"

	"github.com/cabewaldrop/website/pkg/content"
	"github.com/labstack/echo/v4"
)

type HealthCheckResponse struct {
	Status string
}

type IndexParams struct {
	Date string
}

type ContentLink struct {
	Title string
	Link  string
}

type RecipeIndexParams struct {
	Links []ContentLink
	Date  string
}

type RecipeDetailParams struct {
	Recipe content.Recipe
	Date   string
}

type PostDetailParams struct {
	Post content.Post
	Date string
}

type BlogIndexParams struct {
	Date  string
	Links []ContentLink
}

func RegisterRoutes(e *echo.Echo) *echo.Echo {
	now := time.Now().Format("01/02/2006")

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", IndexParams{Date: now})
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

		return c.Render(200, "recipe-index", RecipeIndexParams{Date: now, Links: links})
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

		return c.Render(200, "recipe-detail", RecipeDetailParams{Date: now, Recipe: recipe})
	})

	e.GET("/blog", func(c echo.Context) error {
		posts := content.GetPosts()
		links := []ContentLink{}
		for _, post := range posts {
			links = append(links,
				ContentLink{Title: post.Title, Link: fmt.Sprintf("/blog/%s", post.Slug)})
		}

		return c.Render(200, "blog-index", BlogIndexParams{Date: now, Links: links})
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

		return c.Render(200, "post", PostDetailParams{Date: now, Post: post})
	})

	e.Static("/static", "static")

	return e
}
