package main

import (
	"embed"
	"html/template"
	"io"

	"github.com/cabewaldrop/website/pkg/recipes"
	"github.com/cabewaldrop/website/pkg/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

//go:embed views/*.html
var views embed.FS

const recipeDir = "content/recipes"

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseFS(views, "views/*.html")),
	}
}

func main() {
	err := recipes.LoadRecipes(recipeDir)
	if err != nil {
		log.Fatal().Msgf("Unable to load the recipes. Check yoself before you wreck yoself: %v", err)
	}

	e := echo.New()
	e.Renderer = NewTemplates()
	e.Use(middleware.Logger())

	routes.RegisterRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}
