package recipes_test

import (
	"testing"

	"github.com/cabewaldrop/website/pkg/recipes"
	"github.com/stretchr/testify/assert"
)

func TestLoadRecipes(t *testing.T) {
	recipes.LoadRecipes("../../content/recipes")
	recipe, err := recipes.GetRecipe("coca-cola-carnitas")
	if err != nil {
		t.Fatalf("Unexpected error getting recipe %v", err)
	}

	assert.Equal(t, "coca-cola-carnitas", recipe.Slug)
}
