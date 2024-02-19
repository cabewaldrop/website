package content

import (
	"encoding/json"
	"errors"
	"fmt"

	"io/fs"
	"os"
	"path/filepath"
)

type Recipe struct {
	Title       string
	Image       string
	Ingredients []string
	Steps       []string
	Slug        string
}

var recipes = map[string]Recipe{}

func LoadRecipes(recipeDir string) error {
	cwd, _ := os.Getwd()
	err := filepath.WalkDir(filepath.Join(cwd, recipeDir), LoadRecipe)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to load recipes: %v", err))
	}

	return nil
}

func LoadRecipe(path string, file fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if !file.IsDir() {
		bytes, _ := os.ReadFile(path)
		if err != nil {
			return err
		}

		var recipe Recipe

		err := json.Unmarshal(bytes, &recipe)
		if err != nil {
			return err
		}

		recipes[recipe.Slug] = recipe

	}
	return nil
}

func GetRecipe(slug string) (Recipe, error) {
	if _, ok := recipes[slug]; !ok {
		return Recipe{}, errors.New("Unable to find recipe")
	}

	return recipes[slug], nil
}

func GetRecipes() map[string]Recipe {
	return recipes
}
