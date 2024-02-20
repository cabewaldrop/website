package server

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/cabewaldrop/website/pkg/content"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

const recipeDir = "content/recipes"
const blogDir = "content/blog"

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates(fs embed.FS) *Templates {
	return &Templates{
		templates: template.Must(template.ParseFS(fs, "views/*.html")),
	}
}

// Load all content into a map for easy retrieval.
func loadContent() error {
	err := content.LoadRecipes(recipeDir)
	if err != nil {
		return errors.New(
			fmt.Sprintf("Unable to load the recipes. Check yoself before you wreck yoself: %v", err),
		)
	}

	err = content.LoadPosts(blogDir)
	if err != nil {
		return errors.New(
			fmt.Sprintf("Unable to load the blog posts. Check yoself before you wreck yoself: %v", err),
		)
	}

	return nil
}

// handleError renders an error page named after the status code
// i.e. 404.html, 500.html, etc
func handleError(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	c.Logger().Error(err)
	errorPage := fmt.Sprintf("%d", code)

	if err := c.Render(code, errorPage, nil); err != nil {
		c.Logger().Error(err)
	}
}

func ConfigureServer(e *echo.Echo, fs embed.FS) {
	err := loadContent()
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	templates := NewTemplates(fs)

	e.Renderer = templates
	e.HTTPErrorHandler = handleError

	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())

	RegisterRoutes(e)
}
