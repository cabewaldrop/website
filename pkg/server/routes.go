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

type Card struct {
	Title       string
	URL         string
	Description string
	Image       string
}

type RecipeIndexParams struct {
	Meta  Metadata
	Cards []Card
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
	Cards []Card
}

// Determine if the request is for a partial. If it is return the partial suffix
func partial(c echo.Context) string {
	if c.QueryParam("partial") != "" {
		return "_partial"
	}
	return ""
}

func RegisterRoutes(e *echo.Echo) *echo.Echo {
	now := time.Now().Format("01/02/2006")

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index.html", IndexParams{Meta: Metadata{Date: now, Title: "Cabe Waldrop"}})
	})

	e.File("/favicon.ico", "static/favicon.ico")

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(200, HealthCheckResponse{Status: "ok"})
	})

	e.GET("/recipes", func(c echo.Context) error {
		recipes := content.GetRecipes()
		cards := []Card{}
		for _, recipe := range recipes {
			cards = append(
				cards,
				Card{
					Title:       recipe.Title,
					Description: recipe.Description,
					URL:         fmt.Sprintf("/recipes/%s", recipe.Slug),
					Image:       recipe.Image,
				},
			)
		}

		suffix := partial(c)
		template := fmt.Sprintf("recipe-index.html%s", suffix)
		return c.Render(200, template,
			RecipeIndexParams{
				Meta:  Metadata{Date: now, Title: "Recipes"},
				Cards: cards,
			},
		)
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

		suffix := partial(c)
		template := fmt.Sprintf("recipe-detail.html%s", suffix)
		return c.Render(200, template,
			RecipeDetailParams{
				Meta: Metadata{
					Date:  now,
					Title: slug,
					OpenGraph: OpenGraph{
						Title:       recipe.Title,
						Description: recipe.Description,
						Image:       fmt.Sprintf("%s%s", baseURL, recipe.Image),
					},
				},
				Recipe: recipe,
			},
		)
	})

	e.GET("/blog", func(c echo.Context) error {
		posts := content.GetPosts()
		cards := []Card{}
		for _, post := range posts {
			cards = append(cards,
				Card{
					Title:       post.Title,
					URL:         fmt.Sprintf("/blog/%s", post.Slug),
					Image:       post.Image,
					Description: post.Description,
				},
			)
		}

		suffix := partial(c)
		template := fmt.Sprintf("blog-index.html%s", suffix)
		return c.Render(200, template, BlogIndexParams{Meta: Metadata{Date: now, Title: "Blog"}, Cards: cards})
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

		suffix := partial(c)
		template := fmt.Sprintf("post.html%s", suffix)
		return c.Render(200, template, PostDetailParams{
			Meta: Metadata{
				Date:  now,
				Title: slug,
				OpenGraph: OpenGraph{
					Title:       post.Title,
					Description: post.Description,
					Image:       fmt.Sprintf("%s%s", baseURL, post.Image),
				},
			}, Post: post})
	})

	// Used in dev, but actually intercepted and served by CDN in production
	// See [[statics]] in fly.toml
	e.Static("/static", "static")

	return e
}
