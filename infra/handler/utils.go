package handler

import (
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/env"
	"net/url"
)

func buildRecipeImageUrls(images []domain.RecipeImage) ([]url.URL, error) {
	urls := make([]url.URL, len(images))
	for i, image := range images {
		path, err := url.JoinPath(env.MustGet("BASE_URL"), "images", image.Path)
		if err != nil {
			return nil, err
		}
		imageUrl, err := url.Parse(path)
		if err != nil {
			return nil, err
		}
		urls[i] = *imageUrl
	}
	return urls, nil
}
