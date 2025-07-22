package handler

import (
	"github.com/wolfsblu/go-chef/infra/env"
	"net/url"

	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/config"
)

type urlBuilder struct {
	baseURL string
}

func newURLBuilder() *urlBuilder {
	return &urlBuilder{
		baseURL: env.MustGet("BASE_URL"),
	}
}

// buildImageURL constructs a full URL for an image path
func (b *urlBuilder) buildImageURL(imagePath string) (*url.URL, error) {
	path, err := url.JoinPath(b.baseURL, config.ImagesPathPrefix, imagePath)
	if err != nil {
		return nil, err
	}
	return url.Parse(path)
}

// buildRecipeImageURLs constructs URLs for all recipe images
func (b *urlBuilder) buildRecipeImageURLs(images []domain.RecipeImage) ([]url.URL, error) {
	urls := make([]url.URL, len(images))
	for i, image := range images {
		imageURL, err := b.buildImageURL(image.Path)
		if err != nil {
			return nil, err
		}
		urls[i] = *imageURL
	}
	return urls, nil
}

// buildAPIURL constructs API endpoint URLs
func (b *urlBuilder) buildAPIURL(endpoint string) (*url.URL, error) {
	path, err := url.JoinPath(b.baseURL, config.APIPathPrefix, endpoint)
	if err != nil {
		return nil, err
	}
	return url.Parse(path)
}
