package server

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/cabewaldrop/website/pkg/content"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

const recipeDir = "content/recipes"
const blogDir = "content/blog"

type Templates struct {
	templates map[string]TypedTemplate
}

type TypedTemplate struct {
	template     *template.Template
	rootTemplate string
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates[name].template.ExecuteTemplate(w, t.templates[name].rootTemplate, data)
}

func NewTemplates() *Templates {
	templates := map[string]TypedTemplate{}

	fs, err := os.ReadDir("views")
	if err != nil {
		log.Fatal().Msgf("Unable to open the views directory! Aborting!")
	}

	for _, f := range fs {
		if f.Name() != "base.html" {
			name := f.Name()
			path := fmt.Sprintf("views/%s", name)
			partialName := fmt.Sprintf("%s_partial", name)

			fullTemplate := template.Must(template.ParseFiles(path, "views/base.html", "views/footer.html", "views/navbar.html", "views/head.html", "views/card.html"))
			partialTemplate := template.Must(template.ParseFiles(path, "views/partial.html", "views/card.html"))

			log.Info().Msgf("adding %s", name)
			templates[name] = TypedTemplate{template: fullTemplate, rootTemplate: "base.html"}

			log.Info().Msgf("adding %s", partialName)
			templates[partialName] = TypedTemplate{template: partialTemplate, rootTemplate: "partial.html"}
		}
	}

	return &Templates{
		templates: templates,
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

func ConfigureServer(e *echo.Echo) {
	err := loadContent()
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}

	templates := NewTemplates()

	e.Renderer = templates
	e.HTTPErrorHandler = handleError

	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())

	RegisterRoutes(e)
}
