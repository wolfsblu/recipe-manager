package handler

import (
	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/env"
	"regexp"
)

func NewUploadHandler(recipes *domain.RecipeService) (*tusd.Handler, error) {
	store := filestore.New(env.MustGet("UPLOAD_FOLDER"))
	locker := filelocker.New(env.MustGet("UPLOAD_FOLDER"))

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)
	locker.UseIn(composer)

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              "/api/uploads/",
		StoreComposer:         composer,
		NotifyCreatedUploads:  true,
		NotifyCompleteUploads: true,
		Cors: &tusd.CorsConfig{
			AllowOrigin:      regexp.MustCompile(env.MustGet("CORS_ORIGIN")),
			AllowCredentials: true,
			AllowMethods:     tusd.DefaultCorsConfig.AllowMethods,
			AllowHeaders:     tusd.DefaultCorsConfig.AllowHeaders,
			MaxAge:           tusd.DefaultCorsConfig.MaxAge,
			ExposeHeaders:    tusd.DefaultCorsConfig.ExposeHeaders,
		},
	})

	go registerEventListeners(handler, recipes)

	if err != nil {
		return nil, err
	}
	return handler, nil
}

func registerEventListeners(handler *tusd.Handler, recipes *domain.RecipeService) {
	for {
		select {
		case event := <-handler.CompleteUploads:
			onUploadCompleted(recipes, event)
		case event := <-handler.CreatedUploads:
			onUploadCreated(recipes, event)
		}
	}
}

func onUploadCreated(recipes *domain.RecipeService, event tusd.HookEvent) {
	// If user is not authenticated
	//event.Upload.StopUpload()
}

func onUploadCompleted(recipes *domain.RecipeService, event tusd.HookEvent) {
	// Move file to correct folder
}
