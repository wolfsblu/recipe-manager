package urlbuilder

import (
	"github.com/wolfsblu/go-chef/infra/env"
	"net/url"

	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/config"
)

type Builder struct {
	baseURL string
}

func NewURLBuilder() *Builder {
	return &Builder{
		baseURL: env.MustGet("BASE_URL"),
	}
}

// BuildImageURL constructs a full URL for an image path
func (b *Builder) BuildImageURL(imagePath string) (*url.URL, error) {
	path, err := url.JoinPath(b.baseURL, config.ImagesPathPrefix, imagePath)
	if err != nil {
		return nil, err
	}
	return url.Parse(path)
}

// BuildRecipeImageURLs constructs URLs for all recipe images
func (b *Builder) BuildRecipeImageURLs(images []domain.RecipeImage) ([]url.URL, error) {
	urls := make([]url.URL, len(images))
	for i, image := range images {
		imageURL, err := b.BuildImageURL(image.Path)
		if err != nil {
			return nil, err
		}
		urls[i] = *imageURL
	}
	return urls, nil
}

// BuildAPIURL constructs API endpoint URLs
func (b *Builder) BuildAPIURL(endpoint string) (*url.URL, error) {
	path, err := url.JoinPath(b.baseURL, config.APIPathPrefix, endpoint)
	if err != nil {
		return nil, err
	}
	return url.Parse(path)
}
