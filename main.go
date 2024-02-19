package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/cabewaldrop/website/pkg/recipes"
	"github.com/cabewaldrop/website/pkg/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

//go:embed views
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
		templates: template.Must(template.ParseFS(views, "views/*.html", "views/blog/*.html")),
	}
}

func logRenderedTemplates() {
	templates := strings.Split(NewTemplates().templates.DefinedTemplates(), ",")
	for _, templ := range templates {
		// Get rid of spaces and quotes to make it easier to read
		formatted := strings.ReplaceAll(strings.TrimSpace(templ), "\"", "")

		// We don't care about the file names so ignore them
		if !strings.Contains(formatted, ".html") {
			log.Debug().Msgf("Parsed: %s", formatted)
		}
	}
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

func main() {
	err := recipes.LoadRecipes(recipeDir)
	if err != nil {
		log.Fatal().Msgf("Unable to load the recipes. Check yoself before you wreck yoself: %v", err)
	}

	e := echo.New()
	logRenderedTemplates()
	e.Renderer = NewTemplates()

	e.HTTPErrorHandler = handleError

	e.Use(middleware.Logger())

	routes.RegisterRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}
