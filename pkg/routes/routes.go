package routes

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

type RecipeLink struct {
	Title string
	Link  string
}

type RecipeIndexParams struct {
	Links []RecipeLink
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
	Date string
}

func RegisterRoutes(e *echo.Echo) *echo.Echo {
	now := time.Now().Format("01/02/2006")
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", IndexParams{Date: now})
	})

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(200, HealthCheckResponse{Status: "ok"})
	})

	e.GET("/recipes", func(c echo.Context) error {
		recipes := content.GetRecipes()
		links := []RecipeLink{}
		for _, recipe := range recipes {
			links = append(links,
				RecipeLink{Title: recipe.Title, Link: fmt.Sprintf("/recipes/%s", recipe.Slug)})
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

	e.GET("/api", func(c echo.Context) error {
		return c.Render(200, "federated-autonomy", nil)
	})

	e.Static("/static", "static")

	return e
}
