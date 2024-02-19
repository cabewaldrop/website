package content_test

import (
	"testing"

	"github.com/cabewaldrop/website/pkg/content"
	"github.com/stretchr/testify/assert"
)

func TestLoadRecipes(t *testing.T) {
	content.LoadRecipes("../../content/recipes")
	recipe, err := content.GetRecipe("coca-cola-carnitas")
	if err != nil {
		t.Fatalf("Unexpected error getting recipe %v", err)
	}

	assert.Equal(t, "coca-cola-carnitas", recipe.Slug)
}
