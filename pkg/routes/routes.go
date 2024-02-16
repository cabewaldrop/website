package routes

import (
	"time"

	"github.com/cabewaldrop/website/pkg/recipes"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type HealthCheckResponse struct {
	Status string
}

type IndexParams struct {
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

	e.GET("/recipe/:slug", func(c echo.Context) error {
		slug := c.Param("slug")
		recipe, err := recipes.GetRecipe(slug)
		if err != nil {
			//TODO
		}

		log.Info().Msgf("Recipe is: %v", recipe)
		return c.Render(200, "recipe-detail", recipe)
	})

	e.Static("/images", "images")

	return e
}
